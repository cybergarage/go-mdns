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
mdnsd is a Multicast DNS responder.

	NAME
	mdnsd

	SYNOPSIS
	mdnsd [OPTIONS]

	DESCRIPTION
	mdnsd is under development. The responder side of go-mdns is not
	implemented yet, so the command only listens for the mDNS messages on
	the link. It registers no service, and it answers no query.

	Use mdnslookup to browse and resolve the services which the other
	responders advertise.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cybergarage/go-logger/log"
)

func main() {
	verbose := flag.Bool("v", false, "Enable verbose output")
	flag.Parse()

	if *verbose {
		log.SetSharedLogger(log.NewStdoutLogger(log.LevelTrace))
	}

	fmt.Fprintln(os.Stderr, "mdnsd is under development: it registers no service and answers no query.")
	fmt.Fprintln(os.Stderr, "Use mdnslookup to browse and resolve the services on the link.")

	server := NewServer()

	if *verbose {
		server.RegisterMessageHandler(server.MessageReceived)
	}

	if err := server.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	if err := server.Stop(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
