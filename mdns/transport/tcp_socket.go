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
	"net"
	"strconv"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
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

	boundAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(ifaddr, strconv.Itoa(port)))
	if err != nil {
		return err
	}

	sock.Listener, err = net.ListenTCP("tcp", boundAddr)
	if err != nil {
		return err
	}

	f, err := sock.Listener.File()
	if err != nil {
		return err
	}

	defer f.Close()

	err = sock.SetReuseAddr(f, true)
	if err != nil {
		return err
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
func (sock *TCPSocket) ReadMessage(conn net.Conn) (*protocol.Message, error) {
	remoteAddr := conn.RemoteAddr()

	reader := bufio.NewReader(conn)
	msg, err := protocol.NewMessageWithReader(reader)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = msg.From.ParseString(remoteAddr.String())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return msg, nil
}

// SendMessage sends a message to the destination address.
func (sock *TCPSocket) SendMessage(addr string, port int, msg *protocol.Message, timeout time.Duration) (int, error) {
	conn, nWrote, err := sock.dialAndWriteBytes(addr, port, msg.Bytes(), timeout)
	if conn != nil {
		conn.Close()
	}
	return nWrote, err
}

// PostMessage sends a message to the destination address.
func (sock *TCPSocket) PostMessage(addr string, port int, reqMsg *protocol.Message, timeout time.Duration) (*protocol.Message, error) {
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

// ResponseMessageForRequestMessage sends a specified response message to the request node.
func (sock *TCPSocket) ResponseMessageForRequestMessage(reqMsg *protocol.Message, resMsg *protocol.Message, timeout time.Duration) error {
	dstAddr := reqMsg.From.IP.String()
	dstPort := reqMsg.From.Port
	_, err := sock.SendMessage(dstAddr, dstPort, resMsg, timeout)
	return err
}

// ResponseMessageToConnection sends a response message to the specified connection.
func (sock *TCPSocket) ResponseMessageToConnection(conn *net.TCPConn, resMsg *protocol.Message) error {
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
