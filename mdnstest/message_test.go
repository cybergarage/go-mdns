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
	"github.com/cybergarage/go-mdns/mdns/protocol"
)

//go:embed log/matter01.log
var matter01 string

//go:embed log/matter02.log
var matter02 string

//go:embed log/googlecast01.log
var googlecast01 string

// 4.3.1.13. Examples
// dns-sd -R DD200C20D25AE5F7 _matterc._udp,_S3,_L840,_CM . 11111 D=840 CM=2
//
//go:embed log/matter-spec-120-4.3.1.13.log
var matterSpec12043113 string

func TestResponseMessage(t *testing.T) {
	type answer struct {
		name string
	}
	tests := []struct {
		name    string
		msgLogs string
		answers []answer
	}{
		/*
			{
				"matter01",
				matter01,
				[]answer{
					{"_services._dns-sd._udp.local"},
				},
			},
			{
				"matter02",
				matter02,
				[]answer{
					{"_services._dns-sd._udp.local"},
				},
			},
			{
				"googlecast01",
				googlecast01,
				[]answer{
					{"_services._dns-sd._udp.local"},
				},
			},
		*/
		{
			"matter 120 4.3.1.13",
			matterSpec12043113,
			[]answer{
				{"_services._dns-sd"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msgBytes, err := log.DecodeHexLog(strings.Split(test.msgLogs, "\n"))
			if err != nil {
				t.Error(err)
				return
			}

			msg, err := protocol.NewMessageWithBytes(msgBytes)
			if err != nil {
				t.Error(err)
				return
			}

			for _, answer := range test.answers {
				if !msg.Answers.HasResourceRecord(answer.name) {
					t.Errorf("answer (%s) not found", answer.name)
				}
			}
			t.Log(msg.String())
		})
	}
}
