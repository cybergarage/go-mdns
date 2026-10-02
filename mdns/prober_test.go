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
	"errors"
	"net"
	"testing"
	"time"

	"github.com/cybergarage/go-mdns/mdns/dns"
)

// startTestServer starts a server, and skips the test when the mDNS sockets
// cannot be bound here.
func startTestServer(t *testing.T, opts ...ServerOption) *Server {
	t.Helper()
	server := NewServer(opts...)
	if err := server.Start(); err != nil {
		t.Skipf("the server cannot bind the mDNS sockets here: %v", err)
	}
	t.Cleanup(func() { _ = server.Stop() })
	return server
}

// testConflictService returns a service with names no other test uses.
func testConflictService() *LocalService {
	svc := testLocalService()
	svc.Instance = "go-mdns-probe-test"
	svc.Host = "go-mdns-probe-test"
	svc.Addresses = []net.IP{net.IPv4(192, 0, 2, 1)}
	return svc
}

// srvResponse builds a response, as it arrives from the network, which
// claims the instance name of svc for port.
func srvResponse(t *testing.T, svc *LocalService, port uint16) dns.Message {
	t.Helper()
	srv, err := dns.NewSRVResourceRecord(svc.FullName(), 0, 0, port, "other-node.local", HostRecordTTL)
	if err != nil {
		t.Fatal(err)
	}
	return reparse(t, dns.NewResponseMessage(dns.WithMessageAnswers(srv)))
}

// addrResponse builds a response which claims the host name of svc for ip.
func addrResponse(t *testing.T, svc *LocalService, ip net.IP) dns.Message {
	t.Helper()
	addr, err := dns.NewAddressResourceRecord(svc.HostName(), ip, HostRecordTTL)
	if err != nil {
		t.Fatal(err)
	}
	return reparse(t, dns.NewResponseMessage(dns.WithMessageAnswers(addr)))
}

// probeQuery builds a probe of another node, as it arrives from the network,
// which proposes the records of svc for its instance name, with port in the
// SRV record.
func probeQuery(t *testing.T, svc *LocalService, port uint16) dns.Message {
	t.Helper()
	srv, err := dns.NewSRVResourceRecord(svc.FullName(), 0, 0, port, svc.HostName(), HostRecordTTL)
	if err != nil {
		t.Fatal(err)
	}
	txt, err := dns.NewTXTResourceRecord(svc.FullName(), svc.TXT, OtherRecordTTL)
	if err != nil {
		t.Fatal(err)
	}
	q := dns.NewQuestion(dns.WithQuestionName(svc.FullName()), dns.WithQuestionType(dns.ANY), dns.WithQuestionClass(dns.IN))
	msg := dns.NewRequestMessage(dns.WithMessageQuestions(q), dns.WithMessageNameServers(srv, txt))
	parsed, err := dns.NewMessageWithBytes(msg.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	return parsed
}

func TestConflictError(t *testing.T) {
	svc := testLocalService()
	var err error = &ConflictError{Service: svc, Name: svc.FullName()}
	if !errors.Is(err, ErrConflict) {
		t.Error("errors.Is(err, ErrConflict) = false")
	}
	var conflict *ConflictError
	if !errors.As(err, &conflict) || !conflict.IsInstanceConflict() || conflict.IsHostConflict() {
		t.Errorf("an instance conflict: As %v, instance %v, host %v", conflict, conflict.IsInstanceConflict(), conflict.IsHostConflict())
	}
	conflict = &ConflictError{Service: svc, Name: "b75afb458ecd6d6f.local"}
	if conflict.IsInstanceConflict() || !conflict.IsHostConflict() {
		t.Error("a host conflict is not reported as one")
	}
}

func TestProbeMessage(t *testing.T) {
	svc := testLocalService()
	p := newProbe(svc, []string{svc.FullName(), svc.HostName()})
	msg, err := probeMessage(p, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	msg = reparse(t, msg)
	if !msg.IsQuery() {
		t.Fatal("a probe is not a query")
	}
	questions := msg.Questions()
	if len(questions) != 2 {
		t.Fatalf("%d questions, want 2", len(questions))
	}
	for _, q := range questions {
		if q.Type() != dns.ANY || !q.UnicastResponse() {
			t.Errorf("question %s: type %v, QU %v; want ANY with QU", q.Name(), q.Type(), q.UnicastResponse())
		}
	}
	// SRV, TXT and the two addresses.
	authorities := msg.NameServers()
	if len(authorities) != 4 {
		t.Fatalf("%d authority records, want 4", len(authorities))
	}
	for _, rr := range authorities {
		if rr.UnicastResponse() {
			t.Errorf("%s: the cache-flush bit is set in the Authority section", rr.Name())
		}
	}

	// The later probes ask for multicast responses, and only the names
	// being claimed are probed.
	p = newProbe(svc, []string{svc.FullName()})
	msg, err = probeMessage(p, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	msg = reparse(t, msg)
	if len(msg.Questions()) != 1 || msg.Questions()[0].UnicastResponse() {
		t.Errorf("questions %v, want one without QU", msg.Questions())
	}
	if types := recordTypes(msg.NameServers()); len(types) != 2 {
		t.Errorf("authority record types %v, want SRV and TXT", types)
	}
}

func TestCompareProbeRecords(t *testing.T) {
	newSRV := func(port uint16) dns.ResourceRecord {
		srv, err := dns.NewSRVResourceRecord("a._http._tcp.local", 0, 0, port, "host.local", 120)
		if err != nil {
			t.Fatal(err)
		}
		return srv
	}
	newA := func(last byte) dns.ResourceRecord {
		a, err := dns.NewAResourceRecord("host.local", net.IPv4(192, 0, 2, last), 120)
		if err != nil {
			t.Fatal(err)
		}
		return a
	}

	for _, tc := range []struct {
		name         string
		ours, theirs []dns.ResourceRecord
		want         int
	}{
		{"lower data loses", []dns.ResourceRecord{newSRV(1)}, []dns.ResourceRecord{newSRV(2)}, -1},
		{"higher data wins", []dns.ResourceRecord{newSRV(2)}, []dns.ResourceRecord{newSRV(1)}, 1},
		{"identical", []dns.ResourceRecord{newSRV(1)}, []dns.ResourceRecord{newSRV(1)}, 0},
		// A (1) sorts before SRV (33), so the A records are compared first.
		{"lower type loses", []dns.ResourceRecord{newA(9)}, []dns.ResourceRecord{newSRV(1)}, -1},
		{"sorted before compared", []dns.ResourceRecord{newA(2), newA(1)}, []dns.ResourceRecord{newA(1), newA(3)}, -1},
		{"more records win", []dns.ResourceRecord{newA(1), newA(2)}, []dns.ResourceRecord{newA(1)}, 1},
	} {
		got := compareProbeRecords(tc.ours, tc.theirs)
		if (got < 0) != (tc.want < 0) || (got > 0) != (tc.want > 0) {
			t.Errorf("%s: compareProbeRecords() = %d, want the sign of %d", tc.name, got, tc.want)
		}
	}
}

func TestIsOwnRecord(t *testing.T) {
	svc := testLocalService()
	records, err := svc.records(nil, OtherRecordTTL, HostRecordTTL)
	if err != nil {
		t.Fatal(err)
	}
	localAddrs := svc.Addresses
	res := reparse(t, dns.NewResponseMessage(dns.WithMessageAnswers(records.all()...)))
	for _, rr := range res.Answers() {
		switch rr.Type() {
		case dns.SRV, dns.TXT, dns.A, dns.AAAA:
			if !isOwnRecord(svc, rr, localAddrs) {
				t.Errorf("%s %v: not reported as an own record", rr.Name(), rr.Type())
			}
		}
	}

	other := testLocalService()
	other.Port = 1
	other.TXT = []string{"D=1"}
	other.Addresses = []net.IP{net.IPv4(192, 0, 2, 99)}
	records, err = other.records(nil, OtherRecordTTL, HostRecordTTL)
	if err != nil {
		t.Fatal(err)
	}
	res = reparse(t, dns.NewResponseMessage(dns.WithMessageAnswers(records.all()...)))
	for _, rr := range res.Answers() {
		switch rr.Type() {
		case dns.SRV, dns.TXT, dns.A:
			if isOwnRecord(svc, rr, localAddrs) {
				t.Errorf("%s %v: a record of another node is reported as an own record", rr.Name(), rr.Type())
			}
		}
	}
}

func TestResponderProbes(t *testing.T) {
	r := newResponder()
	svc := testLocalService()
	if names := r.namesToProbe(svc); len(names) != 2 {
		t.Fatalf("namesToProbe() of a new service = %v, want the instance and the host", names)
	}

	first := newProbe(svc, r.namesToProbe(svc))
	r.addProbe(first)
	// A later registration of the instance replaces the earlier probe.
	second := newProbe(copyLocalService(svc), r.namesToProbe(svc))
	r.addProbe(second)
	select {
	case <-first.canceled:
		if !errors.Is(first.cancelErr, ErrDeregistered) {
			t.Errorf("the replaced probe ended with %v", first.cancelErr)
		}
	default:
		t.Fatal("the replaced probe is not canceled")
	}
	if len(r.localServices()) != 0 {
		t.Fatal("a service being probed is published")
	}

	if !r.establish(second) {
		t.Fatal("establish() = false")
	}
	if r.establish(first) {
		t.Fatal("establish() of a canceled probe = true")
	}
	if names := r.namesToProbe(svc); len(names) != 0 {
		t.Errorf("namesToProbe() of a published service = %v, want none", names)
	}
	other := testLocalService()
	other.Instance = "other"
	if names := r.namesToProbe(other); len(names) != 1 || names[0] != other.FullName() {
		t.Errorf("namesToProbe() of a service on a held host = %v, want only the instance", names)
	}

	// Deregistration cancels a probe.
	p := newProbe(other, r.namesToProbe(other))
	r.addProbe(p)
	if _, ok := r.deregister(other); ok {
		t.Error("deregister() of a service being probed reports a published service")
	}
	select {
	case <-p.canceled:
	default:
		t.Fatal("deregister() did not cancel the probe")
	}

	// A published service is probed again without being answered for, and
	// restored when the server stops.
	published := r.localServices()[0]
	p, ok := r.beginReprobe(published)
	if !ok || len(r.localServices()) != 0 || len(p.names) != 2 {
		t.Fatalf("beginReprobe() = %v, %v with %d published", p, ok, len(r.localServices()))
	}
	r.restore(p)
	if len(r.localServices()) != 1 || len(r.probes()) != 0 {
		t.Fatal("restore() did not publish the service again")
	}
}

func TestCheckProbeTiebreak(t *testing.T) {
	server := NewServer()
	svc := testConflictService()
	p := newProbe(svc, []string{svc.FullName()})
	server.responder.addProbe(p)

	lost := func() bool {
		select {
		case <-p.lost:
			return true
		default:
			return false
		}
	}

	// Its own probe received back is not a simultaneous probe.
	own, err := probeMessage(p, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	server.checkProbe(reparse(t, own))
	if lost() {
		t.Fatal("the own probe made the probe defer")
	}

	server.checkProbe(probeQuery(t, svc, uint16(svc.Port-1))) // nolint: gosec
	if lost() {
		t.Fatal("a probe which compares lower made the probe defer")
	}
	server.checkProbe(probeQuery(t, svc, uint16(svc.Port+1))) // nolint: gosec
	if !lost() {
		t.Fatal("a probe which compares higher did not make the probe defer")
	}
}

func TestCheckResponseWhileProbing(t *testing.T) {
	server := NewServer()
	svc := testConflictService()
	p := newProbe(svc, []string{svc.FullName(), svc.HostName()})
	server.responder.addProbe(p)

	conflict := func() string {
		select {
		case name := <-p.conflict:
			return name
		default:
			return ""
		}
	}

	records, err := svc.records(nil, OtherRecordTTL, HostRecordTTL)
	if err != nil {
		t.Fatal(err)
	}
	server.checkResponse(reparse(t, dns.NewResponseMessage(dns.WithMessageAnswers(records.all()...))))
	if name := conflict(); name != "" {
		t.Fatalf("its own records conflict with %s", name)
	}

	server.checkResponse(srvResponse(t, svc, uint16(svc.Port+1))) // nolint: gosec
	if name := conflict(); name != svc.FullName() {
		t.Fatalf("an SRV record of another node: conflict %q, want %q", name, svc.FullName())
	}
	server.checkResponse(addrResponse(t, svc, net.IPv4(192, 0, 2, 200)))
	if name := conflict(); name != svc.HostName() {
		t.Fatalf("an address of another node: conflict %q, want %q", name, svc.HostName())
	}

	// A goodbye record claims nothing.
	goodbye, err := dns.NewSRVResourceRecord(svc.FullName(), 0, 0, 1, "other-node.local", 0)
	if err != nil {
		t.Fatal(err)
	}
	server.checkResponse(reparse(t, dns.NewResponseMessage(dns.WithMessageAnswers(goodbye))))
	if name := conflict(); name != "" {
		t.Fatalf("a goodbye record conflicts with %s", name)
	}
}

// feed passes msg to the server until stop is closed, as another node which
// keeps answering does.
func feed(server *Server, msg dns.Message, stop <-chan struct{}) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			_, _ = server.MessageReceived(msg)
		}
	}
}

func TestServerRegisterConflict(t *testing.T) {
	server := startTestServer(t)
	svc := testConflictService()

	stop := make(chan struct{})
	go feed(server, srvResponse(t, svc, uint16(svc.Port+1)), stop) // nolint: gosec
	err := server.Register(context.Background(), svc)
	close(stop)

	var conflict *ConflictError
	if !errors.As(err, &conflict) || !conflict.IsInstanceConflict() {
		t.Fatalf("Register() = %v, want an instance conflict", err)
	}
	if n := len(server.LocalServices()); n != 0 {
		t.Fatalf("%d services published after a conflict", n)
	}

	// The application registers the service again with a new name.
	svc.Instance = "go-mdns-probe-test-2"
	if err := server.Register(context.Background(), svc); err != nil {
		t.Fatalf("Register() with a new name = %v", err)
	}
}

func TestServerRegisterCanceled(t *testing.T) {
	server := startTestServer(t)

	t.Run("context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if err := server.Register(ctx, testConflictService()); !errors.Is(err, context.Canceled) {
			t.Fatalf("Register() = %v, want context.Canceled", err)
		}
	})

	t.Run("deregister", func(t *testing.T) {
		svc := testConflictService()
		errs := make(chan error, 1)
		go func() { errs <- server.Register(context.Background(), svc) }()
		time.Sleep(probeInterval)
		if err := server.Deregister(svc); err != nil {
			t.Fatal(err)
		}
		if err := <-errs; !errors.Is(err, ErrDeregistered) {
			t.Fatalf("Register() = %v, want ErrDeregistered", err)
		}
	})

	t.Run("stop", func(t *testing.T) {
		svc := testConflictService()
		errs := make(chan error, 1)
		go func() { errs <- server.Register(context.Background(), svc) }()
		time.Sleep(probeInterval)
		if err := server.Stop(); err != nil {
			t.Fatal(err)
		}
		if err := <-errs; !errors.Is(err, ErrNotRunning) {
			t.Fatalf("Register() = %v, want ErrNotRunning", err)
		}
		if n := len(server.LocalServices()); n != 0 {
			t.Fatalf("%d services published after Stop", n)
		}
	})
}

func TestServerConflictAfterPublished(t *testing.T) {
	conflicts := make(chan *ConflictError, 1)
	server := startTestServer(t, WithServerConflictHandler(func(err *ConflictError) {
		conflicts <- err
	}))
	svc := testConflictService()
	if err := server.Register(context.Background(), svc); err != nil {
		t.Fatal(err)
	}

	// A single conflicting record, such as a stale one, makes the server
	// probe the service again, and nobody holds the name now.
	server.MessageReceived(addrResponse(t, svc, net.IPv4(192, 0, 2, 200))) // nolint: errcheck
	if n := len(server.LocalServices()); n != 0 {
		t.Fatal("a service is answered for while it is probed again")
	}
	deadline := time.Now().Add(3 * time.Second)
	for len(server.LocalServices()) == 0 && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
	}
	if n := len(server.LocalServices()); n != 1 {
		t.Fatal("the service is not published again after the probe")
	}

	// Another node which keeps holding the name withdraws the service.
	stop := make(chan struct{})
	defer close(stop)
	go feed(server, addrResponse(t, svc, net.IPv4(192, 0, 2, 200)), stop)
	select {
	case conflict := <-conflicts:
		if !conflict.IsHostConflict() {
			t.Errorf("conflict %v, want a host conflict", conflict)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the conflict handler is not called")
	}
	if n := len(server.LocalServices()); n != 0 {
		t.Fatalf("%d services published after a conflict", n)
	}
}

func TestServerRestartProbesAgain(t *testing.T) {
	server := startTestServer(t)
	svc := testConflictService()
	if err := server.Register(context.Background(), svc); err != nil {
		t.Fatal(err)
	}
	if err := server.Restart(); err != nil {
		t.Fatal(err)
	}
	if n := len(server.LocalServices()); n != 0 {
		t.Fatal("a service is answered for before it is probed again")
	}
	// A server stopped while it probes keeps the service for the next
	// start.
	time.Sleep(probeInterval)
	if err := server.Restart(); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(3 * time.Second)
	for len(server.LocalServices()) == 0 && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
	}
	if n := len(server.LocalServices()); n != 1 {
		t.Fatal("the service is not published again after a restart")
	}
}
