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
	"bytes"
	"strings"
)

// Writer represents a record writer.
type Writer struct {
	*bytes.Buffer
}

// NewWriter returns a new writer instance.
func NewWriter() *Writer {
	return &Writer{
		Buffer: &bytes.Buffer{},
	}
}

// WriteHeader writes a header.
func (writer *Writer) WriteHeader(header *Header) error {
	return writer.WriteBytes(header.Bytes())
}

// WriteByte writes a byte value.
func (writer *Writer) WriteUint8(v uint8) error {
	return writer.WriteByte(v)
}

// WriteUint16 writes a uint16 value.
func (writer *Writer) WriteUint16(v uint16) error {
	return writer.WriteBytes([]byte{byte(v >> 8), byte(v)})
}

// WriteUint32 writes a uint32 value.
func (writer *Writer) WriteUint32(v uint32) error {
	return writer.WriteBytes([]byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)})
}

// WriteBytes writes bytes.
func (writer *Writer) WriteBytes(v []byte) error {
	_, err := writer.Write(v)
	return err
}

// WriteType writes a type.
func (writer *Writer) WriteType(v Type) error {
	return writer.WriteUint16(uint16(v))
}

// WriteClass writes a class.
func (writer *Writer) WriteClass(v Class) error {
	return writer.WriteUint16(uint16(v))
}

// WriteTTL writes a TTL.
func (writer *Writer) WriteTTL(v uint) error {
	return writer.WriteUint32(uint32(v))
}

// WriteData writes data.
func (writer *Writer) WriteData(v []byte) error {
	if err := writer.WriteUint16(uint16(len(v))); err != nil {
		return err
	}
	return writer.WriteBytes(v)
}

// WriteString writes a string with a length.
func (writer *Writer) WriteString(v string) error {
	if err := writer.WriteUint8(uint8(len(v))); err != nil {
		return err
	}
	_, err := writer.Buffer.WriteString(v)
	return err
}

// WriteName writes a name.
func (writer *Writer) WriteName(name string) error {
	labels := strings.Split(name, ".")
	for _, label := range labels {
		if err := writer.WriteString(label); err != nil {
			return err
		}
	}
	return writer.WriteByte(0)
}
