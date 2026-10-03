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
	"reflect"
)

// ExtensionConfig represents a cofiguration for extended specifications.
type ExtensionConfig struct {
	autoPortBindingEnabled bool
}

// NewDefaultExtensionConfig returns a default configuration.
func NewDefaultExtensionConfig() *ExtensionConfig {
	conf := &ExtensionConfig{
		autoPortBindingEnabled: false,
	}
	return conf
}

// SetConfig sets all flags.
func (conf *ExtensionConfig) SetConfig(newConfig *ExtensionConfig) {
	conf.autoPortBindingEnabled = newConfig.autoPortBindingEnabled
}

// SetAutoPortBindingEnabled sets a flag for TCP functions.
func (conf *ExtensionConfig) SetAutoPortBindingEnabled(flag bool) {
	conf.autoPortBindingEnabled = flag
}

// AutoPortBindingEnabled returns true whether the TCP function is enabled, otherwise false.
func (conf *ExtensionConfig) AutoPortBindingEnabled() bool {
	return conf.autoPortBindingEnabled
}

// Equals returns true whether the specified other class is same, otherwise false.
func (conf *ExtensionConfig) Equals(otherConf *ExtensionConfig) bool {
	return reflect.DeepEqual(conf, otherConf)
}
