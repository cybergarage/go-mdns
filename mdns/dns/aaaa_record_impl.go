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
	"net"
	"strings"
)

// aaaaRecord represents a AAAA record.
type aaaaRecord struct {
	*record
}

// NewAAAARecord returns a new AAAA record instance.
func NewAAAARecord() AAAARecord {
	return &aaaaRecord{
		record: newResourceRecord(),
	}
}

// newAAAARecordWithResourceRecord returns a new AAAA record instance.
func newAAAARecordWithResourceRecord(res *record) AAAARecord {
	return &aaaaRecord{
		record: res,
	}
}

// Address returns the resource ip address.
func (a *aaaaRecord) Address() net.IP {
	if len(a.data) != 16 {
		return nil
	}
	var ipstr strings.Builder
	for n, b := range a.data {
		if (n != 0) && ((n % 2) == 0) {
			ipstr.WriteString(":")
		}
		ipstr.WriteString(fmt.Sprintf("%02x", b))
	}
	return net.ParseIP(ipstr.String())
}

// Content returns a string representation to the record data.
func (a *aaaaRecord) Content() string {
	return a.Address().String()
}

// Equal returns true if this record is equal to  the specified resource record. otherwise false.
func (a *aaaaRecord) Equal(other Record) bool {
	return EqualContent(a, other)
}
