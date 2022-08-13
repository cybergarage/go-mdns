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

func TestMulticastSocketBindWithInterface(t *testing.T) {
	sock := NewMulticastSocket()

	ifis, err := GetAvailableInterfaces()
	if err != nil {
		t.Error(err)
		return
	}

	for _, ifi := range ifis {
		ifaddrs, err := GetInterfaceAddresses(ifi)
		if err != nil {
			t.Error(err)
			continue
		}
		for _, ifaddr := range ifaddrs {
			err = sock.Bind(ifi, ifaddr)
			if err != nil {
				t.Error(err)
			}

			err = sock.Close()
			if err != nil {
				t.Error(err)
			}
		}
	}
}
