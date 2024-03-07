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

import "fmt"

// SRVRecord represents a SRV record.
// RFC 2782: A DNS RR for specifying the location of services (DNS SRV).
// https://www.rfc-editor.org/rfc/rfc2782
type SRVRecord struct {
	*record
	service  string
	proto    string
	name     string
	priority uint16
	weight   uint16
	port     uint16
	target   string
}

// NewSRVRecord returns a new SRV record instance.
func NewSRVRecord() *SRVRecord {
	return &SRVRecord{
		record:   newResourceRecord(),
		service:  "",
		proto:    "",
		name:     "",
		priority: 0,
		weight:   0,
		port:     0,
		target:   "",
	}
}

// newSRVRecordWithResourceRecord returns a new SRV record instance.
func newSRVRecordWithResourceRecord(res *record) (*SRVRecord, error) {
	srv := &SRVRecord{
		record:   res,
		service:  "",
		proto:    "",
		name:     "",
		priority: 0,
		weight:   0,
		port:     0,
		target:   "",
	}
	return srv, srv.parseResourceRecord()
}

func (srv *SRVRecord) parseResourceRecord() error {
	if len(srv.data) == 0 {
		return nil
	}

	var err error

	reader := NewReaderWithBytes(srv.data)

	srv.priority, err = reader.ReadUint16()
	if err != nil {
		return err
	}

	srv.weight, err = reader.ReadUint16()
	if err != nil {
		return err
	}

	srv.port, err = reader.ReadUint16()
	if err != nil {
		return err
	}

	srv.target, err = reader.ReadString()
	if err != nil {
		return err
	}

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

// Priority returns the resource priority.
func (srv *SRVRecord) Priority() uint {
	return uint(srv.priority)
}

// Weight returns the resource weight.
func (srv *SRVRecord) Weight() uint {
	return uint(srv.weight)
}

// Port returns the resource port.
func (srv *SRVRecord) Port() uint {
	return uint(srv.port)
}

// Target returns the resource target.
func (srv *SRVRecord) Target() string {
	return srv.target
}

// Content returns a string representation to the record data.
func (srv *SRVRecord) Content() string {
	return fmt.Sprintf("%d %d %d %s", srv.priority, srv.weight, srv.port, srv.target)
}
