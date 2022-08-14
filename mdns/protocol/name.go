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
	"io"
	"strings"

	"github.com/cybergarage/go-mdns/mdns/encoding"
)

const (
	nameSep               = "."
	nameIsCompressionMask = 0xC0
	nameLenMask           = 0x3F
)

func parseName(reader io.Reader) (string, error) {
	name := ""
	nextNameLenBuf := make([]byte, 1)
	_, err := reader.Read(nextNameLenBuf)
	for err == nil {
		nextNameLen := encoding.BytesToInteger(nextNameLenBuf)
		// Note: Compression names are not supported, and so
		// the alias name are indistinguishable from root names yet.
		if (nextNameLen & nameIsCompressionMask) == nameIsCompressionMask {
			// Skips a remain compression offset bit
			_, err := reader.Read(nextNameLenBuf)
			if err != nil {
				return "", err
			}
			return "", nil
		}
		nextNameLen &= nameLenMask
		if nextNameLen == 0 {
			break
		}
		nextName := make([]byte, nextNameLen)
		_, err = reader.Read(nextName)
		if err != nil {
			return "", err
		}
		if 0 < len(name) {
			name += nameSep
		}
		name += string(nextName)
		_, err = reader.Read(nextNameLenBuf)
	}
	if err != nil {
		return "", err
	}
	return name, nil
}

func nameToBytes(name string) []byte {
	bytes := []byte{}
	tokens := strings.Split(name, nameSep)
	for _, token := range tokens {
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
