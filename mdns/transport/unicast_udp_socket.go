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
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"syscall"

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

	var boundAddr *net.UDPAddr
	host := ifaddr
	zone := ""
	if strings.Contains(ifaddr, "%") {
		host, zone, _ = strings.Cut(ifaddr, "%")
	}
	if ip := net.ParseIP(host); ip != nil {
		if zone == "" && ip.To4() == nil && ip.IsLinkLocalUnicast() && ifi != nil {
			zone = ifi.Name
		}
		boundAddr = &net.UDPAddr{IP: ip, Port: port, Zone: zone}
	} else {
		var err error
		boundAddr, err = net.ResolveUDPAddr("udp", net.JoinHostPort(ifaddr, strconv.Itoa(port)))
		if err != nil {
			return err
		}
	}

	listenConfig := net.ListenConfig{ // nolint: exhaustruct
		Control: func(network, address string, c syscall.RawConn) error {
			var ctrlErr error
			if err := c.Control(func(fd uintptr) {
				ctrlErr = sock.SetReuseAddrFd(fd, true)
			}); err != nil {
				return err
			}
			return ctrlErr
		},
	}

	pc, err := listenConfig.ListenPacket(context.Background(), "udp", boundAddr.String())
	if err != nil {
		return err
	}

	conn, ok := pc.(*net.UDPConn)
	if !ok {
		_ = pc.Close()
		return fmt.Errorf("invalid udp packet connection: %T", pc)
	}

	sock.Conn = conn

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
