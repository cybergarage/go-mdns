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
	"bytes"
	"errors"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/cybergarage/go-mdns/mdns/dns"
)

func TestNewTCPSocket(t *testing.T) {
	NewTCPSocket()
}

type tcpReaderConn struct {
	net.Conn
	reader io.Reader
}

func (conn tcpReaderConn) Read(b []byte) (int, error) {
	return conn.reader.Read(b)
}

func TestTCPReadMessageSizeLimit(t *testing.T) {
	input := tcpReaderConn{reader: bytes.NewReader(make([]byte, maxTCPMessageSize+1))}
	if msg, err := NewTCPSocket().ReadMessage(input); msg != nil || !errors.Is(err, ErrInvalid) {
		t.Fatalf("ReadMessage = %v, %v; want a size error", msg, err)
	}
}

func TestUnicastTCPRequestAndShutdown(t *testing.T) {
	server := NewUnicastServer()
	server.SetTCPEnabled(true)
	server.SetMessageProcessor(func(msg dns.Message) (dns.Message, error) {
		return dns.NewResponseMessage(dns.WithMessageID(msg.ID())), nil
	})
	if err := server.Start(nil, "127.0.0.1", 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { server.Stop() })
	addr := server.TCPSocket.Listener.Addr().(*net.TCPAddr)
	req := dns.NewRequestMessage(dns.WithMessageID(123))
	reply, err := NewTCPSocket().PostMessage(addr.IP.String(), addr.Port, req, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if !reply.IsResponse() || reply.ID() != req.ID() {
		t.Fatalf("unexpected response: %v", reply)
	}

	idle, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		t.Fatal(err)
	}
	defer idle.Close()
	if err := idle.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := server.Stop(); err != nil {
		t.Fatal(err)
	}
	if _, err := idle.Read(make([]byte, 1)); err == nil || errors.Is(err, os.ErrDeadlineExceeded) {
		t.Fatalf("read after Stop = %v, want connection close", err)
	}
}
