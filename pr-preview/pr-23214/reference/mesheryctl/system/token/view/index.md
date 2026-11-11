# mesheryctl-system-token-view

LLMS index: [llms.txt](/pr-preview/pr-23214/llms.txt)

---

# mesheryctl system token view

View token

## Synopsis

View a specific token in meshery config
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system token view [token-name] [flags]

</div>
</div>
</pre> 

## Examples

View a specific token
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system token view [token-name]

</div>
</div>
</pre> 

View token of current context
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system token view

</div>
</div>
</pre> 

View all available tokens
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system token view --all

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --all    set the flag to view all the tokens.
  -h, --help   help for view

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string    path to config file (default "~/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

Go back to [command reference index](/pr-preview/pr-23214/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/pr-preview/pr-23214/project/contributing/cli/cli/#documenting) for guidance.
