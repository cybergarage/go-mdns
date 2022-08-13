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
	"errors"
	"net"
	"strings"
)

const (
	libvirtInterfaceName = "virbr0"
)

// IsIPv6Interface returns true whether the specified interface has a IPv6 address.
func IsIPv6Address(addr string) bool {
	if len(addr) == 0 {
		return false
	}
	if 0 <= strings.Index(addr, ":") {
		return true
	}
	return false
}

// IsIPv4Address returns true whether the specified address is a IPv4 address.
func IsIPv4Address(addr string) bool {
	if len(addr) == 0 {
		return false
	}
	return !IsIPv6Address(addr)
}

// IsIPv6Interface returns true whether the specified address is a IPv6 address.
func IsIPv6Interface(ifi *net.Interface) bool {
	addr, err := GetInterfaceAddress(ifi)
	if err != nil {
		return false
	}
	return IsIPv6Address(addr)
}

// IsIPv4Interface returns true whether the specified address is a IPv4 address.
func IsIPv4Interface(ifi *net.Interface) bool {
	addr, err := GetInterfaceAddress(ifi)
	if err != nil {
		return false
	}
	return IsIPv4Address(addr)
}

// IsLoopbackAddress returns true whether the specified address is a loopback addresses.
func IsLoopbackAddress(addr string) bool {
	localAddrs := []string{
		"127.0.0.1",
		"::1",
	}
	for _, localAddr := range localAddrs {
		if localAddr == addr {
			return true
		}
	}
	return false
}

// IsCommunicableAddress returns true whether the address is a effective address to commnicate with other nodes, othwise false.
func IsCommunicableAddress(addr string) bool {
	if len(addr) == 0 {
		return false
	}
	if IsLoopbackAddress(addr) {
		return false
	}
	return true
}

// IsBridgeInterface returns true when the specified interface is a bridege interface, otherwise false.
func IsBridgeInterface(ifi *net.Interface) bool {
	return ifi.Name == libvirtInterfaceName
}

// IsVirtualInterface returns true when the specified interface is a virtual interface, otherwise false.
func IsVirtualInterface(ifi *net.Interface) bool {
	if strings.HasPrefix(ifi.Name, "utun") { // macOS
		return true
	}
	if strings.HasPrefix(ifi.Name, "llw") { // VirtualBox
		return true
	}
	if strings.HasPrefix(ifi.Name, "awdl") { // AirDrop (macOS)
		return true
	}
	if strings.HasPrefix(ifi.Name, "en6") { // iPhone-USB (macOS)
		return true
	}
	return false
}

// GetInterfaceAddress returns a IPv4 or IPv6 address of the specivied interface.
func GetInterfaceAddress(ifi *net.Interface) (string, error) {
	addrs, err := ifi.Addrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		addrStr := addr.String()
		saddr := strings.Split(addrStr, "/")
		if len(saddr) < 2 {
			continue
		}
		return saddr[0], nil
	}
	return "", errors.New(errorAvailableAddressNotFound)
}

// GetAvailableInterfaces returns all available interfaces in the node.
func GetAvailableInterfaces() ([]*net.Interface, error) {
	useIfs := make([]*net.Interface, 0)
	localIfs, err := net.Interfaces()
	if err != nil {
		return useIfs, err
	}

	for n := range localIfs {
		localIf := localIfs[n]
		if (localIf.Flags & net.FlagLoopback) != 0 {
			continue
		}
		if (localIf.Flags & net.FlagUp) == 0 {
			continue
		}
		if (localIf.Flags & net.FlagMulticast) == 0 {
			continue
		}
		if IsBridgeInterface(&localIf) {
			continue
		}
		if IsVirtualInterface(&localIf) {
			continue
		}
		_, addrErr := GetInterfaceAddress(&localIf)
		if addrErr != nil {
			continue
		}

		useIf := localIf
		useIfs = append(useIfs, &useIf)
	}

	if len(useIfs) == 0 {
		return useIfs, errors.New(errorAvailableInterfaceFound)
	}

	return useIfs, err
}

// GetAvailableAddresses returns all available IPv4 addresses in the node.
func GetAvailableAddresses() ([]string, error) {
	addrs := make([]string, 0)
	ifis, err := GetAvailableInterfaces()
	if err != nil {
		return addrs, err
	}
	for _, ifi := range ifis {
		addr, err := GetInterfaceAddress(ifi)
		if err != nil {
			continue
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

func getMatchAddressBlockCount(ifAddr string, targetAddr string) int {
	const addrSep = "."
	targetAddrs := strings.Split(targetAddr, addrSep)
	ifAddrs := strings.Split(ifAddr, addrSep)
	if len(targetAddrs) != len(ifAddrs) {
		return -1
	}
	addrSize := len(targetAddrs)
	for n := 0; n < len(targetAddrs); n++ {
		if targetAddrs[n] != ifAddrs[n] {
			return n
		}
	}
	return addrSize
}

// GetAvailableInterfaceForAddr returns an interface of the specified address.
func GetAvailableInterfaceForAddr(fromAddr string) (*net.Interface, error) {
	ifis, err := GetAvailableInterfaces()
	if err != nil {
		return nil, err
	}

	switch len(ifis) {
	case 0:
		return nil, errors.New(errorAvailableInterfaceFound)
	case 1:
		return ifis[0], nil
	}

	ifAddrs := make([]string, len(ifis))
	for n := 0; n < len(ifAddrs); n++ {
		ifAddrs[n], _ = GetInterfaceAddress(ifis[n])
	}

	selIf := ifis[0]
	selIfMatchBlocks := getMatchAddressBlockCount(fromAddr, ifAddrs[0])
	for n := 0; n < len(ifAddrs); n++ {
		matchBlocks := getMatchAddressBlockCount(fromAddr, ifAddrs[n])
		if matchBlocks < selIfMatchBlocks {
			continue
		}
		selIf = ifis[n]
		selIfMatchBlocks = matchBlocks
	}

	return selIf, nil
}
