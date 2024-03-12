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

package dns

import "fmt"

// Type represents a message type.
type Type uint

const (
	A     Type = 0x0001
	NS    Type = 0x0002
	CNAME Type = 0x0005
	TXT   Type = 0x0010
	SRV   Type = 0x0021
	OPT   Type = 0x0029
	PTR   Type = 0x000C
	HINFO Type = 0x000D
	MX    Type = 0x000F
	AAAA  Type = 0x001C
	AXFR  Type = 0x00FC
	NSEC  Type = 0x002F
	ANY   Type = 0x00FF
)

// String returns the string of the type.
func (t Type) String() string {
	switch t {
	case A:
		return "A"
	case NS:
		return "NS"
	case CNAME:
		return "CNAME"
	case TXT:
		return "TXT"
	case SRV:
		return "SRV"
	case OPT:
		return "OPT"
	case PTR:
		return "PTR"
	case HINFO:
		return "HINFO"
	case MX:
		return "MX"
	case AAAA:
		return "AAAA"
	case AXFR:
		return "AXFR"
	case NSEC:
		return "NSEC"
	case ANY:
		return "ANY"
	}
	return fmt.Sprintf("%04X(%d)", uint16(t), uint16(t))
}
