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

	"github.com/cybergarage/go-mdns/mdns/encoding"
)

func parseTxt(reader io.Reader) ([]string, error) {
	attrs := []string{}
	attrLenByte := make([]byte, 1)
	_, err := reader.Read(attrLenByte)
	for err == nil {
		attrLen := encoding.BytesToInteger(attrLenByte)
		if attrLen == 0 {
			break
		}
		attrBytes := make([]byte, attrLen)
		_, err = reader.Read(attrBytes)
		if err != nil {
			return nil, err
		}
		attrs = append(attrs, string(attrBytes))
		_, err = reader.Read(attrLenByte)
	}
	if err != nil {
		return nil, err
	}
	return attrs, nil
}

func txtToBytes(attrs []string) []byte {
	bytes := []byte{}
	for _, attr := range attrs {
		attrLen := byte(len(attr) & 0xFF)
		bytes = append(bytes, attrLen)
		bytes = append(bytes, []byte(attr)...)
	}
	bytes = append(bytes, 0x00)
	return bytes
}
