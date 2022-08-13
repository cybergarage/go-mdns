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

// RFC 6765: 11. Discovery of Browsing and Registration Domains (Domain Enumeration).
const (
	// A list of domains recommended for browsing.
	RecommendedBrowsingService = "b._dns-sd._udp"
	// A single recommended default domain for browsing.
	DefaultBrowsingService = "db._dns-sd._udp"
	// A list of domains recommended for registering services using Dynamic Update.
	RecommendedRegisteringService = "r._dns-sd._udp"
	// A single recommended default domain for registering services.
	DefaultRegisteringService = "dr._dns-sd._udp"
	// The "legacy browsing" or "automatic browsing" domain(s).
	AutomaticBrowsingService = "lb._dns-sd._udp"
)

// Query represents a question query.
type Query struct {
	Service string
	Domain  string
}

// NewQueryWithService returns a new query instance with the specified service.
func NewQueryWithService(service string) *Query {
	return &Query{
		Service: service,
		Domain:  DefaultDomain,
	}
}

// String returns the string representation.
func (q *Query) String() string {
	return strings.Join([]string{q.Service, q.Domain}, nameSep)
}
