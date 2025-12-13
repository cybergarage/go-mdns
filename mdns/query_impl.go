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
	subtype  string
	services []string
	domain   string
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
		q.services = []string{service}
	}
}

// WithQueryServices sets the service name of the query.
func WithQueryServices(services ...string) QueryOption {
	return func(q *queryImp) {
		q.services = services
	}
}

// WithQueryDomain sets the domain name of the query.
func WithQueryDomain(domain string) QueryOption {
	return func(q *queryImp) {
		q.domain = domain
	}
}

// NewQuery returns a new query instance with the specified options.
func NewQuery(opts ...QueryOption) Query {
	q := &queryImp{
		subtype:  "",
		services: []string{DefaultQueryService},
		domain:   DefaultQueryDomain,
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

// Services returns the service names of the query.
func (q *queryImp) Services() []string {
	return q.services
}

// Domain returns the domain name of the query.
func (q *queryImp) Domain() string {
	return q.domain
}

// String returns the string representation of the query.
func (q *queryImp) String() string {
	labels := []string{}
	if 0 < len(q.subtype) {
		labels = append(labels, q.subtype)
	}
	if 0 < len(q.services) {
		labels = append(labels, q.services...)
	}
	labels = append(labels, q.domain)
	return dns.NewNameWithStrings(labels...)
}
