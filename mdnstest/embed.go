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

package mdnstest

import (
	_ "embed"
)

//go:embed dumps/google-cast-01.dump
var googlecast01 string

//go:embed dumps/google-cast-02.dump
var googlecast02 string

//go:embed dumps/google-cast-03.dump
var googlecast03 string

// 4.3.1.13. Examples
// dns-sd -R DD200C20D25AE5F7 _matterc._udp,_S3,_L840,_CM . 11111 D=840 CM=2
//
//go:embed dumps/matter-spec-120-4.3.1.13-dns-sd-01.dump
var matterSpec12043113DNSSD string

// 4.3.1.13. Examples
// avahi-publish-service --subtype=_S3._sub._matterc._udp --subtype=_L840._sub._matterc._udp DD200C20D25AE5F7 --subtype=_CM._sub._matterc._udp _matterc._udp 11111 D=840 CM=2
//
//go:embed dumps/matter-spec-120-4.3.1.13-avahi-01.dump
var matterSpec12043113Avahi01 string

//go:embed dumps/matter-spec-120-4.3.1.13-avahi-02.dump
var matterSpec12043113Avahi02 string

//go:embed dumps/matter-service-01.dump
var matterService01 string

//go:embed dumps/matter-service-02.dump
var matterService02 string

//go:embed dumps/chip-tool-query-01.dump
var chipToolQuery01 string
