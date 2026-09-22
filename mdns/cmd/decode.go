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
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/cybergarage/go-logger/log/hexdump"
	"github.com/cybergarage/go-mdns/mdns/dns"
	"github.com/spf13/cobra"
)

const (
	rawParamStr = "raw"
)

var decodeCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "decode <file>",
	Short: "Decode a recorded mDNS message",
	Long: `Decode a recorded mDNS message, and print its records. The file is a hex dump
log, such as the dumps which "monitor --hex" prints, or a raw message file when
--raw is set.

The command needs no network, so it can be used to analyze a message which was
captured elsewhere.`,
	Example: `  mdnslookup decode mdnstest/dumps/matter-answer-01.dump
  mdnslookup decode --format json mdnstest/dumps/matter-query-01.dump
  mdnslookup decode --raw message.bin`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		raw, err := cmd.Flags().GetBool(rawParamStr)
		if err != nil {
			return err
		}

		fileBytes, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}

		msgBytes := fileBytes
		if !raw {
			msgBytes, err = decodeDumpBytes(fileBytes)
			if err != nil {
				return err
			}
		}

		msg, err := dns.NewMessageWithBytes(msgBytes)
		if err != nil {
			return err
		}

		return outputMessage(msg)
	},
}

// decodeDumpBytes decodes the specified file contents as a hex dump log, or as
// a hex string when the file holds no dump log.
func decodeDumpBytes(fileBytes []byte) ([]byte, error) {
	lines := strings.Split(string(fileBytes), "\n")

	msgBytes, err := hexdump.DecodeHexdumpLogs(lines)
	if err == nil && 0 < len(msgBytes) {
		return msgBytes, nil
	}

	hexStr := strings.Join(strings.Fields(string(fileBytes)), "")
	msgBytes, hexErr := hex.DecodeString(hexStr)
	if hexErr == nil && 0 < len(msgBytes) {
		return msgBytes, nil
	}

	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("no message is found")
}

func init() {
	decodeCmd.Flags().Bool(rawParamStr, false, "read the file as a raw message instead of a hex dump log")
	rootCmd.AddCommand(decodeCmd)
}
