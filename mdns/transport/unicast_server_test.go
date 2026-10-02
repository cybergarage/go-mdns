// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"testing"

	"github.com/cybergarage/go-mdns/mdns/dns"
)

func TestNewUnicastServer(t *testing.T) {
	NewUnicastServer()
}

// TestUnicastServerSetMessageProcessorWhileReceiving sets the processor
// while the server hands it received messages; run with -race.
func TestUnicastServerSetMessageProcessorWhileReceiving(t *testing.T) {
	server := NewUnicastServer()
	processor := func(dns.Message) (dns.Message, error) { return nil, nil }
	msg := dns.NewRequestMessage()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for range 100 {
			handleUnicastUDPRequestMessage(server, msg)
		}
	}()
	for range 100 {
		server.SetMessageProcessor(processor)
	}
	<-done
}
