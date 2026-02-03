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
	"net"
)

// aRecord represents a A record.
type aRecord struct {
	*record
}

// NewARecord returns a new A record instance.
func NewARecord(res *record) ARecord {
	return &aRecord{
		record: newRecord(),
	}
}

// newARecordWithResourceRecord returns a new A record instance.
func newARecordWithResourceRecord(res *record) ARecord {
	return &aRecord{
		record: res,
	}
}

// Address returns the resource ip address.
func (a *aRecord) Address() net.IP {
	if len(a.data) < 4 {
		return nil
	}
	return net.IPv4(a.data[0], a.data[1], a.data[2], a.data[3])
}

// Content returns a string representation to the record data.
func (a *aRecord) Content() string {
	return a.Address().String()
}

// Equal returns true if this record is equal to  the specified resource record. otherwise false.
func (a *aRecord) Equal(other Record) bool {
	return EqualContent(a, other)
}
