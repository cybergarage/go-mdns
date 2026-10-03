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
	mdnsd publishes a DNS-SD service on the link: it probes the instance
	name and the host name, announces the service, answers the queries for
	it, and withdraws it when it is stopped. With --verbose it also prints
	the queries it receives.

	mdnsd does not rename the service: it exits with an error when another
	node on the link holds the instance name or the host name, when it
	starts or later.

	Use mdnslookup to browse and resolve the services which the other
	responders advertise.

	OPTIONS
	--name string         service instance name
	--service string      service type, such as _http._tcp
	--port int            service port
	--host string         host name, such as myhost or myhost.local (default: this host)
	--subtype string      subtype label, such as _printer; can be repeated
	--txt string          TXT string, such as key=value; can be repeated
	--address string      address the host name resolves to; can be repeated
	                      (default: the addresses of the interface a query arrives on)
	-i, --interface name  network interface to use; can be repeated or comma separated
	                      (default: all available interfaces)
	--family string       address family to use: all|ipv4|ipv6 (default "all")
	-v, --verbose         print the received queries
	--version             print the version and exit

	A flag is given with one dash or two, such as -name or --name.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns"
)

const (
	programName = "mdnsd"

	familyAll  = "all"
	familyIPv4 = "ipv4"
	familyIPv6 = "ipv6"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var verbose, version bool
	flag.BoolVar(&verbose, "v", false, "Print the received queries")
	flag.BoolVar(&verbose, "verbose", false, "Print the received queries")
	flag.BoolVar(&version, "version", false, "Print the version and exit")
	name := flag.String("name", "", "Service instance name")
	service := flag.String("service", "", "Service type, such as _http._tcp")
	port := flag.Int("port", 0, "Service port")
	host := flag.String("host", "", "Host name, such as myhost or myhost.local (default: this host)")
	family := flag.String("family", familyAll, "Address family to use: all|ipv4|ipv6")
	var subtypes, txt, ifnames []string
	var addrs []net.IP
	flag.Func("subtype", "Subtype label, such as _printer; can be repeated", func(v string) error {
		subtypes = append(subtypes, v)
		return nil
	})
	flag.Func("txt", "TXT string, such as key=value; can be repeated", func(v string) error {
		txt = append(txt, v)
		return nil
	})
	flag.Func("address", "Address the host name resolves to; can be repeated (default: the addresses of the interface a query arrives on)", func(v string) error {
		ip := net.ParseIP(v)
		if ip == nil {
			return fmt.Errorf("invalid address: %s", v)
		}
		addrs = append(addrs, ip)
		return nil
	})
	addInterfaces := func(v string) error {
		for ifname := range strings.SplitSeq(v, ",") {
			if ifname = strings.TrimSpace(ifname); ifname != "" {
				ifnames = append(ifnames, ifname)
			}
		}
		return nil
	}
	flag.Func("i", "Network interface to use; can be repeated or comma separated (default: all available interfaces)", addInterfaces)
	flag.Func("interface", "Network interface to use; can be repeated or comma separated (default: all available interfaces)", addInterfaces)
	flag.Parse()

	if version {
		fmt.Printf("%s version %s\n", programName, mdns.Version)
		return nil
	}

	if verbose {
		log.SetSharedLogger(log.NewStdoutLogger(log.LevelTrace))
	}

	opts, err := serverOptions(ifnames, *family)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	conflicts := make(chan *mdns.ConflictError, 1)
	opts = append(opts, mdns.WithServerConflictHandler(func(err *mdns.ConflictError) {
		select {
		case conflicts <- err:
		default:
		}
	}))
	server := NewServer(opts...)

	if verbose {
		server.RegisterMessageHandler(server.MessageReceived)
	}

	if err := server.Start(); err != nil {
		return err
	}

	if *service != "" {
		svc := &mdns.LocalService{
			Instance:  *name,
			Service:   *service,
			Subtypes:  subtypes,
			Host:      *host,
			Port:      *port,
			TXT:       txt,
			Addresses: addrs,
		}
		if svc.Host == "" {
			svc.Host = localHostName()
		}
		if err := server.Register(ctx, svc); err != nil {
			_ = server.Stop()
			return err
		}
		fmt.Fprintf(os.Stderr, "Publishing %s on %s:%d\n", svc.FullName(), svc.HostName(), svc.Port)
	}

	var conflict error
	select {
	case <-ctx.Done():
	case err := <-conflicts:
		conflict = err
	}
	// The default handling returns, so a second signal ends a stop which
	// hangs.
	stop()

	if err := server.Stop(); err != nil {
		return err
	}
	return conflict
}

// localHostName returns the first label of this host's name.
func localHostName() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return "mdnsd"
	}
	return strings.Split(host, ".")[0]
}

// serverOptions returns the server options of the interfaces and the address
// family to use.
func serverOptions(ifnames []string, family string) ([]mdns.ServerOption, error) {
	opts := []mdns.ServerOption{}
	if 0 < len(ifnames) {
		ifis := make([]*net.Interface, 0, len(ifnames))
		for _, ifname := range ifnames {
			ifi, err := net.InterfaceByName(ifname)
			if err != nil {
				return nil, fmt.Errorf("interface %s: %w", ifname, err)
			}
			ifis = append(ifis, ifi)
		}
		opts = append(opts, mdns.WithServerInterfaces(ifis...))
	}
	switch family {
	case familyAll, "":
	case familyIPv4:
		opts = append(opts, mdns.WithServerIPv6Enabled(false))
	case familyIPv6:
		opts = append(opts, mdns.WithServerIPv4Enabled(false))
	default:
		return nil, fmt.Errorf("invalid address family: %s", family)
	}
	return opts, nil
}
