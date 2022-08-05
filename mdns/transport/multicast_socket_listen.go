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

import (
	"fmt"
	"net"
	"strconv"
)

// Listen listens the Ethonet multicast address with the specified interface.
func (sock *MulticastSocket) Listen(ifi *net.Interface) error {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(MulticastIPv4Address, strconv.Itoa(UDPPort)))
	if err != nil {
		return err
	}

	sock.Conn, err = net.ListenMulticastUDP("udp", ifi, addr)
	if err != nil {
		return fmt.Errorf("%w (%s)", err, ifi.Name)
	}

	return nil
}
