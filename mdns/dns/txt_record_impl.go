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
	"strings"
)

// txtRecord represents a TXT record.
type txtRecord struct {
	*record
	strs []string
}

// NewTXTRecord returns a new TXT record instance.
func NewTXTRecord() TXTRecord {
	return &txtRecord{
		record: newRecord(),
		strs:   []string{},
	}
}

// newTXTRecordWithResourceRecord returns a new TXT record instance.
func newTXTRecordWithResourceRecord(res *record) (TXTRecord, error) {
	txt := &txtRecord{
		record: res,
		strs:   []string{},
	}
	return txt, txt.parseResourceRecord()
}

func (txt *txtRecord) parseResourceRecord() error {
	var err error
	if len(txt.data) == 0 {
		return nil
	}
	reader := NewReaderWithBytes(txt.data)
	txt.strs, err = reader.ReadStrings()
	return err
}

// Strings returns the resource attribute strings.
func (txt *txtRecord) Strings() []string {
	return txt.strs
}

// Attributes returns the resource attribute strings.
func (txt *txtRecord) Attributes() (Attributes, error) {
	return NewAttributesFromStrings(txt.strs)
}

// LookupAttribute returns the attribute with the specified name.
func (txt *txtRecord) LookupAttribute(name string) (Attribute, bool) {
	attrs, err := txt.Attributes()
	if err != nil {
		return nil, false
	}
	return attrs.LookupAttribute(name)
}

// HasAttribute returns true if this instance has the specified attribute.
func (txt *txtRecord) HasAttribute(name string) bool {
	attrs, err := txt.Attributes()
	if err != nil {
		return false
	}
	return attrs.HasAttribute(name)
}

// Content returns a string representation to the record data.
func (txt *txtRecord) Content() string {
	return strings.Join(txt.strs, " ")
}

// Equal returns true if this record is equal to  the specified resource record. otherwise false.
func (txt *txtRecord) Equal(other Record) bool {
	return EqualContent(txt, other)
}
