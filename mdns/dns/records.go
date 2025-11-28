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
	"strconv"
	"strings"
)

// Records is a list of Record.
type Records []Record

// HasRecord returns true if the resource record of the specified name is included in the list. otherwise false.
func (records Records) HasRecord(name string) bool {
	_, ok := records.LookupRecordForName(name)
	return ok
}

// LookupRecordForName returns the resource record of the specified name.
func (records Records) LookupRecordForName(name string) (Record, bool) {
	lookupRecords := records.LookupRecordsForName(name)
	if len(lookupRecords) == 0 {
		return nil, false
	}
	return lookupRecords[0], true
}

// LookupRecordForNamePrefix returns the resource record of the specified name prefix.
func (records Records) LookupRecordForNamePrefix(prefix string) (Record, bool) {
	lookupRecords := records.LookupRecordsForNamePrefix(prefix)
	if len(lookupRecords) == 0 {
		return nil, false
	}
	return lookupRecords[0], true
}

// LookupRecordForNameSuffix returns the resource record of the specified name suffix.
func (records Records) LookupRecordForNameSuffix(suffix string) (Record, bool) {
	lookupRecords := records.LookupRecordsForNameSuffix(suffix)
	if len(lookupRecords) == 0 {
		return nil, false
	}
	return lookupRecords[0], true
}

// LookupRecordForType returns the resource record of the specified type.
func (records Records) LookupRecordForType(t Type) (Record, bool) {
	lookupRecords := records.LookupRecordsForType(t)
	if len(lookupRecords) == 0 {
		return nil, false
	}
	return lookupRecords[0], true
}

// LookupRecordsForNamePrefix returns the resource records of the specified name prefix.
func (records Records) LookupRecordsForNamePrefix(prefix string) Records {
	lookupRecords := Records{}
	for _, record := range records {
		if record.HasNamePrefix(prefix) {
			lookupRecords = append(lookupRecords, record)
		}
	}
	return lookupRecords
}

// LookupRecordsForNameSuffix returns the resource records of the specified name suffix.
func (records Records) LookupRecordsForNameSuffix(suffix string) Records {
	lookupRecords := Records{}
	for _, record := range records {
		if record.HasNameSuffix(suffix) {
			lookupRecords = append(lookupRecords, record)
		}
	}
	return lookupRecords
}

// LookupRecordsForName returns the resource records of the specified name.
func (records Records) LookupRecordsForName(name string) Records {
	lookupRecords := Records{}
	for _, record := range records {
		if record.IsName(name) {
			lookupRecords = append(lookupRecords, record)
		}
	}
	return lookupRecords
}

// LookupRecordForType returns the resource records of the specified type.
func (records Records) LookupRecordsForType(t Type) []Record {
	resRecords := []Record{}
	for _, record := range records {
		if record.Type() == t {
			resRecords = append(resRecords, record)
		}
	}
	return resRecords
}

// LookupARecords returns the A records.
func (records Records) LookupARecords() []*ARecord {
	resRecords := []*ARecord{}
	for _, record := range records {
		if aRecord, ok := record.(*ARecord); ok {
			resRecords = append(resRecords, aRecord)
		}
	}
	return resRecords
}

// LookupAAAARecords returns the AAAA records.
func (records Records) LookupAAAARecords() []*AAAARecord {
	resRecords := []*AAAARecord{}
	for _, record := range records {
		if aaaaRecord, ok := record.(*AAAARecord); ok {
			resRecords = append(resRecords, aaaaRecord)
		}
	}
	return resRecords
}

// LookupPTRRecords returns the PTR records.
func (records Records) LookupPTRRecords() []*PTRRecord {
	resRecords := []*PTRRecord{}
	for _, record := range records {
		if ptrRecord, ok := record.(*PTRRecord); ok {
			resRecords = append(resRecords, ptrRecord)
		}
	}
	return resRecords
}

// LookupSRVRecords returns the SRV records.
func (records Records) LookupSRVRecords() []*SRVRecord {
	resRecords := []*SRVRecord{}
	for _, record := range records {
		if srvRecord, ok := record.(*SRVRecord); ok {
			resRecords = append(resRecords, srvRecord)
		}
	}
	return resRecords
}

// LookupTXTRecords returns the TXT records.
func (records Records) LookupTXTRecords() []*TXTRecord {
	resRecords := []*TXTRecord{}
	for _, record := range records {
		if txtRecord, ok := record.(*TXTRecord); ok {
			resRecords = append(resRecords, txtRecord)
		}
	}
	return resRecords
}

// String returns the string representation.
func (records Records) String() string {
	type record []string

	lines := []record{}
	for _, r := range records {
		record := record{
			r.Name(),
			r.Type().String(),
			r.Content(),
		}
		lines = append(lines, record)
	}

	maxRecordLen := []int{0, 0, 0}
	for _, r := range lines {
		for n, s := range r {
			if maxRecordLen[n] < len(s) {
				maxRecordLen[n] = len(s)
			}
		}
	}

	var str strings.Builder
	for n, r := range lines {
		for n, s := range r {
			sfmt := "%-" + strconv.Itoa(maxRecordLen[n]) + "s"
			str.WriteString(fmt.Sprintf(sfmt, s))
			if n < len(r)-1 {
				str.WriteString(" ")
			}
		}
		if n < len(lines)-1 {
			str.WriteString("\n")
		}
	}

	return str.String()
}
