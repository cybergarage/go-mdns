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

// ReadReader represents a read bytes reader.
type ReadReader struct {
	*bytes.Buffer
}

// NewReadReaderWithBytes returns a new reader instance with the specified bytes.
func NewReadReaderWithBytes(b []byte) *ReadReader {
	reader := &ReadReader{
		Buffer: bytes.NewBuffer(b),
	}
	return reader
}

// Read overwrites the io.Reader interface.
func (reader *ReadReader) Read(p []byte) (int, error) {
	n, err := reader.Buffer.Read(p)
	if err != nil {
		return n, err
	}
	return n, nil
}

// Skip skips the specified bytes.
func (reader *ReadReader) Skip(n int) error {
	if b := reader.Buffer.Next(n); len(b) != n {
		return io.EOF
	}
	return nil
}
