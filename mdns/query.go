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

	"github.com/cybergarage/go-mdns/mdns/dns"
)

// Class represents a DNS class.
type Class = dns.Class

// RFC 6762 - Multicast DNS.
const (
	// LocalDomain is the local domain name.
	// 4.1.3. Domain Names.
	LocalDomain = "local"
	// QU is the unicast response class.
	// 5.4. Questions Requesting Unicast Responses
	// 18.12. Repurposing of Top Bit of qclass in Question Section.
	QU Class = 0x8000
	// Subtype is the service subtype label.
	// 7.1. Selective Instance Enumeration (Subtypes).
	Subtype = "_sub"
)

// RFC 6763 - DNS-Based Service Discovery.
const (
	// ServiceTypeEnumerationName is the service type for service enumeration.
	// 7.2. Browsing for Services.
	ServiceTypeEnumerationName = "_services._dns-sd._udp"
)

// RFC 6765: 11. Discovery of Browsing and Registration Domains (Domain Enumeration).
const (
	// RecommendedBrowsingService is a list of domains recommended for browsing.
	RecommendedBrowsingService = "b._dns-sd._udp"
	// DefaultBrowsingService is a single recommended default domain for browsing.
	DefaultBrowsingService = "db._dns-sd._udp"
	// RecommendedRegisteringService is a list of domains recommended for registering services using Dynamic Update.
	RecommendedRegisteringService = "r._dns-sd._udp"
	// DefaultRegisteringService is a single recommended default domain for registering services.
	DefaultRegisteringService = "dr._dns-sd._udp"
	// AutomaticBrowsingService is the "legacy browsing" or "automatic browsing" domain(s).
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
	// Service returns the service name of the query.
	Service() string
	// Domain returns the domain name of the query.
	Domain() string
	// MessageHandler returns the message handler of the query if set.
	MessageHandler() (MessageHandler, bool)
	// String returns the string representation of the query.
	String() string
}
