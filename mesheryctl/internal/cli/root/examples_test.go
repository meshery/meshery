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

// exampleFlagPattern matches a flag token such as --token, --page=2, -s or
// -Xauth.json that starts a whitespace-separated word in an example line.
// For a shorthand with an attached value (-Xauth.json) only the flag letter
// (-X) is captured, so the value is not mistaken for part of the flag.
var exampleFlagPattern = regexp.MustCompile(`(?:^|\s)(--[a-zA-Z][a-zA-Z0-9-]*|-[a-zA-Z])(?:=\S*|[^\s=][^\s]*)?(?=\s|$)`)

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
				t.Errorf("example for %q does not resolve to a valid command: %v\n    %s",
					cmd.CommandPath(), err, line)
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

// TestExampleResolutionWithLeadingFlags pins the behavior of resolving an
// example command when root-level flags (e.g. --verbose) precede the
// subcommand, so TestExamplesUseKnownFlags checks flags against the right
// target command in that case too.
func TestExampleResolutionWithLeadingFlags(t *testing.T) {
	line := "mesheryctl --verbose system config gke --token auth.json"
	words := strings.Fields(strings.TrimPrefix(line, "mesheryctl "))
	target, _, err := RootCmd.Find(words)
	if err != nil {
		t.Fatalf("expected command to resolve: %v", err)
	}
	if target.Name() != "gke" {
		t.Fatalf("expected target command 'gke', got %q", target.Name())
	}
	for _, match := range exampleFlagPattern.FindAllStringSubmatch(line, -1) {
		flag := match[1]
		if !commandAcceptsFlag(target, flag) {
			t.Errorf("expected target %q to accept flag %q", target.CommandPath(), flag)
		}
	}
}
