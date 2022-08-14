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
	"io"

	"github.com/cybergarage/go-mdns/mdns/encoding"
)

// ResourceRecord represents a resource record interface.
type ResourceRecord interface {
	// Name returns the resource record name.
	Name() string
	// Type returns the resource record type.
	Type() Type
	// Class returns the resource record class.
	Class() Class
	// Type returns the cache flush flag.
	CacheFlush() bool
	// TTL returns the TTL second.
	TTL() uint
	// Data returns the resource record data.
	Data() []byte
	// Bytes returns the binary representation.
	Bytes() []byte
}

// resourceRecord represents a resource record.
type resourceRecord struct {
	name       string
	typ        Type
	cacheFlush bool
	class      Class
	ttl        uint
	data       []byte
}

// newResourceRecordWithReader returns a new question innstance with the specified reader.
func newResourceRecordWithReader(reader io.Reader) (ResourceRecord, error) {
	res := &resourceRecord{
		name:       "",
		typ:        0,
		cacheFlush: false,
		class:      0,
		ttl:        0,
		data:       nil,
	}

	if err := res.parse(reader); err != nil {
		return nil, err
	}

	switch res.Type() {
	case PTR:
		return NewPTRRecord(res), nil
	case SRV:
		return NewSRVRecord(res), nil
	case A:
		return NewARecord(res), nil
	case AAAA:
		return NewAAAARecord(res), nil
	}

	return res, nil
}

// Name returns the resource record name.
func (res *resourceRecord) Name() string {
	return res.name
}

// Type returns the resource record type.
func (res *resourceRecord) Type() Type {
	return res.typ
}

// Class returns the resource record class.
func (res *resourceRecord) Class() Class {
	return res.class
}

// Type returns the cache flush flag.
func (res *resourceRecord) CacheFlush() bool {
	return res.cacheFlush
}

// TTL returns the TTL second.
func (res *resourceRecord) TTL() uint {
	return res.ttl
}

// Data returns the resource record data.
func (res *resourceRecord) Data() []byte {
	return res.data
}

// parse parses the specified reader.
func (res *resourceRecord) parse(reader io.Reader) error {
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
	res.typ = Type(encoding.BytesToInteger(typeBytes))

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
	res.class = Class(class & classMask)

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
func (res *resourceRecord) Bytes() []byte {
	bytes := nameToBytes(res.name)

	typeBytes := make([]byte, 2)
	bytes = append(bytes, encoding.IntegerToBytes(uint(res.typ), typeBytes)...)

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
