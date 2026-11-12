# mesheryctl-system-provider

LLMS index: [llms.txt](/docs-preview/pr-preview/pr-23341/llms.txt)

---

# mesheryctl system provider

Switch between providers

## Synopsis

Enforce a provider. Choose between available Meshery providers
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system provider [flags]

</div>
</div>
</pre> 

## Examples

To view provider
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system provider view

</div>
</div>
</pre> 

To list all available providers
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system provider list

</div>
</div>
</pre> 

To set a provider
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system provider set [provider]

</div>
</div>
</pre> 

To switch provider and redeploy Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system provider switch [provider]

</div>
</div>
</pre> 

To clear the configured provider
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system provider reset

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for provider

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

Go back to [command reference index](/docs-preview/pr-preview/pr-23341/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/docs-preview/pr-preview/pr-23341/project/contributing/cli/cli/#documenting) for guidance.

---

Section pages:

- [mesheryctl-system-provider-list](/docs-preview/pr-preview/pr-23341/reference/mesheryctl/system/provider/list/)
- [mesheryctl-system-provider-reset](/docs-preview/pr-preview/pr-23341/reference/mesheryctl/system/provider/reset/)
- [mesheryctl-system-provider-set](/docs-preview/pr-preview/pr-23341/reference/mesheryctl/system/provider/set/)
- [mesheryctl-system-provider-switch](/docs-preview/pr-preview/pr-23341/reference/mesheryctl/system/provider/switch/)
- [mesheryctl-system-provider-view](/docs-preview/pr-preview/pr-23341/reference/mesheryctl/system/provider/view/)
