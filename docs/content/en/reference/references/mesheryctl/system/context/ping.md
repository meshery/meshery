---
title: mesheryctl-system-context-ping
display_title: false
command: system
subcommand: context
categories: [mesheryctl-system]
---

# mesheryctl system context ping

Check connectivity and token validity for a Meshery context

## Synopsis

Check whether the Meshery Server for a context is reachable and whether the stored authentication token is still valid.

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system context ping [context-name | --context context-name] [flags]

</div>
</div>
</pre> 

## Examples

Ping the current context
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system context ping

</div>
</div>
</pre> 

Ping a specified context
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system context ping context-name

</div>
</div>
</pre> 

Ping a specified context using the --context flag
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system context ping --context context-name

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --context string   Ping the given context
  -h, --help             help for ping

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -v, --verbose         verbose output
  -y, --yes             (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
