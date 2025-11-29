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

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cybergarage/go-logger/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(dumpCmd)
}

var dumpCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "dump",
	Short: "Dump mDNS messages.",
	Long:  "Dump mDNS messages.",
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose := viper.GetBool(VerboseParamStr)
		if verbose {
			debug := viper.GetBool(DebugParamStr)
			enableStdoutVerbose(true, debug)
		}

		client := NewClient()
		client.SetListener(client)

		err := client.Start()
		if err != nil {
			return err
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
		if code != 0 {
			return fmt.Errorf("dump failed with code (%d)", code)
		}

		return nil
	},
}
