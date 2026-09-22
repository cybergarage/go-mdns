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

/*
Package dns implements the DNS message encoding and decoding which the mDNS
messages are built on (RFC 1035), with the Multicast DNS extensions of RFC 6762.

The package reads the A, AAAA, PTR, SRV, TXT and NSEC records, and it resolves
the compression pointers of the names. A record of an unknown type keeps its raw
data, which [Record.Data] returns.

The API of this package is not stable until v1.0.0.
*/
package dns
