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

package protocol

import (
	"bytes"
	"testing"
)

func TestQuestion(t *testing.T) {
	t.Run("ParseQuestion", func(t *testing.T) {
		tests := []struct {
			query    []byte
			expected string
		}{
			{
				query:    []byte{0x02, 0x6c, 0x62, 0x07, 0x5f, 0x64, 0x6e, 0x73, 0x2d, 0x73, 0x64, 0x04, 0x5f, 0x75, 0x64, 0x70, 0x05, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x00, 0x00, 0x0c, 0x80, 0x01},
				expected: "lb._dns-sd._udp.local",
			},
		}
		for _, test := range tests {
			q, err := NewQuestionWithReader(bytes.NewReader(test.query))
			if err != nil {
				t.Error(err)
			}
			if q.DomainName != test.expected {
				t.Errorf("%s != %s", q.DomainName, test.expected)
			}
		}
	})
}
