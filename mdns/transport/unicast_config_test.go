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
	"testing"
	"time"
)

func TestNewDefaultUnicastConfig(t *testing.T) {
	NewDefaultUnicastConfig()
}

func TestUnicastConfigEquals(t *testing.T) {
	conf01 := NewDefaultUnicastConfig()
	conf02 := NewDefaultUnicastConfig()

	// Testing Set*()

	if !conf01.Equals(conf02) {
		t.Errorf("%v != %v", conf01, conf02)
	}

	conf01.SetTCPEnabled(true)
	conf02.SetTCPEnabled(false)
	if conf01.Equals(conf02) {
		t.Errorf("%v == %v", conf01, conf02)
	}

	conf01.SetConnectionTimeout(time.Microsecond)
	conf02.SetConnectionTimeout(time.Second)
	if conf01.Equals(conf02) {
		t.Errorf("%v == %v", conf01, conf02)
	}

	// Testing SetConfig()

	conf03 := NewDefaultUnicastConfig()
	conf03.SetConfig(conf01)
	if !conf01.Equals(conf03) {
		t.Errorf("%v != %v", conf01, conf03)
	}
	if conf02.Equals(conf03) {
		t.Errorf("%v == %v", conf01, conf02)
	}

	conf03.SetConfig(conf02)
	if !conf02.Equals(conf03) {
		t.Errorf("%v != %v", conf01, conf03)
	}
}
