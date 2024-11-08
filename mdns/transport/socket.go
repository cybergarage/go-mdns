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
	"net"
	"os"
	"strconv"
	"syscall"
)

// A Socket represents a socket.
type Socket struct {
	BoundInterface *net.Interface
	BoundPort      int
	BoundAddress   string
}

// NewSocket returns a new UDPSocket.
func NewSocket() *Socket {
	sock := &Socket{
		BoundInterface: nil,
		BoundPort:      0,
		BoundAddress:   "",
	}
	sock.Close()
	return sock
}

// Close initialize this socket.
func (sock *Socket) Close() {
	sock.BoundInterface = nil
	sock.BoundAddress = ""
	sock.BoundPort = 0
}

// SetBoundStatus sets the bound interface, port, and address.
func (sock *Socket) SetBoundStatus(i *net.Interface, addr string, port int) {
	sock.BoundInterface = i
	sock.BoundAddress = addr
	sock.BoundPort = port
}

// IsBound returns true whether the socket is bound, otherwise false.
func (sock *Socket) IsBound() bool {
	return sock.BoundPort != 0
}

// GetBoundPort returns the bound port.
func (sock *Socket) GetBoundPort() (int, error) {
	if !sock.IsBound() {
		return 0, errors.New(errorSocketClosed)
	}
	return sock.BoundPort, nil
}

// GetBoundInterface returns the bound interface.
func (sock *Socket) GetBoundInterface() (*net.Interface, error) {
	if !sock.IsBound() {
		return nil, errors.New(errorSocketClosed)
	}
	return sock.BoundInterface, nil
}

// GetBoundAddr returns the bound address.
func (sock *Socket) GetBoundAddr() (string, error) {
	if !sock.IsBound() {
		return "", errors.New(errorSocketClosed)
	}

	return sock.BoundAddress, nil
}

// GetBoundIPAddr returns the bound address.
func (sock *Socket) GetBoundIPAddr() (string, error) {
	port, err := sock.GetBoundPort()
	if err != nil {
		return "", err
	}

	addr, err := sock.GetBoundAddr()
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
