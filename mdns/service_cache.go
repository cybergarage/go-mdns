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
	"sync"
	"time"
)

// serviceCacheEntry represents a cached service with its expiration time.
type serviceCacheEntry struct {
	service   Service
	expiredAt time.Time
}

// serviceCache holds the discovered services until their records expire.
//
// RFC 6762: 10. Resource Record TTL Values and Cache Coherency
// A cached record is held until its TTL elapses, and a responder which is going
// away announces its records with a TTL of zero (10.1).
type serviceCache struct {
	mutex   sync.Mutex
	entries map[string]*serviceCacheEntry
}

// newServiceCache returns a blank service cache.
func newServiceCache() *serviceCache {
	return &serviceCache{
		mutex:   sync.Mutex{},
		entries: map[string]*serviceCacheEntry{},
	}
}

// Update updates the cache by the specified service, and returns the event
// which the change caused. The second return value is false when the service
// does not change the cache.
func (cache *serviceCache) Update(service Service) (ServiceEvent, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	key := service.FullName()
	if len(key) == 0 {
		return ServiceEvent{}, false // nolint: exhaustruct
	}

	ttl := service.TTL()

	// RFC 6762: 10.1. Goodbye Packets
	// A record with a TTL of zero means that the record is no longer valid.
	if ttl <= 0 {
		if _, ok := cache.entries[key]; !ok {
			return ServiceEvent{}, false // nolint: exhaustruct
		}
		delete(cache.entries, key)
		return ServiceEvent{
			Type:    ServiceRemoved,
			Service: service,
		}, true
	}

	entry, ok := cache.entries[key]
	if !ok {
		cache.entries[key] = &serviceCacheEntry{
			service:   service,
			expiredAt: time.Now().Add(ttl),
		}
		return ServiceEvent{
			Type:    ServiceAdded,
			Service: service,
		}, true
	}

	isUpdated := !entry.service.Equal(service)
	entry.service = service
	entry.expiredAt = time.Now().Add(ttl)

	if !isUpdated {
		return ServiceEvent{}, false // nolint: exhaustruct
	}

	return ServiceEvent{
		Type:    ServiceUpdated,
		Service: service,
	}, true
}

// RemoveExpired removes the expired services, and returns their events.
func (cache *serviceCache) RemoveExpired(now time.Time) []ServiceEvent {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	events := []ServiceEvent{}
	for key, entry := range cache.entries {
		if now.Before(entry.expiredAt) {
			continue
		}
		delete(cache.entries, key)
		events = append(events, ServiceEvent{
			Type:    ServiceRemoved,
			Service: entry.service,
		})
	}

	return events
}

// Services returns the cached services.
func (cache *serviceCache) Services() []Service {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	services := make([]Service, 0, len(cache.entries))
	for _, entry := range cache.entries {
		services = append(services, entry.service)
	}

	return services
}
