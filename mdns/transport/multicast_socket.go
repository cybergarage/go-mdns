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

	"github.com/cybergarage/go-mdns/mdns/dns"
)

// A MulticastSocket represents a socket.
type MulticastSocket struct {
	*UDPSocket
}

// NewMulticastSocket returns a new MulticastSocket.
func NewMulticastSocket() *MulticastSocket {
	sock := &MulticastSocket{
		UDPSocket: NewUDPSocket(dns.TransportUDPGroup),
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
		return errAvailableAddressNotFound
	}

	sock.SetListenStatus(ifi, ifaddr, Port)

	if err != nil {
		return fmt.Errorf("%w (%s)", err, ifi.Name)
	}

	sock.Conn.SetReadBuffer(sock.GetReadBufferSize())

	rawConn, err := sock.Conn.SyscallConn()
	if err != nil {
		return err
	}

	var ctrlErr error
	err = rawConn.Control(func(fd uintptr) {
		if err := sock.SetReuseAddrFd(fd, true); err != nil {
			ctrlErr = err
			return
		}
		if err := sock.SetMulticastLoopFd(fd, ifaddr, true); err != nil {
			ctrlErr = err
			return
		}
		if err := sock.SetMulticastHopsFd(fd, ifaddr, 255); err != nil {
			ctrlErr = err
			return
		}
	})
	if err != nil {
		return err
	}
	if ctrlErr != nil {
		return ctrlErr
	}

	return nil
}

// AnnounceMessage announces the message to the listening multicast address.
func (sock *MulticastSocket) AnnounceMessage(msg dns.Message) error {
	addr, err := sock.ListenAddr()
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
