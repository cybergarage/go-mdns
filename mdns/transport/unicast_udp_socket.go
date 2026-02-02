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
	"net"
	"strconv"

	"github.com/cybergarage/go-mdns/mdns/dns"
)

// A UnicastUDPSocket represents a socket.
type UnicastUDPSocket struct {
	*UDPSocket
}

// NewUnicastUDPSocket returns a new UnicastUDPSocket.
func NewUnicastUDPSocket() *UnicastUDPSocket {
	sock := &UnicastUDPSocket{
		UDPSocket: NewUDPSocket(),
	}
	return sock
}

// Bind binds to Echonet multicast address.
func (sock *UnicastUDPSocket) Bind(ifi *net.Interface, ifaddr string, port int) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	boundAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(ifaddr, strconv.Itoa(port)))
	if err != nil {
		return err
	}

	sock.Conn, err = net.ListenUDP("udp", boundAddr)
	if err != nil {
		return err
	}

	f, err := sock.Conn.File()
	if err != nil {
		return err
	}

	defer f.Close()

	err = sock.SetReuseAddr(f, true)
	if err != nil {
		sock.Close()
		return err
	}

	sock.SetListenStatus(ifi, ifaddr, port)

	return nil
}

// AnnounceMessage announces the message to the listening multicast address.
func (sock *UnicastUDPSocket) AnnounceMessage(msg dns.Message) error {
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

// ResponseMessageForRequestMessage sends a specified response message to the request node.
func (sock *UnicastUDPSocket) ResponseMessageForRequestMessage(reqMsg dns.Message, resMsg dns.Message) error {
	dstAddr := reqMsg.From().IP().String()
	dstPort := reqMsg.From().Port()
	_, err := sock.SendMessage(dstAddr, dstPort, resMsg)
	return err
}
