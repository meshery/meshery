// Copyright 2019 The Meshery Authors
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

package cmd

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh]",
	Short: "Output shell completion code for the specified shell (bash or zsh)",
	Long:  "Output shell completion code for the specified shell (bash or zsh).",
	Example: `  # bash <= 3.2
  source /dev/stdin <<< "$(mesheryctl completion bash)"
  # bash >= 4.0
  source <(mesheryctl completion bash)
  # bash <= 3.2 on osx
  brew install bash-completion # ensure you have bash-completion 1.3+
  mesheryctl completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl
  # bash >= 4.0 on osx
  brew install bash-completion@2
  mesheryctl completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl
  # zsh
  source <(mesheryctl completion zsh)
  # zsh on osx / oh-my-zsh
  mesheryctl completion zsh > "${fpath[1]}/_mesheryctl"`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh"},
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := getCompletion(args[0], cmd.Parent())
		if err != nil {
			return err
		}

		fmt.Print(out)
		return nil
	},
}

func getCompletion(sh string, parent *cobra.Command) (string, error) {
	var err error
	var buf bytes.Buffer

	switch sh {
	case "bash":
		err = parent.GenBashCompletion(&buf)
	case "zsh":
		err = parent.GenZshCompletion(&buf)
	default:
		err = errors.New("unsupported shell type (must be bash or zsh): " + sh)
	}

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
