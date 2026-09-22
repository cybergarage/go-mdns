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
	"github.com/cybergarage/go-mdns/mdns/dns"
)

type queryImp struct {
	name            string
	subtype         string
	service         string
	domain          string
	typ             Type
	unicastResponse bool
	handler         MessageHandler
}

// QueryOption represents a query option.
type QueryOption func(*queryImp)

// WithQueryName sets the full question name of the query. The name is used as
// it is, and the subtype, the service and the domain options are ignored. Use
// it to query a service instance name or a host name directly.
func WithQueryName(name string) QueryOption {
	return func(q *queryImp) {
		q.name = name
	}
}

// WithQuerySubtype sets the subtype of the query.
func WithQuerySubtype(subtype string) QueryOption {
	return func(q *queryImp) {
		q.subtype = subtype
	}
}

// WithQueryService sets the service name of the query.
func WithQueryService(service string) QueryOption {
	return func(q *queryImp) {
		q.service = service
	}
}

// WithQueryDomain sets the domain name of the query.
func WithQueryDomain(domain string) QueryOption {
	return func(q *queryImp) {
		q.domain = domain
	}
}

// WithQueryType sets the question record type of the query.
//
// RFC 6763: 4.1. Structured Service Instance Names
// A service type is browsed with PTR, a service instance is resolved with SRV
// and TXT, and a host name is resolved with A and AAAA.
func WithQueryType(t Type) QueryOption {
	return func(q *queryImp) {
		q.typ = t
	}
}

// WithQueryUnicastResponse sets the unicast response bit (QU) of the query.
//
// RFC 6762: 5.4. Questions Requesting Unicast Responses
// A querier sets the unicast response bit when it has not yet joined the
// multicast group, such as the first query after a node starts. The bit is not
// set by default because a multicast response lets the other nodes on the link
// update their caches.
func WithQueryUnicastResponse(flag bool) QueryOption {
	return func(q *queryImp) {
		q.unicastResponse = flag
	}
}

// WithQueryMessageHandler sets the message handler of the query.
func WithQueryMessageHandler(handler MessageHandler) QueryOption {
	return func(q *queryImp) {
		q.handler = handler
	}
}

// NewQuery returns a new query instance with the specified options.
func NewQuery(opts ...QueryOption) Query {
	q := &queryImp{
		name:            "",
		subtype:         "",
		service:         "",
		domain:          DefaultQueryDomain,
		typ:             DefaultQueryType,
		unicastResponse: false,
		handler:         nil,
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

// Name returns the full question name of the query.
func (q *queryImp) Name() string {
	if 0 < len(q.name) {
		return q.name
	}
	labels := []string{}
	if 0 < len(q.subtype) {
		labels = append(labels, q.subtype, Subtype)
	}
	if 0 < len(q.service) {
		labels = append(labels, q.service)
	}
	labels = append(labels, q.domain)
	return dns.NewNameWithStrings(labels...)
}

// Subtype returns the subtype of the query.
func (q *queryImp) Subtype() string {
	return q.subtype
}

// Service returns the service name of the query.
func (q *queryImp) Service() string {
	return q.service
}

// Domain returns the domain name of the query.
func (q *queryImp) Domain() string {
	return q.domain
}

// Type returns the question record type of the query.
func (q *queryImp) Type() Type {
	return q.typ
}

// UnicastResponse returns true if the query requests a unicast response.
func (q *queryImp) UnicastResponse() bool {
	return q.unicastResponse
}

// MessageHandler returns the message handler of the query if set.
func (q *queryImp) MessageHandler() (MessageHandler, bool) {
	if q.handler == nil {
		return nil, false
	}
	return q.handler, true
}

// String returns the string representation of the query.
func (q *queryImp) String() string {
	return q.Name()
}
