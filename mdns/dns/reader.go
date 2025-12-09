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
	"errors"
	"fmt"
	"io"

	"github.com/cybergarage/go-mdns/mdns/encoding"
)

// Reader represents a record reader.
type Reader struct {
	cmpBytes   []byte
	buffer     []byte
	bufferSize int
	offset     int
}

// NewReaderWithBytes returns a new reader instance with the specified bytes.
func NewReaderWithBytes(b []byte) *Reader {
	return &Reader{
		cmpBytes:   b,
		buffer:     b,
		bufferSize: len(b),
		offset:     0,
	}
}

// SetCompressionBytes sets the compression bytes.
func (reader *Reader) SetCompressionBytes(b []byte) {
	reader.cmpBytes = b
}

// CompressionBytes returns the compression bytes.
func (reader *Reader) CompressionBytes() []byte {
	return reader.cmpBytes
}

// Read overwrites the io.Reader interface.
func (reader *Reader) Read(p []byte) (int, error) {
	if reader.bufferSize <= reader.offset {
		return 0, io.EOF
	}
	n := copy(p, reader.buffer[reader.offset:])
	reader.offset += n
	if n == 0 {
		return 0, io.EOF
	}
	return n, nil
}

// ReadUint8 returns a uint8 from the reader.
func (reader *Reader) ReadUint8() (uint8, error) {
	if reader.bufferSize < (reader.offset + 1) {
		return 0, io.EOF
	}
	v := uint8(reader.buffer[reader.offset])
	reader.offset++
	return v, nil
}

// ReadUint16 returns a uint16 from the reader.
func (reader *Reader) ReadUint16() (uint16, error) {
	if reader.bufferSize < (reader.offset + 2) {
		return 0, io.EOF
	}
	v := encoding.BytesToInteger(reader.buffer[reader.offset : reader.offset+2])
	reader.offset += 2
	return uint16(v), nil
}

// ReadUint32 returns a uint32 from the reader.
func (reader *Reader) ReadUint32() (uint32, error) {
	if reader.bufferSize < (reader.offset + 4) {
		return 0, io.EOF
	}
	v := encoding.BytesToInteger(reader.buffer[reader.offset : reader.offset+4])
	reader.offset += 4
	return uint32(v), nil
}

// ReadString returns a string from the reader.
func (reader *Reader) ReadString() (string, error) {
	l, err := reader.ReadUint8()
	if err != nil {
		return "", err
	}
	if l == 0 {
		return "", nil
	}
	strBytes := make([]byte, l)
	_, err = reader.Read(strBytes)
	if err != nil {
		return "", err
	}
	return string(strBytes), nil
}

// ReadStrings returns strings from the reader.
func (reader *Reader) ReadStrings() ([]string, error) {
	strs := make([]string, 0)
	str, err := reader.ReadString()
	for err == nil {
		if len(str) == 0 {
			break
		}
		strs = append(strs, str)
		str, err = reader.ReadString()
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	return strs, nil
}

// ReadName returns a name from the reader with the read reader.
func (reader *Reader) ReadName() (string, error) {
	nameLenIsCompressed := func(l uint8) bool {
		return (l & nameIsCompressionMask) == nameIsCompressionMask
	}

	name := ""
	nextNameLen, err := reader.ReadUint8()
	for err == nil {
		if nameLenIsCompressed(nextNameLen) {
			// RFC1035: 4.1.4. Message compression
			cmpOffsetByte1 := int(nextNameLen)
			cmpOffsetByte2, err := reader.ReadUint8()
			if err != nil {
				return "", err
			}
			cmpOffset := ((cmpOffsetByte1 & int(nameLenMask)) << 8) | int(cmpOffsetByte2)
			cmpBytes := reader.CompressionBytes()
			if cmpOffset < 0 || len(cmpBytes) <= cmpOffset {
				return "", fmt.Errorf("invalid compression offset : %d", cmpOffset)
			}
			cmpReader := NewReaderWithBytes(cmpBytes[cmpOffset:])
			cmpReader.SetCompressionBytes(cmpBytes)
			nextName, err := cmpReader.ReadName()
			if err != nil {
				return "", err
			}
			// Per RFC1035, a compression pointer terminates the current name.
			return AppendName(name, nextName), nil
		}
		if nextNameLen == 0 {
			break
		}
		nextNameBytes := make([]byte, nextNameLen)
		_, err = reader.Read(nextNameBytes)
		if err != nil {
			return "", err
		}
		nextName := string(nextNameBytes)
		name = AppendName(name, nextName)
		nextNameLen, err = reader.ReadUint8()
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return name, nil
}
