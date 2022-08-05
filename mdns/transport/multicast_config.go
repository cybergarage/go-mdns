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

// MulticastConfig represents a cofiguration for extended specifications.
type MulticastConfig struct {
	EachInterfaceBindingEnabled bool
}

// NewDefaultMulticastConfig returns a default configuration.
func NewDefaultMulticastConfig() *MulticastConfig {
	conf := &MulticastConfig{
		EachInterfaceBindingEnabled: true,
	}
	return conf
}

// SetConfig sets all flags.
func (conf *MulticastConfig) SetConfig(newMulticastConfig *MulticastConfig) {
	conf.EachInterfaceBindingEnabled = newMulticastConfig.EachInterfaceBindingEnabled
}

// SetEachInterfaceBindingEnabled sets a flag for binding functions.
func (conf *MulticastConfig) SetEachInterfaceBindingEnabled(flag bool) {
	conf.EachInterfaceBindingEnabled = flag
}

// IsEachInterfaceBindingEnabled returns true whether the binding functions is enabled, otherwise false.
func (conf *MulticastConfig) IsEachInterfaceBindingEnabled() bool {
	return conf.EachInterfaceBindingEnabled
}

// Equals returns true whether the specified other class is same, otherwise false.
func (conf *MulticastConfig) Equals(otherConf *MulticastConfig) bool {
	return conf.EachInterfaceBindingEnabled == otherConf.EachInterfaceBindingEnabled
}
