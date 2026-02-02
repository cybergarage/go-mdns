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
	"bufio"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns/dns"
)

// A TCPSocket represents a socket for TCP.
type TCPSocket struct {
	*Socket

	Listener *net.TCPListener
	readBuf  []byte
}

// NewTCPSocket returns a new TCPSocket.
func NewTCPSocket() *TCPSocket {
	sock := &TCPSocket{
		Socket:   NewSocket(),
		Listener: nil,
		readBuf:  make([]byte, MaxPacketSize),
	}
	return sock
}

// Bind binds to Echonet multicast address.
func (sock *TCPSocket) Bind(ifi *net.Interface, ifaddr string, port int) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	var boundAddr *net.TCPAddr
	host := ifaddr
	zone := ""
	if strings.Contains(ifaddr, "%") {
		host, zone, _ = strings.Cut(ifaddr, "%")
	}
	if ip := net.ParseIP(host); ip != nil {
		if zone == "" && ip.To4() == nil && ip.IsLinkLocalUnicast() && ifi != nil {
			zone = ifi.Name
		}
		boundAddr = &net.TCPAddr{IP: ip, Port: port, Zone: zone}
	} else {
		var err error
		boundAddr, err = net.ResolveTCPAddr("tcp", net.JoinHostPort(ifaddr, strconv.Itoa(port)))
		if err != nil {
			return err
		}
	}

	sock.Listener, err = net.ListenTCP("tcp", boundAddr)
	if err != nil {
		return err
	}

	rawConn, err := sock.Listener.SyscallConn()
	if err != nil {
		return err
	}

	var ctrlErr error
	err = rawConn.Control(func(fd uintptr) {
		if err := sock.SetReuseAddrFd(fd, true); err != nil {
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

	sock.SetListenStatus(ifi, ifaddr, port)

	return nil
}

// Close closes the current opened socket.
func (sock *TCPSocket) Close() error {
	if sock.Listener == nil {
		return nil
	}

	sock.Listener.SetDeadline(time.Now().Add(-time.Second))
	err := sock.Listener.Close()
	if err != nil {
		return err
	}

	sock.Listener = nil

	return nil
}

// ReadMessage reads a message from the current opened socket.
func (sock *TCPSocket) ReadMessage(conn net.Conn) (dns.Message, error) {
	bytes, err := io.ReadAll(bufio.NewReader(conn))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	addr, err := dns.NewAddrFromString(
		conn.RemoteAddr().String(),
		dns.WithAddrTransport(dns.TransportTCP),
	)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	msg, err := dns.NewMessageWithBytes(
		bytes,
		dns.WithMessageFrom(addr),
	)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return msg, nil
}

// SendMessage sends a message to the destination address.
func (sock *TCPSocket) SendMessage(addr string, port int, msg dns.Message, timeout time.Duration) (int, error) {
	conn, nWrote, err := sock.dialAndWriteBytes(addr, port, msg.Bytes(), timeout)
	if conn != nil {
		conn.Close()
	}
	return nWrote, err
}

// PostMessage sends a message to the destination address.
func (sock *TCPSocket) PostMessage(addr string, port int, reqMsg dns.Message, timeout time.Duration) (dns.Message, error) {
	conn, _, err := sock.dialAndWriteBytes(addr, port, reqMsg.Bytes(), timeout)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	err = conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		conn.Close()
		return nil, err
	}

	return sock.ReadMessage(conn)
}

// responseToConnection sends a response message to the specified connection.
func (sock *TCPSocket) responseToConnection(conn *net.TCPConn, resMsg dns.Message) error {
	_, err := sock.writeBytesToConnection(conn, resMsg.Bytes())
	return err
}

// writeBytesToConnection sends the specified bytes to the specified connection.
func (sock *TCPSocket) writeBytesToConnection(conn *net.TCPConn, b []byte) (int, error) {
	var nWrote int
	nWrote, err := conn.Write(b)
	if err != nil {
		log.Error(err)
		return nWrote, err
	}
	return nWrote, nil
}

// dialAndWriteBytes sends the specified bytes to the specified destination.
func (sock *TCPSocket) dialAndWriteBytes(addr string, port int, b []byte, timeout time.Duration) (*net.TCPConn, int, error) {
	toAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	conn, err := net.DialTCP("tcp", nil, toAddr)
	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	err = conn.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		conn.Close()
		return nil, 0, err
	}

	nWrote, err := sock.writeBytesToConnection(conn, b)
	if err != nil {
		conn.Close()
	}

	return conn, nWrote, err
}
