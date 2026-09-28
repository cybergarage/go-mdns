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

package dns

import (
	"net"
	"reflect"
	"strings"
	"testing"
)

func TestResourceRecordBuilders(t *testing.T) {
	ptr, err := NewPTRResourceRecord("_http._tcp.local", "web._http._tcp.local", 4500)
	if err != nil {
		t.Fatal(err)
	}
	srv, err := NewSRVResourceRecord("web._http._tcp.local", 0, 0, 8080, "host.local", 120)
	if err != nil {
		t.Fatal(err)
	}
	srv.SetUnicastResponse(true)
	txt, err := NewTXTResourceRecord("web._http._tcp.local", []string{"path=/", "v=1"}, 4500)
	if err != nil {
		t.Fatal(err)
	}
	a, err := NewAResourceRecord("host.local", net.IPv4(192, 168, 0, 10), 120)
	if err != nil {
		t.Fatal(err)
	}
	aaaa, err := NewAAAAResourceRecord("host.local", net.ParseIP("fe80::1"), 120)
	if err != nil {
		t.Fatal(err)
	}

	msg := NewResponseMessage(WithMessageID(0x1234), WithMessageAnswers(ptr, srv, txt, a, aaaa))
	parsed, err := NewMessageWithBytes(msg.Bytes())
	if err != nil {
		t.Fatalf("NewMessageWithBytes(...) error = %v", err)
	}
	if parsed.ID() != 0x1234 || !parsed.IsResponse() || !parsed.AA() {
		t.Fatalf("header: ID %#x, response %v, AA %v", parsed.ID(), parsed.IsResponse(), parsed.AA())
	}
	answers := parsed.Answers()
	if len(answers) != 5 {
		t.Fatalf("%d answers, want 5", len(answers))
	}

	if got, ok := answers[0].(PTRRecord); !ok || got.DomainName() != "web._http._tcp.local" || got.TTL() != 4500 {
		t.Errorf("PTR = %v", answers[0])
	}
	gotSRV, ok := answers[1].(SRVRecord)
	if !ok || gotSRV.Port() != 8080 || gotSRV.Target() != "host.local" || !gotSRV.UnicastResponse() {
		t.Errorf("SRV = %v, cache flush %v", answers[1], answers[1].UnicastResponse())
	}
	if got, ok := answers[2].(TXTRecord); !ok || !reflect.DeepEqual(got.Strings(), []string{"path=/", "v=1"}) {
		t.Errorf("TXT = %v", answers[2])
	}
	if got, ok := answers[3].(ARecord); !ok || !got.Address().Equal(net.IPv4(192, 168, 0, 10)) {
		t.Errorf("A = %v", answers[3])
	}
	if got, ok := answers[4].(AAAARecord); !ok || !got.Address().Equal(net.ParseIP("fe80::1")) {
		t.Errorf("AAAA = %v", answers[4])
	}
}

func TestResourceRecordBuildersRejectInvalidInput(t *testing.T) {
	if _, err := NewPTRResourceRecord("a.local", "bad..name", 1); err == nil {
		t.Error("NewPTRResourceRecord with an empty label = nil error")
	}
	if _, err := NewPTRResourceRecord("a.local", strings.Repeat("x", 64)+".local", 1); err == nil {
		t.Error("NewPTRResourceRecord with a 64-byte label = nil error")
	}
	if _, err := NewTXTResourceRecord("a.local", []string{strings.Repeat("x", 256)}, 1); err == nil {
		t.Error("NewTXTResourceRecord with a 256-byte string = nil error")
	}
	if _, err := NewAResourceRecord("a.local", net.ParseIP("fe80::1"), 1); err == nil {
		t.Error("NewAResourceRecord with an IPv6 address = nil error")
	}
	if _, err := NewAAAAResourceRecord("a.local", net.IPv4(1, 2, 3, 4), 1); err == nil {
		t.Error("NewAAAAResourceRecord with an IPv4 address = nil error")
	}
}

func TestEmptyTXTResourceRecord(t *testing.T) {
	txt, err := NewTXTResourceRecord("a.local", nil, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(txt.Data(), []byte{0}) {
		t.Fatalf("empty TXT data = %v, want a single empty string", txt.Data())
	}
}
