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
)

// TXTRecord represents a TXT record.
type TXTRecord struct {
	*Record
	attrs Attributes
}

// NewTXTRecord returns a new TXT record innstance.
func NewTXTRecord() *TXTRecord {
	return &TXTRecord{
		Record: newResourceRecord(),
		attrs:  Attributes{},
	}
}

// newTXTRecordWithResourceRecord returns a new TXT record innstance.
func newTXTRecordWithResourceRecord(res *Record) (*TXTRecord, error) {
	txt := &TXTRecord{
		Record: res,
		attrs:  Attributes{},
	}
	return txt, txt.parseResourceRecord()
}

func (txt *TXTRecord) parseResourceRecord() error {
	var err error
	if len(txt.data) == 0 {
		return nil
	}
	reader := NewReaderWithReader(bytes.NewReader(txt.data))
	txt.attrs, err = reader.ReadAttributes()
	return err
}

// Attributes returns the resource attribute strings.
func (txt *TXTRecord) Attributes() Attributes {
	return txt.attrs
}

// GetAttribute returns the attribute with the specified name.
func (txt *TXTRecord) GetAttribute(name string) *Attribute {
	return txt.attrs.GetAttribute(name)
}

// HasAttribute returns true if this instance has the specified attribute.
func (txt *TXTRecord) HasAttribute(name string) bool {
	return txt.attrs.HasAttribute(name)
}
