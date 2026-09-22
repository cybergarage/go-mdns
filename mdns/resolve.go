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
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/cybergarage/go-mdns/mdns/dns"
)

const (
	// DefaultResolveGracePeriod is the duration to wait for the additional
	// responses after the first answer is received. A responder may answer
	// the records of a name with more than one message, such as a host which
	// holds both an IPv4 and an IPv6 address.
	DefaultResolveGracePeriod = time.Duration(250) * time.Millisecond
)

// newRequestWithQuestions returns a request message which asks the specified
// types of a name.
func newRequestWithQuestions(name string, types ...Type) Message {
	questions := make([]dns.Question, 0, len(types))
	for _, t := range types {
		questions = append(questions,
			dns.NewQuestion(
				dns.WithQuestionName(name),
				dns.WithQuestionType(t),
				dns.WithQuestionClass(dns.IN),
			))
	}
	return dns.NewRequestMessage(dns.WithMessageQuestions(questions...))
}

// queryUntilAnswered announces the query message, and calls the response
// handler for every answer. The handler returns true when the answer completes
// the query, and the query then waits for the grace period to collect the
// remaining answers before it returns.
func (client *clientImpl) queryUntilAnswered(ctx context.Context, queryMsg Message, onAnswer func(Message) bool) error {
	answeredCtx, answered := context.WithCancel(ctx)
	defer answered()

	var once sync.Once

	responseHandler := func(resMsg dns.Message) {
		if !resMsg.IsResponse() {
			return
		}
		if !queryMsg.IsQueryAnswer(resMsg) {
			return
		}
		if !onAnswer(resMsg) {
			return
		}
		// The query is answered, but a responder may send the remaining
		// records with another message, so the grace period collects them
		// before the query returns.
		once.Do(func() {
			go func() {
				timer := time.NewTimer(client.resolveGracePeriod)
				defer timer.Stop()
				select {
				case <-ctx.Done():
				case <-timer.C:
				}
				answered()
			}()
		})
	}

	client.RegisterMessageHandler(responseHandler)
	defer client.UnRegisterMessageHandler(responseHandler)

	return client.queryUntilDone(answeredCtx, queryMsg)
}

// Resolve resolves the specified service instance name, such as
// "instance._matter._tcp.local", and returns the service with its host, port,
// TXT attributes and addresses.
//
// RFC 6763: 5. Service Instance Resolution
// A service instance is resolved by querying its SRV and TXT records, and the
// address records of the SRV target name are queried when a responder does not
// send them with the answer.
func (client *clientImpl) Resolve(ctx context.Context, name string) (Service, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("service name is %w", ErrInvalid)
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, client.queryTimeout)
		defer cancel()
	}

	queryMsg := newRequestWithQuestions(name, SRV, TXT)

	var mutex sync.Mutex
	var resolved Service

	onAnswer := func(resMsg Message) bool {
		service, err := NewService(
			WithServiceMessage(resMsg),
		)
		if err != nil {
			return false
		}
		if !strings.EqualFold(service.FullName(), name) {
			return false
		}

		mutex.Lock()
		defer mutex.Unlock()

		// A later answer is kept only when it holds more information, so
		// that an answer without the SRV record does not drop the host.
		switch {
		case resolved == nil:
			resolved = service
		case len(resolved.Host()) == 0 && 0 < len(service.Host()):
			resolved = service
		case len(resolved.Addresses()) == 0 && 0 < len(service.Addresses()):
			resolved = service
		}

		return 0 < len(resolved.Host())
	}

	if err := client.queryUntilAnswered(ctx, queryMsg, onAnswer); err != nil {
		return nil, err
	}

	mutex.Lock()
	service := resolved
	mutex.Unlock()

	if service == nil {
		return nil, fmt.Errorf("service (%s) is %w", name, ErrNotFound)
	}

	if 0 < len(service.Addresses()) || len(service.Host()) == 0 {
		return service, nil
	}

	// The responder did not send the address records of the SRV target, so
	// the host is resolved with another query.
	addrs, err := client.LookupHost(ctx, service.Host())
	if err != nil {
		return service, nil // nolint: nilerr
	}

	if impl, ok := service.(*serviceImpl); ok {
		ips := make([]net.IP, 0, len(addrs))
		for _, addr := range addrs {
			ips = append(ips, addr.IP)
		}
		impl.setAddresses(ips)
	}

	return service, nil
}

// LookupHost resolves the specified host name, such as "device.local", and
// returns its addresses. The IPv6 scoped addressing zone is set for the
// link-local addresses.
//
// RFC 6762: 1. Introduction
// Resolving a host name is the core of Multicast DNS, and it does not require
// the DNS-SD service records.
func (client *clientImpl) LookupHost(ctx context.Context, host string) ([]*net.IPAddr, error) {
	if len(host) == 0 {
		return nil, fmt.Errorf("host name is %w", ErrInvalid)
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, client.queryTimeout)
		defer cancel()
	}

	queryMsg := newRequestWithQuestions(host, A, AAAA)

	var mutex sync.Mutex
	addrs := []*net.IPAddr{}

	onAnswer := func(resMsg Message) bool {
		mutex.Lock()
		defer mutex.Unlock()

		for _, record := range resMsg.ResourceRecordSet() {
			if !record.IsName(host) {
				continue
			}
			// RFC 6762: 10.1. Goodbye Packets
			if record.TTL() == 0 {
				continue
			}
			addrRecord, ok := record.(interface{ Address() net.IP })
			if !ok {
				continue
			}
			ip := addrRecord.Address()
			if ip == nil {
				continue
			}

			zone := ""
			if ip.To4() == nil && ip.IsLinkLocalUnicast() {
				if from := resMsg.From(); from != nil {
					if ifi := from.Interface(); ifi != nil {
						zone = ifi.Name
					}
				}
			}

			isKnown := false
			for _, addr := range addrs {
				if addr.IP.Equal(ip) && addr.Zone == zone {
					isKnown = true
					break
				}
			}
			if isKnown {
				continue
			}

			addrs = append(addrs, &net.IPAddr{IP: ip, Zone: zone})
		}

		return 0 < len(addrs)
	}

	if err := client.queryUntilAnswered(ctx, queryMsg, onAnswer); err != nil {
		return nil, err
	}

	mutex.Lock()
	defer mutex.Unlock()

	if len(addrs) == 0 {
		return nil, fmt.Errorf("host (%s) is %w", host, ErrNotFound)
	}

	return addrs, nil
}
