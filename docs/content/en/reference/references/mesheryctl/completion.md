---
title: mesheryctl-completion
display_title: false
command: completion
subcommand: nil
---

# mesheryctl completion

Generate shell completion scripts

## Synopsis

Output shell completion code
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl completion [bash|zsh|fish]

</div>
</div>
</pre> 

## Examples

### bash <= 3.2
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
source /dev/stdin &lt;&lt;&lt; "$(mesheryctl completion bash)"

</div>
</div>
</pre> 

bash <= 3.2 on osx
ensure you have bash-completion 1.3+
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
brew install bash-completion 

</div>
</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl completion bash &gt; $(brew --prefix)/etc/bash_completion.d/mesheryctl

</div>
</div>
</pre> 

### bash >= 4.0
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
source &lt;(mesheryctl completion bash)

</div>
</div>
</pre> 

bash >= 4.0 on osx
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
brew install bash-completion@2

</div>
</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl completion bash &gt; $(brew --prefix)/etc/bash_completion.d/mesheryctl

</div>
</div>
</pre> 

### zsh
If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:
Might need to start a new shell for this setup to take effect.
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
$ echo "autoload -U compinit; compinit" &gt;&gt; ~/.zshrc

</div>
</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
source &lt;(mesheryctl completion zsh)

</div>
</div>
</pre> 

zsh on osx / oh-my-zsh
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
COMPLETION_DIR=$(echo $fpath | grep -o '[^ ]*completions' | grep -v cache) &amp;&amp; mkdir -p $COMPLETION_DIR &amp;&amp; mesheryctl completion zsh &gt; "${COMPLETION_DIR}/_mesheryctl"

</div>
</div>
</pre> 

### fish:
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl completion fish | source

</div>
</div>
</pre> 

To load fish shell completions for each session, execute once:
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl completion fish &gt; ~/.config/fish/completions/mesheryctl.fish

</div>
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
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
