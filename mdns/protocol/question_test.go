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
	"testing"
)

func TestQuestion(t *testing.T) {
	t.Run("ParseQuestion", func(t *testing.T) {
		tests := []struct {
			query           []byte
			expectedName    string
			expectedType    Type
			expectedUnicast bool
			expectedClass   Class
		}{
			{
				query:           []byte{0x02, 0x6c, 0x62, 0x07, 0x5f, 0x64, 0x6e, 0x73, 0x2d, 0x73, 0x64, 0x04, 0x5f, 0x75, 0x64, 0x70, 0x05, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x00, 0x00, 0x0c, 0x80, 0x01},
				expectedName:    "lb._dns-sd._udp.local",
				expectedType:    PTR,
				expectedUnicast: true,
				expectedClass:   IN,
			},
		}
		for _, test := range tests {
			q, err := NewQuestionWithReader(NewReaderWithBytes(test.query))
			if err != nil {
				t.Error(err)
			}
			// Checks each field
			if q.Name() != test.expectedName {
				t.Errorf("%s != %s", q.Name(), test.expectedName)
			}
			if q.Type() != test.expectedType {
				t.Errorf("%2X != %2X", q.Type(), test.expectedType)
			}
			if q.UnicastResponse() != test.expectedUnicast {
				t.Errorf("%t != %t", q.UnicastResponse(), test.expectedUnicast)
			}
			if q.Class() != test.expectedClass {
				t.Errorf("%2X != %2X", q.Class(), test.expectedClass)
			}
		}
	})
}
