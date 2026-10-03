// Copyright (C) 2026 The go-mdns Authors All rights reserved.
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
	"github.com/cybergarage/go-mdns/mdns/transport"
)

// TTLs of the records a responder publishes (RFC 6762, 10): the records
// which name a host or depend on its addresses expire sooner than the rest.
const (
	HostRecordTTL  = 120
	OtherRecordTTL = 4500
)

// LocalService is a DNS-SD service which this node publishes through a
// Server (RFC 6763).
//
// For example, a Matter device in commissioning mode publishes
//
//	&mdns.LocalService{
//		Instance: "665F6E75B5D3A9C2",
//		Service:  "_matterc._udp",
//		Subtypes: []string{"_L3840", "_S15", "_V65521", "_CM"},
//		Host:     "B75AFB458ECD6D6F",
//		Port:     5540,
//		TXT:      []string{"D=3840", "CM=1", "VP=65521+32769"},
//	}
type LocalService struct {
	// Instance is the instance name, such as "My Printer" or
	// "665F6E75B5D3A9C2".
	Instance string
	// Service is the service type, such as "_http._tcp".
	Service string
	// Domain is the domain; empty means "local".
	Domain string
	// Subtypes are the subtype labels the service can also be browsed by,
	// such as "_printer" (RFC 6763, 7.1). A label is given without the
	// "._sub" suffix and the service type: "_printer", not
	// "_printer._sub._http._tcp".
	Subtypes []string
	// Host is the host name which the SRV record points to, such as
	// "host". The domain is added unless it is given, so "host.local" is
	// the same host; a name in another domain is not valid.
	Host string
	// Port is the service port.
	Port int
	// TXT holds the TXT record strings, such as "key=value".
	TXT []string
	// Addresses are the addresses the host name resolves to. Empty means
	// the addresses of the interface a query arrives on, or which an
	// announcement is sent from.
	Addresses []net.IP
}

// Validate reports whether the service can be published.
func (svc *LocalService) Validate() error {
	if svc == nil {
		return fmt.Errorf("%w: service", dns.ErrNil)
	}
	for name, v := range map[string]string{"instance": svc.Instance, "service": svc.Service, "host": svc.Host} {
		if v == "" {
			return fmt.Errorf("%w: %s name is empty", dns.ErrInvalid, name)
		}
	}
	labels := dns.SplitName(svc.Service)
	if len(labels) != 2 || !strings.HasPrefix(labels[0], "_") || (labels[1] != "_udp" && labels[1] != "_tcp") {
		return fmt.Errorf("%w: service type %q is not \"_name._udp\" or \"_name._tcp\"", dns.ErrInvalid, svc.Service)
	}
	if svc.Port <= 0 || 0xFFFF < svc.Port {
		return fmt.Errorf("%w: port %d", dns.ErrInvalid, svc.Port)
	}
	if host := svc.hostLabel(); host == "" || strings.Contains(host, dns.LabelSeparator) {
		return fmt.Errorf("%w: host %q is not a single label with or without the domain %q, such as \"host\" or \"host.%s\"", dns.ErrInvalid, svc.Host, svc.domain(), svc.domain())
	}
	for _, st := range svc.Subtypes {
		if err := validateSubtype(st); err != nil {
			return err
		}
	}
	if _, err := svc.records(nil, OtherRecordTTL, HostRecordTTL); err != nil {
		return err
	}
	return nil
}

// maxLabelLength is the longest label a DNS name can hold (RFC 1035, 2.3.4).
const maxLabelLength = 63

// validateSubtype reports whether st is a subtype label, such as "_printer",
// rather than a subtype name, such as "_printer._sub._http._tcp" as
// avahi-publish-service takes it, which would be published under a doubled
// name (RFC 6763, 7.1).
func validateSubtype(st string) error {
	switch {
	case st == "":
		return fmt.Errorf("%w: subtype is empty", dns.ErrInvalid)
	case strings.Contains(st, dns.LabelSeparator):
		label, _, _ := strings.Cut(st, dns.LabelSeparator)
		return fmt.Errorf("%w: subtype %q is not a single label, such as %q", dns.ErrInvalid, st, label)
	case maxLabelLength < len(st):
		return fmt.Errorf("%w: subtype %q is longer than %d bytes", dns.ErrInvalid, st, maxLabelLength)
	}
	return nil
}

func (svc *LocalService) domain() string {
	if svc.Domain == "" {
		return LocalDomain
	}
	return svc.Domain
}

// ServiceName returns the service type with its domain, such as
// "_http._tcp.local".
func (svc *LocalService) ServiceName() string {
	return dns.NewNameWithStrings(svc.Service, svc.domain())
}

// FullName returns the service instance name, such as
// "My Printer._http._tcp.local".
func (svc *LocalService) FullName() string {
	return dns.NewNameWithStrings(svc.Instance, svc.ServiceName())
}

// HostName returns the host name with its domain, such as "host.local".
// The domain is added to Host unless Host already ends with it, so "host",
// "host.local" and "host.local." all give "host.local".
func (svc *LocalService) HostName() string {
	return dns.NewNameWithStrings(svc.hostLabel(), svc.domain())
}

// hostLabel returns Host without the service's domain and a trailing dot.
func (svc *LocalService) hostLabel() string {
	host := strings.TrimSuffix(svc.Host, dns.LabelSeparator)
	suffix := dns.LabelSeparator + svc.domain()
	if len(suffix) < len(host) && strings.EqualFold(host[len(host)-len(suffix):], suffix) {
		host = host[:len(host)-len(suffix)]
	}
	return host
}

// SubtypeNames returns the names the service is browsed by as a subtype,
// such as "_printer._sub._http._tcp.local".
func (svc *LocalService) SubtypeNames() []string {
	names := make([]string, 0, len(svc.Subtypes))
	for _, st := range svc.Subtypes {
		names = append(names, dns.NewNameWithStrings(st, Subtype, svc.ServiceName()))
	}
	return names
}

// serviceTypeEnumerationName returns the name which lists the service types
// of the service's domain, "_services._dns-sd._udp.local".
func (svc *LocalService) serviceTypeEnumerationName() string {
	return dns.NewNameWithStrings(ServiceTypeEnumerationName, svc.domain())
}

// addresses returns the service's own addresses, or those of ifi when it
// has none.
func (svc *LocalService) addresses(ifi *net.Interface) []net.IP {
	if 0 < len(svc.Addresses) {
		return svc.Addresses
	}
	return interfaceAddresses(ifi)
}

// localServiceRecords are the records which publish a service.
type localServiceRecords struct {
	// ptr points the service type to the instance, and subtypes point each
	// subtype name to it. enum points the type enumeration name to the
	// service type. They are shared records: other nodes publish the same
	// names.
	ptr      dns.ResourceRecord
	subtypes []dns.ResourceRecord
	enum     dns.ResourceRecord
	// srv, txt and addrs are unique to this node, and carry the cache-flush
	// bit (RFC 6762, 10.2).
	srv   dns.ResourceRecord
	txt   dns.ResourceRecord
	addrs []dns.ResourceRecord
}

func (r *localServiceRecords) all() []dns.ResourceRecord {
	records := make([]dns.ResourceRecord, 0, 5+len(r.subtypes)+len(r.addrs))
	records = append(records, r.ptr)
	records = append(records, r.subtypes...)
	records = append(records, r.enum, r.srv, r.txt)
	return append(records, r.addrs...)
}

// records builds the records of the service for interface ifi with the
// given TTLs; a zero TTL builds the goodbye records (RFC 6762, 10.1).
func (svc *LocalService) records(ifi *net.Interface, otherTTL, hostTTL uint) (*localServiceRecords, error) {
	var err error
	r := &localServiceRecords{
		ptr:      nil,
		subtypes: make([]dns.ResourceRecord, 0, len(svc.Subtypes)),
		enum:     nil,
		srv:      nil,
		txt:      nil,
		addrs:    []dns.ResourceRecord{},
	}

	if r.ptr, err = dns.NewPTRResourceRecord(svc.ServiceName(), svc.FullName(), otherTTL); err != nil {
		return nil, err
	}
	for _, name := range svc.SubtypeNames() {
		st, err := dns.NewPTRResourceRecord(name, svc.FullName(), otherTTL)
		if err != nil {
			return nil, err
		}
		r.subtypes = append(r.subtypes, st)
	}
	if r.enum, err = dns.NewPTRResourceRecord(svc.serviceTypeEnumerationName(), svc.ServiceName(), otherTTL); err != nil {
		return nil, err
	}

	port := uint16(svc.Port) // nolint: gosec // Validate checks the port range.
	if r.srv, err = dns.NewSRVResourceRecord(svc.FullName(), 0, 0, port, svc.HostName(), hostTTL); err != nil {
		return nil, err
	}
	r.srv.SetUnicastResponse(true)
	if r.txt, err = dns.NewTXTResourceRecord(svc.FullName(), svc.TXT, otherTTL); err != nil {
		return nil, err
	}
	r.txt.SetUnicastResponse(true)
	for _, ip := range svc.addresses(ifi) {
		addr, err := dns.NewAddressResourceRecord(svc.HostName(), ip, hostTTL)
		if err != nil {
			return nil, err
		}
		addr.SetUnicastResponse(true)
		r.addrs = append(r.addrs, addr)
	}
	return r, nil
}

// interfaceAddresses returns the unicast addresses of ifi, or of every
// interface transport.GetAvailableInterfaces considers usable when ifi is
// nil. Using the same selection as the one the servers bind to keeps this
// fallback from resurfacing an address on an interface (such as Docker's
// docker0) that the bound sockets themselves already exclude.
func interfaceAddresses(ifi *net.Interface) []net.IP {
	var ifis []net.Interface
	if ifi != nil {
		ifis = []net.Interface{*ifi}
	} else if all, err := transport.GetAvailableInterfaces(); err == nil {
		ifis = make([]net.Interface, 0, len(all))
		for _, i := range all {
			ifis = append(ifis, *i)
		}
	}
	ips := make([]net.IP, 0, 2*len(ifis))
	for _, i := range ifis {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.IsMulticast() || ipnet.IP.IsUnspecified() {
				continue
			}
			ips = append(ips, ipnet.IP)
		}
	}
	return ips
}
