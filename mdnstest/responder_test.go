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

package mdnstest

import (
	"context"
	"testing"
	"time"

	"github.com/cybergarage/go-mdns/mdns"
)

// TestResponder registers a Matter commissionable service with a server and
// finds it with a client over the network, by its service type and by one
// of its subtypes.
func TestResponder(t *testing.T) {
	server := mdns.NewServer()
	if err := server.Start(); err != nil {
		t.Skipf("the server cannot bind the mDNS sockets here: %v", err)
	}
	defer server.Stop()

	svc := &mdns.LocalService{
		Instance: "665F6E75B5D3A9C2",
		Service:  "_matterc._udp",
		Subtypes: []string{"_L3840", "_S15", "_V65521", "_CM"},
		Host:     "B75AFB458ECD6D6F",
		Port:     5540,
		TXT:      []string{"D=3840", "CM=1", "VP=65521+32769"},
	}
	if err := server.Register(context.Background(), svc); err != nil {
		t.Fatal(err)
	}

	client := mdns.NewClient()
	if err := client.Start(); err != nil {
		t.Skipf("the client cannot bind the mDNS sockets here: %v", err)
	}
	defer client.Stop()

	for _, query := range []mdns.Query{
		mdns.NewQuery(mdns.WithQueryService("_matterc._udp")),
		mdns.NewQuery(mdns.WithQueryService("_matterc._udp"), mdns.WithQuerySubtype("_L3840")),
	} {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		services, err := client.Query(ctx, query)
		cancel()
		if err != nil {
			t.Fatalf("Query(%s) error = %v", query, err)
		}
		var found mdns.Service
		for _, s := range services {
			if s.FullName() == svc.FullName() {
				found = s
			}
		}
		if found == nil {
			t.Fatalf("Query(%s) did not find %s among %d services", query, svc.FullName(), len(services))
		}
		if found.Port() != svc.Port {
			t.Errorf("Query(%s): port %d, want %d", query, found.Port(), svc.Port)
		}
		if attr, ok := found.LookupResourceAttribute("VP"); !ok || attr.Value() != "65521+32769" {
			t.Errorf("Query(%s): TXT VP = %v, %v", query, attr, ok)
		}
		if len(found.Addresses()) == 0 {
			t.Errorf("Query(%s): no address", query)
		}
	}
}
