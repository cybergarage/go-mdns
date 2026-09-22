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

package main

import (
	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns"
	"github.com/cybergarage/go-mdns/mdns/dns"
)

// Server is a Multicast DNS responder.
//
// The responder side of go-mdns is under development, so the server only
// listens for the messages on the link. It registers no service, and it answers
// no query.
type Server struct {
	*mdns.Server
}

func NewServer() *Server {
	server := &Server{
		Server: mdns.NewServer(),
	}
	return server
}

// MessageReceived is called when a message is received on the link.
func (server *Server) MessageReceived(msg dns.Message) {
	log.Infof("%s", msg.String())
}
