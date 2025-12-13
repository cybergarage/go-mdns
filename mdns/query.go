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
	"time"
)

const (
	// 4.1.3. Domain Names
	LocalDomain = "local"
	// Subtype is the service subtype label.
	// 7.1. Selective Instance Enumeration (Subtypes)
	Subtype = "_sub"
)

// RFC 6763: 7.2. Browsing for Services.
const (
	ServiceTypeEnumerationName = "_services._dns-sd._udp"
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

const (
	// DefaultQueryService is the default service for mDNS queries.
	DefaultQueryService = ServiceTypeEnumerationName
	// DefaultQueryDomain is the default domain for mDNS queries.
	DefaultQueryDomain = LocalDomain
	// DefaultQueryTimeout is the default timeout duration for mDNS queries.
	DefaultQueryTimeout = time.Duration(5) * time.Second
)

// Query represents a question query.
type Query interface {
	// Subtype returns the subtype of the query.
	Subtype() string
	// Services returns the service names of the query.
	Services() []string
	// Domain returns the domain name of the query.
	Domain() string
	// String returns the string representation of the query.
	String() string
}
