// Copyright (C) 2022 Satoshi Konno All rights reserved.
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

import "bytes"

// TXTRecord represents a TXT record.
type TXTRecord struct {
	*Record
}

// NewTXTRecord returns a new TXT record innstance.
func NewTXTRecord() *TXTRecord {
	return &TXTRecord{
		Record: newResourceRecord(),
	}
}

// newTXTRecordWithResourceRecord returns a new TXT record innstance.
func newTXTRecordWithResourceRecord(res *Record) *TXTRecord {
	return &TXTRecord{
		Record: res,
	}
}

// Attributes returns the resource attribute strings.
func (txt *TXTRecord) Attributes() []string {
	attrs, err := parseTxt(bytes.NewReader(txt.data))
	if err != nil {
		return []string{}
	}
	return attrs
}
