// Copyright (C) 2026 The go-mdns Authors All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mdns

import (
	"net"
	"slices"
	"strings"
	"sync"

	"github.com/cybergarage/go-mdns/mdns/dns"
	"github.com/cybergarage/go-mdns/mdns/transport"
)

// legacyUnicastTTL caps the TTL of the records in a reply to a legacy
// unicast query (RFC 6762, 6.7).
const legacyUnicastTTL = 10

// responder holds the services this node publishes and answers the queries
// for them. A service whose names are being probed is held as a probe, and
// it is not answered for until the probe finds its names unique.
type responder struct {
	mutex    sync.RWMutex
	services []*LocalService
	probing  []*probe
}

func newResponder() *responder {
	return &responder{
		mutex:    sync.RWMutex{},
		services: []*LocalService{},
		probing:  []*probe{},
	}
}

// register adds svc, replacing a service with the same instance name.
func (r *responder) register(svc *LocalService) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.registerLocked(svc)
}

func (r *responder) registerLocked(svc *LocalService) {
	for i, s := range r.services {
		if strings.EqualFold(s.FullName(), svc.FullName()) {
			r.services[i] = svc
			return
		}
	}
	r.services = append(r.services, svc)
}

// deregister removes the service with the instance name of svc, and
// returns the removed one. A probe of the instance is canceled.
func (r *responder) deregister(svc *LocalService) (*LocalService, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.cancelProbesLocked(svc.FullName(), ErrDeregistered)
	for i, s := range r.services {
		if strings.EqualFold(s.FullName(), svc.FullName()) {
			r.services = append(r.services[:i], r.services[i+1:]...)
			return s, true
		}
	}
	return nil, false
}

// hostHeldLocked reports whether a published service holds host.
func (r *responder) hostHeldLocked(host string) bool {
	for _, s := range r.services {
		if strings.EqualFold(s.HostName(), host) {
			return true
		}
	}
	return false
}

// namesToProbe returns the names of svc which this node does not hold yet:
// a name which a published service holds is not probed again, such as when
// a service is registered again to update its TXT strings.
func (r *responder) namesToProbe(svc *LocalService) []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	names := []string{}
	held := slices.ContainsFunc(r.services, func(s *LocalService) bool {
		return strings.EqualFold(s.FullName(), svc.FullName())
	})
	if !held {
		names = append(names, svc.FullName())
	}
	if !r.hostHeldLocked(svc.HostName()) {
		names = append(names, svc.HostName())
	}
	return names
}

func (r *responder) cancelProbesLocked(fullName string, err error) {
	r.probing = slices.DeleteFunc(r.probing, func(p *probe) bool {
		if strings.EqualFold(p.svc.FullName(), fullName) {
			p.cancel(err)
			return true
		}
		return false
	})
}

// addProbe adds p, canceling an earlier probe of the same instance, which a
// later registration replaces.
func (r *responder) addProbe(p *probe) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.cancelProbesLocked(p.svc.FullName(), ErrDeregistered)
	r.probing = append(r.probing, p)
}

// removeProbe removes p.
func (r *responder) removeProbe(p *probe) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.probing = slices.DeleteFunc(r.probing, func(q *probe) bool { return q == p })
}

// establish publishes the service of p, whose names were found unique. It
// reports false when p was canceled meanwhile.
func (r *responder) establish(p *probe) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if !slices.Contains(r.probing, p) {
		return false
	}
	r.probing = slices.DeleteFunc(r.probing, func(q *probe) bool { return q == p })
	r.registerLocked(p.svc)
	return true
}

// beginReprobe stops answering for svc, this very copy, and returns a probe
// of its names, such as after a conflicting record is received (RFC 6762, 9)
// or when the server starts again. It reports false when svc is not
// published.
func (r *responder) beginReprobe(svc *LocalService) (*probe, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if !slices.Contains(r.services, svc) {
		return nil, false
	}
	r.services = slices.DeleteFunc(r.services, func(s *LocalService) bool { return s == svc })
	names := []string{svc.FullName()}
	if !r.hostHeldLocked(svc.HostName()) {
		names = append(names, svc.HostName())
	}
	p := newProbe(svc, names)
	r.probing = append(r.probing, p)
	return p, true
}

// restore puts the service of p back without publishing it again, such as
// when the server stops while it probes a registered service: the service
// is probed again when the server starts.
func (r *responder) restore(p *probe) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if !slices.Contains(r.probing, p) {
		return
	}
	r.probing = slices.DeleteFunc(r.probing, func(q *probe) bool { return q == p })
	r.registerLocked(p.svc)
}

// probes returns the probes in progress.
func (r *responder) probes() []*probe {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return slices.Clone(r.probing)
}

// allServices returns the published services and the services being
// probed.
func (r *responder) allServices() []*LocalService {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	services := slices.Clone(r.services)
	for _, p := range r.probing {
		services = append(services, p.svc)
	}
	return services
}

// ifRegistered calls fn while svc, this very copy, is registered, and
// reports whether it did. fn runs under the read lock, so a service being
// deregistered or replaced waits for it, and nothing fn sends can follow
// the goodbye records of a deregistration.
func (r *responder) ifRegistered(svc *LocalService, fn func()) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if slices.Contains(r.services, svc) {
		fn()
		return true
	}
	return false
}

func (r *responder) localServices() []*LocalService {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	services := make([]*LocalService, len(r.services))
	copy(services, r.services)
	return services
}

// isLegacyUnicastQuery reports whether query was sent from a port other
// than 5353, by a resolver which is not a full mDNS querier (RFC 6762,
// 6.7).
func isLegacyUnicastQuery(query dns.Message) bool {
	from := query.From()
	return from != nil && from.Port() != transport.Port
}

// answer builds the response to query, which was received on ifi. It
// returns nil when no registered service answers any question. The second
// result reports whether the answers include a shared record, which a
// multicast response delays (RFC 6762, 6).
func (r *responder) answer(query dns.Message, ifi *net.Interface) (dns.Message, bool) {
	if query == nil || !query.IsQuery() || query.ResponseCode() != 0 {
		return nil, false
	}

	legacy := isLegacyUnicastQuery(query)
	otherTTL, hostTTL := uint(OtherRecordTTL), uint(HostRecordTTL)
	if legacy {
		otherTTL, hostTTL = legacyUnicastTTL, legacyUnicastTTL
	}

	b := newResponseBuilder(query.Answers())
	for _, svc := range r.localServices() {
		records, err := svc.records(ifi, otherTTL, hostTTL)
		if err != nil {
			continue
		}
		if legacy {
			// A legacy resolver does not know the cache-flush bit.
			for _, rr := range records.all() {
				rr.SetUnicastResponse(false)
			}
		}
		for _, q := range query.Questions() {
			b.answerQuestion(q, records)
		}
	}
	if len(b.answers) == 0 {
		return nil, false
	}

	opts := []dns.MessageOption{
		dns.WithMessageAnswers(b.answers...),
		dns.WithMessageAdditions(b.additions...),
	}
	if legacy {
		opts = append(opts,
			dns.WithMessageID(query.ID()),
			dns.WithMessageQuestions(query.Questions()...),
		)
	}
	return dns.NewResponseMessage(opts...), b.shared
}

// announcement builds the unsolicited response which publishes svc on ifi
// (RFC 6762, 8.3), or which withdraws it when goodbye is set (10.1).
func announcement(svc *LocalService, ifi *net.Interface, goodbye bool) (dns.Message, error) {
	otherTTL, hostTTL := uint(OtherRecordTTL), uint(HostRecordTTL)
	if goodbye {
		otherTTL, hostTTL = 0, 0
	}
	records, err := svc.records(ifi, otherTTL, hostTTL)
	if err != nil {
		return nil, err
	}
	return dns.NewResponseMessage(dns.WithMessageAnswers(records.all()...)), nil
}

// responseBuilder collects the answers and the additional records of a
// response without duplicates, leaving out the answers the querier
// already knows.
type responseBuilder struct {
	known     dns.ResourceRecordSet
	answers   []dns.ResourceRecord
	additions []dns.ResourceRecord
	shared    bool
}

func newResponseBuilder(known dns.ResourceRecordSet) *responseBuilder {
	return &responseBuilder{
		known:     known,
		answers:   []dns.ResourceRecord{},
		additions: []dns.ResourceRecord{},
		shared:    false,
	}
}

func sameRecord(a, b dns.ResourceRecord) bool {
	return a.IsName(b.Name()) && a.Type() == b.Type() && strings.EqualFold(a.Content(), b.Content())
}

func containsRecord(records []dns.ResourceRecord, rr dns.ResourceRecord) bool {
	for _, r := range records {
		if sameRecord(r, rr) {
			return true
		}
	}
	return false
}

// isKnownAnswer reports whether the query listed rr as a known answer with
// at least half of its TTL left (RFC 6762, 7.1).
func (b *responseBuilder) isKnownAnswer(rr dns.ResourceRecord) bool {
	for _, known := range b.known {
		if sameRecord(known, rr) && rr.TTL()/2 <= known.TTL() {
			return true
		}
	}
	return false
}

func (b *responseBuilder) addAnswer(rr dns.ResourceRecord, shared bool) {
	if containsRecord(b.answers, rr) || b.isKnownAnswer(rr) {
		return
	}
	b.answers = append(b.answers, rr)
	if shared {
		b.shared = true
	}
}

func (b *responseBuilder) addAddition(rr dns.ResourceRecord) {
	if containsRecord(b.answers, rr) || containsRecord(b.additions, rr) || b.isKnownAnswer(rr) {
		return
	}
	b.additions = append(b.additions, rr)
}

func matchesQuestion(q dns.Question, rr dns.ResourceRecord) bool {
	return rr.IsName(q.Name()) && (q.Type() == dns.ANY || q.Type() == rr.Type())
}

// answerQuestion adds the records of one service which answer q, and the
// additional records a querier needs next (RFC 6763, 12).
func (b *responseBuilder) answerQuestion(q dns.Question, records *localServiceRecords) {
	instance := false
	host := false

	for _, ptr := range append([]dns.ResourceRecord{records.ptr}, records.subtypes...) {
		if matchesQuestion(q, ptr) {
			b.addAnswer(ptr, true)
			instance = true
		}
	}
	if matchesQuestion(q, records.enum) {
		b.addAnswer(records.enum, true)
	}
	if matchesQuestion(q, records.srv) {
		b.addAnswer(records.srv, false)
		host = true
	}
	if matchesQuestion(q, records.txt) {
		b.addAnswer(records.txt, false)
	}
	for _, addr := range records.addrs {
		if matchesQuestion(q, addr) {
			b.addAnswer(addr, false)
			host = true
		}
	}

	if instance {
		b.addAddition(records.srv)
		b.addAddition(records.txt)
		host = true
	}
	if host {
		for _, addr := range records.addrs {
			b.addAddition(addr)
		}
	}
}
