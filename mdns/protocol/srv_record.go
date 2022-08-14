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
	"github.com/cybergarage/go-mdns/mdns/encoding"
)

// SRVRecord represents a SRV record.
// RFC 2782: A DNS RR for specifying the location of services (DNS SRV).
// https://www.rfc-editor.org/rfc/rfc2782
type SRVRecord struct {
	*resourceRecord
}

// NewSRVRecord returns a new SRV record innstance.
func NewSRVRecord(res *resourceRecord) *SRVRecord {
	return &SRVRecord{
		resourceRecord: res,
	}
}

// Priority returns the resource priority.
func (srv *SRVRecord) Priority() uint {
	if len(srv.data) < 2 {
		return 0
	}
	return encoding.BytesToInteger(srv.data[0:2])
}

// Weight returns the resource weight.
func (srv *SRVRecord) Weight() uint {
	if len(srv.data) < 4 {
		return 0
	}
	return encoding.BytesToInteger(srv.data[2:4])
}

// Port returns the resource port.
func (srv *SRVRecord) Port() uint {
	if len(srv.data) < 6 {
		return 0
	}
	return encoding.BytesToInteger(srv.data[4:6])
}

// Target returns the resource domain name of the target host.
func (srv *SRVRecord) Target() string {
	if len(srv.data) < 7 {
		return ""
	}
	return string(srv.data[7:])
}
