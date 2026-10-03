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

package dns

import (
	"fmt"
)

// Type represents the type of a resource record or a question (RFC 1035,
// 3.2.2 and 3.2.3).
type Type uint

// The record and question types.
const (
	// A is a host's IPv4 address (RFC 1035).
	A Type = 0x0001
	// NS is an authoritative name server (RFC 1035).
	NS Type = 0x0002
	// CNAME is the canonical name of an alias (RFC 1035).
	CNAME Type = 0x0005
	// TXT holds text strings, the attributes of a DNS-SD service (RFC 6763, 6).
	TXT Type = 0x0010
	// SRV is the host and the port of a service (RFC 2782).
	SRV Type = 0x0021
	// OPT is the EDNS(0) pseudo record (RFC 6891).
	OPT Type = 0x0029
	// PTR points to a domain name, a service instance in DNS-SD (RFC 6763, 4).
	PTR Type = 0x000C
	// HINFO is the host information (RFC 1035).
	HINFO Type = 0x000D
	// MX is a mail exchange (RFC 1035).
	MX Type = 0x000F
	// AAAA is a host's IPv6 address (RFC 3596).
	AAAA Type = 0x001C
	// AXFR asks for a zone transfer (RFC 1035).
	AXFR Type = 0x00FC
	// NSEC lists the record types a name has, which mDNS uses for the
	// negative responses (RFC 6762, 6.1).
	NSEC Type = 0x002F
	// ANY asks for the records of every type (RFC 1035).
	ANY Type = 0x00FF
)

// Equal returns true if the type matches the specified one.
func (t Type) Equal(other Type) bool {
	if t == ANY || other == ANY {
		return true
	}
	return t == other
}

// String returns the string of the type.
func (t Type) String() string {
	switch t {
	case A:
		return "A"
	case NS:
		return "NS"
	case CNAME:
		return "CNAME"
	case TXT:
		return "TXT"
	case SRV:
		return "SRV"
	case OPT:
		return "OPT"
	case PTR:
		return "PTR"
	case HINFO:
		return "HINFO"
	case MX:
		return "MX"
	case AAAA:
		return "AAAA"
	case AXFR:
		return "AXFR"
	case NSEC:
		return "NSEC"
	case ANY:
		return "ANY"
	}
	return fmt.Sprintf("%04X(%d)", uint16(t), uint16(t))
}
