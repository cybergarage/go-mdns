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
	"os"
	"os/signal"
	"syscall"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns/dns"
	"github.com/spf13/cobra"
)

const (
	hexParamStr       = "hex"
	queriesParamStr   = "queries"
	responsesParamStr = "responses"
)

var monitorCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "monitor",
	Short: "Watch the mDNS messages on the link",
	Long: `Listen on the mDNS multicast address, and print every message which is sent
on the link. The monitor runs until it is interrupted.`,
	Example: `  mdnslookup monitor
  mdnslookup monitor --responses --format json
  mdnslookup monitor --hex`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		hexEnabled, err := cmd.Flags().GetBool(hexParamStr)
		if err != nil {
			return err
		}
		queriesOnly, err := cmd.Flags().GetBool(queriesParamStr)
		if err != nil {
			return err
		}
		responsesOnly, err := cmd.Flags().GetBool(responsesParamStr)
		if err != nil {
			return err
		}

		client, stop, err := startClient()
		if err != nil {
			return err
		}
		defer stop()

		client.RegisterMessageHandler(func(msg dns.Message) {
			if queriesOnly && !msg.IsQuery() {
				return
			}
			if responsesOnly && !msg.IsResponse() {
				return
			}
			if err := outputMessage(msg); err != nil {
				log.Error(err)
			}
			if hexEnabled {
				log.HexInfo(msg.Bytes())
			}
		})

		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		<-ctx.Done()

		return nil
	},
}

func init() {
	monitorCmd.Flags().Bool(hexParamStr, false, "print the raw bytes of the messages")
	monitorCmd.Flags().Bool(queriesParamStr, false, "print only the query messages")
	monitorCmd.Flags().Bool(responsesParamStr, false, "print only the response messages")
	rootCmd.AddCommand(monitorCmd)
}
