// Copyright (C) 2026 The go-mdns Authors All rights reserved.
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
	"net"
)

// CacheFlush is the cache-flush bit of a resource record class. It shares
// the top bit of the class with QU, which carries a different meaning in a
// question (RFC 6762, 10.2 and 18.13). A record sets it with
// SetUnicastResponse(true).
const CacheFlush = QU

// Maximum lengths of the encoded names and strings (RFC 1035, 3.1 and 3.3).
const (
	maxLabelLength  = 63
	maxStringLength = 255
)

func newResourceRecord(name string, typ Type, ttl uint, data []byte) *record {
	r := newRecord()
	r.name = name
	r.typ = typ
	r.class = IN
	r.ttl = ttl
	r.data = data
	return r
}

func encodeName(name string) ([]byte, error) {
	for _, label := range SplitName(name) {
		if len(label) == 0 || maxLabelLength < len(label) {
			return nil, fmt.Errorf("%w: label %q in %q", ErrInvalid, label, name)
		}
	}
	w := NewWriter()
	if err := w.WriteName(name); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// NewPTRResourceRecord returns a PTR record named name which points to
// domainName.
func NewPTRResourceRecord(name string, domainName string, ttl uint) (PTRRecord, error) {
	data, err := encodeName(domainName)
	if err != nil {
		return nil, err
	}
	return newPTRRecordWithResourceRecord(newResourceRecord(name, PTR, ttl, data))
}

// NewSRVResourceRecord returns a SRV record named name, a service instance
// name such as "instance._http._tcp.local", which points to port on target.
func NewSRVResourceRecord(name string, priority, weight, port uint16, target string, ttl uint) (SRVRecord, error) {
	targetBytes, err := encodeName(target)
	if err != nil {
		return nil, err
	}
	w := NewWriter()
	for _, v := range []uint16{priority, weight, port} {
		if err := w.WriteUint16(v); err != nil {
			return nil, err
		}
	}
	if err := w.WriteBytes(targetBytes); err != nil {
		return nil, err
	}
	return newSRVRecordWithResourceRecord(newResourceRecord(name, SRV, ttl, w.Bytes()))
}

// NewTXTResourceRecord returns a TXT record named name which holds strs,
// such as "key=value" attributes. A TXT record without a string holds a
// single empty string (RFC 6763, 6.1).
func NewTXTResourceRecord(name string, strs []string, ttl uint) (TXTRecord, error) {
	if len(strs) == 0 {
		strs = []string{""}
	}
	w := NewWriter()
	for _, s := range strs {
		if maxStringLength < len(s) {
			return nil, fmt.Errorf("%w: TXT string of %d bytes", ErrInvalid, len(s))
		}
		if err := w.WriteString(s); err != nil {
			return nil, err
		}
	}
	return newTXTRecordWithResourceRecord(newResourceRecord(name, TXT, ttl, w.Bytes()))
}

// NewAResourceRecord returns an A record named name for the IPv4 address
// ip.
func NewAResourceRecord(name string, ip net.IP, ttl uint) (ARecord, error) {
	ip4 := ip.To4()
	if ip4 == nil {
		return nil, fmt.Errorf("%w: %s is not an IPv4 address", ErrInvalid, ip)
	}
	return newARecordWithResourceRecord(newResourceRecord(name, A, ttl, []byte(ip4))), nil
}

// NewAAAAResourceRecord returns an AAAA record named name for the IPv6
// address ip.
func NewAAAAResourceRecord(name string, ip net.IP, ttl uint) (AAAARecord, error) {
	if ip.To4() != nil || ip.To16() == nil {
		return nil, fmt.Errorf("%w: %s is not an IPv6 address", ErrInvalid, ip)
	}
	return newAAAARecordWithResourceRecord(newResourceRecord(name, AAAA, ttl, []byte(ip.To16()))), nil
}

// NewAddressResourceRecord returns an A record for an IPv4 address and an
// AAAA record for an IPv6 one.
func NewAddressResourceRecord(name string, ip net.IP, ttl uint) (ResourceRecord, error) {
	if ip.To4() != nil {
		return NewAResourceRecord(name, ip, ttl)
	}
	return NewAAAAResourceRecord(name, ip, ttl)
}
