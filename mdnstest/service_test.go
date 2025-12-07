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
	_ "embed"
	"strings"
	"testing"

	"github.com/cybergarage/go-logger/log/hexdump"
	"github.com/cybergarage/go-mdns/mdns"
	"github.com/cybergarage/go-mdns/mdns/dns"
)

func TestServiceMessages(t *testing.T) {
	type expected struct {
		name   string
		domain string
	}
	tests := []struct {
		name string
		dump string
		exp  expected
	}{
		{
			"matter 120 4.3.1.13/dns-sd",
			matterSpec12043113DNSSD,
			expected{
				name:   "DD200C20D25AE5F7._matterc._udp",
				domain: "local",
			},
		},
		// {
		// 	"matter 120 4.3.1.13/avahi#01",
		// 	matterSpec12043113Avahi01,
		// 	expected{
		// 		name:   "DD200C20D25AE5F7._matterc._udp.local",
		// 		domain: "local",
		// 	},
		// },
		{
			"matter 120 4.3.1.13/avahi#02",
			matterSpec12043113Avahi02,
			expected{
				name:   "DD200C20D25AE5F7._matterc._udp",
				domain: "local",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msgBytes, err := hexdump.DecodeHexdumpLogs(strings.Split(test.dump, "\n"))
			if err != nil {
				t.Error(err)
				return
			}

			msg, err := dns.NewMessageWithBytes(msgBytes)
			if err != nil {
				t.Error(err)
				return
			}

			t.Logf("\n%s", msg.String())

			service, err := mdns.NewService(
				mdns.WithServiceMessage(msg),
			)
			if err != nil {
				t.Error(err)
				return
			}

			t.Logf("\n%s", service.String())

			name := service.Name()
			if name != test.exp.name {
				t.Errorf("service name mismatch: expected %s, got %s", test.exp.name, name)
			}

			domain := service.Domain()
			if domain != test.exp.domain {
				t.Errorf("service domain mismatch: expected %s, got %s", test.exp.domain, domain)
			}
		})
	}
}
