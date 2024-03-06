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

// ResourceRecords is a list of ResourceRecord.
type ResourceRecords []ResourceRecord

// HasResourceRecord returns true if the resource record of the specified name is included in the list. otherwise false.
func (records ResourceRecords) HasResourceRecord(name string) bool {
	_, ok := records.LookupResourceRecordForName(name)
	return ok
}

// LookupResourceRecordForName returns the resource record of the specified name.
func (records ResourceRecords) LookupResourceRecordForName(name string) (ResourceRecord, bool) {
	lookupRecords := records.LookupResourceRecordsForName(name)
	if len(lookupRecords) == 0 {
		return nil, false
	}
	return lookupRecords[0], true
}

// LookupResourceRecordForType returns the resource record of the specified type.
func (records ResourceRecords) LookupResourceRecordForType(t Type) (ResourceRecord, bool) {
	lookupRecords := records.LookupResourceRecordsForType(t)
	if len(lookupRecords) == 0 {
		return nil, false
	}
	return lookupRecords[0], true
}

// LookupResourceRecordsForName returns the resource records of the specified name.
func (records ResourceRecords) LookupResourceRecordsForName(name string) ResourceRecords {
	resRecords := ResourceRecords{}
	for _, record := range records {
		if record.IsName(name) {
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

// LookupARecords returns the A records.
func (records ResourceRecords) LookupARecords() []*ARecord {
	resRecords := []*ARecord{}
	for _, record := range records {
		if aRecord, ok := record.(*ARecord); ok {
			resRecords = append(resRecords, aRecord)
		}
	}
	return resRecords
}

// LookupAAAARecords returns the AAAA records.
func (records ResourceRecords) LookupAAAARecords() []*AAAARecord {
	resRecords := []*AAAARecord{}
	for _, record := range records {
		if aaaaRecord, ok := record.(*AAAARecord); ok {
			resRecords = append(resRecords, aaaaRecord)
		}
	}
	return resRecords
}

// LookupPTRRecords returns the PTR records.
func (records ResourceRecords) LookupPTRRecords() []*PTRRecord {
	resRecords := []*PTRRecord{}
	for _, record := range records {
		if ptrRecord, ok := record.(*PTRRecord); ok {
			resRecords = append(resRecords, ptrRecord)
		}
	}
	return resRecords
}

// LookupSRVRecords returns the SRV records.
func (records ResourceRecords) LookupSRVRecords() []*SRVRecord {
	resRecords := []*SRVRecord{}
	for _, record := range records {
		if srvRecord, ok := record.(*SRVRecord); ok {
			resRecords = append(resRecords, srvRecord)
		}
	}
	return resRecords
}

// LookupTXTRecords returns the TXT records.
func (records ResourceRecords) LookupTXTRecords() []*TXTRecord {
	resRecords := []*TXTRecord{}
	for _, record := range records {
		if txtRecord, ok := record.(*TXTRecord); ok {
			resRecords = append(resRecords, txtRecord)
		}
	}
	return resRecords
}
