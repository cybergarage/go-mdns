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

// services represents a service array.
type services struct {
	services []*Service
}

// newServices returns a blank service array.
func newServices() *services {
	return &services{
		services: []*Service{},
	}
}

// Services returns the sercice array.
func (services *services) Services() []*Service {
	return services.services
}

// HasService returns true if the specified service is already added, otherwise false.
func (services *services) HasService(service *Service) bool {
	return false
}

// AddService adds the specified service into th service array.
func (services *services) AddService(service *Service) {
	services.services = append(services.services, service)
}
