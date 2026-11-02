// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package root

import (
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// exampleFlagPattern matches a flag token such as --token, --page=2 or -s
// that starts a whitespace-separated word in an example line.
var exampleFlagPattern = regexp.MustCompile(`(?:^|\s)(--[a-zA-Z][a-zA-Z0-9-]*|-[a-zA-Z])(?:=|\s|$)`)

// TestExamplesUseKnownFlags makes sure every flag used in a command's Example
// text is a flag that the example's target command actually accepts. The
// examples are rendered into `--help` and into the reference docs, so a stale
// flag there sends users straight into an "unknown flag" error.
func TestExamplesUseKnownFlags(t *testing.T) {
	walkCommands(RootCmd, func(cmd *cobra.Command) {
		for _, line := range strings.Split(cmd.Example, "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "mesheryctl ") {
				continue
			}

			words := strings.Fields(strings.TrimPrefix(line, "mesheryctl "))
			target, _, err := RootCmd.Find(words)
			if err != nil {
				// The example does not resolve to a command; that is
				// outside the scope of this test.
				continue
			}

			for _, match := range exampleFlagPattern.FindAllStringSubmatch(line, -1) {
				flag := match[1]
				if !commandAcceptsFlag(target, flag) {
					t.Errorf("example for %q uses %s, which %q does not accept:\n    %s",
						cmd.CommandPath(), flag, target.CommandPath(), line)
				}
			}
		}
	})
}

// walkCommands calls fn for cmd and every command below it.
func walkCommands(cmd *cobra.Command, fn func(*cobra.Command)) {
	fn(cmd)
	for _, child := range cmd.Commands() {
		walkCommands(child, fn)
	}
}

// commandAcceptsFlag reports whether cmd accepts flag, counting its own flags,
// flags inherited from parent commands, and cobra's built-in --help/-h.
func commandAcceptsFlag(cmd *cobra.Command, flag string) bool {
	if flag == "--help" || flag == "-h" {
		return true
	}

	if strings.HasPrefix(flag, "--") {
		return cmd.Flag(strings.TrimPrefix(flag, "--")) != nil
	}

	shorthand := strings.TrimPrefix(flag, "-")
	return cmd.Flags().ShorthandLookup(shorthand) != nil ||
		cmd.InheritedFlags().ShorthandLookup(shorthand) != nil
}
