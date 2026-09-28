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

package mdns

import (
	"errors"
	"math/rand/v2"
	"net"
	"slices"
	"sync"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns/dns"
	"github.com/cybergarage/go-mdns/mdns/transport"
)

// Timing of the announcements and of the delayed responses (RFC 6762, 6
// and 8.3).
const (
	announceCount    = 2
	announceInterval = time.Second
	minSharedDelay   = 20 * time.Millisecond
	maxSharedDelay   = 120 * time.Millisecond
)

// Server represents a server node instance.
//
// A server is a responder. It publishes the services registered with
// Register: it answers the queries for them, announces them when they are
// registered or the server starts, and withdraws them when they are
// deregistered or the server stops. It also passes every received message to
// the registered message handlers.
//
// The responder does not probe for name conflicts yet (RFC 6762, 8.1 and
// 9): a registered name is assumed to be unique on the link.
type Server struct {
	sync.Mutex
	*transport.MessageManager
	*serviceSet
	*msgHandler
	*responder
	running bool
	done    chan struct{}
	wg      sync.WaitGroup
}

// NewServer returns a new server instance.
func NewServer() *Server {
	server := &Server{
		Mutex:          sync.Mutex{},
		MessageManager: transport.NewMessageManager(),
		serviceSet:     newServiceSet(),
		msgHandler:     newMessageHandler(),
		responder:      newResponder(),
		running:        false,
		done:           nil,
		wg:             sync.WaitGroup{},
	}
	server.SetMessageProcessor(server.MessageReceived)
	return server
}

// Start starts the server instance, and announces the registered services.
func (server *Server) Start() error {
	if err := server.Stop(); err != nil {
		return err
	}
	if err := server.MessageManager.Start(); err != nil {
		return err
	}
	server.Lock()
	server.running = true
	server.done = make(chan struct{})
	server.Unlock()
	for _, svc := range server.localServices() {
		server.announce(svc)
	}
	return nil
}

// Stop withdraws the registered services and stops the server instance.
func (server *Server) Stop() error {
	server.Lock()
	running := server.running
	if running {
		server.running = false
		close(server.done)
	}
	server.Unlock()
	if running {
		server.wg.Wait()
		for _, svc := range server.localServices() {
			server.sendAnnouncement(svc, true)
		}
	}
	return server.MessageManager.Stop()
}

// Restart restarts the server instance.
func (server *Server) Restart() error {
	if err := server.Stop(); err != nil {
		return err
	}
	return server.Start()
}

// Register publishes svc. A service with the same instance name is
// replaced. The server keeps a copy of svc, so later changes to svc are
// published only by registering it again.
func (server *Server) Register(svc *LocalService) error {
	if err := svc.Validate(); err != nil {
		return err
	}
	copied := copyLocalService(svc)
	server.responder.register(copied)
	server.announce(copied)
	return nil
}

// Deregister withdraws the service with the instance name of svc, sending
// goodbye records if the server is running.
func (server *Server) Deregister(svc *LocalService) error {
	if svc == nil {
		return errors.New("mdns: nil service")
	}
	removed, ok := server.responder.deregister(svc)
	if !ok {
		return nil
	}
	server.Lock()
	running := server.running
	server.Unlock()
	if running {
		server.sendAnnouncement(removed, true)
	}
	return nil
}

// LocalServices returns the registered services.
func (server *Server) LocalServices() []*LocalService {
	return server.localServices()
}

// announce sends the announcements of svc in the background while the
// server is running.
func (server *Server) announce(svc *LocalService) {
	server.Lock()
	defer server.Unlock()
	if !server.running {
		return
	}
	done := server.done
	server.wg.Go(func() {
		for i := range announceCount {
			if 0 < i {
				select {
				case <-done:
					return
				case <-time.After(announceInterval):
				}
			}
			// A service deregistered or replaced since stops being
			// announced; otherwise a late announcement would follow its
			// goodbye records and publish it again.
			announced := server.ifRegistered(svc, func() {
				server.sendAnnouncement(svc, false)
			})
			if !announced {
				return
			}
		}
	})
}

// sendAnnouncement sends the records of svc, or its goodbye records, from
// every multicast socket with the addresses of that socket's interface.
func (server *Server) sendAnnouncement(svc *LocalService, goodbye bool) {
	for _, ms := range server.MessageManager.MulticastManager.Servers {
		ifi, err := ms.MulticastSocket.ListenInterface()
		if err != nil {
			ifi = nil
		}
		msg, err := announcement(svc, ifi, goodbye)
		if err != nil {
			log.Warnf("mdns: build the announcement of %s: %s", svc.FullName(), err)
			continue
		}
		if err := ms.AnnounceMessage(msg); err != nil {
			log.Debugf("mdns: announce %s: %s", svc.FullName(), err)
		}
	}
}

// MessageReceived passes a query to the message handlers, and returns the
// response to it for the registered services, or nil.
func (server *Server) MessageReceived(msg dns.Message) (dns.Message, error) {
	if msg.IsResponse() {
		return nil, nil
	}

	server.processMessageHandlers(msg)

	var ifi *net.Interface
	if from := msg.From(); from != nil {
		ifi = from.Interface()
	}
	res, shared := server.answer(msg, ifi)
	if res == nil {
		return nil, nil
	}
	// A multicast response with shared records is delayed, so that the
	// responses of other nodes to the same query do not collide
	// (RFC 6762, 6).
	if shared && !isLegacyUnicastQuery(msg) && !msg.IsQueryWithUnicastResponse() {
		delay := minSharedDelay + rand.N(maxSharedDelay-minSharedDelay) // nolint: gosec
		time.Sleep(delay)
	}
	return res, nil
}

func copyLocalService(svc *LocalService) *LocalService {
	copied := *svc
	copied.Subtypes = slices.Clone(svc.Subtypes)
	copied.TXT = slices.Clone(svc.TXT)
	copied.Addresses = slices.Clone(svc.Addresses)
	return &copied
}
