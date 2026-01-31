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
	"regexp"
	"strconv"
	"strings"
)

// RecordSet represents a set of resource records.
type RecordSet []Record

// HasRecord returns true if the resource record of the specified name is included in the list. otherwise false.
func (records RecordSet) HasRecord(name string) bool {
	_, ok := records.LookupRecordByName(name)
	return ok
}

// LookupRecordByName returns the resource record of the specified name.
func (records RecordSet) LookupRecordByName(name string) (Record, bool) {
	lookupRecords := records.LookupRecordSetByName(name)
	if len(lookupRecords) == 0 {
		return nil, false
	}
	return lookupRecords[0], true
}

// LookupRecordByNameRegex returns the resource record of the specified name.
func (records RecordSet) LookupRecordByNameRegex(re *regexp.Regexp) (Record, bool) {
	lookupRecords := records.LookupRecordSetByNameRegex(re)
	if len(lookupRecords) == 0 {
		return nil, false
	}
	return lookupRecords[0], true
}

// LookupRecordByNamePrefix returns the resource record of the specified name prefix.
func (records RecordSet) LookupRecordByNamePrefix(prefix string) (Record, bool) {
	re := regexp.MustCompile("^" + regexp.QuoteMeta(prefix))
	return records.LookupRecordByNameRegex(re)
}

// LookupRecordByNameSuffix returns the resource record of the specified name suffix.
func (records RecordSet) LookupRecordByNameSuffix(suffix string) (Record, bool) {
	re := regexp.MustCompile(regexp.QuoteMeta(suffix) + "$")
	return records.LookupRecordByNameRegex(re)
}

// LookupRecordByType returns the resource record of the specified type.
func (records RecordSet) LookupRecordByType(t Type) (Record, bool) {
	lookupRecords := records.LookupRecordSetByType(t)
	if len(lookupRecords) == 0 {
		return nil, false
	}
	return lookupRecords[0], true
}

// LookupRecordSetByName returns the resource records of the specified name.
func (records RecordSet) LookupRecordSetByName(name string) RecordSet {
	lookupRecords := RecordSet{}
	for _, record := range records {
		if record.IsName(name) {
			lookupRecords = append(lookupRecords, record)
		}
	}
	return lookupRecords
}

// LookupRecordSetByNameRegex returns the resource records of the specified name regex.
func (records RecordSet) LookupRecordSetByNameRegex(re *regexp.Regexp) RecordSet {
	lookupRecords := RecordSet{}
	for _, record := range records {
		if re.MatchString(record.Name()) {
			lookupRecords = append(lookupRecords, record)
		}
	}
	return lookupRecords
}

// LookupRecordSetByNamePrefix returns the resource records of the specified name prefix.
func (records RecordSet) LookupRecordSetByNamePrefix(prefix string) RecordSet {
	re := regexp.MustCompile("^" + regexp.QuoteMeta(prefix))
	return records.LookupRecordSetByNameRegex(re)
}

// LookupRecordSetByNameSuffix returns the resource records of the specified name suffix.
func (records RecordSet) LookupRecordSetByNameSuffix(suffix string) RecordSet {
	re := regexp.MustCompile(regexp.QuoteMeta(suffix) + "$")
	return records.LookupRecordSetByNameRegex(re)
}

// LookupRecordSetByType returns the resource records of the specified type.
func (records RecordSet) LookupRecordSetByType(t Type) []Record {
	resRecords := []Record{}
	for _, record := range records {
		if record.Type() == t {
			resRecords = append(resRecords, record)
		}
	}
	return resRecords
}

// LookupARecordSet returns the A records.
func (records RecordSet) LookupARecordSet() []ARecord {
	resRecords := []ARecord{}
	for _, record := range records {
		if aRecord, ok := record.(ARecord); ok {
			resRecords = append(resRecords, aRecord)
		}
	}
	return resRecords
}

// LookupAAAARecordSet returns the AAAA records.
func (records RecordSet) LookupAAAARecordSet() []AAAARecord {
	resRecords := []AAAARecord{}
	for _, record := range records {
		if aaaaRecord, ok := record.(AAAARecord); ok {
			resRecords = append(resRecords, aaaaRecord)
		}
	}
	return resRecords
}

// LookupPTRRecordSet returns the PTR records.
func (records RecordSet) LookupPTRRecordSet() []PTRRecord {
	resRecords := []PTRRecord{}
	for _, record := range records {
		if ptrRecord, ok := record.(PTRRecord); ok {
			resRecords = append(resRecords, ptrRecord)
		}
	}
	return resRecords
}

// LookupSRVRecordSet returns the SRV records.
func (records RecordSet) LookupSRVRecordSet() []SRVRecord {
	resRecords := []SRVRecord{}
	for _, record := range records {
		if srvRecord, ok := record.(SRVRecord); ok {
			resRecords = append(resRecords, srvRecord)
		}
	}
	return resRecords
}

// LookupTXTRecordSet returns the TXT records.
func (records RecordSet) LookupTXTRecordSet() []TXTRecord {
	resRecords := []TXTRecord{}
	for _, record := range records {
		if txtRecord, ok := record.(TXTRecord); ok {
			resRecords = append(resRecords, txtRecord)
		}
	}
	return resRecords
}

// Equal returns true if the record sets are equal. otherwise false.
func (records RecordSet) Equal(other RecordSet) bool {
	if len(records) != len(other) {
		return false
	}
	for _, record := range records {
		hasRecord := false
		name := record.Name()
		otherRecordSet := other.LookupRecordSetByName(name)
		if len(otherRecordSet) == 0 {
			name = record.Name()
			return false
		}
		for _, otherRecord := range otherRecordSet {
			if record.Equal(otherRecord) {
				hasRecord = true
				break
			}
		}
		if !hasRecord {
			return false
		}
	}
	return true
}

// String returns the string representation.
func (records RecordSet) String() string {
	type record []string

	lines := make([]record, 0, len(records))
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
