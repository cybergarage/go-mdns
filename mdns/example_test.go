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

package mdns_test

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/cybergarage/go-mdns/mdns"
)

// Browse reports the services of a service type as they are added, updated and
// removed. It blocks until the context is done.
func ExampleClient_Browse() {
	client := mdns.NewClient()
	if err := client.Start(); err != nil {
		return
	}
	defer client.Stop()

	query := mdns.NewQuery(
		mdns.WithQueryService("_matterc._udp"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client.Browse(ctx, query, func(event mdns.ServiceEvent) {
		service := event.Service
		fmt.Printf("%s %s (%s:%d)\n", event.Type, service.FullName(), service.Host(), service.Port())
	})
}

// A subtype browses only the instances which are registered with it, such as
// the Matter discriminator subtypes.
func ExampleClient_Browse_subtype() {
	client := mdns.NewClient()
	if err := client.Start(); err != nil {
		return
	}
	defer client.Stop()

	query := mdns.NewQuery(
		mdns.WithQuerySubtype("_S3"),
		mdns.WithQueryService("_matterc._udp"),
	)

	client.Browse(context.Background(), query, func(event mdns.ServiceEvent) {
		fmt.Println(event.String())
	})
}

// Query returns the services which answered before the context is done.
func ExampleClient_Query() {
	client := mdns.NewClient()
	if err := client.Start(); err != nil {
		return
	}
	defer client.Stop()

	// RFC 6763 (7.2) enumerates the advertised service types.
	query := mdns.NewQuery(
		mdns.WithQueryService(mdns.ServiceTypeEnumerationName),
	)

	services, err := client.Query(context.Background(), query)
	if err != nil {
		return
	}

	for _, service := range services {
		fmt.Println(service.FullName())
	}
}

// Resolve returns the host, the port, the addresses and the TXT attributes of a
// service instance.
func ExampleClient_Resolve() {
	client := mdns.NewClient()
	if err := client.Start(); err != nil {
		return
	}
	defer client.Stop()

	service, err := client.Resolve(context.Background(), "DD200C20D25AE5F7._matterc._udp.local")
	if err != nil {
		return
	}

	// The IPv6 scoped addressing zone is set for the link-local addresses,
	// so the address can be used to connect to the service.
	for _, addr := range service.Addrs() {
		fmt.Println(addr.String())
	}

	if attr, ok := service.LookupResourceAttribute("CM"); ok {
		fmt.Printf("commissioning mode: %s\n", attr.Value())
	}
}

// LookupHost resolves a host name without the DNS-SD service records.
func ExampleClient_LookupHost() {
	client := mdns.NewClient()
	if err := client.Start(); err != nil {
		return
	}
	defer client.Stop()

	addrs, err := client.LookupHost(context.Background(), "macmini.local")
	if err != nil {
		return
	}

	for _, addr := range addrs {
		fmt.Println(addr.String())
	}
}

// The client listens on all available interfaces by default.
func ExampleNewClient_options() {
	ifi, err := net.InterfaceByName("en0")
	if err != nil {
		return
	}

	client := mdns.NewClient(
		mdns.WithClientInterfaces(ifi),
		mdns.WithClientIPv4Enabled(false),
		mdns.WithClientQueryTimeout(10*time.Second),
	)

	if err := client.Start(); err != nil {
		return
	}
	defer client.Stop()
}
