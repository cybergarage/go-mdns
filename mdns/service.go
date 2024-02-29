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
	"fmt"
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
	protocol.Attributes
}

// NewService returns a new service instance.
func NewService(name, domain string, port uint) *Service {
	return &Service{
		Name:       name,
		Domain:     domain,
		Host:       "",
		AddrV4:     nil,
		AddrV6:     nil,
		Port:       port,
		Attributes: protocol.Attributes{},
	}
}

// NewServiceWithMessage returns a new service instance.
func NewServiceWithMessage(msg *Message) (*Service, error) {
	srv := NewService("", "", 0)
	srv.Update(msg)
	return srv, nil
}

// Update updates the service data by the specified message.
func (srv *Service) Update(msg *Message) {
	for _, record := range msg.ResourceRecords() {
		switch rr := record.(type) {
		case *protocol.PTRRecord:
			srv.Name = rr.DomainName()
		case *protocol.SRVRecord:
			host := rr.Target()
			if 0 < len(host) {
				srv.Host = host
			}
			port := rr.Port()
			if 0 < port {
				srv.Port = port
			}
		case *protocol.TXTRecord:
			srv.Attributes = append(srv.Attributes, rr.Attributes()...)
		case *protocol.ARecord:
			ip := rr.Address()
			if ip != nil {
				srv.AddrV4 = ip
			}
		case *protocol.AAAARecord:
			ip := rr.Address()
			if ip != nil {
				srv.AddrV6 = ip
			}
		}
	}
}

// Equal returns true if the header is same as the specified header, otherwise false.
func (srv *Service) Equal(other *Service) bool {
	if other == nil {
		return false
	}
	if srv.Name != other.Name {
		return false
	}
	if srv.Host != other.Host {
		return false
	}
	if srv.Domain != other.Domain {
		return false
	}
	return true
}

// String returns the string representation.
func (srv *Service) String() string {
	return fmt.Sprintf(
		"%s (%s:%d, %s:%d)",
		strings.Join([]string{srv.Name, srv.Host, srv.Domain}, nameSep),
		srv.AddrV4,
		srv.Port,
		srv.AddrV6,
		srv.Port,
	)
}
