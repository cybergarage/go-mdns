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
	"time"

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
	ifi    *net.Interface
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

// WithServiceInterface returns a service option with the specified network interface.
func WithServiceInterface(ifi *net.Interface) ServiceOptions {
	return func(srv *serviceImpl) error {
		srv.ifi = ifi
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
		ifi:     nil,
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

// FullName returns the service name with its domain.
func (srv *serviceImpl) FullName() string {
	return dns.NewNameWithStrings(srv.name, srv.domain)
}

// TTL returns the shortest TTL of the service records.
//
// RFC 6762: 10. Resource Record TTL Values and Cache Coherency
// The service is cached until the shortest TTL of its records elapses, and a
// zero TTL means that the service is going away (10.1).
func (srv *serviceImpl) TTL() time.Duration {
	if srv.Message == nil {
		return 0
	}

	fullName := srv.FullName()

	ttl := uint(0)
	hasRecord := false
	for _, record := range srv.Message.ResourceRecordSet() {
		isServiceRecord := record.IsName(fullName)
		if !isServiceRecord && 0 < len(srv.host) {
			isServiceRecord = record.IsName(srv.host)
		}
		if !isServiceRecord {
			continue
		}
		recordTTL := record.TTL()
		if recordTTL == 0 {
			return 0
		}
		if !hasRecord || recordTTL < ttl {
			ttl = recordTTL
			hasRecord = true
		}
	}

	if !hasRecord {
		return 0
	}

	return time.Duration(ttl) * time.Second
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

// Interface returns the network interface which the service was discovered on.
func (srv *serviceImpl) Interface() *net.Interface {
	if srv.ifi != nil {
		return srv.ifi
	}
	if srv.Message == nil {
		return nil
	}
	from := srv.Message.From()
	if from == nil {
		return nil
	}
	return from.Interface()
}

// Addrs returns the service addresses with the service port.
//
// RFC 4007: IPv6 Scoped Address Architecture
// A link-local address is ambiguous without its zone, and it cannot be used to
// connect to the service. The interface which the response was received on is
// used as the zone, because a link-local address is only reachable on the link
// which advertised it.
func (srv *serviceImpl) Addrs() []*net.UDPAddr {
	zone := ""
	if ifi := srv.Interface(); ifi != nil {
		zone = ifi.Name
	}

	addrs := make([]*net.UDPAddr, 0, len(srv.addrs))
	for _, ip := range srv.addrs {
		udpAddr := &net.UDPAddr{ // nolint: exhaustruct,exhaustruct_v5
			IP:   ip,
			Port: srv.port,
		}
		if ip.To4() == nil && ip.IsLinkLocalUnicast() {
			udpAddr.Zone = zone
		}
		addrs = append(addrs, udpAddr)
	}

	return addrs
}

// setAddresses sets the service addresses.
func (srv *serviceImpl) setAddresses(ips []net.IP) {
	srv.addrs = ips
}

// parseMessage updates the service data by the specified message.
func (srv *serviceImpl) parseMessage(msg Message) error {
	srv.Message = msg

	records := msg.ResourceRecordSet()

	// The SRV and TXT records are parsed first because the SRV target name is
	// required to select the address records which belong to this service.
	for _, record := range records {
		err := srv.parseRecord(record)
		if err != nil {
			return err
		}
	}

	srv.parseAddressRecords(records)

	return nil
}

// parseAddressRecords updates the service addresses by the specified records.
func (srv *serviceImpl) parseAddressRecords(records ResourceRecordSet) {
	appendAddress := func(record dns.Record) bool {
		addrRecord, ok := record.(interface{ Address() net.IP })
		if !ok {
			return false
		}
		ip := addrRecord.Address()
		if ip == nil {
			return false
		}
		for _, addr := range srv.addrs {
			if addr.Equal(ip) {
				return false
			}
		}
		srv.addrs = append(srv.addrs, ip)
		return true
	}

	// RFC 6763: 5. Service Instance Resolution
	// The SRV record target names the host, and only the address records of
	// that host belong to this service. A single mDNS response may carry the
	// records of multiple services and hosts, so the address records must not
	// be collected without matching the target name.
	if 0 < len(srv.host) {
		for _, record := range records {
			if !record.IsName(srv.host) {
				continue
			}
			appendAddress(record)
		}
		return
	}

	// No SRV record is included in the message. The address records are
	// collected as they are, because there is no target name to match.
	for _, record := range records {
		appendAddress(record)
	}
}

// ResourceRecordSet returns the service resource records.
func (srv *serviceImpl) ResourceRecordSet() ResourceRecordSet {
	return srv.Message.ResourceRecordSet()
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

	switch rr := record.(type) {
	case dns.PTRRecord:
		// RFC 6763: 4.1. Structured Service Instance Names
		// The data of a PTR record is the service instance name, so a
		// response which holds no SRV record still names the service.
		err := parseNameDomain(rr.DomainName())
		if err != nil {
			return err
		}
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

	thisRRSet := srv.ResourceRecordSet()
	otherRRSet := other.ResourceRecordSet()
	if thisRRSet != nil && otherRRSet != nil {
		return thisRRSet.Equal(otherRRSet)
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
	str := dns.NewNameWithStrings(srv.name, srv.host, srv.domain)
	if len(str) == 0 {
		for _, record := range srv.ResourceRecordSet() {
			name := record.Name()
			if len(name) == 0 {
				continue
			}
			str = name
			break
		}
	}
	addrs := []string{}
	for _, addr := range srv.Addrs() {
		host := addr.IP.String()
		if 0 < len(addr.Zone) {
			host += "%" + addr.Zone
		}
		addrs = append(addrs, host)
	}
	if len(addrs) == 0 {
		from := srv.From()
		if from != nil {
			addrs = append(addrs, from.IP().String())
		}
	}
	if 0 < len(addrs) {
		str += " [" + strings.Join(addrs, ", ") + "]"
	}
	return str
}
