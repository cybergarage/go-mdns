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
	subtype string
	service string
	domain  string
	handler MessageHandler
}

// QueryOption represents a query option.
type QueryOption func(*queryImp)

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

// WithQueryMessageHandler sets the message handler of the query.
func WithQueryMessageHandler(handler MessageHandler) QueryOption {
	return func(q *queryImp) {
		q.handler = handler
	}
}

// NewQuery returns a new query instance with the specified options.
func NewQuery(opts ...QueryOption) Query {
	q := &queryImp{
		subtype: "",
		service: "",
		domain:  DefaultQueryDomain,
		handler: nil,
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
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

// MessageHandler returns the message handler of the query if set.
func (q *queryImp) MessageHandler() (MessageHandler, bool) {
	if q.handler == nil {
		return nil, false
	}
	return q.handler, true
}

// String returns the string representation of the query.
func (q *queryImp) String() string {
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
