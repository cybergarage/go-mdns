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
)

// serviceSet represents a service array.
// The service array is updated from the message handler goroutines, and it is
// read from the caller goroutines. Thus all accessors are synchronized.
type serviceSet struct {
	mutex    sync.RWMutex
	services []Service
}

// newServiceSet returns a blank service array.
func newServiceSet() *serviceSet {
	return &serviceSet{
		mutex:    sync.RWMutex{},
		services: []Service{},
	}
}

// Services returns the service array.
func (services *serviceSet) Services() []Service {
	services.mutex.RLock()
	defer services.mutex.RUnlock()
	copiedServices := make([]Service, len(services.services))
	copy(copiedServices, services.services)
	return copiedServices
}

// HasService returns true if the specified service is included in the service array.
func (services *serviceSet) HasService(targetService Service) bool {
	services.mutex.RLock()
	defer services.mutex.RUnlock()
	return services.hasService(targetService)
}

func (services *serviceSet) hasService(targetService Service) bool {
	for _, service := range services.services {
		if service.Equal(targetService) {
			return true
		}
	}
	return false
}

// AddService adds the specified service into the service array.
func (services *serviceSet) AddService(service Service) bool {
	services.mutex.Lock()
	defer services.mutex.Unlock()
	if services.hasService(service) {
		return false
	}
	services.services = append(services.services, service)
	return true
}

// AddServices adds the specified services into the service array.
func (services *serviceSet) AddServices(newServiceSet []Service) int {
	addedCount := 0
	for _, service := range newServiceSet {
		if services.AddService(service) {
			addedCount++
		}
	}
	return addedCount
}

// Clear removes all services from the service array.
func (services *serviceSet) Clear() {
	services.mutex.Lock()
	defer services.mutex.Unlock()
	services.services = []Service{}
}
