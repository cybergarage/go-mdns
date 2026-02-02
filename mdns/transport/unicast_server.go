// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"

	"github.com/cybergarage/go-mdns/mdns/dns"
)

// A UnicastServer represents a unicast server.
type UnicastServer struct {
	*UnicastConfig
	*Server

	TCPSocket  *UnicastTCPSocket
	TCPChannel chan any
	UDPSocket  *UnicastUDPSocket
	UDPChannel chan any
	processor  dns.MessageProcessor
}

// NewUnicastServer returns a new UnicastServer.
func NewUnicastServer() *UnicastServer {
	server := &UnicastServer{
		UnicastConfig: NewDefaultUnicastConfig(),
		Server:        NewServer(),
		TCPSocket:     NewUnicastTCPSocket(),
		TCPChannel:    nil,
		UDPSocket:     NewUnicastUDPSocket(),
		UDPChannel:    nil,
		processor:     nil,
	}
	return server
}

// SetMessageProcessor sets the message processor.
func (server *UnicastServer) SetMessageProcessor(processor dns.MessageProcessor) {
	server.processor = processor
}

// SendMessage send a message to the destination address.
func (server *UnicastServer) SendMessage(addr string, port int, msg dns.Message) (int, error) {
	if server.TCPEnabled() {
		n, err := server.TCPSocket.SendMessage(addr, port, msg, server.ConnectionTimeout())
		if err == nil {
			return n, nil
		}
	}
	return server.UDPSocket.SendMessage(addr, port, msg)
}

// AnnounceMessage sends a message to the multicast address.
func (server *UnicastServer) AnnounceMessage(msg dns.Message) error {
	return server.UDPSocket.AnnounceMessage(msg)
}

// Start starts this server.
func (server *UnicastServer) Start(ifi *net.Interface, ifaddr string, port int) error {
	err := server.UDPSocket.Bind(ifi, ifaddr, port)
	if err != nil {
		server.TCPSocket.Close()
		return err
	}
	server.UDPChannel = make(chan any)
	go handleUnicastUDPConnection(server, server.UDPChannel)

	if server.TCPEnabled() {
		err := server.TCPSocket.Bind(ifi, ifaddr, port)
		if err != nil {
			return err
		}
		server.TCPChannel = make(chan any)
		go handleUnicastTCPListener(server, server.TCPChannel)
	}

	server.TCPSocket.SetListenStatus(ifi, ifaddr, port)
	server.UDPSocket.SetListenStatus(ifi, ifaddr, port)

	return nil
}

// Stop stops this server.
func (server *UnicastServer) Stop() error {
	var lastErr error

	close(server.UDPChannel)
	err := server.UDPSocket.Close()
	if err != nil {
		lastErr = err
	}

	if server.TCPEnabled() {
		close(server.TCPChannel)
		err := server.TCPSocket.Close()
		if err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func handleUnicastUDPRequestMessage(server *UnicastServer, reqMsg dns.Message) {
	if server.processor == nil {
		return
	}

	resMsg, err := server.processor(reqMsg)
	if err != nil || resMsg == nil {
		return
	}

	server.UDPSocket.ResponseMessageForRequestMessage(reqMsg, resMsg)
}

func handleUnicastUDPConnection(server *UnicastServer, cancel chan any) {
	for {
		select {
		case <-cancel:
			return
		default:
			reqMsg, err := server.UDPSocket.ReadMessage()
			if err != nil {
				return
			}

			go handleUnicastUDPRequestMessage(server, reqMsg)
		}
	}
}

func handleUnicastTCPConnection(server *UnicastServer, conn *net.TCPConn) {
	defer conn.Close()

	reqMsg, err := server.TCPSocket.ReadMessage(conn)
	if err != nil {
		return
	}

	if server.processor == nil {
		return
	}

	resMsg, err := server.processor(reqMsg)
	if err != nil || resMsg == nil {
		return
	}

	server.TCPSocket.ResponseMessageToConnection(conn, resMsg)
}

func handleUnicastTCPListener(server *UnicastServer, cancel chan any) {
	for {
		select {
		case <-cancel:
			return
		default:
			tcpConn, err := server.TCPSocket.Listener.AcceptTCP()
			if err != nil {
				return
			}
			go handleUnicastTCPConnection(server, tcpConn)
		}
	}
}
