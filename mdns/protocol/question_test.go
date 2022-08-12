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
	tests := []struct {
		query    []byte
		expected string
	}{
		{
			query:    []byte{},
			expected: "",
		},
	}
	// 02 6c 62 07 5f 64 6e 73 2d 73 64 04 5f 75 64 70 05 6c 6f 63 61 6c 00 00 0c 80 01
	// lb._dns-sd._udp.local
	for _, test := range tests {
	}
}
