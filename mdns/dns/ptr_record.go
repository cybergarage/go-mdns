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

// PTRRecord represents a PTR record.
type PTRRecord struct {
	*record
	domainName string
}

// NewPTRRecord returns a new PTR record instance.
func NewPTRRecord() *PTRRecord {
	return &PTRRecord{
		record:     newResourceRecord(),
		domainName: "",
	}
}

// newPTRRecordWithResourceRecord returns a new PTR record instance.
func newPTRRecordWithResourceRecord(res *record) (*PTRRecord, error) {
	ptr := &PTRRecord{
		record:     res,
		domainName: "",
	}
	return ptr, ptr.parseResourceRecord()
}

func (ptr *PTRRecord) parseResourceRecord() error {
	if len(ptr.data) == 0 {
		return nil
	}
	var err error
	reader := NewReaderWithBytes(ptr.data)
	reader.SetCompressionReader(ptr.reader.CompressionReader())
	ptr.domainName, err = reader.ReadName()
	return err
}

// DomainName returns the resource domain name.
func (ptr *PTRRecord) DomainName() string {
	return ptr.domainName
}

// Content returns a string representation to the record data.
func (ptr *PTRRecord) Content() string {
	return ptr.DomainName()
}
