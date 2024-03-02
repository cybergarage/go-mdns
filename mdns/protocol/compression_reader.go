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
	"io"
)

// CompressionReader represents a read bytes reader.
type CompressionReader struct {
	bytes []byte
	*bytes.Buffer
}

// NewCompressionReaderWithBytes returns a new reader instance with the specified bytes.
func NewCompressionReaderWithBytes(b []byte) *CompressionReader {
	reader := &CompressionReader{
		bytes:  b,
		Buffer: bytes.NewBuffer(b),
	}
	return reader
}

// Read overwrites the io.Reader interface.
func (reader *CompressionReader) Read(p []byte) (int, error) {
	n, err := reader.Buffer.Read(p)
	if err != nil {
		return n, err
	}
	return n, nil
}

// Skip skips the specified bytes.
func (reader *CompressionReader) Skip(n int) error {
	if b := reader.Buffer.Next(n); len(b) != n {
		return io.EOF
	}
	return nil
}

// Bytes returns the read bytes.
func (reader *CompressionReader) Bytes() []byte {
	return reader.bytes
}
