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
mdns-dump is a dumpr utility for mDNS protocol.

	NAME
	mdns-dump

	SYNOPSIS
	mdns-dump [OPTIONS]

	mdns-dump is a dumpr utility for mDNS protocol.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/cybergarage/go-logger/log"
)

func main() {
	log.SetSharedLogger(log.NewStdoutLogger(log.LevelTrace))

	client := NewClient()
	client.SetListener(client)

	err := client.Start()
	if err != nil {
		return
	}

	defer client.Stop()

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM)

	exitCh := make(chan int)

	go func() {
		for {
			s := <-sigCh
			switch s {
			case syscall.SIGINT, syscall.SIGTERM:
				if err := client.Stop(); err != nil {
					log.Errorf("couldn't be terminated (%s)", err.Error())
					os.Exit(1)
				}
				exitCh <- 0
			}
		}
	}()

	code := <-exitCh

	os.Exit(code)
}
