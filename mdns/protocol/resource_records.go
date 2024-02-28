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

// ResourceRecords is a list of ResourceRecord.
type ResourceRecords []ResourceRecord

// LookupResourceRecordForName returns the resource record of the specified name.
func (records ResourceRecords) LookupResourceRecordForName(name string) (ResourceRecord, bool) {
	for _, record := range records {
		if record.Name() == name {
			return record, true
		}
	}
	return nil, false
}

// LookupResourceRecordForType returns the resource record of the specified type.
func (records ResourceRecords) LookupResourceRecordForType(t Type) (ResourceRecord, bool) {
	for _, record := range records {
		if record.Type() == t {
			return record, true
		}
	}
	return nil, false
}

// LookupResourceRecordsForName returns the resource records of the specified name.
func (records ResourceRecords) LookupResourceRecordsForName(name string) ResourceRecords {
	resRecords := ResourceRecords{}
	for _, record := range records {
		if record.Name() == name {
			resRecords = append(resRecords, record)
		}
	}
	return resRecords
}

// LookupResourceRecordForType returns the resource records of the specified type.
func (records ResourceRecords) LookupResourceRecordsForType(t Type) []ResourceRecord {
	resRecords := []ResourceRecord{}
	for _, record := range records {
		if record.Type() == t {
			resRecords = append(resRecords, record)
		}
	}
	return resRecords
}

// HasResourceRecord returns true if the resource record of the specified name is included in the list. otherwise false.
func (records ResourceRecords) HasResourceRecord(name string) bool {
	_, ok := records.LookupResourceRecordForName(name)
	return ok
}
