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
	"errors"
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
	r := newResourceRecord()
	if err := r.Parse(reader); err != nil {
		return nil, err
	}

	switch r.Type() {
	case PTR:
		return newPTRRecordWithResourceRecord(r), nil
	case SRV:
		return newSRVRecordWithResourceRecord(r), nil
	case TXT:
		return newTXTRecordWithResourceRecord(r), nil
	case A:
		return newARecordWithResourceRecord(r), nil
	case AAAA:
		return newAAAARecordWithResourceRecord(r), nil
	}

	return r, nil
}

// SetName sets the specified name.
func (r *Record) SetName(name string) *Record {
	r.name = name
	return r
}

// SetUnicastResponse sets the specified unicast response flag.
func (r *Record) SetUnicastResponse(enabled bool) *Record {
	r.unicastResponse = enabled
	return r
}

// SetType sets the specified resource record type.
func (r *Record) SetType(typ Type) *Record {
	r.typ = typ
	return r
}

// SetClass sets the specified resource record class.
func (r *Record) SetClass(cls Class) *Record {
	r.class = cls
	return r
}

// SetCacheFlush sets the specified cache flush flag.
func (r *Record) SetCacheFlush(enabled bool) *Record {
	r.cacheFlush = enabled
	return r
}

// SetTTL returns the specified TTL second.
func (r *Record) SetTTL(ttl uint) *Record {
	r.ttl = ttl
	return r
}

// SetData returns the specified resource record data.
func (r *Record) SetData(b []byte) *Record {
	r.data = b
	return r
}

// Name returns the resource record name.
func (r *Record) Name() string {
	return r.name
}

// Type returns the resource record type.
func (r *Record) Type() Type {
	return r.typ
}

// UnicastResponse returns the unicast response flag.
func (r *Record) UnicastResponse() bool {
	return r.cacheFlush
}

// Class returns the resource record class.
func (r *Record) Class() Class {
	return r.class
}

// CacheFlush returns the cache flush flag.
func (r *Record) CacheFlush() bool {
	return r.cacheFlush
}

// TTL returns the TTL second.
func (r *Record) TTL() uint {
	return r.ttl
}

// Data returns the resource record data.
func (r *Record) Data() []byte {
	return r.data
}

// Parse parses the specified reader.
func (r *Record) Parse(reader io.Reader) error {
	var err error

	// Parses domain names
	r.name, err = parseName(reader)
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
	r.unicastResponse = false
	if (typ & unicastResponseMask) != 0 {
		r.unicastResponse = true
	}
	r.typ = typ & typeMask

	// Parses class type
	classBytes := make([]byte, 2)
	_, err = reader.Read(classBytes)
	if err != nil {
		return err
	}
	cls := encoding.BytesToInteger(classBytes)
	r.cacheFlush = false
	if (cls & cacheFlushMask) != 0 {
		r.cacheFlush = true
	}
	r.class = Class(cls & classMask)

	// Parses TTL
	ttlBytes := make([]byte, 4)
	_, err = reader.Read(ttlBytes)
	if err != nil {
		if errors.Is(err, io.EOF) { // QR == 0
			return nil
		}
		return err
	}
	r.ttl = encoding.BytesToInteger(ttlBytes)

	// Parses data
	dataLenBytes := make([]byte, 2)
	_, err = reader.Read(dataLenBytes)
	if err != nil {
		return err
	}
	dataLen := encoding.BytesToInteger(dataLenBytes)
	r.data = make([]byte, dataLen)
	if 0 < dataLen {
		_, err = reader.Read(r.data)
		if err != nil {
			return err
		}
	}

	return nil
}

// Bytes returns the binary representation.
func (r *Record) Bytes() []byte {
	bytes := nameToBytes(r.name)

	typeBytes := make([]byte, 2)
	typ := r.typ
	if r.unicastResponse {
		typ |= unicastResponseMask
	}
	bytes = append(bytes, encoding.IntegerToBytes(uint(typ), typeBytes)...)

	classBytes := make([]byte, 2)
	cls := r.class
	if r.cacheFlush {
		cls |= cacheFlushMask
	}
	bytes = append(bytes, encoding.IntegerToBytes(uint(cls), classBytes)...)

	ttlBytes := make([]byte, 4)
	bytes = append(bytes, encoding.IntegerToBytes(r.ttl, ttlBytes)...)

	dataLenBytes := make([]byte, 2)
	bytes = append(bytes, encoding.IntegerToBytes(uint(len(r.data)), dataLenBytes)...)
	bytes = append(bytes, r.data...)

	return bytes
}

// Equal returns true if this record is equal to  the specified resource record. otherwise false.
func (r *Record) Equal(other ResourceRecord) bool {
	return bytes.Equal(r.Bytes(), other.Bytes())
}
