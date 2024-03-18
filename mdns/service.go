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
	host   string
	addrs  []net.IP
	port   uint
	dns.Attributes
}

// NewService returns a new service instance.
func NewService(name, domain string, port uint) *Service {
	return &Service{
		Message:    nil,
		name:       name,
		domain:     domain,
		host:       "",
		addrs:      []net.IP{},
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

// Host returns the service host.
func (srv *Service) Host() string {
	return srv.host
}

// Port returns the service port.
func (srv *Service) Port() int {
	return int(srv.port)
}

// Addresses returns the service addresses.
func (srv *Service) Addresses() []net.IP {
	return srv.addrs
}

// parseMessage updates the service data by the specified message.
func (srv *Service) parseMessage(msg *Message) error {
	srv.Message = msg

	for _, record := range msg.ResourceRecords() {
		switch rr := record.(type) {
		case *dns.SRVRecord:
			host := rr.Target()
			if 0 < len(host) {
				srv.host = host
			}
			port := rr.Port()
			if 0 < port {
				srv.port = port
			}
		case *dns.TXTRecord:
			srv.name = rr.Name()
			attrs, err := rr.Attributes()
			if err == nil {
				srv.Attributes = append(srv.Attributes, attrs...)
			}
		case *dns.ARecord:
			ip := rr.Address()
			if ip != nil {
				srv.addrs = append(srv.addrs, ip)
			}
		case *dns.AAAARecord:
			ip := rr.Address()
			if ip != nil {
				srv.addrs = append(srv.addrs, ip)
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

	if len(srv.addrs) != len(other.addrs) {
		return false
	}

	equalAddrCount := 0
	for n, addr := range srv.addrs {
		if addr.Equal(other.addrs[n]) {
			equalAddrCount++
		}
	}
	if equalAddrCount != len(srv.addrs) {
		return false
	}

	if srv.name != other.name {
		return false
	}
	if srv.host != other.host {
		return false
	}
	if srv.port != other.port {
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
		"%s (%s:%d)",
		strings.Join([]string{srv.name, srv.host, srv.domain}, nameSep),
		srv.host,
		srv.port,
	)
}
