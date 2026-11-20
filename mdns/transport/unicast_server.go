// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns/dns"
)

const maxUnicastTCPConnections = 64

// A UnicastServer represents a unicast server.
type UnicastServer struct {
	*UnicastConfig
	*Server

	TCPSocket  *UnicastTCPSocket
	TCPChannel chan any
	UDPSocket  *UnicastUDPSocket
	UDPChannel chan any

	processorMutex sync.RWMutex
	processor      dns.MessageProcessor
}

// NewUnicastServer returns a new UnicastServer.
func NewUnicastServer() *UnicastServer {
	server := &UnicastServer{
		UnicastConfig:  NewDefaultUnicastConfig(),
		Server:         NewServer(),
		TCPSocket:      NewUnicastTCPSocket(),
		TCPChannel:     nil,
		UDPSocket:      NewUnicastUDPSocket(),
		UDPChannel:     nil,
		processorMutex: sync.RWMutex{},
		processor:      nil,
	}
	return server
}

// SetMessageProcessor sets the message processor. It may be called while
// the server receives messages.
func (server *UnicastServer) SetMessageProcessor(processor dns.MessageProcessor) {
	server.processorMutex.Lock()
	defer server.processorMutex.Unlock()
	server.processor = processor
}

// messageProcessor returns the message processor.
func (server *UnicastServer) messageProcessor() dns.MessageProcessor {
	server.processorMutex.RLock()
	defer server.processorMutex.RUnlock()
	return server.processor
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
	if server.TCPEnabled() {
		err := server.TCPSocket.Bind(ifi, ifaddr, port)
		if err != nil {
			server.UDPSocket.Close()
			server.TCPSocket.Close()
			return err
		}
		server.TCPChannel = make(chan any)
	}

	server.TCPSocket.SetListenStatus(ifi, ifaddr, port)
	server.UDPSocket.SetListenStatus(ifi, ifaddr, port)
	server.UDPChannel = make(chan any)
	go handleUnicastUDPConnection(server, server.UDPChannel)
	if server.TCPChannel != nil {
		go handleUnicastTCPListener(server, server.TCPSocket.Listener, server.TCPChannel)
	}

	return nil
}

// Stop stops this server.
func (server *UnicastServer) Stop() error {
	var lastErr error

	if server.UDPChannel != nil {
		close(server.UDPChannel)
		server.UDPChannel = nil
	}
	err := server.UDPSocket.Close()
	if err != nil {
		lastErr = err
	}

	if server.TCPChannel != nil {
		close(server.TCPChannel)
		server.TCPChannel = nil
		err := server.TCPSocket.Close()
		if err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func handleUnicastUDPRequestMessage(server *UnicastServer, reqMsg dns.Message) {
	processor := server.messageProcessor()
	if processor == nil {
		return
	}

	resMsg, err := processor(reqMsg)
	if err != nil {
		log.Error(err)
	}
	if resMsg == nil {
		return
	}

	if err := server.UDPSocket.responseForRequest(reqMsg, resMsg); err != nil {
		log.Error(err)
	}
}

func handleUnicastUDPConnection(server *UnicastServer, cancel chan any) {
	for {
		select {
		case <-cancel:
			return
		default:
			reqMsg, err := server.UDPSocket.ReadMessage()
			if err != nil {
				if errors.Is(err, errInvalidDatagram) {
					continue
				}
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

	processor := server.messageProcessor()
	if processor == nil {
		return
	}

	resMsg, err := processor(reqMsg)
	if err != nil {
		log.Error(err)
		return
	}
	if resMsg == nil {
		return
	}

	if err := server.TCPSocket.responseToConnection(conn, resMsg); err != nil {
		log.Error(err)
	}
}

func handleUnicastTCPListener(server *UnicastServer, listener *net.TCPListener, cancel <-chan any) {
	ctx, closeConnections := context.WithCancel(context.Background())
	defer closeConnections()
	active := make(chan struct{}, maxUnicastTCPConnections)
	timeout := server.ConnectionTimeout()
	if timeout <= 0 {
		timeout = DefaultConnectTimeout
	}

	for {
		select {
		case <-cancel:
			return
		default:
			tcpConn, err := listener.AcceptTCP()
			if err != nil {
				return
			}
			if err := tcpConn.SetDeadline(time.Now().Add(timeout)); err != nil {
				tcpConn.Close()
				continue
			}
			select {
			case active <- struct{}{}:
				stop := context.AfterFunc(ctx, func() { tcpConn.Close() })
				go func() {
					defer func() {
						stop()
						<-active
					}()
					handleUnicastTCPConnection(server, tcpConn)
				}()
			default:
				tcpConn.Close()
			}
		}
	}
}
