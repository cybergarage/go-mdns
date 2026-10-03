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
	"context"
	"errors"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/cybergarage/go-mdns/mdns/dns"
	"github.com/cybergarage/go-mdns/mdns/transport"
)

func testLocalService() *LocalService {
	return &LocalService{
		Instance:  "665F6E75B5D3A9C2",
		Service:   "_matterc._udp",
		Subtypes:  []string{"_L3840", "_S15", "_V65521", "_CM"},
		Host:      "B75AFB458ECD6D6F",
		Port:      5540,
		TXT:       []string{"D=3840", "CM=1", "VP=65521+32769"},
		Addresses: []net.IP{net.IPv4(192, 168, 0, 10), net.ParseIP("fe80::1")},
	}
}

func testResponder(t *testing.T) *responder {
	t.Helper()
	svc := testLocalService()
	if err := svc.Validate(); err != nil {
		t.Fatal(err)
	}
	r := newResponder()
	r.register(svc)
	return r
}

// testQuery builds a query as it arrives from the network: encoded, parsed
// back, and tagged with its source port.
func testQuery(t *testing.T, name string, typ dns.Type, fromPort int, known ...dns.ResourceRecord) dns.Message {
	t.Helper()
	q := dns.NewQuestion(dns.WithQuestionName(name), dns.WithQuestionType(typ), dns.WithQuestionClass(dns.IN))
	msg := dns.NewRequestMessage(dns.WithMessageQuestions(q), dns.WithMessageAnswers(known...))
	from := dns.NewAddr(dns.WithAddrIP(net.IPv4(192, 168, 0, 20)), dns.WithAddrPort(fromPort))
	parsed, err := dns.NewMessageWithBytes(msg.Bytes(), dns.WithMessageFrom(from))
	if err != nil {
		t.Fatal(err)
	}
	return parsed
}

// reparse encodes res and parses it back, as a querier receives it.
func reparse(t *testing.T, res dns.Message) dns.Message {
	t.Helper()
	parsed, err := dns.NewMessageWithBytes(res.Bytes())
	if err != nil {
		t.Fatalf("parse the response: %v", err)
	}
	return parsed
}

func recordTypes(records dns.ResourceRecordSet) []dns.Type {
	types := make([]dns.Type, 0, len(records))
	for _, r := range records {
		types = append(types, r.Type())
	}
	return types
}

func TestResponderAnswersBrowse(t *testing.T) {
	r := testResponder(t)
	for _, name := range []string{"_matterc._udp.local", "_L3840._sub._matterc._udp.local", "_MATTERC._UDP.LOCAL"} {
		res, shared := r.answer(testQuery(t, name, dns.PTR, transport.Port), nil)
		if res == nil {
			t.Fatalf("%s: no response", name)
		}
		if !shared {
			t.Errorf("%s: a PTR answer is shared", name)
		}
		res = reparse(t, res)
		if res.ID() != 0 || !res.IsResponse() || !res.AA() || len(res.Questions()) != 0 {
			t.Errorf("%s: header ID %d, response %v, AA %v, %d questions", name, res.ID(), res.IsResponse(), res.AA(), len(res.Questions()))
		}
		if got := recordTypes(res.Answers()); !reflect.DeepEqual(got, []dns.Type{dns.PTR}) {
			t.Errorf("%s: answers %v, want [PTR]", name, got)
		}
		if got := recordTypes(res.Additions()); !reflect.DeepEqual(got, []dns.Type{dns.SRV, dns.TXT, dns.A, dns.AAAA}) {
			t.Errorf("%s: additions %v, want [SRV TXT A AAAA]", name, got)
		}

		// A browser resolves the whole service from the one response.
		svc, err := NewService(WithServiceMessage(res))
		if err != nil {
			t.Fatal(err)
		}
		if svc.Port() != 5540 || svc.Host() != "B75AFB458ECD6D6F.local" || len(svc.Addresses()) != 2 {
			t.Errorf("%s: resolved port %d, host %q, addresses %v", name, svc.Port(), svc.Host(), svc.Addresses())
		}
		if attr, ok := svc.LookupResourceAttribute("D"); !ok || attr.Value() != "3840" {
			t.Errorf("%s: TXT D = %v, %v", name, attr, ok)
		}
	}
}

func TestResponderAnswersResolveAndLookup(t *testing.T) {
	r := testResponder(t)

	res, shared := r.answer(testQuery(t, "665F6E75B5D3A9C2._matterc._udp.local", dns.SRV, transport.Port), nil)
	if res == nil || shared {
		t.Fatalf("SRV: response %v, shared %v", res, shared)
	}
	res = reparse(t, res)
	if got := recordTypes(res.Answers()); !reflect.DeepEqual(got, []dns.Type{dns.SRV}) {
		t.Errorf("SRV: answers %v", got)
	}
	if got := recordTypes(res.Additions()); !reflect.DeepEqual(got, []dns.Type{dns.A, dns.AAAA}) {
		t.Errorf("SRV: additions %v, want the host addresses", got)
	}
	if !res.Answers()[0].UnicastResponse() {
		t.Error("SRV: a unique record carries the cache-flush bit")
	}

	res, _ = r.answer(testQuery(t, "665F6E75B5D3A9C2._matterc._udp.local", dns.TXT, transport.Port), nil)
	if res == nil || !reflect.DeepEqual(recordTypes(reparse(t, res).Answers()), []dns.Type{dns.TXT}) {
		t.Errorf("TXT: response %v", res)
	}

	res, _ = r.answer(testQuery(t, "B75AFB458ECD6D6F.local", dns.A, transport.Port), nil)
	if res == nil {
		t.Fatal("A: no response")
	}
	res = reparse(t, res)
	if got := recordTypes(res.Answers()); !reflect.DeepEqual(got, []dns.Type{dns.A}) {
		t.Errorf("A: answers %v", got)
	}
	if got := recordTypes(res.Additions()); !reflect.DeepEqual(got, []dns.Type{dns.AAAA}) {
		t.Errorf("A: additions %v, want the AAAA record", got)
	}

	res, _ = r.answer(testQuery(t, "665F6E75B5D3A9C2._matterc._udp.local", dns.ANY, transport.Port), nil)
	if res == nil || !reflect.DeepEqual(recordTypes(reparse(t, res).Answers()), []dns.Type{dns.SRV, dns.TXT}) {
		t.Errorf("ANY on the instance: response %v", res)
	}
}

func TestResponderAnswersServiceTypeEnumeration(t *testing.T) {
	res, _ := testResponder(t).answer(testQuery(t, "_services._dns-sd._udp.local", dns.PTR, transport.Port), nil)
	if res == nil {
		t.Fatal("no response")
	}
	answers := reparse(t, res).Answers()
	ptr, ok := answers[0].(dns.PTRRecord)
	if len(answers) != 1 || !ok || ptr.DomainName() != "_matterc._udp.local" {
		t.Fatalf("answers = %v, want a PTR to _matterc._udp.local", answers)
	}
}

func TestResponderIgnoresOtherNames(t *testing.T) {
	r := testResponder(t)
	for _, q := range []struct {
		name string
		typ  dns.Type
	}{
		{"_matter._tcp.local", dns.PTR},
		{"_L1._sub._matterc._udp.local", dns.PTR},
		{"other._matterc._udp.local", dns.SRV},
		{"665F6E75B5D3A9C2._matterc._udp.local", dns.A},
	} {
		if res, _ := r.answer(testQuery(t, q.name, q.typ, transport.Port), nil); res != nil {
			t.Errorf("%s %s: answered %v", q.name, q.typ, res)
		}
	}
}

func TestResponderSuppressesKnownAnswers(t *testing.T) {
	r := testResponder(t)
	svc := testLocalService()

	fresh, _ := dns.NewPTRResourceRecord(svc.ServiceName(), svc.FullName(), OtherRecordTTL)
	if res, _ := r.answer(testQuery(t, svc.ServiceName(), dns.PTR, transport.Port, fresh), nil); res != nil {
		t.Errorf("answered a PTR the querier already holds: %v", res)
	}

	stale, _ := dns.NewPTRResourceRecord(svc.ServiceName(), svc.FullName(), OtherRecordTTL/2-1)
	if res, _ := r.answer(testQuery(t, svc.ServiceName(), dns.PTR, transport.Port, stale), nil); res == nil {
		t.Error("did not refresh a PTR whose known TTL is below half")
	}
}

func TestResponderAnswersLegacyUnicast(t *testing.T) {
	query := testQuery(t, "_matterc._udp.local", dns.PTR, 50000)
	res, _ := testResponder(t).answer(query, nil)
	if res == nil {
		t.Fatal("no response")
	}
	res = reparse(t, res)
	if res.ID() != query.ID() || len(res.Questions()) != 1 {
		t.Errorf("legacy reply: ID %d (query %d), %d questions", res.ID(), query.ID(), len(res.Questions()))
	}
	for _, rr := range res.ResourceRecordSet() {
		if legacyUnicastTTL < rr.TTL() || rr.UnicastResponse() {
			t.Errorf("legacy reply record %s %s: TTL %d, cache flush %v", rr.Name(), rr.Type(), rr.TTL(), rr.UnicastResponse())
		}
	}
}

func TestAnnouncementAndGoodbye(t *testing.T) {
	svc := testLocalService()
	msg, err := announcement(svc, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	msg = reparse(t, msg)
	// PTR, 4 subtype PTRs, the enumeration PTR, SRV, TXT, A and AAAA.
	if n := len(msg.Answers()); n != 10 {
		t.Fatalf("announcement has %d records, want 10", n)
	}
	for _, rr := range msg.Answers() {
		want := uint(OtherRecordTTL)
		if rr.Type() == dns.SRV || rr.Type() == dns.A || rr.Type() == dns.AAAA {
			want = HostRecordTTL
		}
		if rr.TTL() != want {
			t.Errorf("announced %s %s with TTL %d, want %d", rr.Name(), rr.Type(), rr.TTL(), want)
		}
	}

	bye, err := announcement(svc, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	for _, rr := range reparse(t, bye).Answers() {
		if rr.TTL() != 0 {
			t.Errorf("goodbye %s %s has TTL %d, want 0", rr.Name(), rr.Type(), rr.TTL())
		}
	}
}

func TestLocalServiceValidate(t *testing.T) {
	for name, mutate := range map[string]func(*LocalService){
		"instance": func(s *LocalService) { s.Instance = "" },
		"host":     func(s *LocalService) { s.Host = "" },
		"service":  func(s *LocalService) { s.Service = "matterc._udp" },
		"proto":    func(s *LocalService) { s.Service = "_matterc._sctp" },
		"port":     func(s *LocalService) { s.Port = 70000 },
		"txt":      func(s *LocalService) { s.TXT = []string{string(make([]byte, 256))} },
		"address":  func(s *LocalService) { s.Addresses = []net.IP{{1, 2, 3}} },
		// avahi-publish-service takes the subtype name, which would be
		// published as "_S15._sub._matterc._udp._sub._matterc._udp.local".
		"subtype name":  func(s *LocalService) { s.Subtypes = []string{"_S15._sub._matterc._udp"} },
		"empty subtype": func(s *LocalService) { s.Subtypes = []string{""} },
		"host domain":   func(s *LocalService) { s.Host = "B75AFB458ECD6D6F.example.com" },
		"host only dot": func(s *LocalService) { s.Host = ".local" },
		"long subtype":  func(s *LocalService) { s.Subtypes = []string{"_" + strings.Repeat("x", 63)} },
	} {
		svc := testLocalService()
		mutate(svc)
		if err := svc.Validate(); err == nil {
			t.Errorf("%s: Validate() = nil, want an error", name)
		}
	}
	var nilService *LocalService
	if err := nilService.Validate(); err == nil {
		t.Error("Validate() on nil = nil, want an error")
	}
}

func TestLocalServiceHostName(t *testing.T) {
	for _, tc := range []struct {
		host, domain, want string
	}{
		{"B75AFB458ECD6D6F", "", "B75AFB458ECD6D6F.local"},
		{"B75AFB458ECD6D6F.local", "", "B75AFB458ECD6D6F.local"},
		{"B75AFB458ECD6D6F.local.", "", "B75AFB458ECD6D6F.local"},
		{"B75AFB458ECD6D6F.LOCAL", "", "B75AFB458ECD6D6F.local"},
		{"printer", "example.local", "printer.example.local"},
		{"printer.example.local", "example.local", "printer.example.local"},
	} {
		svc := testLocalService()
		svc.Host = tc.host
		svc.Domain = tc.domain
		if got := svc.HostName(); got != tc.want {
			t.Errorf("Host %q, Domain %q: HostName() = %q, want %q", tc.host, tc.domain, got, tc.want)
		}
		if err := svc.Validate(); err != nil {
			t.Errorf("Host %q, Domain %q: Validate() = %v", tc.host, tc.domain, err)
		}
	}
}

func TestServerRegisterAndDeregister(t *testing.T) {
	server := NewServer()
	if err := server.Register(context.Background(), testLocalService()); !errors.Is(err, ErrNotRunning) {
		t.Fatalf("Register() before Start = %v, want ErrNotRunning", err)
	}

	server = startTestServer(t)
	svc := testLocalService()
	if err := server.Register(context.Background(), svc); err != nil {
		t.Fatal(err)
	}
	// The server keeps its own copy.
	svc.Port = 1
	svc.TXT[0] = "D=1"
	registered := server.LocalServices()
	if len(registered) != 1 || registered[0].Port != 5540 || registered[0].TXT[0] != "D=3840" {
		t.Fatalf("LocalServices() = %+v, want the service as registered", registered)
	}

	// Registering the same instance again replaces it without probing its
	// names, which the server already holds.
	updated := testLocalService()
	updated.TXT = []string{"D=3840", "CM=0"}
	start := time.Now()
	if err := server.Register(context.Background(), updated); err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); probeInterval <= elapsed {
		t.Errorf("an update took %s, want no probe", elapsed)
	}
	if registered := server.LocalServices(); len(registered) != 1 || registered[0].TXT[1] != "CM=0" {
		t.Fatalf("LocalServices() after an update = %+v", registered)
	}

	if err := server.Deregister(testLocalService()); err != nil {
		t.Fatal(err)
	}
	if n := len(server.LocalServices()); n != 0 {
		t.Fatalf("%d services after Deregister, want 0", n)
	}
	if err := server.Register(context.Background(), &LocalService{}); err == nil {
		t.Fatal("Register() of an empty service = nil error")
	}
}

// TestResponderIfRegistered guards against a regression where the second
// announcement of a service deregistered right after it was registered
// followed its goodbye records and published it again.
func TestResponderIfRegistered(t *testing.T) {
	r := newResponder()
	svc := testLocalService()
	r.register(svc)

	called := false
	if !r.ifRegistered(svc, func() { called = true }) || !called {
		t.Fatal("ifRegistered did not run for a registered service")
	}

	replacement := testLocalService()
	r.register(replacement)
	if r.ifRegistered(svc, func() { t.Error("ran for a replaced copy") }) {
		t.Error("ifRegistered reported a replaced copy as registered")
	}

	r.deregister(replacement)
	if r.ifRegistered(replacement, func() { t.Error("ran for a deregistered service") }) {
		t.Error("ifRegistered reported a deregistered service as registered")
	}
}
