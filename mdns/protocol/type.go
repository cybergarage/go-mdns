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

package protocol

type Type uint

const (
	A     Type = 0x0001
	NS    Type = 0x0002
	CNAME Type = 0x0005
	SRV   Type = 0x0021
	OPT   Type = 0x0029
	PTR   Type = 0x000C
	HINFO Type = 0x000D
	MX    Type = 0x000F
	AXFR  Type = 0x00FC
	ANY   Type = 0x00FF
)
