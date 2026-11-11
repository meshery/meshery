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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type flagDoc struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Usage       string `yaml:"usage"`
	Example     string `yaml:"example"`
}

type cmdDoc struct {
	Name        string             `yaml:"name"`
	Description string             `yaml:"description"`
	Usage       string             `yaml:"usage"`
	Flags       map[string]flagDoc `yaml:"flags"`
	Subcommands map[string]cmdDoc  `yaml:"subcommands"`
}

// parseFlagTokens extracts flag identifiers (e.g. "--token", "-t") from a documented flag name string.
func parseFlagTokens(nameStr string) []string {
	var tokens []string
	for _, part := range strings.Split(nameStr, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "-") {
			fields := strings.Fields(part)
			flag := fields[0]
			if idx := strings.Index(flag, "="); idx != -1 {
				flag = flag[:idx]
			}
			tokens = append(tokens, flag)
		}
	}
	return tokens
}

// TestCmdsYmlDocumentedFlags ensures every flag documented in docs/data/mesheryctlcommands/cmds.yml
// is actually accepted by its corresponding Cobra command.
func TestCmdsYmlDocumentedFlags(t *testing.T) {
	cmdsPath := filepath.Join("..", "..", "..", "..", "docs", "data", "mesheryctlcommands", "cmds.yml")
	data, err := os.ReadFile(cmdsPath)
	if err != nil {
		t.Fatalf("failed to read cmds.yml: %v", err)
	}

	var rootMap map[string]cmdDoc
	err = yaml.Unmarshal(data, &rootMap)
	if err != nil {
		t.Fatalf("failed to unmarshal cmds.yml: %v", err)
	}

	checkFlags := func(cmdWords []string, flags map[string]flagDoc) {
		target, _, err := RootCmd.Find(cmdWords)
		if err != nil {
			t.Errorf("command %q does not resolve: %v", strings.Join(cmdWords, " "), err)
			return
		}
		for _, f := range flags {
			tokens := parseFlagTokens(f.Name)
			for _, token := range tokens {
				if !commandAcceptsFlag(target, token) {
					t.Errorf("command %q documents flag %q in cmds.yml, but the command does not accept it",
						target.CommandPath(), token)
				}
			}
		}
	}

	for key, cmd := range rootMap {
		var baseWords []string
		if key != "global" && strings.TrimSpace(cmd.Name) != "" {
			baseWords = strings.Fields(cmd.Name)
		}

		if len(cmd.Flags) > 0 {
			checkFlags(baseWords, cmd.Flags)
		}

		for _, sub := range cmd.Subcommands {
			subWords := append([]string{}, baseWords...)
			if strings.TrimSpace(sub.Name) != "" {
				subWords = append(subWords, strings.Fields(sub.Name)...)
			}
			if len(sub.Flags) > 0 {
				checkFlags(subWords, sub.Flags)
			}
		}
	}
}
