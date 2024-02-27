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

	"github.com/cybergarage/go-mdns/mdns/encoding"
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

// ReadString returns a string from the reader.
func (reader *Reader) ReadString() (string, error) {
	lenByte := make([]byte, 1)
	_, err := reader.Read(lenByte)
	if err != nil {
		return "", err
	}
	strBytes := make([]byte, int(lenByte[0]))
	_, err = reader.Read(strBytes)
	if err != nil {
		return "", err
	}
	return string(strBytes), nil
}

// ReadNameWith returns a name from the reader with the read reader.
func (reader *Reader) ReadNameWith(readReader *CompressionReader) (string, error) {
	nameLenIsCompressed := func(l uint) bool {
		return (l & nameIsCompressionMask) == nameIsCompressionMask
	}

	name := ""
	nextNameLenBuf := make([]byte, 1)
	_, err := reader.Read(nextNameLenBuf)
	for err == nil {
		nextNameField := encoding.BytesToInteger(nextNameLenBuf)
		if nameLenIsCompressed(nextNameField) {
			if readReader == nil {
				return "", ErrNilReader
			}
			remainNameOffsetBuf := make([]byte, 1)
			_, err := reader.Read(remainNameOffsetBuf)
			if err != nil {
				return "", err
			}
			remainNameOffset := encoding.BytesToInteger(remainNameOffsetBuf)
			nameOffset := int(((nextNameField & nameLenMask) << 8) + remainNameOffset)
			if err := readReader.Skip(nameOffset); err != nil {
				return "", err
			}
			return NewReaderWithReader(readReader).ReadNameWith(nil)
		}

		nextNameLen := int(nextNameField & nameLenMask)
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

// CompressionReader returns a read reader instance.
func (reader *Reader) CompressionReader() *CompressionReader {
	return NewCompressionReaderWithBytes(reader.Bytes)
}
