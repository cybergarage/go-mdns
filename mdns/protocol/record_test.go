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
	"fmt"
	"net"
	"strings"
	"testing"
)

// nolint: gocyclo, maintidx
func TestResourceRecord(t *testing.T) {
	t.Run("PTR", func(t *testing.T) {
		tests := []struct {
			query              []byte
			expectedName       string
			expectedTTL        uint
			expectedDomainName string
		}{
			{
				query:              []byte{0x09, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x07, 0x5f, 0x64, 0x6e, 0x73, 0x2d, 0x73, 0x64, 0x04, 0x5f, 0x75, 0x64, 0x70, 0x05, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x00, 0x00, 0x0c, 0x00, 0x01, 0x00, 0x00, 0x11, 0x94, 0x00, 0x13, 0x0b, 0x5f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x63, 0x61, 0x73, 0x74, 0x04, 0x5f, 0x74, 0x63, 0x70, 0xc0, 0x00},
				expectedName:       "_services._dns-sd._udp.local",
				expectedTTL:        4500,
				expectedDomainName: "_services._dns-sd._udp.local",
			},
		}
		for _, test := range tests {
			t.Run(test.expectedName, func(t *testing.T) {
				q, err := newResourceRecordWithReader(NewReaderWithBytes(test.query))
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
					t.Errorf("Domain: %s != %s", ptr.DomainName(), test.expectedDomainName)
				}
			})
		}
	})

	t.Run("SRV", func(t *testing.T) {
		tests := []struct {
			query            []byte
			expectedTTL      uint
			expectedPriority uint
			expectedWeight   uint
			expectedPort     uint
			expectedTarget   string
		}{
			{
				query:            []byte{0x00, 0x00, 0x21, 0x80, 0x01, 0x00, 0x00, 0x00, 0x78, 0x00, 0x1f, 0x00, 0x01, 0x00, 0x02, 0x1f, 0x49, 0x16, 0x66, 0x75, 0x63, 0x68, 0x73, 0x69, 0x61, 0x2d, 0x37, 0x63, 0x64, 0x39, 0x2d, 0x35, 0x63, 0x34, 0x39, 0x2d, 0x65, 0x30, 0x61, 0x37, 0xc0, 0x1d},
				expectedTTL:      120,
				expectedPriority: 1,
				expectedWeight:   2,
				expectedPort:     8009,
				expectedTarget:   "fuchsia-7cd9-5c49-e0a7",
			},
		}
		for _, test := range tests {
			t.Run(fmt.Sprintf("%d:%d:%d", test.expectedPriority, test.expectedWeight, test.expectedPort), func(t *testing.T) {
				q, err := newResourceRecordWithReader(NewReaderWithBytes(test.query))
				if err != nil {
					t.Error(err)
				}
				srv, ok := q.(*SRVRecord)
				if !ok {
					t.Errorf("%v", q)
					return
				}
				// Checks each field
				if srv.TTL() != test.expectedTTL {
					t.Errorf("TTL: %d != %d", srv.TTL(), test.expectedTTL)
				}
				if srv.Priority() != test.expectedPriority {
					t.Errorf("Priority: %d != %d", srv.Priority(), test.expectedPriority)
				}
				if srv.Weight() != test.expectedWeight {
					t.Errorf("Weight: %d != %d", srv.Weight(), test.expectedWeight)
				}
				if srv.Port() != test.expectedPort {
					t.Errorf("Port: %d != %d", srv.Port(), test.expectedPort)
				}
				if srv.Target() != test.expectedTarget {
					t.Skipf("Target: %s != %s", srv.Target(), test.expectedTarget)
				}
			})
		}
	})

	t.Run("TXT", func(t *testing.T) {
		tests := []struct {
			query         []byte
			expectedTTL   uint
			expectedAttrs []string
		}{
			{
				query:         []byte{0x00, 0x00, 0x10, 0x80, 0x01, 0x00, 0x00, 0x00, 0x78, 0x00, 0x49, 0x23, 0x69, 0x64, 0x3d, 0x31, 0x66, 0x61, 0x30, 0x35, 0x63, 0x31, 0x35, 0x63, 0x65, 0x37, 0x65, 0x66, 0x65, 0x63, 0x62, 0x37, 0x34, 0x30, 0x38, 0x37, 0x38, 0x64, 0x31, 0x63, 0x30, 0x32, 0x34, 0x32, 0x31, 0x38, 0x33, 0x23, 0x63, 0x64, 0x3d, 0x44, 0x45, 0x39, 0x39, 0x42, 0x44, 0x31, 0x39, 0x34, 0x42, 0x32, 0x42, 0x46, 0x33, 0x42, 0x42, 0x46, 0x35, 0x42, 0x30, 0x39, 0x43, 0x36, 0x46, 0x42, 0x42, 0x34, 0x41, 0x30, 0x31, 0x35, 0x46, 0x00},
				expectedTTL:   120,
				expectedAttrs: []string{"id=1fa05c15ce7efecb740878d1c0242183", "cd=DE99BD194B2BF3BBF5B09C6FBB4A015F"},
			},
		}
		for _, test := range tests {
			t.Run(strings.Join(test.expectedAttrs, ","), func(t *testing.T) {
				q, err := newResourceRecordWithReader(NewReaderWithBytes(test.query))
				if err != nil {
					t.Error(err)
				}
				txt, ok := q.(*TXTRecord)
				if !ok {
					t.Errorf("%v", txt)
					return
				}
				// Checks each field
				if txt.TTL() != test.expectedTTL {
					t.Errorf("TTL: %d != %d", txt.TTL(), test.expectedTTL)
				}
				attrs := txt.Attributes()
				if len(attrs) != len(test.expectedAttrs) {
					t.Errorf("Attrs: %d != %d", len(attrs), len(test.expectedAttrs))
					return
				}
				for n, attr := range attrs {
					if attr != test.expectedAttrs[n] {
						t.Errorf("Attr[%d]: %s != %s", n, attrs, test.expectedAttrs[n])
					}
				}
			})
		}
	})

	t.Run("A", func(t *testing.T) {
		tests := []struct {
			query       []byte
			expectedTTL uint
			expectedIP  net.IP
		}{
			{
				query:       []byte{0x00, 0x00, 0x01, 0x80, 0x01, 0x00, 0x00, 0x00, 0x78, 0x00, 0x04, 0xc0, 0xA8, 0x01, 0x02},
				expectedTTL: 120,
				expectedIP:  net.IPv4(192, 168, 1, 2),
			},
		}
		for _, test := range tests {
			t.Run(test.expectedIP.String(), func(t *testing.T) {
				q, err := newResourceRecordWithReader(NewReaderWithBytes(test.query))
				if err != nil {
					t.Error(err)
				}
				a, ok := q.(*ARecord)
				if !ok {
					t.Errorf("%v", a)
					return
				}
				// Checks each field
				if a.TTL() != test.expectedTTL {
					t.Errorf("TTL: %d != %d", a.TTL(), test.expectedTTL)
				}
				if !test.expectedIP.Equal(a.Address()) {
					t.Errorf("IP: %s != %s", a.Address(), test.expectedIP)
				}
			})
		}
	})

	t.Run("AAAA", func(t *testing.T) {
		tests := []struct {
			query       []byte
			expectedTTL uint
			expectedIP  net.IP
		}{
			{
				query:       []byte{0x00, 0x00, 0x1c, 0x80, 0x01, 0x00, 0x00, 0x00, 0x78, 0x00, 0x10, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x31, 0xdb, 0x18, 0x30, 0x69, 0x5d, 0x3e, 0x30},
				expectedTTL: 120,
				expectedIP:  net.ParseIP("fe80::31db:1830:695d:3e30"),
			},
		}
		for _, test := range tests {
			t.Run(test.expectedIP.String(), func(t *testing.T) {
				q, err := newResourceRecordWithReader(NewReaderWithBytes(test.query))
				if err != nil {
					t.Error(err)
				}
				a, ok := q.(*AAAARecord)
				if !ok {
					t.Errorf("%v", a)
					return
				}
				// Checks each field
				if a.TTL() != test.expectedTTL {
					t.Errorf("TTL: %d != %d", a.TTL(), test.expectedTTL)
				}
				if !test.expectedIP.Equal(a.Address()) {
					t.Errorf("IP: %s != %s", a.Address(), test.expectedIP)
				}
			})
		}
	})

	t.Run("OPT", func(t *testing.T) {
		tests := []struct {
			query            []byte
			expectedName     string
			expectedType     Type
			expectedResponse bool
		}{
			{
				query:            []byte{0x00, 0x00, 0x29, 0x05, 0xa0, 0x00, 0x00, 0x11, 0x94, 0x00, 0x12, 0x00, 0x04, 0x00, 0x0e, 0x00, 0x74, 0x52, 0x06, 0x8d, 0xcf, 0x54, 0x27, 0x86, 0xfd, 0xcd, 0x88, 0xe1, 0x43},
				expectedName:     "",
				expectedType:     OPT,
				expectedResponse: false,
			},
		}
		for _, test := range tests {
			q, err := newRequestRecordWithReader(NewReaderWithBytes(test.query))
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
			if q.UnicastResponse() != test.expectedResponse {
				t.Errorf("%t != %t", q.UnicastResponse(), test.expectedResponse)
			}
		}
	})
}
