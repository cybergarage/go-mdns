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

// Records is a list of Record.
type Records []Record

// HasRecord returns true if the resource record of the specified name is included in the list. otherwise false.
func (records Records) HasRecord(name string) bool {
	_, ok := records.LookupRecordByName(name)
	return ok
}

// LookupRecordByName returns the resource record of the specified name.
func (records Records) LookupRecordByName(name string) (Record, bool) {
	lookupRecords := records.LookupRecordsByName(name)
	if len(lookupRecords) == 0 {
		return nil, false
	}
	return lookupRecords[0], true
}

// LookupRecordByNameRegex returns the resource record of the specified name.
func (records Records) LookupRecordByNameRegex(re *regexp.Regexp) (Record, bool) {
	lookupRecords := records.LookupRecordsByNameRegex(re)
	if len(lookupRecords) == 0 {
		return nil, false
	}
	return lookupRecords[0], true
}

// LookupRecordByNamePrefix returns the resource record of the specified name prefix.
func (records Records) LookupRecordByNamePrefix(prefix string) (Record, bool) {
	re := regexp.MustCompile("^" + regexp.QuoteMeta(prefix))
	return records.LookupRecordByNameRegex(re)
}

// LookupRecordByNameSuffix returns the resource record of the specified name suffix.
func (records Records) LookupRecordByNameSuffix(suffix string) (Record, bool) {
	re := regexp.MustCompile(regexp.QuoteMeta(suffix) + "$")
	return records.LookupRecordByNameRegex(re)
}

// LookupRecordByType returns the resource record of the specified type.
func (records Records) LookupRecordByType(t Type) (Record, bool) {
	lookupRecords := records.LookupRecordsByType(t)
	if len(lookupRecords) == 0 {
		return nil, false
	}
	return lookupRecords[0], true
}

// LookupRecordsByName returns the resource records of the specified name.
func (records Records) LookupRecordsByName(name string) Records {
	lookupRecords := Records{}
	for _, record := range records {
		if record.IsName(name) {
			lookupRecords = append(lookupRecords, record)
		}
	}
	return lookupRecords
}

// LookupRecordsByNameRegex returns the resource records that match the specified name regular expression.
func (records Records) LookupRecordsByNameRegex(re *regexp.Regexp) Records {
	lookupRecords := Records{}
	for _, record := range records {
		if re.MatchString(record.Name()) {
			lookupRecords = append(lookupRecords, record)
		}
	}
	return lookupRecords
}

// LookupRecordsByNamePrefix returns the resource records of the specified name prefix.
func (records Records) LookupRecordsByNamePrefix(prefix string) Records {
	re := regexp.MustCompile("^" + regexp.QuoteMeta(prefix))
	return records.LookupRecordsByNameRegex(re)
}

// LookupRecordsByNameSuffix returns the resource records of the specified name suffix.
func (records Records) LookupRecordsByNameSuffix(suffix string) Records {
	re := regexp.MustCompile(regexp.QuoteMeta(suffix) + "$")
	return records.LookupRecordsByNameRegex(re)
}

// LookupRecordsByType returns the resource records of the specified type.
func (records Records) LookupRecordsByType(t Type) []Record {
	resRecords := []Record{}
	for _, record := range records {
		if record.Type() == t {
			resRecords = append(resRecords, record)
		}
	}
	return resRecords
}

// LookupARecords returns the A records.
func (records Records) LookupARecords() []ARecord {
	resRecords := []ARecord{}
	for _, record := range records {
		if aRecord, ok := record.(ARecord); ok {
			resRecords = append(resRecords, aRecord)
		}
	}
	return resRecords
}

// LookupAAAARecords returns the AAAA records.
func (records Records) LookupAAAARecords() []AAAARecord {
	resRecords := []AAAARecord{}
	for _, record := range records {
		if aaaaRecord, ok := record.(AAAARecord); ok {
			resRecords = append(resRecords, aaaaRecord)
		}
	}
	return resRecords
}

// LookupPTRRecords returns the PTR records.
func (records Records) LookupPTRRecords() []PTRRecord {
	resRecords := []PTRRecord{}
	for _, record := range records {
		if ptrRecord, ok := record.(PTRRecord); ok {
			resRecords = append(resRecords, ptrRecord)
		}
	}
	return resRecords
}

// LookupSRVRecords returns the SRV records.
func (records Records) LookupSRVRecords() []SRVRecord {
	resRecords := []SRVRecord{}
	for _, record := range records {
		if srvRecord, ok := record.(SRVRecord); ok {
			resRecords = append(resRecords, srvRecord)
		}
	}
	return resRecords
}

// LookupTXTRecords returns the TXT records.
func (records Records) LookupTXTRecords() []TXTRecord {
	resRecords := []TXTRecord{}
	for _, record := range records {
		if txtRecord, ok := record.(TXTRecord); ok {
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
