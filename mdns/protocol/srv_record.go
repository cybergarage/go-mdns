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
	*Record
	service  string
	proto    string
	name     string
	ttl      uint
	class    Class
	priority uint
	weight   uint
	port     uint
	target   string
}

// NewSRVRecord returns a new SRV record innstance.
func NewSRVRecord() *SRVRecord {
	return &SRVRecord{
		Record:   newResourceRecord(),
		service:  "",
		proto:    "",
		name:     "",
		ttl:      0,
		class:    0,
		priority: 0,
		weight:   0,
		port:     0,
		target:   "",
	}
}

// newSRVRecordWithResourceRecord returns a new SRV record innstance.
func newSRVRecordWithResourceRecord(res *Record) (*SRVRecord, error) {
	srv := &SRVRecord{
		Record:   res,
		service:  "",
		proto:    "",
		name:     "",
		ttl:      0,
		class:    0,
		priority: 0,
		weight:   0,
		port:     0,
		target:   "",
	}
	return srv, srv.parseResourceRecord()
}

func (srv *SRVRecord) parseResourceRecord() error {
	return nil
}

// Service returns the service name.
func (srv *SRVRecord) Service() string {
	return srv.service
}

// Proto returns the protocol name.
func (srv *SRVRecord) Proto() string {
	return srv.proto
}

// Name returns the resource name.
func (srv *SRVRecord) Name() string {
	return srv.name
}

// TTL returns the resource TTL.
func (srv *SRVRecord) TTL() uint {
	return srv.ttl
}

// Class returns the resource class.
func (srv *SRVRecord) Class() Class {
	return srv.class
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
	return srv.weight
}

// Port returns the resource port.
func (srv *SRVRecord) Port() uint {
	return srv.port
}

// Target returns the resource target.
func (srv *SRVRecord) Target() string {
	return srv.target
}
