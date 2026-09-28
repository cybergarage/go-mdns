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

package transport

import (
	"net"
	"testing"

	"github.com/cybergarage/go-mdns/mdns/dns"
)

func TestUnicastResponseAddr(t *testing.T) {
	query := func(port int, zone string, qu bool, v6 bool) dns.Message {
		cls := dns.IN
		if qu {
			cls |= dns.QU
		}
		q := dns.NewQuestion(dns.WithQuestionName("_http._tcp.local"), dns.WithQuestionType(dns.PTR), dns.WithQuestionClass(cls))
		ip := net.IPv4(192, 168, 0, 20)
		if v6 {
			ip = net.ParseIP("fe80::20")
		}
		from := dns.NewAddr(dns.WithAddrIP(ip), dns.WithAddrPort(port), dns.WithAddrZone(zone))
		msg, err := dns.NewMessageWithBytes(dns.NewRequestMessage(dns.WithMessageQuestions(q)).Bytes(), dns.WithMessageFrom(from))
		if err != nil {
			t.Fatal(err)
		}
		return msg
	}

	if _, _, ok := unicastResponseAddr(query(Port, "", false, false)); ok {
		t.Error("a multicast query from port 5353 is answered by unicast")
	}
	if addr, port, ok := unicastResponseAddr(query(50000, "", false, false)); !ok || addr != "192.168.0.20" || port != 50000 {
		t.Errorf("legacy unicast query: (%q, %d, %v), want the source address and port", addr, port, ok)
	}
	if addr, port, ok := unicastResponseAddr(query(Port, "", true, false)); !ok || addr != "192.168.0.20" || port != Port {
		t.Errorf("QU query: (%q, %d, %v), want the source on port 5353", addr, port, ok)
	}
	if addr, _, ok := unicastResponseAddr(query(Port, "en0", true, true)); !ok || addr != "fe80::20%en0" {
		t.Errorf("QU query from a link-local address: %q, %v, want the zone kept", addr, ok)
	}
	// A received IPv4 address carries its interface as the zone too, which
	// must not reach the destination address.
	if addr, _, ok := unicastResponseAddr(query(Port, "eth0", true, false)); !ok || addr != "192.168.0.20" {
		t.Errorf("QU query from IPv4 with an interface zone: %q, %v, want no zone", addr, ok)
	}
}
