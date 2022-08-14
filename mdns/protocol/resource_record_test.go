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
	t.Run("ParseAuthoritative", func(t *testing.T) {
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
