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
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docCmd = &cobra.Command{ // nolint:exhaustruct,exhaustruct_v5
	Use:   "doc",
	Short: "Generate markdown documentation to stdout",
	Long: `Generate the markdown reference of the commands to stdout, as one page in
which each command links to the sections of the others.`,
	// The reference is generated for the repository, not for the users.
	Hidden: true,
	Args:   cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		md, err := generateMarkdown(rootCmd)
		if err != nil {
			return err
		}
		fmt.Fprint(cmd.OutOrStdout(), md)
		return nil
	},
}

// generateMarkdown returns the markdown reference of root and its available
// subcommands as one page.
func generateMarkdown(root *cobra.Command) (string, error) {
	var buf bytes.Buffer
	var appendMarkdown func(cmd *cobra.Command) error
	appendMarkdown = func(cmd *cobra.Command) error {
		if err := doc.GenMarkdownCustom(cmd, &buf, sectionLink); err != nil {
			return err
		}
		for _, c := range cmd.Commands() {
			if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
				continue
			}
			if err := appendMarkdown(c); err != nil {
				return err
			}
		}
		return nil
	}
	if err := appendMarkdown(root); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// sectionLink returns the link to the section of a command in the page,
// such as "#mdnslookup-browse" for "mdnslookup_browse.md", which cobra
// gives as the file of the command.
func sectionLink(name string) string {
	anchor := strings.TrimSuffix(name, ".md")
	anchor = strings.ReplaceAll(anchor, "_", "-")
	return "#" + strings.ToLower(anchor)
}

func init() {
	rootCmd.AddCommand(docCmd)
}
