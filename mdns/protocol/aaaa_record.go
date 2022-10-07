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

import (
	"fmt"
	"net"
)

// AAAARecord represents a AAAA record.
type AAAARecord struct {
	*Record
}

// NewAAAARecord returns a new AAAA record innstance.
func NewAAAARecord() *AAAARecord {
	return &AAAARecord{
		Record: newResourceRecord(),
	}
}

// newAAAARecordWithResourceRecord returns a new AAAA record innstance.
func newAAAARecordWithResourceRecord(res *Record) *AAAARecord {
	return &AAAARecord{
		Record: res,
	}
}

// IP returns the resource ip address.
func (a *AAAARecord) IP() net.IP {
	if len(a.data) != 16 {
		return nil
	}
	ipstr := ""
	for n, b := range a.data {
		if (n != 0) && ((n % 2) == 0) {
			ipstr += ":"
		}
		ipstr += fmt.Sprintf("%02x", b)
	}
	return net.ParseIP(ipstr)
}
