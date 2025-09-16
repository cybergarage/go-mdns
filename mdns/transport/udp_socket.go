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
	"io"
	"net"
	"strconv"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns/dns"
)

// A UDPSocket represents a socket for UDP.
type UDPSocket struct {
	*Socket
	Conn           *net.UDPConn
	ReadBufferSize int
	ReadBuffer     []byte
}

// NewUDPSocket returns a new UDPSocket.
func NewUDPSocket() *UDPSocket {
	sock := &UDPSocket{
		Socket:         NewSocket(),
		Conn:           nil,
		ReadBufferSize: MaxPacketSize,
		ReadBuffer:     make([]byte, 0),
	}
	sock.SetReadBufferSize(MaxPacketSize)
	return sock
}

// SetReadBufferSize sets the read buffer size.
func (sock *UDPSocket) SetReadBufferSize(n int) {
	sock.ReadBufferSize = n
	sock.ReadBuffer = make([]byte, n)
}

// GetReadBufferSize returns the read buffer size.
func (sock *UDPSocket) GetReadBufferSize() int {
	return sock.ReadBufferSize
}

// Close closes the current opened socket.
func (sock *UDPSocket) Close() error {
	if sock.Conn == nil {
		return nil
	}

	// FIXME : sock.Conn.Close() hung up on darwin
	/*
		err := sock.Conn.Close()
		if err != nil {
			return err
		}
	*/
	go sock.Conn.Close()
	time.Sleep(time.Millisecond * 100)

	sock.Conn = nil

	return nil
}

// SendMessage sends the message to the destination address.
func (sock *UDPSocket) SendMessage(toAddr string, toPort int, msg *dns.Message) (int, error) {
	toUDPAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(toAddr, strconv.Itoa(toPort)))
	if err != nil {
		return 0, err
	}

	fromAddr, _ := sock.GetBoundAddr()
	fromPort, _ := sock.GetBoundPort()
	log.Debugf("SEND %s -> %s", net.JoinHostPort(fromAddr, strconv.Itoa(fromPort)), net.JoinHostPort(toAddr, strconv.Itoa(toPort)))

	return sock.Conn.WriteToUDP(msg.Bytes(), toUDPAddr)
}

// ReadMessage reads a message from the current opened socket.
func (sock *UDPSocket) ReadMessage() (*dns.Message, error) {
	if sock.Conn == nil {
		return nil, fmt.Errorf("%w: %s", io.EOF, errorSocketClosed)
	}

	n, fromAddr, err := sock.Conn.ReadFromUDP(sock.ReadBuffer)
	if err != nil {
		return nil, err
	}

	toAddr, _ := sock.GetBoundAddr()
	toPort, _ := sock.GetBoundPort()

	log.Debugf("RECV %s -> %s", net.JoinHostPort(fromAddr.IP.String(), strconv.Itoa(fromAddr.Port)), net.JoinHostPort(toAddr, strconv.Itoa(toPort)))

	msg, err := dns.NewMessageWithBytes(sock.ReadBuffer[:n])
	if err != nil {
		log.Debugf("Failed to parse DNS message: %s", err)
		log.HexDebug(sock.ReadBuffer[:n])
		return nil, err
	}

	return msg, nil
}
