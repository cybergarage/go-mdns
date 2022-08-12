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

type ResourceType = Type
type ResourceClass = Class

// Resource represents a question.
type Resource struct {
	Name       string
	Type       ResourceType
	CacheFlush bool
	Class      ResourceClass
	TTL        uint
	Data       []byte
}

// NewResource returns a new question innstance.
func NewResource() *Resource {
	return &Resource{
		Name:       "",
		Type:       0,
		CacheFlush: false,
		Class:      0,
		TTL:        0,
		Data:       nil,
	}
}

// NewResourceWithReader returns a new question innstance with the specified reader.
func NewResourceWithReader(reader io.Reader) (*Resource, error) {
	res := NewResource()
	return res, res.Parse(reader)
}

// Parse parses the specified reader.
func (res *Resource) Parse(reader io.Reader) error {
	var err error

	// Parses domain names
	res.Name, err = parseName(reader)
	if err != nil {
		return err
	}

	// Parses query type
	typeBytes := make([]byte, 2)
	_, err = reader.Read(typeBytes)
	if err != nil {
		return err
	}
	res.Type = ResourceType(encoding.BytesToInteger(typeBytes))

	// Parses class type
	classBytes := make([]byte, 2)
	_, err = reader.Read(classBytes)
	if err != nil {
		return err
	}
	class := encoding.BytesToInteger(classBytes)
	res.CacheFlush = false
	if (class & cacheFlushMask) != 0 {
		res.CacheFlush = true
	}
	res.Class = ResourceClass(class & classMask)

	// Parses TTL
	ttlBytes := make([]byte, 4)
	_, err = reader.Read(ttlBytes)
	if err != nil {
		return err
	}
	res.TTL = encoding.BytesToInteger(ttlBytes)

	// Parses data
	dataLenBytes := make([]byte, 2)
	_, err = reader.Read(dataLenBytes)
	if err != nil {
		return err
	}
	dataLen := encoding.BytesToInteger(dataLenBytes)
	res.Data = make([]byte, dataLen)
	_, err = reader.Read(res.Data)
	if err != nil {
		return err
	}

	return nil
}

// Bytes returns the binary representation.
func (res *Resource) Bytes() []byte {
	bytes := nameToBytes(res.Name)

	typeBytes := make([]byte, 2)
	bytes = append(bytes, encoding.IntegerToBytes(uint(res.Type), typeBytes)...)

	classBytes := make([]byte, 2)
	class := res.Class
	if res.CacheFlush {
		class |= cacheFlushMask
	}
	bytes = append(bytes, encoding.IntegerToBytes(uint(class), classBytes)...)

	ttlBytes := make([]byte, 4)
	bytes = append(bytes, encoding.IntegerToBytes(res.TTL, ttlBytes)...)

	dataLenBytes := make([]byte, 2)
	bytes = append(bytes, encoding.IntegerToBytes(uint(len(res.Data)), dataLenBytes)...)
	bytes = append(bytes, res.Data...)

	return bytes
}
