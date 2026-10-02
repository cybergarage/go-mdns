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
Discovery (RFC 6763) client and server.

# Status

[Client] is a querier: it browses, resolves and looks up the services and the
hosts which the other responders advertise. [Server] is a responder: it
publishes the services registered with it, after it probes their names for
conflicts with the other nodes on the link (RFC 6762, 8.1 and 9).

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

# Publishing

[Server.Register] publishes a [LocalService]. The server probes its instance
name and host name, announces it, answers the queries for its type,
subtypes, instance and host, and sends goodbye records when it is
deregistered or the server stops. The server must be started first.

	server := mdns.NewServer()
	if err := server.Start(); err != nil {
		return err
	}
	defer server.Stop()

	err := server.Register(ctx, &mdns.LocalService{
		Instance: "665F6E75B5D3A9C2",
		Service:  "_matterc._udp",
		Subtypes: []string{"_L3840", "_S15"},
		Host:     "B75AFB458ECD6D6F",
		Port:     5540,
		TXT:      []string{"D=3840", "CM=1"},
	})

# Name conflicts

The server does not rename a service whose name another node holds: the
application chooses the new name, since its naming rules depend on the
service, such as a new random instance name for a Matter commissionable node.
[Server.Register] returns a [*ConflictError] which wraps [ErrConflict] when
the probe finds a name held, and the service is not published:

	var conflict *mdns.ConflictError
	if errors.As(err, &conflict) && conflict.IsInstanceConflict() {
		// Register the service again with a new instance name.
	}

A published service whose names another node claims later is probed again,
and when the other node holds them, it is withdrawn and reported to the
handler given by [WithServerConflictHandler].

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
