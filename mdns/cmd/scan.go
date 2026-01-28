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
	"fmt"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "scan",
	Short: "Scan for mDNS devices.",
	Long:  "Scan for mDNS devices.",
	RunE: func(cmd *cobra.Command, args []string) error {

		client := NewClient()
		err := client.Start()
		if err != nil {
			return err
		}
		defer client.Stop()

		msgHandler := mdns.MessageHandler(func(msg mdns.Message) {
			log.Debugf("%s\n", msg.String())
		})

		query := mdns.NewQuery(
			mdns.WithQueryService(mdns.DefaultQueryService),
			mdns.WithQueryDomain(mdns.DefaultQueryDomain),
			mdns.WithQueryMessageHandler(msgHandler),
		)

		services, err := client.Query(context.Background(), query)
		if err != nil {
			return err
		}

		for n, srv := range services {
			fmt.Printf("[%d] %s\n", n, srv.String())
			fmt.Printf("%s\n", srv.ResourceRecords().String())
		}

		return nil
	},
}
