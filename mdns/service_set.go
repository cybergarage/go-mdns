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

// serviceSet represents a service array.
type serviceSet struct {
	services []Service
}

// newServiceSet returns a blank service array.
func newServiceSet() *serviceSet {
	return &serviceSet{
		services: []Service{},
	}
}

// Services returns the sercice array.
func (services *serviceSet) Services() []Service {
	return services.services
}

func (services *serviceSet) HasService(targetService Service) bool {
	for _, service := range services.services {
		if service.Equal(targetService) {
			return true
		}
	}
	return false
}

// AddService adds the specified service into th service array.
func (services *serviceSet) AddService(service Service) bool {
	if services.HasService(service) {
		return false
	}
	services.services = append(services.services, service)
	return true
}

// AddServices adds the specified services into th service array.
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
	services.services = []Service{}
}
