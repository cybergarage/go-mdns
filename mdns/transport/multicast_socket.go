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
)

// A MulticastSocket represents a socket.
type MulticastSocket struct {
	*UDPSocket
}

// NewMulticastSocket returns a new MulticastSocket.
func NewMulticastSocket() *MulticastSocket {
	sock := &MulticastSocket{
		UDPSocket: NewUDPSocket(),
	}
	return sock
}

// Bind binds to the Echonet multicast address with the specified interface.
func (sock *MulticastSocket) Bind(ifi *net.Interface) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	err = sock.Listen(ifi)
	if err != nil {
		return fmt.Errorf("%w (%s)", err, ifi.Name)
	}

	sock.Conn.SetReadBuffer(sock.GetReadBufferSize())

	f, err := sock.Conn.File()
	if err != nil {
		return err
	}

	defer f.Close()

	err = sock.SetReuseAddr(f, true)
	if err != nil {
		return err
	}

	err = sock.SetMulticastLoop(f, true)
	if err != nil {
		return err
	}

	sock.SetBoundStatus(ifi, MulticastAddress, UDPPort)

	return nil
}
