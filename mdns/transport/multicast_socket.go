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
	"errors"
	"fmt"
	"net"

	"github.com/cybergarage/go-mdns/mdns/protocol"
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
func (sock *MulticastSocket) Bind(ifi *net.Interface, ifaddr string) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	switch {
	case IsIPv4Address(ifaddr):
		err = sock.Listen(ifi, MulticastIPv4Address, Port)
	case IsIPv6Address(ifaddr):
		err = sock.Listen(ifi, MulticastIPv6Address, Port)
	default:
		return errors.New(errorAvailableAddressNotFound)
	}

	sock.SetBoundStatus(ifi, ifaddr, Port)

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

	err = sock.SetMulticastLoop(f, ifaddr, true)
	if err != nil {
		return err
	}

	return nil
}

// AnnounceMessage announces the message to the bound multicast address.
func (sock *MulticastSocket) AnnounceMessage(msg *protocol.Message) error {
	addr, err := sock.GetBoundAddr()
	if err != nil {
		return err
	}
	toAddr := MulticastIPv4Address
	if IsIPv6Address(addr) {
		toAddr = MulticastIPv6Address
	}
	_, err = sock.SendMessage(toAddr, Port, msg)
	return err
}
