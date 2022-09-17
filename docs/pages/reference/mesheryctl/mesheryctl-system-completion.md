---
layout: default
title: mesheryctl-system-completion
permalink: reference/mesheryctl/system/completion
redirect_from: reference/mesheryctl/system/completion/
type: reference
display-title: "false"
language: en
command: system
subcommand: completion
---

# mesheryctl system completion

Output shell completion code

## Synopsis

Output shell completion code

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system completion [bash|zsh|fish]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>
  # bash <= 3.2

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  source /dev/stdin <<< "$(mesheryctl system completion bash)"

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  # bash >= 4.0

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  source <(mesheryctl system completion bash)

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  # bash <= 3.2 on osx

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  brew install bash-completion # ensure you have bash-completion 1.3+

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  mesheryctl system completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  # bash >= 4.0 on osx

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  brew install bash-completion@2

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  mesheryctl system completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  # zsh

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  # If shell completion is not already enabled in your environment you will need

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  # to enable it.  You can execute the following once:

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  # Might need to start a new shell for this setup to take effect.

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  source <(mesheryctl system completion zsh)

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  # zsh on osx / oh-my-zsh

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  mesheryctl system completion zsh > "${fpath[1]}/_mesheryctl"

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  # fish:

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  mesheryctl system completion fish | source

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  # To load fish shell completions for each session, execute once:

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
  mesheryctl system completion fish > ~/.config/fish/completions/mesheryctl.fish

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for completion

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string    path to config file (default "/home/runner/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/) 
