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
	"fmt"
	"strings"
	"unicode"

	"github.com/cybergarage/go-mdns/mdns/encoding"
)

// recordOptions represents a record option.
type recordOptions func(*record)

// record represents a base record.
type record struct {
	reader          *Reader
	name            string
	unicastResponse bool
	typ             Type
	class           Class
	ttl             uint
	data            []byte
	cmpBytes        []byte
}

// newRecord returns a new base record instance.
func newRecord(opts ...recordOptions) *record {
	r := &record{
		reader:          nil,
		name:            "",
		unicastResponse: false,
		typ:             0,
		class:           0,
		ttl:             0,
		data:            nil,
		cmpBytes:        nil,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// newRecordWithReader returns a new base record instance with the specified reader.
func newRecordWithReader(reader *Reader) *record {
	r := newRecord()
	r.reader = reader
	return r
}

// NewRequestRecordWithReader returns a new request resource record instance with the specified reader.
func NewRequestRecordWithReader(reader *Reader) (*record, error) {
	r := newRecordWithReader(reader)
	return r, r.parseSection(reader)
}

// NewResourceRecordWithReader returns a new resource record instance with the specified reader.
func NewResourceRecordWithReader(reader *Reader) (ResourceRecord, error) {
	r := newRecordWithReader(reader)
	err := r.parseResourceRecord(reader)
	if err != nil {
		return nil, err
	}
	r.SetCompressionBytes(reader.CompressionBytes())

	switch r.Type() {
	case PTR:
		return newPTRRecordWithResourceRecord(r)
	case SRV:
		return newSRVRecordWithResourceRecord(r)
	case TXT:
		return newTXTRecordWithResourceRecord(r)
	case A:
		return newARecordWithResourceRecord(r), nil
	case AAAA:
		return newAAAARecordWithResourceRecord(r), nil
	case NSEC:
		return newNSECRecordWithResourceRecord(r)
	}

	return r, nil
}

// Reader returns a record reader.
func (r *record) Reader() (*Reader, error) {
	if r.reader == nil {
		return nil, ErrNilReader
	}
	return r.reader, nil
}

// SetName sets the specified name.
func (r *record) SetName(name string) Record {
	r.name = name
	return r
}

// SetUnicastResponse sets the specified unicast response flag.
func (r *record) SetUnicastResponse(enabled bool) Record {
	r.unicastResponse = enabled
	return r
}

// SetType sets the specified resource record type.
func (r *record) SetType(typ Type) Record {
	r.typ = typ
	return r
}

// SetClass sets the specified resource record class.
func (r *record) SetClass(cls Class) Record {
	r.class = cls
	return r
}

// SetTTL returns the specified TTL second.
func (r *record) SetTTL(ttl uint) Record {
	r.ttl = ttl
	return r
}

// SetData returns the specified record data.
func (r *record) SetData(b []byte) Record {
	r.data = b
	return r
}

// Name returns the resource record name.
func (r *record) Name() string {
	return r.name
}

// IsName returns true if the resource record name is the specified name.
func (r *record) IsName(name string) bool {
	// RFC1035: 2.3.3. Character Case
	return strings.EqualFold(r.name, name)
}

// HasNamePrefix returns true if the resource record name has the specified prefix.
func (r *record) HasNamePrefix(prefix string) bool {
	return strings.HasPrefix(r.name, prefix)
}

// HasNameSuffix returns true if the resource record name has the specified suffix.
func (r *record) HasNameSuffix(suffix string) bool {
	return strings.HasSuffix(r.name, suffix)
}

// Type returns the resource record type.
func (r *record) Type() Type {
	return r.typ
}

// UnicastResponse returns the unicast response flag.
func (r *record) UnicastResponse() bool {
	return r.unicastResponse
}

// Class returns the resource record class.
func (r *record) Class() Class {
	return r.class
}

// TTL returns the TTL second.
func (r *record) TTL() uint {
	return r.ttl
}

// Data returns the record data.
func (r *record) Data() []byte {
	return r.data
}

// Content returns a string representation to the record data.
func (r *record) Content() string {
	var c strings.Builder
	for n := range r.data {
		rb := rune(r.data[n])
		if unicode.IsPrint(rb) {
			c.WriteString(fmt.Sprintf("%c", rb))
		} else {
			c.WriteString(LabelSeparator)
		}
	}
	return c.String()
}

func (r *record) parseSection(reader *Reader) error {
	// Parses domain names
	name, err := reader.ReadName()
	if err != nil {
		return err
	}
	r.name = name

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
	if (cls & uint(QU)) != 0 {
		r.unicastResponse = true
	}
	r.class = Class(cls & classMask)

	return nil
}

// parseResourceRecord parses a resource record from the specified reader.
func (r *record) parseResourceRecord(reader *Reader) error {
	var err error

	err = r.parseSection(reader)
	if err != nil {
		return err
	}

	// Parses TTL
	ttl, err := reader.ReadUint32()
	if err != nil {
		return err
	}
	r.ttl = uint(ttl)

	// Parses data
	dataLen, err := reader.ReadUint16()
	if err != nil {
		return err
	}
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
func (r *record) RequestBytes() ([]byte, error) {
	w := NewWriter()
	if err := w.WriteName(r.name); err != nil {
		return nil, err
	}
	if err := w.WriteType(r.typ); err != nil {
		return nil, err
	}
	cls := r.class
	if r.unicastResponse {
		cls |= QU
	}
	if err := w.WriteClass(cls); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

// ResponseBytes returns only the binary representation of the all fields.
func (r *record) ResponseBytes() ([]byte, error) {
	commonBytes, err := r.RequestBytes()
	if err != nil {
		return nil, err
	}

	w := NewWriter()
	if err := w.WriteTTL(r.ttl); err != nil {
		return nil, err
	}
	if err := w.WriteData(r.data); err != nil {
		return nil, err
	}
	return append(commonBytes, w.Bytes()...), nil
}

// Bytes returns the binary representation.
func (r *record) Bytes() ([]byte, error) {
	return r.ResponseBytes()
}

// SetCompressionBytes sets the compression bytes.
func (r *record) SetCompressionBytes(b []byte) {
	r.cmpBytes = b
}

// CompressionBytes returns the compression bytes.
func (r *record) CompressionBytes() []byte {
	return r.cmpBytes
}

// EqualContent returns true if the record contents are equal. otherwise false.
func EqualContent(r1, r2 Record) bool {
	if r1.Type() != r2.Type() {
		return false
	}
	if !strings.EqualFold(r1.Name(), r2.Name()) {
		return false
	}
	if !strings.EqualFold(r1.Content(), r2.Content()) {
		return false
	}
	return true
}

// Equal returns true if this record is equal to  the specified resource record. otherwise false.
func (r *record) Equal(other Record) bool {
	return EqualContent(r, other)
}
