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
	"net"
	"regexp"

	"github.com/cybergarage/go-mdns/mdns/dns"
)

// Attribute represents a TXT record attribute.
type Attribute = dns.Attribute

// Attributes represents a TXT record attributes map.
type Attributes = dns.Attributes

// ResourceRecordSet represents a resource record array.
type ResourceRecordSet = dns.ResourceRecordSet

// ResourceRecord represents a resource record.
type ResourceRecord = dns.ResourceRecord

// Service represents a SRV record.
type Service interface {
	// Name returns the service name.
	Name() string
	// Domain returns the service domain.
	Domain() string
	// Host returns the service host.
	Host() string
	// Port returns the service port.
	Port() int
	// Addresses returns the service addresses.
	Addresses() []net.IP
	// ResourceRecordSet returns the service resource records.
	ResourceRecordSet() ResourceRecordSet
	// ResourceAttributes returns the service TXT attributes.
	ResourceAttributes() Attributes
	// LookupResourceAttribute returns the attribute with the specified name.
	LookupResourceAttribute(name string) (Attribute, bool)
	// LookupRecordByName returns the resource record of the specified name.
	LookupResourceByName(name string) (ResourceRecord, bool)
	// LookupResourceRecordByNameRegex returns the resource record of the specified name regex.
	LookupResourceByNameRegex(re *regexp.Regexp) (ResourceRecord, bool)
	// Equal returns true if the header is same as the specified header, otherwise false.
	Equal(other Service) bool
	// String returns the string representation.
	String() string
	// ServiceHelper returns the service helper.
	ServiceHelper
}

// ServiceHelper represents a service helper interface.
type ServiceHelper interface {
	// LookupResourceByNamePrefix returns the resource record of the specified name prefix.
	LookupResourceByNamePrefix(prefix string) (ResourceRecord, bool)
	// LookupResourceByNameSuffix returns the resource record of the specified name suffix.
	LookupResourceByNameSuffix(suffix string) (ResourceRecord, bool)
}
