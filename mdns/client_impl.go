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
	"sync"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns/dns"
	"github.com/cybergarage/go-mdns/mdns/transport"
)

// clientImpl represents a client node instance.
type clientImpl struct {
	sync.Mutex
	*transport.MessageManager
	*services
	*msgHandler
}

// NewClient returns a new client instance.
func NewClient() Client {
	client := &clientImpl{
		Mutex:          sync.Mutex{},
		MessageManager: transport.NewMessageManager(),
		services:       newServices(),
		msgHandler:     newMessageHandler(),
	}
	client.MessageManager.SetMessageProcessor(
		func(msg dns.Message) (dns.Message, error) {
			client.processMessageHandlers(msg)
			return nil, nil
		})

	return client
}

// Start starts the client instance.
func (client *clientImpl) Start() error {
	if err := client.Stop(); err != nil {
		return err
	}
	return client.MessageManager.Start()
}

// Stop stops the client instance.
func (client *clientImpl) Stop() error {
	return client.MessageManager.Stop()
}

// Restart restarts the client instance.
func (client *clientImpl) Restart() error {
	if err := client.Stop(); err != nil {
		return err
	}
	return client.Start()
}

// Query sends a question message to the multicast address.
func (client *clientImpl) Query(ctx context.Context, q Query) ([]Service, error) {
	client.Lock()
	defer client.Unlock()

	client.services.Clear()

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, DefaultQueryTimeout)
		defer cancel()
	}

	handler, ok := q.MessageHandler()
	if ok {
		client.RegisterMessageHandler(handler)
		defer client.UnRegisterMessageHandler(handler)
	}

	queryMsg := NewRequestWithQuery(q)

	respondServices := newServices()
	queryResponseHandler := func(resMsg dns.Message) {
		if !resMsg.IsResponse() {
			return
		}
		if queryMsg.IsQueryWithUnicastResponse() {
			if resMsg.From().Transport().Is(dns.TransportUDPGroup) {
				return
			}
		}

		newService, err := NewService(
			WithServiceMessage(resMsg),
		)
		if err != nil {
			return
		}

		added := false
		if queryMsg.IsQueryAnswer(resMsg) {
			added = respondServices.AddService(newService)
		}

		log.Debugf("mDNS Service responded: %s (added=%t)", newService.String(), added)
	}
	client.RegisterMessageHandler(queryResponseHandler)
	defer client.UnRegisterMessageHandler(queryResponseHandler)

	err := client.AnnounceMessage(queryMsg)
	if err != nil {
		return []Service{}, err
	}

	<-ctx.Done()

	client.AddServices(respondServices.Services())

	return respondServices.Services(), nil
}
