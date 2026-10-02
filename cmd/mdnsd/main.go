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
	mdnsd publishes a DNS-SD service on the link: it announces the service,
	answers the queries for it, and withdraws it when it is stopped. With
	-v it also prints the queries it receives.

	Use mdnslookup to browse and resolve the services which the other
	responders advertise.

	OPTIONS
	-name string       service instance name
	-service string    service type, such as _http._tcp
	-port int          service port
	-host string       host name without the domain (default: this host)
	-subtype string    subtype label, such as _printer; can be repeated
	-txt string        TXT string, such as key=value; can be repeated
	-v                 print the received queries

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
	"strings"
	"syscall"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns"
)

func main() {
	verbose := flag.Bool("v", false, "Print the received queries")
	name := flag.String("name", "", "Service instance name")
	service := flag.String("service", "", "Service type, such as _http._tcp")
	port := flag.Int("port", 0, "Service port")
	host := flag.String("host", "", "Host name without the domain (default: this host)")
	var subtypes, txt []string
	flag.Func("subtype", "Subtype label, such as _printer; can be repeated", func(v string) error {
		subtypes = append(subtypes, v)
		return nil
	})
	flag.Func("txt", "TXT string, such as key=value; can be repeated", func(v string) error {
		txt = append(txt, v)
		return nil
	})
	flag.Parse()

	if *verbose {
		log.SetSharedLogger(log.NewStdoutLogger(log.LevelTrace))
	}

	server := NewServer()

	if *verbose {
		server.RegisterMessageHandler(server.MessageReceived)
	}

	if *service != "" {
		svc := &mdns.LocalService{
			Instance: *name,
			Service:  *service,
			Subtypes: subtypes,
			Host:     *host,
			Port:     *port,
			TXT:      txt,
		}
		if svc.Host == "" {
			svc.Host = localHostName()
		}
		if err := server.Register(svc); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Publishing %s on %s:%d\n", svc.FullName(), svc.HostName(), svc.Port)
	}

	if err := server.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	<-ctx.Done()
	// The default handling returns, so a second signal ends a stop which
	// hangs.
	stop()

	if err := server.Stop(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// localHostName returns the first label of this host's name.
func localHostName() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return "mdnsd"
	}
	return strings.Split(host, ".")[0]
}
