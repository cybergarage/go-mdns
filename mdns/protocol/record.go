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

// Record represents a resource record.
type Record struct {
	name            string
	unicastResponse bool
	typ             Type
	cacheFlush      bool
	class           Class
	ttl             uint
	data            []byte
}

// newResourceRecord returns a new resource record innstance.
func newResourceRecord() *Record {
	return &Record{
		name:            "",
		unicastResponse: false,
		typ:             0,
		cacheFlush:      false,
		class:           0,
		ttl:             0,
		data:            nil,
	}
}

// newRecordWithReader returns a new resource record innstance with the specified reader.
func newRecordWithReader(reader io.Reader) (ResourceRecord, error) {
	res := newResourceRecord()
	if err := res.Parse(reader); err != nil {
		return nil, err
	}

	switch res.Type() {
	case PTR:
		return newPTRRecordWithResourceRecord(res), nil
	case SRV:
		return newSRVRecordWithResourceRecord(res), nil
	case TXT:
		return newTXTRecordWithResourceRecord(res), nil
	case A:
		return newARecordWithResourceRecord(res), nil
	case AAAA:
		return newAAAARecordWithResourceRecord(res), nil
	}

	return res, nil
}

// SetName sets the specified name.
func (res *Record) SetName(name string) *Record {
	res.name = name
	return res
}

// SetUnicastResponse sets the specified unicast response flag.
func (res *Record) SetUnicastResponse(enabled bool) *Record {
	res.unicastResponse = enabled
	return res
}

// SetType sets the specified resource record type.
func (res *Record) SetType(typ Type) *Record {
	res.typ = typ
	return res
}

// SetClass sets the specified resource record class.
func (res *Record) SetClass(cls Class) *Record {
	res.class = cls
	return res
}

// SetCacheFlush sets the specified cache flush flag.
func (res *Record) SetCacheFlush(enabled bool) *Record {
	res.cacheFlush = enabled
	return res
}

// SetTTL returns the specified TTL second.
func (res *Record) SetTTL(ttl uint) *Record {
	res.ttl = ttl
	return res
}

// SetData returns the specified resource record data.
func (res *Record) SetData(b []byte) *Record {
	res.data = b
	return res
}

// Name returns the resource record name.
func (res *Record) Name() string {
	return res.name
}

// Type returns the resource record type.
func (res *Record) Type() Type {
	return res.typ
}

// UnicastResponse returns the unicast response flag.
func (res *Record) UnicastResponse() bool {
	return res.cacheFlush
}

// Class returns the resource record class.
func (res *Record) Class() Class {
	return res.class
}

// CacheFlush returns the cache flush flag.
func (res *Record) CacheFlush() bool {
	return res.cacheFlush
}

// TTL returns the TTL second.
func (res *Record) TTL() uint {
	return res.ttl
}

// Data returns the resource record data.
func (res *Record) Data() []byte {
	return res.data
}

// Parse parses the specified reader.
func (res *Record) Parse(reader io.Reader) error {
	var err error

	// Parses domain names
	res.name, err = parseName(reader)
	if err != nil {
		return err
	}

	// Parses query type
	typeBytes := make([]byte, 2)
	_, err = reader.Read(typeBytes)
	if err != nil {
		return err
	}
	typ := Type(encoding.BytesToInteger(typeBytes))
	res.unicastResponse = false
	if (typ & unicastResponseMask) != 0 {
		res.unicastResponse = true
	}
	res.typ = typ & (^unicastResponseMask & 0xFFFF)

	// Parses class type
	classBytes := make([]byte, 2)
	_, err = reader.Read(classBytes)
	if err != nil {
		return err
	}
	class := encoding.BytesToInteger(classBytes)
	res.cacheFlush = false
	if (class & cacheFlushMask) != 0 {
		res.cacheFlush = true
	}
	res.class = Class(class & (^classMask & 0xFFFF))

	// Parses TTL
	ttlBytes := make([]byte, 4)
	_, err = reader.Read(ttlBytes)
	if err != nil {
		return err
	}
	res.ttl = encoding.BytesToInteger(ttlBytes)

	// Parses data
	dataLenBytes := make([]byte, 2)
	_, err = reader.Read(dataLenBytes)
	if err != nil {
		return err
	}
	dataLen := encoding.BytesToInteger(dataLenBytes)
	res.data = make([]byte, dataLen)
	_, err = reader.Read(res.data)
	if err != nil {
		return err
	}

	return nil
}

// Bytes returns the binary representation.
func (res *Record) Bytes() []byte {
	bytes := nameToBytes(res.name)

	typeBytes := make([]byte, 2)
	typ := res.typ
	if res.unicastResponse {
		typ |= unicastResponseMask
	}
	bytes = append(bytes, encoding.IntegerToBytes(uint(typ), typeBytes)...)

	classBytes := make([]byte, 2)
	class := res.class
	if res.cacheFlush {
		class |= cacheFlushMask
	}
	bytes = append(bytes, encoding.IntegerToBytes(uint(class), classBytes)...)

	ttlBytes := make([]byte, 4)
	bytes = append(bytes, encoding.IntegerToBytes(res.ttl, ttlBytes)...)

	dataLenBytes := make([]byte, 2)
	bytes = append(bytes, encoding.IntegerToBytes(uint(len(res.data)), dataLenBytes)...)
	bytes = append(bytes, res.data...)

	return bytes
}

// Equal returns true if this record is equal to  the specified resource record. otherwise false.
func (res *Record) Equal(other ResourceRecord) bool {
	return bytes.Equal(res.Bytes(), other.Bytes())
}
