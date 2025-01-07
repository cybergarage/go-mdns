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

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns"
	"github.com/cybergarage/go-mdns/mdns/dns"
)

//go:embed log/service01.log
var service01 string

//go:embed log/service02.log
var service02 string

//go:embed log/googlecast01.log
var googlecast01 string

//go:embed log/googlecast02.log
var googlecast02 string

//go:embed log/googlecast03.log
var googlecast03 string

// 4.3.1.13. Examples
// dns-sd -R DD200C20D25AE5F7 _matterc._udp,_S3,_L840,_CM . 11111 D=840 CM=2
//
//go:embed log/matter-spec-120-4.3.1.13-dns-sd.log
var matterSpec12043113DNSSD string

// 4.3.1.13. Examples
// avahi-publish-service --subtype=_S3._sub._matterc._udp --subtype=_L840._sub._matterc._udp DD200C20D25AE5F7 --subtype=_CM._sub._matterc._udp _matterc._udp 11111 D=840 CM=2
//
//go:embed log/matter-spec-120-4.3.1.13-avahi01.log
var matterSpec12043113Avahi01 string

//go:embed log/matter-spec-120-4.3.1.13-avahi02.log
var matterSpec12043113Avahi02 string

func TestResponseMessages(t *testing.T) {
	type answer struct {
		name string
	}
	tests := []struct {
		name       string
		msgLogs    string
		answers    []answer
		attributes map[string]string
	}{
		{
			"service01",
			service01,
			[]answer{
				{"_services._dns-sd._udp.local"},
			},
			map[string]string{},
		},
		// {
		// 	"service02",
		// 	service02,
		// 	[]answer{
		// 		{"_companion-link._tcp.local"},
		// 	},
		// 	map[string]string{},
		// },
		{
			"googlecast01",
			googlecast01,
			[]answer{
				{"_services._dns-sd._udp.local"},
			},
			map[string]string{},
		},
		{
			"googlecast02",
			googlecast02,
			[]answer{
				{"_googlecast._tcp.local"},
			},
			map[string]string{},
		},
		{
			"googlecast03",
			googlecast03,
			[]answer{
				{"_googlezone._tcp.local"},
			},
			map[string]string{
				"id": "4E50AF186C368EE8A98A648BE272AAD5",
			},
		},
		{
			"matter 120 4.3.1.13/dns-sd",
			matterSpec12043113DNSSD,
			[]answer{
				{"_services._dns-sd"},
			},
			map[string]string{
				"D":  "840",
				"CM": "2",
			},
		},
		{
			"matter 120 4.3.1.13/avahi#01",
			matterSpec12043113Avahi01,
			[]answer{
				{"_matterc._udp.local"},
			},
			map[string]string{
				"D":  "840",
				"CM": "2",
			},
		},
		{
			"matter 120 4.3.1.13/avahi#02",
			matterSpec12043113Avahi02,
			[]answer{
				{"_matterc._udp.local"},
			},
			map[string]string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msgBytes, err := log.DecodeHexLogs(strings.Split(test.msgLogs, "\n"))
			if err != nil {
				t.Error(err)
				return
			}

			msg, err := dns.NewMessageWithBytes(msgBytes)
			if err != nil {
				t.Error(err)
				return
			}

			t.Log("\n" + msg.String())

			srv, err := mdns.NewServiceWithMessage(msg)
			if err != nil {
				t.Error(err)
				return
			}

			if !srv.Equal(srv) {
				t.Error("service not equal")
			}

			for _, answer := range test.answers {
				if msg.HasResourceRecord(answer.name) {
					continue
				}
				if _, ok := msg.LookupResourceRecordForNameSuffix(answer.name); ok {
					continue
				}
				t.Errorf("answer (%s) not found", answer.name)
			}

			for name, value := range test.attributes {
				attr, ok := srv.LookupAttribute(name)
				if !ok {
					t.Errorf("attribute (%s) not found", name)
					continue
				}
				if attr.Value() != value {
					t.Errorf("attribute (%s) value (%s) != (%s)", name, attr.Value(), value)
				}
			}
		})
	}
}
