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
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns/dns"
)

const (
	// browseExpirationCheckInterval is the interval to check the cached
	// services for their expiration.
	browseExpirationCheckInterval = time.Duration(1) * time.Second
)

// Browse browses the services of the specified query, and calls the handler
// when a service is added, updated or removed. Browse blocks until the context
// is done.
//
// RFC 6762: 5.2. Continuous Multicast DNS Querying
// A one-shot query is not enough to follow the services which come and go, so
// the query is retransmitted while browsing, and the responses update the cache
// of the browsed services.
func (client *clientImpl) Browse(ctx context.Context, q Query, handler ServiceHandler) error {
	cache := newServiceCache()

	notify := func(event ServiceEvent) {
		log.Debugf("mDNS service %s", event.String())
		if handler == nil {
			return
		}
		handler(event)
	}

	if queryHandler, ok := q.MessageHandler(); ok {
		client.RegisterMessageHandler(queryHandler)
		defer client.UnRegisterMessageHandler(queryHandler)
	}

	queryMsg := NewRequestWithQuery(q)

	browseHandler := func(resMsg dns.Message) {
		if !resMsg.IsResponse() {
			return
		}
		if !queryMsg.IsQueryAnswer(resMsg) {
			return
		}

		service, err := NewService(
			WithServiceMessage(resMsg),
		)
		if err != nil {
			return
		}

		event, ok := cache.Update(service)
		if !ok {
			return
		}

		if event.Type != ServiceRemoved {
			client.AddService(service)
		}

		notify(event)
	}
	client.RegisterMessageHandler(browseHandler)
	defer client.UnRegisterMessageHandler(browseHandler)

	expirationCtx, stopExpirationChecker := context.WithCancel(ctx)
	defer stopExpirationChecker()

	go func() {
		ticker := time.NewTicker(browseExpirationCheckInterval)
		defer ticker.Stop()
		for {
			select {
			case <-expirationCtx.Done():
				return
			case now := <-ticker.C:
				for _, event := range cache.RemoveExpired(now) {
					notify(event)
				}
			}
		}
	}()

	return client.queryUntilDone(ctx, queryMsg)
}
