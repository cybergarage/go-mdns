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
	"regexp"
	"strings"

	"github.com/cybergarage/go-mdns/mdns/dns"
)

// serviceImpl represents a SRV record.
type serviceImpl struct {
	Message
	name   string
	domain string
	host   string
	addrs  []net.IP
	port   int
	attrs  dns.Attributes
}

// ServiceOptions represents a service option.
type ServiceOptions func(*serviceImpl) error

// WithServiceName returns a service option with the specified name.
func WithServiceName(name string) ServiceOptions {
	return func(srv *serviceImpl) error {
		srv.name = name
		return nil
	}
}

// WithServiceDomain returns a service option with the specified domain.
func WithServiceDomain(domain string) ServiceOptions {
	return func(srv *serviceImpl) error {
		srv.domain = domain
		return nil
	}
}

// WithServiceHost returns a service option with the specified host.
func WithServiceHost(host string) ServiceOptions {
	return func(srv *serviceImpl) error {
		srv.host = host
		return nil
	}
}

// WithServicePort returns a service option with the specified port.
func WithServicePort(port int) ServiceOptions {
	return func(srv *serviceImpl) error {
		srv.port = port
		return nil
	}
}

// WithServiceMessage returns a service option with the specified message.
func WithServiceMessage(msg Message) ServiceOptions {
	return func(srv *serviceImpl) error {
		return srv.parseMessage(msg)
	}
}

// NewService returns a new service instance.
func NewService(opts ...ServiceOptions) (Service, error) {
	return newService(opts...)
}

func newService(opts ...ServiceOptions) (*serviceImpl, error) {
	srv := &serviceImpl{
		Message: nil,
		name:    "",
		domain:  "",
		host:    "",
		addrs:   []net.IP{},
		port:    0,
		attrs:   dns.Attributes{},
	}
	for _, opt := range opts {
		err := opt(srv)
		if err != nil {
			return nil, err
		}
	}
	return srv, nil
}

// Name returns the service name.
func (srv *serviceImpl) Name() string {
	return srv.name
}

// Domain returns the service domain.
func (srv *serviceImpl) Domain() string {
	return srv.domain
}

// Host returns the service host.
func (srv *serviceImpl) Host() string {
	return srv.host
}

// Port returns the service port.
func (srv *serviceImpl) Port() int {
	return int(srv.port)
}

// Addresses returns the service addresses.
func (srv *serviceImpl) Addresses() []net.IP {
	return srv.addrs
}

// parseMessage updates the service data by the specified message.
func (srv *serviceImpl) parseMessage(msg Message) error {
	srv.Message = msg
	for _, record := range msg.ResourceRecords() {
		err := srv.parseRecord(record)
		if err != nil {
			return err
		}
	}
	return nil
}

// ResourceRecords returns the service resource records.
func (srv *serviceImpl) ResourceRecords() ResourceRecords {
	return srv.Message.ResourceRecords()
}

// ResourceAttributes returns the service TXT attributes.
func (srv *serviceImpl) ResourceAttributes() dns.Attributes {
	return srv.attrs
}

// LookupResourceAttribute returns the attribute with the specified name.
func (srv *serviceImpl) LookupResourceAttribute(name string) (Attribute, bool) {
	return srv.attrs.LookupAttribute(name)
}

// LookupResourceByName returns the resource record of the specified name.
func (srv *serviceImpl) LookupResourceByName(name string) (ResourceRecord, bool) {
	return srv.Message.LookupResourceRecordByName(name)
}

// LookupResourceByNameRegex returns the resource record of the specified name regex.
func (srv *serviceImpl) LookupResourceByNameRegex(re *regexp.Regexp) (ResourceRecord, bool) {
	return srv.Message.LookupResourceRecordByNameRegex(re)
}

// LookupResourceByNamePrefix returns the resource record of the specified name prefix.
func (srv *serviceImpl) LookupResourceByNamePrefix(prefix string) (ResourceRecord, bool) {
	return srv.Message.LookupResourceRecordByNamePrefix(prefix)
}

// LookupResourceByNameSuffix returns the resource record of the specified name suffix.
func (srv *serviceImpl) LookupResourceByNameSuffix(suffix string) (ResourceRecord, bool) {
	return srv.Message.LookupResourceRecordByNameSuffix(suffix)
}

func (srv *serviceImpl) parseRecord(record dns.Record) error {
	parseNameDomain := func(fullname string) error {
		if len(fullname) == 0 {
			return nil
		}
		idx := strings.LastIndex(fullname, dns.LabelSeparator)
		if idx == -1 {
			return fmt.Errorf("invalid record name: %s", fullname)
		}
		parts := []string{fullname[:idx], fullname[idx+1:]}
		if len(parts) != 2 {
			return fmt.Errorf("invalid record name: %s", fullname)
		}
		if 0 < len(parts[0]) && len(srv.name) == 0 {
			srv.name = parts[0]
		}
		if 0 < len(parts[1]) && len(srv.domain) == 0 {
			srv.domain = parts[1]
		}
		return nil
	}

	// Handle address records (A/AAAA) via shared Address() method.
	if ar, ok := record.(interface{ Address() net.IP }); ok {
		ip := ar.Address()
		if ip != nil {
			srv.addrs = append(srv.addrs, ip)
		}
	}

	switch rr := record.(type) {
	case dns.SRVRecord:
		err := parseNameDomain(rr.Name())
		if err != nil {
			return err
		}
		host := rr.Target()
		if 0 < len(host) {
			srv.host = host
		}
		port := rr.Port()
		if 0 < port {
			srv.port = int(port)
		}
	case dns.TXTRecord:
		err := parseNameDomain(rr.Name())
		if err != nil {
			return err
		}
		attrs, err := rr.Attributes()
		if err == nil {
			srv.attrs = append(srv.attrs, attrs...)
		}
	}

	return nil
}

// Attributes returns the service TXT attributes.
func (srv *serviceImpl) Attributes() dns.Attributes {
	return srv.attrs
}

// LookupAttribute returns the attribute with the specified name.
func (srv *serviceImpl) LookupAttribute(name string) (Attribute, bool) {
	return srv.attrs.LookupAttribute(name)
}

// Equal returns true if the header is same as the specified header, otherwise false.
func (srv *serviceImpl) Equal(other Service) bool {
	if other == nil {
		return false
	}

	if len(srv.addrs) != len(other.Addresses()) {
		return false
	}

	equalAddrCount := 0
	for n, addr := range srv.addrs {
		if addr.Equal(other.Addresses()[n]) {
			equalAddrCount++
		}
	}
	if equalAddrCount != len(srv.addrs) {
		return false
	}

	if srv.name != other.Name() {
		return false
	}
	if srv.host != other.Host() {
		return false
	}
	if srv.port != other.Port() {
		return false
	}
	if srv.domain != other.Domain() {
		return false
	}

	return true
}

// String returns the string representation.
func (srv *serviceImpl) String() string {
	return fmt.Sprintf(
		"%s (%s:%d)",
		dns.NewNameWithStrings(srv.name, srv.host, srv.domain),
		srv.host,
		srv.port,
	)
}
