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
	"strings"
)

type queryImp struct {
	services []string
	domain   string
}

// QueryOption represents a query option.
type QueryOption func(*queryImp)

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
		services: []string{DefaultBrowsingService},
		domain:   DefaultQueryDomain,
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
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
	var s strings.Builder
	for i, service := range q.services {
		if 0 < i {
			s.WriteString(",")
		}
		s.WriteString(strings.Join([]string{service, q.domain}, queryNameSep))
	}
	return s.String()
}
