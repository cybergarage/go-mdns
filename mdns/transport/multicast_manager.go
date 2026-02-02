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
	"net"

	"github.com/cybergarage/go-mdns/mdns/dns"
)

// A MulticastManager represents a multicast server manager.
type MulticastManager struct {
	Servers   []*MulticastServer
	processor dns.MessageProcessor
}

// NewMulticastManager returns a new MulticastManager.
func NewMulticastManager() *MulticastManager {
	mgr := &MulticastManager{
		Servers:   make([]*MulticastServer, 0),
		processor: nil,
	}
	return mgr
}

// SetMessageProcessor sets a message processor to all servers.
func (mgr *MulticastManager) SetMessageProcessor(processor dns.MessageProcessor) {
	mgr.processor = processor
}

// BoundInterfaces returns the listen interfaces.
func (mgr *MulticastManager) BoundInterfaces() []*net.Interface {
	boundIfs := make([]*net.Interface, 0, len(mgr.Servers))
	for _, server := range mgr.Servers {
		boundIfs = append(boundIfs, server.Interface)
	}
	return boundIfs
}

// AnnounceMessage announces the message from the all bound multicast interfaces.
func (mgr *MulticastManager) AnnounceMessage(msg dns.Message) error {
	var lastErr error
	for _, server := range mgr.Servers {
		if err := server.AnnounceMessage(msg); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// startWithInterface starts this server on the specified interface.
func (mgr *MulticastManager) startWithInterface(ifi *net.Interface, ifaddr string) (*MulticastServer, error) {
	server := NewMulticastServer()
	server.processor = mgr.processor
	if err := server.Start(ifi, ifaddr); err != nil {
		return nil, err
	}
	return server, nil
}

// Start starts servers on the all avairable interfaces.
func (mgr *MulticastManager) Start() error {
	err := mgr.Stop()
	if err != nil {
		return err
	}

	ifis, err := GetAvailableInterfaces()
	if err != nil {
		return err
	}

	for _, ifi := range ifis {
		ifaddrs, err := GetInterfaceAddresses(ifi)
		if err != nil {
			continue
		}
		for _, ifaddr := range ifaddrs {
			server, err := mgr.startWithInterface(ifi, ifaddr)
			if err != nil {
				continue
			}
			server.SetMessageProcessor(mgr.processor)
			mgr.Servers = append(mgr.Servers, server)
		}
	}

	if len(mgr.Servers) == 0 {
		return errAvailableInterfaceFound
	}

	return nil
}

// Stop stops this server.
func (mgr *MulticastManager) Stop() error {
	var lastErr error

	for _, server := range mgr.Servers {
		err := server.Stop()
		if err != nil {
			lastErr = err
		}
	}

	mgr.Servers = make([]*MulticastServer, 0)

	return lastErr
}

// IsRunning returns true whether the local servers are running, otherwise false.
func (mgr *MulticastManager) IsRunning() bool {
	return len(mgr.Servers) != 0
}
