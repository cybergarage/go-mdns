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

var typesCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "types",
	Short: "List the service types which are advertised on the link",
	Long: `List the service types which are advertised on the link.

RFC 6763 (7.2) enumerates the service types by browsing "_services._dns-sd._udp".`,
	Example: `  mdnslookup types
  mdnslookup types --timeout 10s --format json`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		domain, err := cmd.Flags().GetString(domainParamStr)
		if err != nil {
			return err
		}

		client, stop, err := startClient()
		if err != nil {
			return err
		}
		defer stop()

		query := mdns.NewQuery(
			mdns.WithQueryService(mdns.ServiceTypeEnumerationName),
			mdns.WithQueryDomain(domain),
		)

		ctx, cancel := context.WithTimeout(cmd.Context(), queryTimeout())
		defer cancel()

		services, err := client.Query(ctx, query)
		if err != nil {
			return err
		}

		return outputServices(services)
	},
}

func init() {
	typesCmd.Flags().String(domainParamStr, mdns.DefaultQueryDomain, "domain to browse")
	rootCmd.AddCommand(typesCmd)
}
