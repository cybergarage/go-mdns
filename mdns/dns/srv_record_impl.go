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
	"strings"
)

// srvRecord represents a SRV record.
// RFC 2782: A DNS RR for specifying the location of services (DNS SRV).
// https://www.rfc-editor.org/rfc/rfc2782
type srvRecord struct {
	*record
	service  string
	proto    string
	priority uint16
	weight   uint16
	port     uint16
	target   string
}

// NewSRVRecord returns a new SRV record instance.
func NewSRVRecord() SRVRecord {
	return &srvRecord{
		record:   newResourceRecord(),
		service:  "",
		proto:    "",
		priority: 0,
		weight:   0,
		port:     0,
		target:   "",
	}
}

// newSRVRecordWithResourceRecord returns a new SRV record instance.
func newSRVRecordWithResourceRecord(res *record) (SRVRecord, error) {
	srv := &srvRecord{
		record:   res,
		service:  "",
		proto:    "",
		priority: 0,
		weight:   0,
		port:     0,
		target:   "",
	}
	if err := srv.parseName(); err != nil {
		return nil, err
	}
	if err := srv.parseResourceRecord(); err != nil {
		return nil, err
	}
	return srv, nil
}

func (srv *srvRecord) parseName() error {
	name := srv.Name()
	if len(name) == 0 {
		return nil
	}
	names := strings.SplitN(name, ".", 3)
	if len(names) < 3 {
		return fmt.Errorf("%w SRV record name: %s", ErrInvalid, name)
	}
	srv.service = names[0]
	srv.proto = names[1]
	return nil
}

func (srv *srvRecord) parseResourceRecord() error {
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
func (srv *srvRecord) Service() string {
	return srv.service
}

// Proto returns the protocol name.
func (srv *srvRecord) Proto() string {
	return srv.proto
}

// Priority returns the resource priority.
func (srv *srvRecord) Priority() uint {
	return uint(srv.priority)
}

// Weight returns the resource weight.
func (srv *srvRecord) Weight() uint {
	return uint(srv.weight)
}

// Port returns the resource port.
func (srv *srvRecord) Port() uint {
	return uint(srv.port)
}

// Target returns the resource target.
func (srv *srvRecord) Target() string {
	return srv.target
}

// Content returns a string representation to the record data.
func (srv *srvRecord) Content() string {
	if len(srv.data) == 0 {
		return ""
	}
	return fmt.Sprintf("%d %d %d %s", srv.priority, srv.weight, srv.port, srv.target)
}

// Equal returns true if this record is equal to  the specified resource record. otherwise false.
func (srv *srvRecord) Equal(other Record) bool {
	return EqualContent(srv, other)
}
