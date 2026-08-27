---
title: mesheryctl-system-provider
display_title: false
command: system
subcommand: provider
---

# mesheryctl system provider

Switch between providers

## Synopsis

Enforce a provider. Choose between available Meshery providers
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl system provider [flags]</code>
</div>
</pre> 

## Examples

To view provider
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl system provider view</code>
</div>
</pre> 

To list all available providers
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl system provider list</code>
</div>
</pre> 

To set a provider
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl system provider set [provider]</code>
</div>
</pre> 

To switch provider and redeploy Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl system provider switch [provider]</code>
</div>
</pre> 

To clear the configured provider
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl system provider reset</code>
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
      --config string    path to config file (default "/home/runner/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
