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
	"encoding/binary"
	"strings"
	"testing"
	"time"

	"github.com/cybergarage/go-mdns/mdns/dns"
)

// newTestTXTResponse returns a response message which holds a single TXT record
// of the specified service name, TTL and attributes.
func newTestTXTResponse(t *testing.T, name string, ttl uint32, attrs ...string) Message {
	t.Helper()

	encodeName := func(name string) []byte {
		b := []byte{}
		for label := range strings.SplitSeq(name, ".") {
			b = append(b, byte(len(label)))
			b = append(b, []byte(label)...)
		}
		return append(b, 0x00)
	}

	rdata := []byte{}
	for _, attr := range attrs {
		rdata = append(rdata, byte(len(attr)))
		rdata = append(rdata, []byte(attr)...)
	}

	msgBytes := []byte{
		0x00, 0x00, // ID
		0x84, 0x00, // Flags (response, authoritative)
		0x00, 0x00, // QD
		0x00, 0x01, // AN
		0x00, 0x00, // NS
		0x00, 0x00, // AR
	}
	msgBytes = append(msgBytes, encodeName(name)...)
	msgBytes = binary.BigEndian.AppendUint16(msgBytes, uint16(dns.TXT))
	msgBytes = binary.BigEndian.AppendUint16(msgBytes, uint16(dns.IN))
	msgBytes = binary.BigEndian.AppendUint32(msgBytes, ttl)
	msgBytes = binary.BigEndian.AppendUint16(msgBytes, uint16(len(rdata)))
	msgBytes = append(msgBytes, rdata...)

	msg, err := dns.NewMessageWithBytes(msgBytes)
	if err != nil {
		t.Fatal(err)
	}

	return msg
}

func newTestService(t *testing.T, name string, ttl uint32, attrs ...string) Service {
	t.Helper()

	service, err := NewService(
		WithServiceMessage(newTestTXTResponse(t, name, ttl, attrs...)),
	)
	if err != nil {
		t.Fatal(err)
	}

	return service
}

func TestServiceCache(t *testing.T) {
	const serviceName = "DD200C20D25AE5F7._matterc._udp.local"

	cache := newServiceCache()

	// An unknown service is added.
	service := newTestService(t, serviceName, 120, "D=840", "CM=2")
	if service.FullName() != serviceName {
		t.Fatalf("full name %s != %s", service.FullName(), serviceName)
	}
	if service.TTL() != 120*time.Second {
		t.Fatalf("ttl %v != %v", service.TTL(), 120*time.Second)
	}

	event, ok := cache.Update(service)
	if !ok || event.Type != ServiceAdded {
		t.Fatalf("event %v (%t) != added", event.Type, ok)
	}

	// The same service does not raise an event.
	if _, ok := cache.Update(newTestService(t, serviceName, 120, "D=840", "CM=2")); ok {
		t.Errorf("an unchanged service should not raise an event")
	}

	// A changed service is updated.
	event, ok = cache.Update(newTestService(t, serviceName, 120, "D=840", "CM=0"))
	if !ok || event.Type != ServiceUpdated {
		t.Errorf("event %v (%t) != updated", event.Type, ok)
	}

	// The cached service does not expire before its TTL elapses.
	if events := cache.RemoveExpired(time.Now().Add(119 * time.Second)); len(events) != 0 {
		t.Errorf("the service should not expire yet: %v", events)
	}

	// RFC 6762: 10.1. Goodbye Packets
	event, ok = cache.Update(newTestService(t, serviceName, 0, "D=840", "CM=0"))
	if !ok || event.Type != ServiceRemoved {
		t.Errorf("event %v (%t) != removed", event.Type, ok)
	}
	if 0 < len(cache.Services()) {
		t.Errorf("the service should be removed: %v", cache.Services())
	}

	// A goodbye packet of an unknown service is ignored.
	if _, ok := cache.Update(newTestService(t, serviceName, 0)); ok {
		t.Errorf("an unknown goodbye packet should not raise an event")
	}
}

func TestServiceCacheExpiration(t *testing.T) {
	const serviceName = "DD200C20D25AE5F7._matterc._udp.local"

	cache := newServiceCache()

	if _, ok := cache.Update(newTestService(t, serviceName, 120, "D=840")); !ok {
		t.Fatal("the service should be added")
	}

	events := cache.RemoveExpired(time.Now().Add(121 * time.Second))
	if len(events) != 1 {
		t.Fatalf("expired event count %d != 1", len(events))
	}
	if events[0].Type != ServiceRemoved {
		t.Errorf("event %v != removed", events[0].Type)
	}
	if 0 < len(cache.Services()) {
		t.Errorf("the expired service should be removed: %v", cache.Services())
	}
}
