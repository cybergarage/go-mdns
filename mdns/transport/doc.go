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
Package transport binds the mDNS multicast and unicast sockets, and it sends and
receives the messages on them.

A server is bound to every available network interface by default.
[InterfaceSelector] selects the interfaces and the address families to bind, and
the interface which a message was received on is attached to the message, so
that an IPv6 link-local address can be scoped to it.

The API of this package is not stable until v1.0.0.
*/
package transport
