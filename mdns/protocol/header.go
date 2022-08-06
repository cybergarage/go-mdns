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

// RFC 1035: Domain names - implementation and specification
// https://www.rfc-editor.org/rfc/rfc1035.html
// RFC 6762: Multicast DNS
// https://www.rfc-editor.org/rfc/rfc6762.html

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/cybergarage/go-mdns/mdns/encoding"
)

const (
	headerSize = 12
)

// Header represents a protocol header.
type Header struct {
	bytes []byte
}

// NewHeader returns a header instance.
func NewHeader() *Header {
	header := &Header{
		bytes: make([]byte, headerSize),
	}
	return header
}

// NewHeaderWithBytes returns a header instance with the specified bytes.
func NewHeaderWithBytes(bytes []byte) *Header {
	return &Header{
		bytes: bytes,
	}
}

// Parse parses the specified reader.
func (header *Header) Parse(reader io.Reader) error {
	header.bytes = make([]byte, headerSize)
	n, err := reader.Read(header.bytes)
	if err != nil {
		return err
	}
	if n != headerSize {
		return fmt.Errorf(errorHeaderShortLength, n)
	}
	return nil
}

// ID returns the query identifier.
// RFC 6762: 18.1. ID (Query Identifier)
// In multicast query messages, the Query Identifier SHOULD be set to zero on transmission.
// In multicast responses, including unsolicited multicast responses, the Query Identifier MUST be set to zero on transmission, and MUST be ignored on reception.
func (header *Header) ID() uint {
	return encoding.BytesToInteger(header.bytes[:2])
}

// Equals returns true if the header is same as the specified header, otherwise false.
func (header *Header) Equals(other *Header) bool {
	return bytes.Equal(header.bytes, other.bytes)
}

// Copy returns the copy header instance.
func (header *Header) Copy() *Header {
	return NewHeaderWithBytes(header.bytes)
}

// Bytes returns the binary representation.
func (header *Header) Bytes() []byte {
	return header.bytes
}

// String returns the string representation.
func (header *Header) String() string {
	return hex.EncodeToString(header.bytes)
}
