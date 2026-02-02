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
	"net"
)

// Transport represents the transport protocol.
type Transport int

const (
	// TransportUnknown represents an unknown transport protocol.
	TransportUnknown Transport = iota
	// TransportUDP represents the UDP transport protocol.
	TransportUDP
	// TransportTCP represents the TCP transport protocol.
	TransportTCP
	// TransportUDPGroup represents the UDP multicast transport protocol.
	TransportUDPGroup
)

// Addr represents a network address.
type Addr interface {
	// IP returns the IP address.
	IP() net.IP
	// Port returns the port number.
	Port() int
	// Zone returns the IPv6 scoped addressing zone.
	Zone() string
	// String returns the string representation of the address.
	String() string
	// Transport returns the transport protocol.
	Transport() Transport
}
