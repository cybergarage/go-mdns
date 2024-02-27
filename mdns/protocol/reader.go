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

// Reader represents a record reader.
type Reader struct {
	Reader io.Reader
	Bytes  []byte
}

// NewReaderWithReader returns a new reader instance with the specified reader.
func NewReaderWithReader(reader io.Reader) *Reader {
	return &Reader{
		Reader: reader,
		Bytes:  []byte{},
	}
}

// NewReaderWithBytes returns a new reader instance with the specified bytes.
func NewReaderWithBytes(b []byte) *Reader {
	return NewReaderWithReader(bytes.NewReader(b))
}

// Read overwrites the io.Reader interface.
func (reader *Reader) Read(p []byte) (int, error) {
	n, err := reader.Reader.Read(p)
	if err != nil {
		return n, err
	}
	reader.Bytes = append(reader.Bytes, p[:n]...)
	return n, nil
}

// ReadReader returns a read reader instance.
func (reader *Reader) ReadReader() *ReadReader {
	return NewReadReaderWithBytes(reader.Bytes)
}
