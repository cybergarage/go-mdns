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
	"os"
	"strconv"
	"syscall"
)

// A Socket represents a socket.
type Socket struct {
	listenInterface *net.Interface
	listenPort      int
	listenAddress   string
}

// NewSocket returns a new UDPSocket.
func NewSocket() *Socket {
	sock := &Socket{
		listenInterface: nil,
		listenPort:      0,
		listenAddress:   "",
	}
	sock.Close()
	return sock
}

// Close initialize this socket.
func (sock *Socket) Close() {
	sock.listenInterface = nil
	sock.listenAddress = ""
	sock.listenPort = 0
}

// SetListenStatus sets the listening interface, port, and address.
func (sock *Socket) SetListenStatus(i *net.Interface, addr string, port int) {
	sock.listenInterface = i
	sock.listenAddress = addr
	sock.listenPort = port
}

// IsListening returns true whether the socket is listening, otherwise false.
func (sock *Socket) IsListening() bool {
	return sock.listenPort != 0
}

// ListenPort returns the listening port.
func (sock *Socket) ListenPort() (int, error) {
	if !sock.IsListening() {
		return 0, errSocketClosed
	}
	return sock.listenPort, nil
}

// ListenInterface returns the listening interface.
func (sock *Socket) ListenInterface() (*net.Interface, error) {
	if !sock.IsListening() {
		return nil, errSocketClosed
	}
	return sock.listenInterface, nil
}

// ListenAddr returns the listening address.
func (sock *Socket) ListenAddr() (string, error) {
	if !sock.IsListening() {
		return "", errSocketClosed
	}

	return sock.listenAddress, nil
}

// ListenIPAddr returns the listening address.
func (sock *Socket) ListenIPAddr() (string, error) {
	port, err := sock.ListenPort()
	if err != nil {
		return "", err
	}

	addr, err := sock.ListenAddr()
	if err != nil {
		return "", err
	}

	return net.JoinHostPort(addr, strconv.Itoa(port)), nil
}

// SetMulticastLoop sets a flag to IP_MULTICAST_LOOP.
// nolint: nosnakecase
func (sock *Socket) SetMulticastLoop(file *os.File, addr string, flag bool) error {
	fd := file.Fd()

	opt := 0
	if flag {
		opt = 1
	}

	if IsIPv6Address(addr) {
		return syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IPV6, syscall.IPV6_MULTICAST_LOOP, opt)
	}
	return syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_MULTICAST_LOOP, opt)
}

// SetMulticastHops sets the multicast TTL (IPv4) / HopLimit (IPv6).
// RFC 6762 requires IP TTL / IPv6 Hop Limit = 255 for mDNS.
func (sock *Socket) SetMulticastHops(file *os.File, addr string, hops int) error {
	if hops < 0 || 255 < hops {
		return fmt.Errorf("invalid multicast hops: %d", hops)
	}

	fd := file.Fd()
	if IsIPv6Address(addr) {
		return syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IPV6, syscall.IPV6_MULTICAST_HOPS, hops)
	}
	return syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_MULTICAST_TTL, hops)
}
