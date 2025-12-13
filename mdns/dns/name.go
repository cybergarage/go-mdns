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
	"strings"
)

const (
	LabelSeparator        = "."
	nameIsCompressionMask = uint8(0xC0)
	nameLenMask           = uint8(0x3F)
)

// NewNameWithStrings returns a DNS name constructed by joining the given strings with dots.
func NewNameWithStrings(s ...string) string {
	return strings.Join(s, LabelSeparator)
}

// AppendName appends the given name to the base name.
func AppendName(base string, name string) string {
	if len(base) == 0 {
		return name
	}
	return base + LabelSeparator + name
}

// SplitName splits the given name into its labels.
func SplitName(name string) []string {
	return strings.Split(name, LabelSeparator)
}

func nameToBytes(name string) []byte {
	bytes := []byte{}
	tokens := strings.SplitSeq(name, LabelSeparator)
	for token := range tokens {
		if len(token) == 0 {
			continue
		}
		nameLen := byte(len(token) & 0xFF)
		bytes = append(bytes, nameLen)
		bytes = append(bytes, []byte(token)...)
	}
	bytes = append(bytes, 0x00)
	return bytes
}
