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

	"github.com/cybergarage/go-mdns/mdns/dns"
)

// Service represents a SRV record.
type Service struct {
	*Message
	name   string
	domain string
	Host   string
	AddrV4 net.IP
	AddrV6 net.IP
	port   uint
	dns.Attributes
}

// NewService returns a new service instance.
func NewService(name, domain string, port uint) *Service {
	return &Service{
		Message:    nil,
		name:       name,
		domain:     domain,
		Host:       "",
		AddrV4:     nil,
		AddrV6:     nil,
		port:       port,
		Attributes: dns.Attributes{},
	}
}

// NewServiceWithMessage returns a new service instance.
func NewServiceWithMessage(msg *Message) (*Service, error) {
	srv := NewService("", "", 0)
	err := srv.parseMessage(msg)
	if err != nil {
		return nil, err
	}
	return srv, nil
}

// Name returns the service name.
func (srv *Service) Name() string {
	return srv.name
}

// Domain returns the service domain.
func (srv *Service) Domain() string {
	return srv.domain
}

// Port returns the service port.
func (srv *Service) Port() int {
	return int(srv.port)
}

// parseMessage updates the service data by the specified message.
func (srv *Service) parseMessage(msg *Message) error {
	srv.Message = msg

	for _, record := range msg.ResourceRecords() {
		switch rr := record.(type) {
		case *dns.PTRRecord:
			srv.name = rr.DomainName()
		case *dns.SRVRecord:
			host := rr.Target()
			if 0 < len(host) {
				srv.Host = host
			}
			port := rr.Port()
			if 0 < port {
				srv.port = port
			}
		case *dns.TXTRecord:
			attrs, err := rr.Attributes()
			if err == nil {
				srv.Attributes = append(srv.Attributes, attrs...)
			}
		case *dns.ARecord:
			ip := rr.Address()
			if ip != nil {
				srv.AddrV4 = ip
			}
		case *dns.AAAARecord:
			ip := rr.Address()
			if ip != nil {
				srv.AddrV6 = ip
			}
		}
	}

	return nil
}

// Equal returns true if the header is same as the specified header, otherwise false.
func (srv *Service) Equal(other *Service) bool {
	if other == nil {
		return false
	}
	if srv.name != other.name {
		return false
	}
	if srv.Host != other.Host {
		return false
	}
	if srv.domain != other.domain {
		return false
	}
	return true
}

// String returns the string representation.
func (srv *Service) String() string {
	return fmt.Sprintf(
		"%s (%s:%d, %s:%d)",
		strings.Join([]string{srv.name, srv.Host, srv.domain}, nameSep),
		srv.AddrV4,
		srv.port,
		srv.AddrV6,
		srv.port,
	)
}
