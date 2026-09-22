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

package transport

import (
	"net"
)

// InterfaceSelector selects the network interfaces and the addresses which the
// servers are bound to.
type InterfaceSelector struct {
	ifis        []*net.Interface
	ipv4Enabled bool
	ipv6Enabled bool
}

// NewInterfaceSelector returns a new selector which binds all available
// interfaces and address families.
func NewInterfaceSelector() *InterfaceSelector {
	return &InterfaceSelector{
		ifis:        nil,
		ipv4Enabled: true,
		ipv6Enabled: true,
	}
}

// SetInterfaces sets the network interfaces to bind. All available interfaces
// are bound when no interface is set.
func (sel *InterfaceSelector) SetInterfaces(ifis []*net.Interface) {
	sel.ifis = ifis
}

// Interfaces returns the configured network interfaces.
func (sel *InterfaceSelector) Interfaces() []*net.Interface {
	return sel.ifis
}

// SetIPv4Enabled sets whether the IPv4 addresses are bound.
func (sel *InterfaceSelector) SetIPv4Enabled(flag bool) {
	sel.ipv4Enabled = flag
}

// IsIPv4Enabled returns true if the IPv4 addresses are bound.
func (sel *InterfaceSelector) IsIPv4Enabled() bool {
	return sel.ipv4Enabled
}

// SetIPv6Enabled sets whether the IPv6 addresses are bound.
func (sel *InterfaceSelector) SetIPv6Enabled(flag bool) {
	sel.ipv6Enabled = flag
}

// IsIPv6Enabled returns true if the IPv6 addresses are bound.
func (sel *InterfaceSelector) IsIPv6Enabled() bool {
	return sel.ipv6Enabled
}

// BindInterfaces returns the network interfaces to bind.
func (sel *InterfaceSelector) BindInterfaces() ([]*net.Interface, error) {
	if 0 < len(sel.ifis) {
		return sel.ifis, nil
	}
	return GetAvailableInterfaces()
}

// BindAddresses returns the addresses of the specified interface to bind.
func (sel *InterfaceSelector) BindAddresses(ifi *net.Interface) ([]string, error) {
	addrs, err := GetInterfaceAddresses(ifi)
	if err != nil {
		return nil, err
	}

	bindAddrs := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		if IsIPv6Address(addr) {
			if !sel.ipv6Enabled {
				continue
			}
		} else if !sel.ipv4Enabled {
			continue
		}
		bindAddrs = append(bindAddrs, addr)
	}

	if len(bindAddrs) == 0 {
		return nil, errAvailableAddressNotFound
	}

	return bindAddrs, nil
}
