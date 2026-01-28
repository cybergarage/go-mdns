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
mdns-server is a generic server for mDNS protocol.

	NAME
	mdns-server

	SYNOPSIS
	mdns-server [OPTIONS]

	mdns-serveru is a generic server for mDNS protocol.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
package main

import (
	"flag"
	"time"

	"github.com/cybergarage/go-logger/log"
)

func main() {
	verbose := flag.Bool("v", false, "Enable verbose output")
	flag.Parse()

	// Setup logger

	if *verbose {
		log.SetSharedLogger(log.NewStdoutLogger(log.LevelTrace))
	}

	// Start a controller for Echonet Lite node

	server := NewServer()

	if *verbose {
		server.RegisterHandler(server.MessageReceived)
	}

	err := server.Start()
	if err != nil {
		return
	}

	// Wait node responses in the local network

	time.Sleep(time.Second * 1)

	// Output all found nodes

	// Stop the controller

	err = server.Stop()
	if err != nil {
		return
	}
}
