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

package mdnstest

import (
	"context"
	"errors"
	"net"
	"strconv"
	"sync"
	"testing"

	"github.com/cybergarage/go-mdns/mdns"
)

// startProbeServer starts a server, which plays a node on the link.
func startProbeServer(t *testing.T) *mdns.Server {
	t.Helper()
	server := mdns.NewServer()
	if err := server.Start(); err != nil {
		t.Skipf("the server cannot bind the mDNS sockets here: %v", err)
	}
	t.Cleanup(func() { _ = server.Stop() })
	return server
}

// probeTestService returns a service with random names, published with its
// own address so that two servers of this host hold different records.
func probeTestService(t *testing.T, last byte) *mdns.LocalService {
	t.Helper()
	suffix := randomSuffix(t)
	return &mdns.LocalService{
		Instance:  "go-mdns-probe-" + suffix,
		Service:   "_go-mdns-test._tcp",
		Host:      "go-mdns-probe-" + suffix,
		Port:      20000 + int(last),
		TXT:       []string{"node=" + strconv.Itoa(int(last))},
		Addresses: []net.IP{net.IPv4(192, 0, 2, last)},
	}
}

// TestProbeConflict publishes a service with one server, and registers the
// same names with another, which finds them held over the network.
func TestProbeConflict(t *testing.T) {
	first := startProbeServer(t)
	second := startProbeServer(t)

	svc := probeTestService(t, 1)
	if err := first.Register(context.Background(), svc); err != nil {
		t.Fatal(err)
	}

	t.Run("instance", func(t *testing.T) {
		other := probeTestService(t, 2)
		other.Instance = svc.Instance
		err := second.Register(context.Background(), other)
		var conflict *mdns.ConflictError
		if !errors.As(err, &conflict) || !conflict.IsInstanceConflict() {
			t.Fatalf("Register() = %v, want an instance conflict", err)
		}
	})

	t.Run("host", func(t *testing.T) {
		other := probeTestService(t, 2)
		other.Host = svc.Host
		err := second.Register(context.Background(), other)
		var conflict *mdns.ConflictError
		if !errors.As(err, &conflict) || !conflict.IsHostConflict() {
			t.Fatalf("Register() = %v, want a host conflict", err)
		}
	})

	t.Run("unique", func(t *testing.T) {
		if err := second.Register(context.Background(), probeTestService(t, 2)); err != nil {
			t.Fatalf("Register() of unique names = %v", err)
		}
	})
}

// TestSimultaneousProbe registers the same names with two servers at once:
// the tiebreak lets exactly one of them publish the service (RFC 6762, 8.2).
func TestSimultaneousProbe(t *testing.T) {
	servers := []*mdns.Server{startProbeServer(t), startProbeServer(t)}
	base := probeTestService(t, 1)

	errs := make([]error, len(servers))
	var wg sync.WaitGroup
	for i, server := range servers {
		svc := *base
		svc.Port = base.Port + i
		svc.Addresses = []net.IP{net.IPv4(192, 0, 2, byte(10+i))}
		wg.Go(func() {
			errs[i] = server.Register(context.Background(), &svc)
		})
	}
	wg.Wait()

	published := 0
	for i, err := range errs {
		switch {
		case err == nil:
			published++
		case errors.Is(err, mdns.ErrConflict):
		default:
			t.Errorf("server %d: Register() = %v", i, err)
		}
	}
	if published != 1 {
		t.Fatalf("%d servers published the service (errors %v), want 1", published, errs)
	}
}
