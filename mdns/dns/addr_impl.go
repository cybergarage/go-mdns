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

type addr struct {
	ip   net.IP
	port int
	zone string // IPv6 scoped addring zone
}

func newAddr() *addr {
	return &addr{
		ip:   nil,
		port: 0,
		zone: "",
	}
}

// NewAddr returns a new blank addr.
func NewAddr() Addr {
	return newAddr()
}

// NewAddrFromString returns a new addr parsed from the specified addr string.
func NewAddrFromString(addrString string) (Addr, error) {
	addr := newAddr()
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

// String returns the node string representation.
func (addr *addr) String() string {
	return net.JoinHostPort(addr.ip.String(), strconv.Itoa(addr.port))
}
