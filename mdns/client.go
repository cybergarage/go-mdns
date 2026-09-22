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
	"net"
)

// Client represents a client node instance.
type Client interface {
	// Start starts the client.
	Start() error
	// Stop stops the client.
	Stop() error
	// Restart restarts the client.
	Restart() error
	// RegisterMessageHandler adds a message handler to the client.
	RegisterMessageHandler(handler MessageHandler)
	// UnRegisterMessageHandler removes a message handler from the client.
	UnRegisterMessageHandler(handler MessageHandler)
	// Query sends a question message to the multicast address, and returns the
	// services which responded until the context is done. The query is
	// retransmitted while it is waiting, as RFC 6762 (5.2) requires.
	Query(ctx context.Context, query Query) ([]Service, error)
	// Browse browses the services of the query, and calls the handler when a
	// service is added, updated or removed. Browse blocks until the context
	// is done. Use Browse instead of Query to follow the services which come
	// and go, such as the Matter nodes on a link.
	Browse(ctx context.Context, query Query, handler ServiceHandler) error
	// Resolve resolves the specified service instance name, such as
	// "instance._matter._tcp.local", and returns the service with its host,
	// port, TXT attributes and addresses.
	Resolve(ctx context.Context, name string) (Service, error)
	// LookupHost resolves the specified host name, such as "device.local",
	// and returns its addresses. The IPv6 scoped addressing zone is set for
	// the link-local addresses.
	LookupHost(ctx context.Context, host string) ([]*net.IPAddr, error)
	// Services returns the services which the client has discovered.
	Services() []Service
}
