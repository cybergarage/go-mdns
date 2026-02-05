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
	"net"
	"strconv"
	"testing"
)

func TestNewAddr(t *testing.T) {
	NewAddr()
}

func TestNewAddressParse(t *testing.T) {
	testCases := []struct {
		hostWithZone string
		port         int
		expectedIP   string
		expectedZone string
	}{
		{hostWithZone: "192.0.2.1", port: 25, expectedIP: "192.0.2.1", expectedZone: ""},
		{hostWithZone: "2001:db8::1", port: 80, expectedIP: "2001:db8::1", expectedZone: ""},
		{hostWithZone: "fe80::1%en0", port: 5353, expectedIP: "fe80::1", expectedZone: "en0"},
	}

	for _, testCase := range testCases {
		testAddr := net.JoinHostPort(testCase.hostWithZone, strconv.Itoa(testCase.port))

		addr, err := NewAddrFromString(testAddr)
		if err != nil {
			t.Error(err)
			continue
		}

		if addr.IP().String() != testCase.expectedIP {
			t.Errorf("%s != %s", addr.IP().String(), testCase.expectedIP)
		}
		if addr.Zone() != testCase.expectedZone {
			t.Errorf("%s != %s", addr.Zone(), testCase.expectedZone)
		}
		if addr.Port() != testCase.port {
			t.Errorf("%d != %d", addr.Port(), testCase.port)
		}
		if addr.String() != testAddr {
			t.Errorf("%s != %s", addr.String(), testAddr)
		}
	}
}
