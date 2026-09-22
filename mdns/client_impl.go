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
	"net"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns/dns"
	"github.com/cybergarage/go-mdns/mdns/transport"
)

// clientImpl represents a client node instance.
type clientImpl struct {
	*transport.MessageManager
	*serviceSet
	*msgHandler
	queryTimeout     time.Duration
	queryInterval    time.Duration
	maxQueryInterval time.Duration
}

// ClientOption represents a client option.
type ClientOption func(*clientImpl)

// WithClientInterfaces sets the network interfaces which the client listens on.
// All available interfaces are used when no interface is set.
func WithClientInterfaces(ifis ...*net.Interface) ClientOption {
	return func(client *clientImpl) {
		client.SetInterfaces(ifis)
	}
}

// WithClientIPv4Enabled sets whether the client listens on the IPv4 addresses.
func WithClientIPv4Enabled(flag bool) ClientOption {
	return func(client *clientImpl) {
		client.SetIPv4Enabled(flag)
	}
}

// WithClientIPv6Enabled sets whether the client listens on the IPv6 addresses.
func WithClientIPv6Enabled(flag bool) ClientOption {
	return func(client *clientImpl) {
		client.SetIPv6Enabled(flag)
	}
}

// WithClientQueryTimeout sets the default timeout of a query. The timeout is
// used only when the query context has no deadline.
func WithClientQueryTimeout(timeout time.Duration) ClientOption {
	return func(client *clientImpl) {
		client.queryTimeout = timeout
	}
}

// WithClientQueryInterval sets the interval before the first retransmission of
// a query.
//
// RFC 6762: 5.2. Continuous Multicast DNS Querying
// The interval between the queries is doubled, so that the network traffic is
// not increased when a service does not respond.
func WithClientQueryInterval(interval time.Duration) ClientOption {
	return func(client *clientImpl) {
		client.queryInterval = interval
	}
}

// WithClientMaxQueryInterval sets the upper bound of the query interval.
func WithClientMaxQueryInterval(interval time.Duration) ClientOption {
	return func(client *clientImpl) {
		client.maxQueryInterval = interval
	}
}

// NewClient returns a new client instance with the specified options.
func NewClient(opts ...ClientOption) Client {
	client := &clientImpl{
		MessageManager:   transport.NewMessageManager(),
		serviceSet:       newServiceSet(),
		msgHandler:       newMessageHandler(),
		queryTimeout:     DefaultQueryTimeout,
		queryInterval:    DefaultQueryInterval,
		maxQueryInterval: DefaultMaxQueryInterval,
	}

	for _, opt := range opts {
		opt(client)
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

// Services returns the services which the client has discovered.
func (client *clientImpl) Services() []Service {
	return client.serviceSet.Services()
}

// queryUntilDone announces the query message, and retransmits it until the
// context is done.
//
// RFC 6762: 5.2. Continuous Multicast DNS Querying
// The querier sends the first query, and the interval between the successive
// queries is doubled up to the maximum interval.
func (client *clientImpl) queryUntilDone(ctx context.Context, queryMsg Message) error {
	if err := client.AnnounceMessage(queryMsg); err != nil {
		return err
	}

	interval := client.queryInterval
	if interval <= 0 {
		<-ctx.Done()
		return nil
	}

	timer := time.NewTimer(interval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			if err := client.AnnounceMessage(queryMsg); err != nil {
				log.Debugf("mDNS query retransmission failed: %s", err)
			}
			interval *= 2
			if 0 < client.maxQueryInterval && client.maxQueryInterval < interval {
				interval = client.maxQueryInterval
			}
			timer.Reset(interval)
		}
	}
}

// Query sends a question message to the multicast address.
//
// The query is not serialized by a client mutex because it waits for the
// responses until the context is done. Serializing the queries would stop the
// caller from running another query, such as resolving a found service, while
// browsing.
func (client *clientImpl) Query(ctx context.Context, q Query) ([]Service, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, client.queryTimeout)
		defer cancel()
	}

	handler, ok := q.MessageHandler()
	if ok {
		client.RegisterMessageHandler(handler)
		defer client.UnRegisterMessageHandler(handler)
	}

	queryMsg := NewRequestWithQuery(q)

	respondServices := newServiceSet()
	queryResponseHandler := func(resMsg dns.Message) {
		if !resMsg.IsResponse() {
			return
		}
		if queryMsg.IsQueryWithUnicastResponse() {
			if resMsg.From().Transport().Is(dns.TransportMulticast) {
				return
			}
		}

		if !queryMsg.IsQueryAnswer(resMsg) {
			return
		}

		newService, err := NewService(
			WithServiceMessage(resMsg),
		)
		if err != nil {
			return
		}

		added := respondServices.AddService(newService)

		log.Debugf("mDNS Service responded: %s (added=%t)", newService.String(), added)
	}
	client.RegisterMessageHandler(queryResponseHandler)
	defer client.UnRegisterMessageHandler(queryResponseHandler)

	if err := client.queryUntilDone(ctx, queryMsg); err != nil {
		return []Service{}, err
	}

	services := respondServices.Services()
	client.AddServices(services)

	return services, nil
}
