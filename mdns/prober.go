// Copyright (C) 2022 The go-mdns Authors All rights reserved.
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
	"bytes"
	"context"
	"math/rand/v2"
	"net"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns/dns"
)

// Timing of the probes (RFC 6762, 8.1 and 8.2).
const (
	probeCount         = 3
	probeInterval      = 250 * time.Millisecond
	maxProbeDelay      = 250 * time.Millisecond
	probeDeferInterval = time.Second
	// A host which sees fifteen conflicts within ten seconds waits five
	// seconds before each further probe (RFC 6762, 8.1).
	conflictRateCount  = 15
	conflictRateWindow = 10 * time.Second
	conflictRateDelay  = 5 * time.Second
)

// probe is a service whose names are being probed before it is published.
type probe struct {
	svc *LocalService
	// names are the names the probe claims: the instance name, the host
	// name, or both. A name the server already holds for another service
	// is not probed again.
	names []string

	conflict chan string
	lost     chan struct{}

	cancelOnce sync.Once
	canceled   chan struct{}
	cancelErr  error
}

func newProbe(svc *LocalService, names []string) *probe {
	return &probe{
		svc:        svc,
		names:      names,
		conflict:   make(chan string, 1),
		lost:       make(chan struct{}, 1),
		cancelOnce: sync.Once{},
		canceled:   make(chan struct{}),
		cancelErr:  nil,
	}
}

// signalConflict reports that another node holds name.
func (p *probe) signalConflict(name string) {
	select {
	case p.conflict <- name:
	default:
	}
}

// signalLost reports that a simultaneous probe of another node won the
// tiebreak (RFC 6762, 8.2).
func (p *probe) signalLost() {
	select {
	case p.lost <- struct{}{}:
	default:
	}
}

// cancel stops the probe with err.
func (p *probe) cancel(err error) {
	p.cancelOnce.Do(func() {
		p.cancelErr = err
		close(p.canceled)
	})
}

func (p *probe) claims(name string) bool {
	return slices.ContainsFunc(p.names, func(n string) bool {
		return strings.EqualFold(n, name)
	})
}

// probeRecords returns the records of svc which are published on ifi and
// named by one of names, without the cache-flush bit, which the Authority
// section of a probe does not carry.
func probeRecords(svc *LocalService, ifi *net.Interface, names []string) ([]dns.ResourceRecord, error) {
	records, err := svc.records(ifi, OtherRecordTTL, HostRecordTTL)
	if err != nil {
		return nil, err
	}
	unique := append([]dns.ResourceRecord{records.srv, records.txt}, records.addrs...)
	claimed := []dns.ResourceRecord{}
	for _, rr := range unique {
		if slices.ContainsFunc(names, rr.IsName) {
			rr.SetUnicastResponse(false)
			claimed = append(claimed, rr)
		}
	}
	return claimed, nil
}

// probeMessage builds a probe of the names of p for ifi: a query of type ANY
// for every name, with the proposed records in the Authority section
// (RFC 6762, 8.1). The first probe asks for unicast responses.
func probeMessage(p *probe, ifi *net.Interface, first bool) (dns.Message, error) {
	questions := make([]dns.Question, 0, len(p.names))
	for _, name := range p.names {
		q := dns.NewQuestion(
			dns.WithQuestionName(name),
			dns.WithQuestionType(dns.ANY),
			dns.WithQuestionClass(dns.IN),
		)
		q.SetUnicastResponse(first)
		questions = append(questions, q)
	}
	authorities, err := probeRecords(p.svc, ifi, p.names)
	if err != nil {
		return nil, err
	}
	return dns.NewRequestMessage(
		dns.WithMessageQuestions(questions...),
		dns.WithMessageNameServers(authorities...),
	), nil
}

// runProbe probes the names of p until they are found unique, and returns
// nil then. It returns a *ConflictError when another node holds one of
// them, ErrNotRunning when the server stops, the cancellation error of p
// when it is deregistered or replaced, and the error of ctx when ctx is
// done.
func (server *Server) runProbe(ctx context.Context, done <-chan struct{}, p *probe) error {
	if ctx == nil {
		ctx = context.Background()
	}

	// wait waits for d, and reports whether the probe lost a tiebreak
	// meanwhile; a non-nil error ends the probe.
	wait := func(d time.Duration) (bool, error) {
		timer := time.NewTimer(d)
		defer timer.Stop()
		select {
		case <-timer.C:
			return false, nil
		case <-p.lost:
			return true, nil
		case name := <-p.conflict:
			server.recordConflict()
			return false, &ConflictError{Service: copyLocalService(p.svc), Name: name}
		case <-p.canceled:
			return false, p.cancelErr
		case <-done:
			return false, ErrNotRunning
		case <-ctx.Done():
			return false, ctx.Err()
		}
	}

	for {
		if delay := server.conflictDelay(); 0 < delay {
			if _, err := wait(delay); err != nil {
				return err
			}
		}
		// A random delay keeps the hosts which power on together from
		// probing at the same instant (RFC 6762, 8.1).
		lost, err := wait(rand.N(maxProbeDelay)) // nolint: gosec
		for i := 0; !lost && err == nil && i < probeCount; i++ {
			server.sendProbe(p, i == 0)
			lost, err = wait(probeInterval)
		}
		if err != nil {
			return err
		}
		if !lost {
			return nil
		}
		// The probe lost a tiebreak against a simultaneous probe of
		// another node, which takes the name unless it is withdrawn
		// within a second (RFC 6762, 8.2).
		log.Debugf("mdns: probe of %s deferred by a simultaneous probe", p.svc.FullName())
		if _, err := wait(probeDeferInterval); err != nil {
			return err
		}
	}
}

// sendProbe sends a probe of p on every bound interface.
func (server *Server) sendProbe(p *probe, first bool) {
	for _, ms := range server.MessageManager.MulticastManager.Servers {
		ifi, err := ms.MulticastSocket.ListenInterface()
		if err != nil {
			ifi = nil
		}
		msg, err := probeMessage(p, ifi, first)
		if err != nil {
			log.Warnf("mdns: build the probe of %s: %s", p.svc.FullName(), err)
			continue
		}
		if err := ms.AnnounceMessage(msg); err != nil {
			log.Debugf("mdns: probe %s: %s", p.svc.FullName(), err)
		}
	}
}

// recordConflict counts a conflict for the rate limit of the probes.
func (server *Server) recordConflict() {
	server.conflictMutex.Lock()
	defer server.conflictMutex.Unlock()
	now := time.Now()
	server.conflicts = append(server.conflicts, now)
	server.conflicts = slices.DeleteFunc(server.conflicts, func(t time.Time) bool {
		return conflictRateWindow < now.Sub(t)
	})
}

// conflictDelay returns how long a probe waits before it starts, which is
// not zero while the conflicts exceed the rate limit (RFC 6762, 8.1).
func (server *Server) conflictDelay() time.Duration {
	server.conflictMutex.Lock()
	defer server.conflictMutex.Unlock()
	now := time.Now()
	server.conflicts = slices.DeleteFunc(server.conflicts, func(t time.Time) bool {
		return conflictRateWindow < now.Sub(t)
	})
	if len(server.conflicts) < conflictRateCount {
		return 0
	}
	return conflictRateDelay
}

// localAddresses returns the addresses of this node: the addresses of its
// interfaces and the addresses its services are given.
func (server *Server) localAddresses() []net.IP {
	ips := []net.IP{}
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	for _, svc := range server.responder.allServices() {
		ips = append(ips, svc.Addresses...)
	}
	return ips
}

// isOwnRecord reports whether rr, a record named by a name of svc, is one
// this node publishes or would publish for svc: an identical record is not
// a conflict (RFC 6762, 9), and this node receives its own messages back.
func isOwnRecord(svc *LocalService, rr dns.ResourceRecord, localAddrs []net.IP) bool {
	switch rr.Type() {
	case dns.SRV:
		srv, ok := rr.(dns.SRVRecord)
		if !ok {
			return false
		}
		return rr.IsName(svc.FullName()) &&
			strings.EqualFold(strings.TrimSuffix(srv.Target(), "."), svc.HostName()) &&
			srv.Port() == uint(svc.Port) && srv.Priority() == 0 && srv.Weight() == 0 // nolint: gosec
	case dns.TXT:
		txt, err := dns.NewTXTResourceRecord(svc.FullName(), svc.TXT, 0)
		if err != nil {
			return false
		}
		return rr.IsName(svc.FullName()) && bytes.Equal(rr.Data(), txt.Data())
	case dns.A, dns.AAAA:
		if !rr.IsName(svc.HostName()) {
			return false
		}
		ip := net.IP(rr.Data())
		return slices.ContainsFunc(localAddrs, ip.Equal)
	default:
		return false
	}
}

// canonicalRData returns the uncompressed data of rr, which the tiebreak of
// the simultaneous probes compares (RFC 6762, 8.2).
func canonicalRData(rr dns.ResourceRecord) []byte {
	if srv, ok := rr.(dns.SRVRecord); ok {
		canonical, err := dns.NewSRVResourceRecord(
			rr.Name(),
			uint16(srv.Priority()), // nolint: gosec
			uint16(srv.Weight()),   // nolint: gosec
			uint16(srv.Port()),     // nolint: gosec
			srv.Target(),
			0)
		if err == nil {
			return canonical.Data()
		}
	}
	return rr.Data()
}

// compareRecords compares two records by their class, type and data, as the
// tiebreak of the simultaneous probes does (RFC 6762, 8.2).
func compareRecords(a, b dns.ResourceRecord) int {
	if c := int(a.Class()) - int(b.Class()); c != 0 {
		return c
	}
	if c := int(a.Type()) - int(b.Type()); c != 0 {
		return c
	}
	return bytes.Compare(canonicalRData(a), canonicalRData(b))
}

// compareProbeRecords compares the records two hosts propose for a name:
// both are sorted, and the first record which differs decides; a host with
// more records wins when one set is a prefix of the other (RFC 6762, 8.2).
func compareProbeRecords(ours, theirs []dns.ResourceRecord) int {
	ours = slices.Clone(ours)
	theirs = slices.Clone(theirs)
	slices.SortFunc(ours, compareRecords)
	slices.SortFunc(theirs, compareRecords)
	for i := 0; i < len(ours) && i < len(theirs); i++ {
		if c := compareRecords(ours[i], theirs[i]); c != 0 {
			return c
		}
	}
	return len(ours) - len(theirs)
}

// responseRecords returns the records of a response which claim names: the
// answers and the additional records, without the goodbye records.
func responseRecords(msg dns.Message) []dns.ResourceRecord {
	records := []dns.ResourceRecord{}
	for _, rr := range append(msg.Answers(), msg.Additions()...) {
		if rr.TTL() == 0 {
			continue
		}
		records = append(records, rr)
	}
	return records
}

// checkResponse looks for the records of another node which conflict with
// the names this server probes or holds (RFC 6762, 8.1 and 9).
func (server *Server) checkResponse(msg dns.Message) {
	records := responseRecords(msg)
	if len(records) == 0 {
		return
	}

	var localAddrs []net.IP
	conflicts := func(svc *LocalService, name string, types ...dns.Type) bool {
		for _, rr := range records {
			if !rr.IsName(name) {
				continue
			}
			if 0 < len(types) && !slices.Contains(types, rr.Type()) {
				continue
			}
			if localAddrs == nil {
				localAddrs = server.localAddresses()
			}
			if !isOwnRecord(svc, rr, localAddrs) {
				return true
			}
		}
		return false
	}

	// While a name is probed, a record of any type with the name is a
	// conflict (RFC 6762, 8.1).
	for _, p := range server.responder.probes() {
		for _, name := range p.names {
			if conflicts(p.svc, name) {
				log.Infof("mdns: %s is used by another node", name)
				p.signalConflict(name)
			}
		}
	}

	// Once the names are held, a record with the same name and type and
	// different data is a conflict, and the service is probed again
	// (RFC 6762, 9).
	for _, svc := range server.responder.localServices() {
		if conflicts(svc, svc.FullName(), dns.SRV, dns.TXT) || conflicts(svc, svc.HostName(), dns.A, dns.AAAA) {
			log.Infof("mdns: a record of %s conflicts, probing it again", svc.FullName())
			server.reprobe(svc)
		}
	}
}

// checkProbe resolves a simultaneous probe of another node for a name this
// server probes: the host whose proposed records compare lower defers
// (RFC 6762, 8.2).
func (server *Server) checkProbe(msg dns.Message) {
	authorities := msg.NameServers()
	if !msg.IsQuery() || len(authorities) == 0 {
		return
	}
	probes := server.responder.probes()
	if len(probes) == 0 {
		return
	}

	var ifi *net.Interface
	if from := msg.From(); from != nil {
		ifi = from.Interface()
	}
	localAddrs := server.localAddresses()

	for _, p := range probes {
		for _, name := range p.names {
			theirs := []dns.ResourceRecord{}
			own := true
			for _, rr := range authorities {
				if !rr.IsName(name) {
					continue
				}
				theirs = append(theirs, rr)
				if !isOwnRecord(p.svc, rr, localAddrs) {
					own = false
				}
			}
			// The probe is not for this name, or it is this node's own
			// probe received back, or it proposes the same records.
			if len(theirs) == 0 || own {
				continue
			}
			ours, err := probeRecords(p.svc, ifi, []string{name})
			if err != nil {
				continue
			}
			if compareProbeRecords(ours, theirs) < 0 {
				p.signalLost()
			}
		}
	}
}
