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
	"context"

	"github.com/cybergarage/go-mdns/mdns"
	"github.com/spf13/cobra"
)

var resolveCmd = &cobra.Command{ // nolint:exhaustruct,exhaustruct_v5
	Use:   "resolve <instance>",
	Short: "Resolve a service instance to its host, port and attributes",
	Long: `Resolve a service instance name to its host, port, addresses and TXT
attributes.

RFC 6763 (5) resolves an instance by querying its SRV and TXT records, and the
address records of the SRV target are queried when the responder does not send
them with the answer.`,
	Example: `  mdnslookup resolve DD200C20D25AE5F7._matterc._udp.local
  mdnslookup resolve --format json 95BDD2C0BDB33593._matterc._udp.local`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, stop, err := startClient()
		if err != nil {
			return err
		}
		defer stop()

		ctx, cancel := context.WithTimeout(cmd.Context(), queryTimeout())
		defer cancel()

		service, err := client.Resolve(ctx, args[0])
		if err != nil {
			return err
		}

		return outputServices([]mdns.Service{service})
	},
}

func init() {
	rootCmd.AddCommand(resolveCmd)
}
