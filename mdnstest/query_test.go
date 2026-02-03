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
	"github.com/cybergarage/go-mdns/mdns/dns"
)

func TestQueryResponse(t *testing.T) {
	tests := []struct {
		name       string
		queryDump  string
		answerDump string
	}{
		{
			"matterc._udp.local",
			matterQuery01,
			matterAnswer01,
		},
		{
			"_S9._sub._matterc._udp.local",
			matterQuery02,
			matterAnswer02,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			queryBytes, err := hexdump.DecodeHexdumpLogs(strings.Split(test.queryDump, "\n"))
			if err != nil {
				t.Error(err)
				return
			}

			queryMsg, err := dns.NewMessageWithBytes(queryBytes)
			if err != nil {
				t.Error(err)
				return
			}

			answerBytes, err := hexdump.DecodeHexdumpLogs(strings.Split(test.answerDump, "\n"))
			if err != nil {
				t.Error(err)
				return
			}

			answerMsg, err := dns.NewMessageWithBytes(answerBytes)
			if err != nil {
				t.Error(err)
				return
			}

			if !queryMsg.IsQueryAnswer(answerMsg) {
				t.Errorf("Invalid answer message for the query message")
				t.Log("\n=== Query Message ===\n" + queryMsg.String())
				t.Log("\n=== Answer Message ===\n" + answerMsg.String())
				return
			}
		})
	}
}
