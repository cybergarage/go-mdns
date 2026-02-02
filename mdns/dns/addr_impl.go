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
	"net"
	"strconv"
)

// AddrOption is a function type that modifies an addr.
type AddrOption func(*addr)

type addr struct {
	ip        net.IP
	port      int
	zone      string // IPv6 scoped addring zone
	transport Transport
}

func newAddr(opts ...AddrOption) *addr {
	addr := &addr{
		ip:        nil,
		port:      0,
		zone:      "",
		transport: TransportUnknown,
	}
	for _, opt := range opts {
		opt(addr)
	}
	return addr
}

// WithAddrIP returns an AddrOption with the specified IP address.
func WithAddrIP(ip net.IP) AddrOption {
	return func(a *addr) {
		a.ip = ip
	}
}

// WithAddrPort returns an AddrOption with the specified port number.
func WithAddrPort(port int) AddrOption {
	return func(a *addr) {
		a.port = port
	}
}

// WithAddrZone returns an AddrOption with the specified zone string.
func WithAddrZone(zone string) AddrOption {
	return func(a *addr) {
		a.zone = zone
	}
}

// WithAddrTransport returns an AddrOption with the specified transport protocol.
func WithAddrTransport(transport Transport) AddrOption {
	return func(a *addr) {
		a.transport = transport
	}
}

// NewAddr returns a new blank addr.
func NewAddr(opts ...AddrOption) Addr {
	return newAddr(opts...)
}

// NewAddrFromString returns a new addr parsed from the specified addr string.
func NewAddrFromString(addrString string, opts ...AddrOption) (Addr, error) {
	addr := newAddr(opts...)
	err := addr.parseString(addrString)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// ParseString parses the specified addr string.
func (addr *addr) parseString(addrStr string) error {
	hostStr, portStr, err := net.SplitHostPort(addrStr)
	if err != nil {
		return err
	}

	addr.ip = net.ParseIP(hostStr)

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("%s %s", ErrInvalid, addrStr)
	}
	addr.port = port

	return nil
}

// IP returns the IP address.
func (addr *addr) IP() net.IP {
	return addr.ip
}

// Port returns the port number.
func (addr *addr) Port() int {
	return addr.port
}

// Zone returns the zone string.
func (addr *addr) Zone() string {
	return addr.zone
}

// Transport returns the transport protocol.
func (addr *addr) Transport() Transport {
	return addr.transport
}

// String returns the node string representation.
func (addr *addr) String() string {
	return net.JoinHostPort(addr.ip.String(), strconv.Itoa(addr.port))
}
