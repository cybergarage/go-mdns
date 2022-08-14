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
	"encoding/hex"
	"testing"
)

func TestResourceRecord(t *testing.T) {
	t.Run("PTR", func(t *testing.T) {
		tests := []struct {
			query              []byte
			expectedName       string
			expectedTTL        uint
			expectedDomainName string
		}{
			{
				query:              []byte{0x09, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x07, 0x5f, 0x64, 0x6e, 0x73, 0x2d, 0x73, 0x64, 0x04, 0x5f, 0x75, 0x64, 0x70, 0x05, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x00, 0x00, 0x0c, 0x00, 0x01, 0x00, 0x00, 0x11, 0x94, 0x00, 0x13, 0x0b, 0x5f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x63, 0x61, 0x73, 0x74, 0x04, 0x5f, 0x74, 0x63, 0x70, 0xc0, 0x23},
				expectedName:       "_services._dns-sd._udp.local",
				expectedTTL:        4500,
				expectedDomainName: "googlecast._tcp.local",
			},
		}
		for _, test := range tests {
			t.Run(test.expectedName, func(t *testing.T) {
				q, err := newResourceRecordWithReader(bytes.NewReader(test.query))
				if err != nil {
					t.Error(err)
				}
				ptr, ok := q.(*PTRRecord)
				if !ok {
					t.Errorf("%v", q)
					return
				}
				// Checks each field
				if ptr.Name() != test.expectedName {
					t.Errorf("Name: %s != %s", q.Name(), test.expectedName)
				}
				if ptr.TTL() != test.expectedTTL {
					t.Errorf("TTL: %d != %d", q.TTL(), test.expectedTTL)
				}
				if ptr.DomainName() != test.expectedDomainName {
					t.Skipf("Domain: %s != %s", ptr.DomainName(), test.expectedDomainName)
				}
			})
		}
	})

	t.Run("OPT", func(t *testing.T) {
		tests := []struct {
			query         []byte
			expectedName  string
			expectedType  Type
			expectedCache bool
			expectedClass Class
		}{
			{
				query:         []byte{0x00, 0x00, 0x29, 0x05, 0xa0, 0x00, 0x00, 0x11, 0x94, 0x00, 0x12, 0x00, 0x04, 0x00, 0x0e, 0x00, 0x74, 0x52, 0x06, 0x8d, 0xcf, 0x54, 0x27, 0x86, 0xfd, 0xcd, 0x88, 0xe1, 0x43},
				expectedName:  "",
				expectedType:  OPT,
				expectedCache: false,
			},
		}
		for _, test := range tests {
			q, err := newResourceRecordWithReader(bytes.NewReader(test.query))
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
			if q.CacheFlush() != test.expectedCache {
				t.Errorf("%t != %t", q.CacheFlush(), test.expectedCache)
			}
			// Checks all bytes
			if !bytes.Equal(q.Bytes(), test.query) {
				t.Errorf("%s != %s", hex.EncodeToString(q.Bytes()), hex.EncodeToString(test.query))
			}
		}
	})
}
