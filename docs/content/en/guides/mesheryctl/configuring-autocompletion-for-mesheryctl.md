---
title: Configuring Autocompletion for `mesheryctl`
categories: [mesheryctl]
description: Bash, Zsh, Oh My Zsh, and fish autocompletion for `mesheryctl` commands.
---

`mesheryctl completion` prints a completion script for bash, zsh, or fish. Load that script in your shell so that pressing <kbd>Tab</kbd> completes `mesheryctl` commands and flags.

## Bash

Check your bash version with `bash --version`, then use the matching commands.

### Linux, bash 4.0 or later

Load completion in the current shell:

```bash
source <(mesheryctl completion bash)
```

To load it in every new shell, add that line to `~/.bashrc`.

### Linux, bash 3.2 or earlier

```bash
source /dev/stdin <<< "$(mesheryctl completion bash)"
```

### macOS

macOS ships bash 3.2. Install the `bash-completion` package that matches your bash version with Homebrew, then write the script to its completion directory:

```bash
# bash 3.2: needs bash-completion 1.3 or later
brew install bash-completion
# bash 4.2 or later: use bash-completion@2 instead
# brew install bash-completion@2

mesheryctl completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl
```

Follow the "Caveats" that `brew install` prints to source `bash_completion` from your `~/.bash_profile`.

## Zsh

1. If completion is not already turned on in zsh, turn it on in the current shell and in every new one:

    ```bash
    autoload -U compinit; compinit
    echo "autoload -U compinit; compinit" >> ~/.zshrc
    ```

2. Load `mesheryctl` completion in the current shell:

    ```bash
    source <(mesheryctl completion zsh)
    ```

3. Start a new shell for the change to take effect.

### Zsh on macOS and Oh My Zsh

Write the script into a completion directory of your own and add that directory to your `fpath`, so that zsh loads it in every session:

```bash
mkdir -p ~/.zsh/completions
mesheryctl completion zsh > ~/.zsh/completions/_mesheryctl
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
```

The `fpath` line must come before `compinit` in `~/.zshrc`. If you use Oh My Zsh, add it before the line that sources `oh-my-zsh.sh`. Then start a new shell.

## Fish

Load completion in the current shell:

```bash
mesheryctl completion fish | source
```

To load it in every session, write it to fish's completions directory once:

```bash
mesheryctl completion fish > ~/.config/fish/completions/mesheryctl.fish
```

For every `mesheryctl` command and flag, see the [`mesheryctl` command reference]({{< ref "reference/references/mesheryctl/_index.md" >}}).

{{< related-discussions tag="mesheryctl" >}}
