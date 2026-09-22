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

package mdnstest

import (
	"net"
	"strings"
	"testing"

	"github.com/cybergarage/go-logger/log/hexdump"
	"github.com/cybergarage/go-mdns/mdns"
	"github.com/cybergarage/go-mdns/mdns/dns"
)

func newTestServiceFromDump(t *testing.T, dump string, opts ...mdns.ServiceOptions) mdns.Service {
	t.Helper()

	msgBytes, err := hexdump.DecodeHexdumpLogs(strings.Split(dump, "\n"))
	if err != nil {
		t.Fatal(err)
	}

	msg, err := dns.NewMessageWithBytes(msgBytes)
	if err != nil {
		t.Fatal(err)
	}

	opts = append(opts, mdns.WithServiceMessage(msg))
	service, err := mdns.NewService(opts...)
	if err != nil {
		t.Fatal(err)
	}

	return service
}

// The SRV target is a domain name, and the address records of the target host
// are the addresses of the service.
func TestServiceAddresses(t *testing.T) {
	service := newTestServiceFromDump(t, matterAnswer01)

	expectedHost := "84FCE6036F38.local"
	if service.Host() != expectedHost {
		t.Errorf("host %s != %s", service.Host(), expectedHost)
	}

	expectedPort := 5540
	if service.Port() != expectedPort {
		t.Errorf("port %d != %d", service.Port(), expectedPort)
	}

	expectedAddrs := []string{
		"fe80::86fc:e6ff:fe03:6f38",
		"2400:2410:b242:bf00:86fc:e6ff:fe03:6f38",
		"192.168.100.58",
	}

	addrs := service.Addresses()
	if len(addrs) != len(expectedAddrs) {
		t.Fatalf("addresses %v != %v", addrs, expectedAddrs)
	}
	for _, expectedAddr := range expectedAddrs {
		expectedIP := net.ParseIP(expectedAddr)
		found := false
		for _, addr := range addrs {
			if addr.Equal(expectedIP) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%s is not found in %v", expectedAddr, addrs)
		}
	}
}

// RFC 4007: An IPv6 link-local address is ambiguous without its zone, and the
// zone is the interface which the response was received on.
func TestServiceLinkLocalAddrZone(t *testing.T) {
	testIfi := &net.Interface{ // nolint: exhaustruct,exhaustruct_v5
		Index: 1,
		Name:  "en0",
	}

	service := newTestServiceFromDump(t, matterAnswer01, mdns.WithServiceInterface(testIfi))

	if ifi := service.Interface(); ifi == nil || ifi.Name != testIfi.Name {
		t.Fatalf("interface %v != %v", ifi, testIfi)
	}

	zonedAddrCount := 0
	for _, addr := range service.Addrs() {
		if addr.Port != service.Port() {
			t.Errorf("port %d != %d", addr.Port, service.Port())
		}
		isLinkLocal := addr.IP.To4() == nil && addr.IP.IsLinkLocalUnicast()
		switch {
		case isLinkLocal:
			if addr.Zone != testIfi.Name {
				t.Errorf("%s : zone %q != %q", addr.IP, addr.Zone, testIfi.Name)
				continue
			}
			zonedAddrCount++
		default:
			// A global or an IPv4 address is not scoped to an interface.
			if addr.Zone != "" {
				t.Errorf("%s : zone %q should be empty", addr.IP, addr.Zone)
			}
		}
	}

	if zonedAddrCount != 1 {
		t.Errorf("zoned address count %d != 1", zonedAddrCount)
	}
}
