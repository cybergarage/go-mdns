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

package mdns

import (
	"net"
	"strings"

	"github.com/cybergarage/go-mdns/mdns/protocol"
)

// Service represents a SRV record.
type Service struct {
	Name   string
	Domain string
	Host   string
	AddrV4 net.IP
	AddrV6 net.IP
	Port   uint
}

// NewService returns a new service instance.
func NewService(name, domain string, port uint) *Service {
	return &Service{
		Name:   name,
		Domain: domain,
		Host:   "",
		AddrV4: nil,
		AddrV6: nil,
		Port:   port,
	}
}

// NewServiceWithMessage returns a new service instance.
func NewServiceWithMessage(msg *Message) (*Service, error) {
	srv := NewService("", "", 0)
	for _, res := range append(msg.Answers, msg.Additions...) {
		switch rr := res.(type) {
		case *protocol.PTRRecord:
			srv.Name = rr.DomainName()
		case *protocol.SRVRecord:
			srv.Host = rr.Target()
			srv.Port = rr.Port()
		case *protocol.ARecord:
			srv.AddrV4 = rr.IP()
		case *protocol.AAAARecord:
			srv.AddrV6 = rr.IP()
		}
	}
	return srv, nil
}

// String returns the string representation.
func (srv *Service) String() string {
	return strings.Join([]string{srv.Name, srv.Domain}, nameSep)
}
