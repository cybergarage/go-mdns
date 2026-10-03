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

package transport

import "time"

// The mDNS addresses (RFC 6762, 3).
const (
	// Port is the mDNS port.
	Port = 5353
	// UDPPort is the UDP port the servers bind by default.
	UDPPort = Port
	// TCPPort is the TCP port the servers bind by default.
	TCPPort = Port
	// MulticastIPv4Address is the mDNS IPv4 multicast address.
	MulticastIPv4Address = "224.0.0.251"
	// MulticastIPv6Address is the mDNS IPv6 multicast address.
	MulticastIPv6Address = "ff02::fb"
	// MaxPacketSize is the size of a message the sockets read at most, the
	// Ethernet MTU (RFC 6762, 17).
	MaxPacketSize = 1500
)

// The defaults of the socket configuration.
const (
	// DefaultConnectTimeout is the timeout to connect a TCP socket.
	DefaultConnectTimeout = (time.Millisecond * 5000)
	// DefaultRequestTimeout is the timeout to wait for a reply on a socket.
	DefaultRequestTimeout = (time.Millisecond * 5000)
	// DefaultBindRetryCount is how many times a socket is bound again when
	// its port is in use.
	DefaultBindRetryCount = 5
	// DefaultBindRetryWaitTime is the interval of the bind retries.
	DefaultBindRetryWaitTime = (time.Millisecond * 500)
)

// UDPPortRange is how many ports above the configured one a unicast server
// tries when the automatic port binding is enabled.
const UDPPortRange = 100
