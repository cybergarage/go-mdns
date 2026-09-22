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

package dns

import (
	"testing"
)

// FuzzParseMessage checks that the parser does not panic or hang for any input,
// because the parsed messages are received from the network.
func FuzzParseMessage(f *testing.F) {
	seeds := [][]byte{
		{},
		{0x00},
		// A query message for "_services._dns-sd._udp.local".
		{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x09, '_', 's', 'e', 'r', 'v', 'i', 'c', 'e', 's',
			0x07, '_', 'd', 'n', 's', '-', 's', 'd',
			0x04, '_', 'u', 'd', 'p',
			0x05, 'l', 'o', 'c', 'a', 'l',
			0x00, 0x00, 0x0c, 0x00, 0x01,
		},
		// A message which holds compression pointers referring to each other.
		{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0xC0, 0x0E, 0x00, 0x0c, 0x00, 0x01, 0xC0, 0x0C,
		},
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, msgBytes []byte) {
		msg, err := NewMessageWithBytes(msgBytes)
		if err != nil {
			return
		}
		// The parsed message must be printable without panicking.
		_ = msg.String()
		_ = msg.Bytes()
	})
}
