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

func TestRequestMessages(t *testing.T) {
	type query struct {
		name   string
		domain string
	}

	tests := []struct {
		name  string
		dump  string
		query query
	}{
		{
			"chip-tool-query-01",
			chipToolQuery01,
			query{
				name:   "S9._sub._matterc._udp",
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

			t.Log("\n" + msg.String())
		})
	}
}

func TestResponseMessages(t *testing.T) {
	type answer struct {
		name string
	}
	tests := []struct {
		name       string
		dump       string
		answers    []answer
		attributes map[string]string
	}{
		{
			"google-cast-01",
			googlecast01,
			[]answer{
				{"_services._dns-sd._udp.local"},
			},
			map[string]string{},
		},
		{
			"google-cast-02",
			googlecast02,
			[]answer{
				{"_googlecast._tcp.local"},
			},
			map[string]string{},
		},
		{
			"google-cast-03",
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
		{
			"matter service 01",
			matterService01,
			[]answer{
				{"_services._dns-sd._udp.local"},
			},
			map[string]string{},
		},
		{
			"matter service 02",
			matterService02,
			[]answer{
				{"_S9._sub._matterc._udp.local"},
			},
			map[string]string{
				"VP": "5002+5010",
				"DT": "14",
				"T":  "1",
				"CM": "1",
				"RI": "0F0072207A66D38D32D2BF9DF653E9735F86",
				"PH": "33",
				"PI": "",
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

			srv, err := mdns.NewService(
				mdns.WithServiceMessage(msg),
			)
			if err != nil {
				t.Log("\n" + msg.String())
				t.Error(err)
				return
			}

			reportError := func(msg mdns.Message, srv mdns.Service, format string, args ...any) {
				t.Errorf(format, args...)
				t.Log("\n" + msg.String())
				t.Log("\n" + srv.String())
			}

			for _, answer := range test.answers {
				if _, ok := msg.LookupResourceRecordByName(answer.name); ok {
					continue
				}
				if _, ok := msg.LookupResourceRecordByNameSuffix(answer.name); ok {
					continue
				}
				reportError(msg, srv, "answer (%s) not found", answer.name)
				return
			}

			for name, value := range test.attributes {
				attr, ok := srv.LookupResourceAttribute(name)
				if !ok {
					reportError(msg, srv, "attribute (%s) not found", name)
					return
				}
				if attr.Value() != value {
					reportError(msg, srv, "attribute (%s) value (%s) != (%s)", name, attr.Value(), value)
					return
				}
			}
		})
	}
}

func TestEqualMessages(t *testing.T) {
	tests := []struct {
		name string
		dump string
	}{
		{
			"google-cast-01",
			googlecast01,
		},
		{
			"google-cast-02",
			googlecast02,
		},
		{
			"google-cast-03",
			googlecast03,
		},
		{
			"matter 120 4.3.1.13/dns-sd",
			matterSpec12043113DNSSD,
		},
		{
			"matter 120 4.3.1.13/avahi#01",
			matterSpec12043113Avahi01,
		},
		{
			"matter 120 4.3.1.13/avahi#02",
			matterSpec12043113Avahi02,
		},
		{
			"matter service 01",
			matterService01,
		},
		{
			"matter service 02",
			matterService02,
		},
	}

	msgs := make([]dns.Message, 0)
	for _, test := range tests {
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
		msgs = append(msgs, msg)
	}

	for i, msg1 := range msgs {
		for j, msg2 := range msgs {
			equal := msg1.Equal(msg2)
			if i != j {
				continue
			}
			if i == j {
				if !equal {
					t.Errorf("messages[%d] and messages[%d] should be equal", i, j)
				}
			} else {
				if equal {
					t.Errorf("messages[%d] and messages[%d] should not be equal", i, j)
				}
			}
		}
	}
}
