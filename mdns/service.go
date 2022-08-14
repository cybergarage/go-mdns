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

package mdns

import (
	"strings"

	"github.com/cybergarage/go-mdns/mdns/protocol"
)

// Service represents a SRV record.
type Service struct {
	Name   string
	Domain string
	Port   int
}

// NewService returns a new service instance.
func NewService(name, domain string, port int) *Service {
	return &Service{
		Name:   name,
		Domain: domain,
		Port:   port,
	}
}

// NewServiceWithMessage returns a new service instance.
func NewServiceWithMessage(msg *Message) (*Service, error) {
	srv := &Service{
		Name:   "",
		Domain: "",
		Port:   0,
	}

	parseResouce := func(res *protocol.ResourceRecord) {
		switch res.Type {
		case protocol.PTR:
			srv.Name = res.Name
		case protocol.SRV:
		case protocol.TXT:
		case protocol.A:
		}
	}

	for _, answer := range msg.Answers {
		parseResouce(answer)
	}

	for _, add := range msg.Additions {
		parseResouce(add)
	}

	return srv, nil
}

// String returns the string representation.
func (srv *Service) String() string {
	return strings.Join([]string{srv.Name, srv.Domain}, nameSep)
}
