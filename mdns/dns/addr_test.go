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
	testHosts := []string{
		"192.0.2.1",
		"2001:db8::1",
	}
	testPorts := []int{
		25,
		80,
	}

	for n := range testHosts {
		testHost := testHosts[n]
		testPort := testPorts[n]

		testAddr := net.JoinHostPort(testHost, strconv.Itoa(testPort))

		addr, err := NewAddrFromString(testAddr)
		if err != nil {
			t.Error(err)
			continue
		}

		if addr.IP().String() != testHost {
			t.Errorf("%s != %s", addr.IP().String(), testHost)
		}

		if addr.Port() != testPort {
			t.Errorf("%d != %d", addr.Port(), testPort)
		}

		if addr.String() != testAddr {
			t.Errorf("%s != %s", addr.String(), testAddr)
		}
	}
}
