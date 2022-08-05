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
	"github.com/cybergarage/uecho-go/net/echonet/log"
)

const (
	logSocketTypeUDPMulticast = "UM"
	logSocketTypeUDPUnicast   = "UU"
	logSocketTypeTCPUnicast   = "TU"
	logSocketWriteFormat      = "W (%s) : %21s -> %21s (%d) : %s"
	logSocketReadFormat       = "R (%s) : %21s <- %21s (%d) : %s"
)

const (
	logSocketDirectionWrite = 0
	logSocketDirectionRead  = 1
)

func outputSocketLog(logLevel log.Level, socketType string, socketDirection int, msgFrom string, msgTo string, msg string, msgSize int) {
	switch socketDirection {
	case logSocketDirectionWrite:
		{
			log.Output(logLevel, logSocketWriteFormat, socketType, msgFrom, msgTo, msgSize, msg)
		}
	case logSocketDirectionRead:
		{
			log.Output(logLevel, logSocketReadFormat, socketType, msgTo, msgFrom, msgSize, msg)
		}
	}
}
