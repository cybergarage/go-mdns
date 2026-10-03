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
	"context"
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
// Register: it probes their names, answers the queries for them, announces
// them, and withdraws them when they are deregistered or the server stops. It
// also passes every received message to the registered message handlers.
//
// Before a service is published, its instance name and host name are probed
// to check that no other node on the link holds them (RFC 6762, 8.1), and a
// published service whose names another node claims later is probed again
// (RFC 6762, 9). The server does not rename a service on a conflict: Register
// returns a *ConflictError, and a conflict found later is reported to the
// handler given by WithServerConflictHandler, so that the application chooses
// a new name.
type Server struct {
	sync.Mutex
	*transport.MessageManager
	*serviceSet
	*msgHandler
	*responder
	running         bool
	done            chan struct{}
	wg              sync.WaitGroup
	conflictHandler ConflictHandler
	conflictMutex   sync.Mutex
	conflicts       []time.Time
}

// ConflictHandler is called when a published service is withdrawn because
// another node holds one of its names (RFC 6762, 9).
type ConflictHandler func(err *ConflictError)

// ServerOption configures a server.
type ServerOption func(*Server)

// WithServerConflictHandler sets the handler which is called when a published
// service is withdrawn because another node holds one of its names. The
// handler runs on its own goroutine, so it may register the service again
// with a new name.
func WithServerConflictHandler(handler ConflictHandler) ServerOption {
	return func(server *Server) {
		server.conflictHandler = handler
	}
}

// WithServerInterfaces sets the network interfaces the server listens and
// publishes on. All available interfaces are used when none is set.
func WithServerInterfaces(ifis ...*net.Interface) ServerOption {
	return func(server *Server) {
		server.SetInterfaces(ifis)
	}
}

// WithServerIPv4Enabled sets whether the server listens on the IPv4
// addresses.
func WithServerIPv4Enabled(flag bool) ServerOption {
	return func(server *Server) {
		server.SetIPv4Enabled(flag)
	}
}

// WithServerIPv6Enabled sets whether the server listens on the IPv6
// addresses.
func WithServerIPv6Enabled(flag bool) ServerOption {
	return func(server *Server) {
		server.SetIPv6Enabled(flag)
	}
}

// NewServer returns a new server instance.
func NewServer(opts ...ServerOption) *Server {
	server := &Server{
		Mutex:           sync.Mutex{},
		MessageManager:  transport.NewMessageManager(),
		serviceSet:      newServiceSet(),
		msgHandler:      newMessageHandler(),
		responder:       newResponder(),
		running:         false,
		done:            nil,
		wg:              sync.WaitGroup{},
		conflictHandler: nil,
		conflictMutex:   sync.Mutex{},
		conflicts:       []time.Time{},
	}
	for _, opt := range opts {
		opt(server)
	}
	server.SetMessageProcessor(server.MessageReceived)
	return server
}

// Start starts the server instance. The services which stay registered
// after Stop are probed again and published in the background, and a
// conflict is reported to the conflict handler.
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
		server.reprobe(svc)
	}
	return nil
}

// Stop withdraws the published services and stops the server instance. The
// services stay registered, and they are published again by Start.
func (server *Server) Stop() error {
	server.Lock()
	running := server.running
	var published []*LocalService
	if running {
		server.running = false
		close(server.done)
		published = server.localServices()
	}
	server.Unlock()
	if running {
		server.wg.Wait()
		for _, svc := range published {
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

// Register publishes svc, and returns when it is published. A service with
// the same instance name is replaced. The server keeps a copy of svc, so
// later changes to svc are published only by registering it again.
//
// The instance name and the host name of svc are probed first (RFC 6762,
// 8.1), which takes about a second; a name which a published service of this
// server already holds is not probed again, so registering a service again to
// update its TXT strings returns at once. Register returns a *ConflictError,
// which wraps ErrConflict, when another node holds one of the names: the
// service is not published, and the caller registers it again with a new
// name. It returns ErrNotRunning when the server is not running or stops,
// ErrDeregistered when the service is deregistered or registered again
// meanwhile, and the error of ctx when ctx is done.
func (server *Server) Register(ctx context.Context, svc *LocalService) error {
	if err := svc.Validate(); err != nil {
		return err
	}
	copied := copyLocalService(svc)

	server.Lock()
	running, done := server.running, server.done
	server.Unlock()
	if !running {
		return ErrNotRunning
	}

	names := server.responder.namesToProbe(copied)
	if len(names) == 0 {
		server.responder.register(copied)
		server.announce(copied)
		return nil
	}

	p := newProbe(copied, names)
	server.responder.addProbe(p)
	if err := server.runProbe(ctx, done, p); err != nil {
		server.responder.removeProbe(p)
		return err
	}
	if err := server.establish(p); err != nil {
		server.responder.removeProbe(p)
		return err
	}
	server.announce(copied)
	return nil
}

// establish publishes the service of p unless the server stopped or p was
// canceled meanwhile.
func (server *Server) establish(p *probe) error {
	server.Lock()
	defer server.Unlock()
	if !server.running {
		return ErrNotRunning
	}
	if !server.responder.establish(p) {
		if p.cancelErr != nil {
			return p.cancelErr
		}
		return ErrDeregistered
	}
	return nil
}

// reprobe probes the names of svc, a published service, again in the
// background (RFC 6762, 9). The service is published again when its names
// are found unique, and it is withdrawn and reported to the conflict handler
// when another node holds one of them.
func (server *Server) reprobe(svc *LocalService) {
	server.Lock()
	defer server.Unlock()
	if !server.running {
		return
	}
	p, ok := server.responder.beginReprobe(svc)
	if !ok {
		return
	}
	done := server.done
	server.wg.Go(func() {
		err := server.runProbe(context.Background(), done, p)
		switch {
		case err == nil:
			if err := server.establish(p); err == nil {
				server.announce(p.svc)
			} else if errors.Is(err, ErrNotRunning) {
				// The server stopped after the probe: keep the service
				// for the next start.
				server.responder.restore(p)
			}
		case errors.Is(err, ErrNotRunning):
			server.responder.restore(p)
		default:
			server.responder.removeProbe(p)
			var conflict *ConflictError
			if errors.As(err, &conflict) {
				log.Warnf("mdns: %s is withdrawn: %s", svc.FullName(), err)
				if handler := server.conflictHandler; handler != nil {
					go handler(conflict)
				}
			}
		}
	})
}

// Deregister withdraws the service with the instance name of svc, sending
// goodbye records if it is published and the server is running. A probe of
// the service is canceled, and its Register returns ErrDeregistered.
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

// LocalServices returns the published services. A service whose names are
// being probed is not included.
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
		server.checkResponse(msg)
		return nil, nil
	}

	server.processMessageHandlers(msg)
	server.checkProbe(msg)

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
