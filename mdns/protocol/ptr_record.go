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

import "bytes"

// PTRRecord represents a PTR record.
type PTRRecord struct {
	*Record
	domainName string
}

// NewPTRRecord returns a new PTR record innstance.
func NewPTRRecord() *PTRRecord {
	return &PTRRecord{
		Record:     newResourceRecord(),
		domainName: "",
	}
}

// newPTRRecordWithResourceRecord returns a new PTR record innstance.
func newPTRRecordWithResourceRecord(res *Record) (*PTRRecord, error) {
	ptr := &PTRRecord{
		Record:     res,
		domainName: "",
	}
	return ptr, ptr.parseResourceRecord()
}

func (ptr *PTRRecord) parseResourceRecord() error {
	name, err := NewReaderWithReader(bytes.NewReader(ptr.data)).ReadNameWith(ptr.reader.CompressionReader())
	if err != nil {
		return err
	}
	ptr.domainName = name
	return nil
}

// DomainName returns the resource domain name.
func (ptr *PTRRecord) DomainName() string {
	return ptr.domainName
}
