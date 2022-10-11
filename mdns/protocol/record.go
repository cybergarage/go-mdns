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
		class:           0,
		ttl:             0,
		data:            nil,
	}
}

// newRecordWithReader returns a new record innstance with the specified reader.
func newRecordWithReader(reader io.Reader) (*Record, error) {
	return newRequestRecordWithReader(reader)
}

// newResourceRecordWithReader returns a new resource record innstance with the specified reader.
func newResourceRecordWithReader(reader io.Reader) (ResourceRecord, error) {
	return newResponseResourceRecordWithReader(reader)
}

// newRequestRecordWithReader returns a new request resource record innstance with the specified reader.
func newRequestRecordWithReader(reader io.Reader) (*Record, error) {
	r := newResourceRecord()
	return r, r.ParseRequest(reader)
}

// newRequestResourceRecordWithReader returns a new request resource record innstance with the specified reader.
func newRequestResourceRecordWithReader(reader io.Reader) (ResourceRecord, error) {
	r, err := newRequestRecordWithReader(reader)
	if err != nil {
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

// newResponseRecordWithReader returns a new response resource record innstance with the specified reader.
func newResponseRecordWithReader(reader io.Reader) (*Record, error) {
	r := newResourceRecord()
	return r, r.ParseResponse(reader)
}

// newResponseResourceRecordWithReader returns a new response resource record innstance with the specified reader.
func newResponseResourceRecordWithReader(reader io.Reader) (ResourceRecord, error) {
	r, err := newResponseRecordWithReader(reader)
	if err != nil {
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
	return r.unicastResponse
}

// Class returns the resource record class.
func (r *Record) Class() Class {
	return r.class
}

// TTL returns the TTL second.
func (r *Record) TTL() uint {
	return r.ttl
}

// Data returns the resource record data.
func (r *Record) Data() []byte {
	return r.data
}

// ParseRequest parses a request record from the specified reader.
func (r *Record) ParseRequest(reader io.Reader) error {
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
	r.typ = Type(encoding.BytesToInteger(typeBytes))

	// Parses class type
	classBytes := make([]byte, 2)
	_, err = reader.Read(classBytes)
	if err != nil {
		return err
	}
	cls := encoding.BytesToInteger(classBytes)
	r.unicastResponse = false
	if (cls & unicastResponseMask) != 0 {
		r.unicastResponse = true
	}
	r.class = Class(cls & classMask)

	return nil
}

// ParseResponse parses a response record from the specified reader.
func (r *Record) ParseResponse(reader io.Reader) error {
	var err error

	err = r.ParseRequest(reader)
	if err != nil {
		return err
	}

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

// RequestBytes returns only the binary representation of the request fields.
func (r *Record) RequestBytes() []byte {
	bytes := nameToBytes(r.name)

	typeBytes := make([]byte, 2)
	typ := r.typ
	bytes = append(bytes, encoding.IntegerToBytes(uint(typ), typeBytes)...)

	classBytes := make([]byte, 2)
	cls := r.class
	if r.unicastResponse {
		cls |= cacheFlushMask
	}
	bytes = append(bytes, encoding.IntegerToBytes(uint(cls), classBytes)...)

	return bytes
}

// ResponseBytes returns only the binary representation of the all fields.
func (r *Record) ResponseBytes() []byte {
	bytes := r.RequestBytes()

	ttlBytes := make([]byte, 4)
	bytes = append(bytes, encoding.IntegerToBytes(r.ttl, ttlBytes)...)

	dataLenBytes := make([]byte, 2)
	bytes = append(bytes, encoding.IntegerToBytes(uint(len(r.data)), dataLenBytes)...)
	bytes = append(bytes, r.data...)

	return bytes
}

// Bytes returns the binary representation.
func (r *Record) Bytes() []byte {
	return r.ResponseBytes()
}

// Equal returns true if this record is equal to  the specified resource record. otherwise false.
func (r *Record) Equal(other ResourceRecord) bool {
	return bytes.Equal(r.Bytes(), other.Bytes())
}
