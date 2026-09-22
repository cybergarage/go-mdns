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

/*
Package mdns implements a Multicast DNS (RFC 6762) and DNS-Based Service
Discovery (RFC 6763) client.

# Status

The package is a client (querier). It browses, resolves and looks up the
services and the hosts which the other responders advertise. The server
(responder) side is under development: [Server] only listens for the messages,
and it registers no service and answers no query yet.

# Browsing

[Client.Browse] keeps querying a service type, and reports a service when it is
added, updated or removed, so that a caller can follow the services which come
and go. It blocks until the context is done.

	client := mdns.NewClient()
	if err := client.Start(); err != nil {
		return err
	}
	defer client.Stop()

	query := mdns.NewQuery(
		mdns.WithQueryService("_matterc._udp"),
	)

	err := client.Browse(ctx, query, func(event mdns.ServiceEvent) {
		log.Printf("%s %s", event.Type, event.Service.FullName())
	})

[Client.Query] sends a one-shot query instead, and returns the services which
answered before the context is done.

# Resolving

[Client.Resolve] resolves a service instance name to its host, port, addresses
and TXT attributes, and [Client.LookupHost] resolves a host name to its
addresses.

An IPv6 link-local address is ambiguous without its zone (RFC 4007), so
[Service.Addrs] returns the addresses with the interface which the service was
discovered on. Use them to connect to a service which advertises only a
link-local address.
*/
package mdns
