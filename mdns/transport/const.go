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
	"time"
)

const (
	Port             = 3610
	UDPPort          = Port
	TCPPort          = Port
	MulticastAddress = "224.0.23.0"
	MaxPacketSize    = 1024
)

//  Extension for Echonet Lite
const (
	UDPPortRange = 100
)

const (
	MessageFormatSpecified = 0x01
	MessageFormatArbitrary = 0x02
)

const (
	DefaultConnectimeTimeOut = (time.Millisecond * 5000)
	DefaultRequestTimeout    = (time.Millisecond * 5000)
	DefaultBindRetryCount    = 5
	DefaultBindRetryWaitTime = (time.Millisecond * 500)
)
