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
	"io"

	"github.com/cybergarage/go-mdns/mdns/encoding"
)

// Reader represents a record reader.
type Reader struct {
	Buffer     []byte
	bufferSize int
	offset     int
}

// NewReaderWithBytes returns a new reader instance with the specified bytes.
func NewReaderWithBytes(b []byte) *Reader {
	return &Reader{
		Buffer:     b,
		bufferSize: int(len(b)),
		offset:     0,
	}
}

// Read overwrites the io.Reader interface.
func (reader *Reader) Read(p []byte) (int, error) {
	if reader.bufferSize < (reader.offset + len(p)) {
		return 0, io.EOF
	}
	copy(p, reader.Buffer[reader.offset:])
	reader.offset += len(p)
	return len(p), nil
}

// ReadUint8 returns a uint8 from the reader.
func (reader *Reader) ReadUint8() (uint8, error) {
	if reader.bufferSize < (reader.offset + 1) {
		return 0, io.EOF
	}
	v := uint8(reader.Buffer[reader.offset])
	reader.offset++
	return v, nil
}

// ReadUint16 returns a uint16 from the reader.
func (reader *Reader) ReadUint16() (uint16, error) {
	if reader.bufferSize < (reader.offset + 2) {
		return 0, io.EOF
	}
	v := encoding.BytesToInteger(reader.Buffer[reader.offset : reader.offset+2])
	reader.offset += 2
	return uint16(v), nil
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
			remainNameOffset, err := reader.ReadUint8()
			if err != nil {
				return "", err
			}
			nameOffset := (int(nextNameLen & ^nameIsCompressionMask) << 8) + int(remainNameOffset)

			compReader := reader.CompressionReader()
			if err := compReader.Skip(nameOffset); err != nil {
				return "", err
			}

			nextCompReader := NewReaderWithReader(compReader)
			nextCompReader.SetCompressionReader(NewCompressionReaderWithBytes(compReader.Bytes()))
			return nextCompReader.ReadName()
		}
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
		nextNameLen, err = reader.ReadUint8()
	}
	if err != nil {
		return "", err
	}
	return name, nil
}

// SetCompressionReader sets a read reader instance.
func (reader *Reader) SetCompressionReader(cmpReader *CompressionReader) *Reader {
	reader.rootCmpReader = cmpReader
	return reader
}

// CompressionReader returns a read reader instance.
func (reader *Reader) CompressionReader() *CompressionReader {
	if reader.rootCmpReader != nil {
		return reader.rootCmpReader
	}
	return NewCompressionReaderWithBytes(reader.Buffer)
}
