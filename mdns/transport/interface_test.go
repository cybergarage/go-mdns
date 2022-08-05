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
	"testing"
)

func TestGetAvailableInterfaces(t *testing.T) {
	ifs, err := GetAvailableInterfaces()
	if err != nil {
		t.Error(err)
	}
	if len(ifs) == 0 {
		t.Errorf("available interface is not found")
	}
}

func TestGetAvailableAddresses(t *testing.T) {
	addrs, err := GetAvailableAddresses()
	if err != nil {
		t.Error(err)
	}
	if len(addrs) == 0 {
		t.Errorf("available address is not found")
	}
}

func TestIPv6Addresses(t *testing.T) {
	goodAddrs := []string{
		"::1",
		"fe80::1875:6549:801:d487",
	}

	for n, addr := range goodAddrs {
		if !IsIPv6Address(addr) {
			t.Errorf("[%d] %s", n, addr)
		}
	}

	badAddrs := []string{
		"",
		"127.0.0.1",
		"192.168.0.1",
	}

	for n, addr := range badAddrs {
		if IsIPv6Address(addr) {
			t.Errorf("[%d] %s", n, addr)
		}
	}
}

func TestIPv4Addresses(t *testing.T) {
	goodAddrs := []string{
		"127.0.0.1",
		"192.168.0.1",
	}

	for n, addr := range goodAddrs {
		if !IsIPv4Address(addr) {
			t.Errorf("[%d] %s", n, addr)
		}
	}

	barAddrs := []string{
		"",
		"::1",
		"fe80::1875:6549:801:d487",
	}

	for n, addr := range barAddrs {
		if IsIPv4Address(addr) {
			t.Errorf("[%d] %s", n, addr)
		}
	}
}

func TestLocalAddresses(t *testing.T) {
	goodAddrs := []string{
		"127.0.0.1",
		"::1",
	}

	for n, addr := range goodAddrs {
		if !IsLoopbackAddress(addr) {
			t.Errorf("[%d] %s", n, addr)
		}
	}

	badAddrs := []string{
		"",
		"192.168.0.1",
		"fe80::1875:6549:801:d487",
	}

	for n, addr := range badAddrs {
		if IsLoopbackAddress(addr) {
			t.Errorf("[%d] %s", n, addr)
		}
	}
}

func TestCommunicableAddresses(t *testing.T) {
	goodAddrs := []string{
		"192.168.0.1",
		"fe80::1875:6549:801:d487",
	}

	for n, addr := range goodAddrs {
		if !IsCommunicableAddress(addr) {
			t.Errorf("[%d] %s", n, addr)
		}
	}

	badAddrs := []string{
		"",
		"127.0.0.1",
		"::1",
	}

	for n, addr := range badAddrs {
		if IsCommunicableAddress(addr) {
			t.Errorf("[%d] %s", n, addr)
		}
	}
}
