// Copyright 2020 Layer5, Inc.
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

package system

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const example = `  # bash <= 3.2
  source /dev/stdin <<< "$(mesheryctl system completion bash)"

  # bash >= 4.0
  source <(mesheryctl system completion bash)

  # bash <= 3.2 on osx
  brew install bash-completion # ensure you have bash-completion 1.3+
  mesheryctl system completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl

  # bash >= 4.0 on osx
  brew install bash-completion@2
  mesheryctl system completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl

  # zsh

  # If shell completion is not already enabled in your environment you will need
  # to enable it.  You can execute the following once:
  # Might need to start a new shell for this setup to take effect.
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  source <(mesheryctl system completion zsh)

  # zsh on osx / oh-my-zsh
  mesheryctl system completion zsh > "${fpath[1]}/_mesheryctl"

  # fish:
  mesheryctl system completion fish | source
  # To load fish shell completions for each session, execute once:
  mesheryctl system completion fish > ~/.config/fish/completions/mesheryctl.fish`

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish]",
	Short:                 "Output shell completion code",
	Long:                  "Output shell completion code",
	Example:               example,
	Args:                  cobra.ExactValidArgs(1),
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish"},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		default:
			return errors.New("shell not supported")
		}
	},
}
