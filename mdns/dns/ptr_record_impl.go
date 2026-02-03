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

// ptrRecord represents a PTR record.
type ptrRecord struct {
	*record
	domainName string
}

// NewPTRRecord returns a new PTR record instance.
func NewPTRRecord() PTRRecord {
	return &ptrRecord{
		record:     newRecord(),
		domainName: "",
	}
}

// newPTRRecordWithResourceRecord returns a new PTR record instance.
func newPTRRecordWithResourceRecord(res *record) (PTRRecord, error) {
	ptr := &ptrRecord{
		record:     res,
		domainName: "",
	}
	return ptr, ptr.parseResourceRecord()
}

func (ptr *ptrRecord) parseResourceRecord() error {
	if len(ptr.data) == 0 {
		return nil
	}
	var err error
	reader := NewReaderWithBytes(ptr.data)
	reader.SetCompressionBytes(ptr.CompressionBytes())
	ptr.domainName, err = reader.ReadName()
	return err
}

// DomainName returns the resource domain name.
func (ptr *ptrRecord) DomainName() string {
	return ptr.domainName
}

// Content returns a string representation to the record data.
func (ptr *ptrRecord) Content() string {
	return ptr.DomainName()
}

// Equal returns true if this record is equal to  the specified resource record. otherwise false.
func (ptr *ptrRecord) Equal(other Record) bool {
	return EqualContent(ptr, other)
}
