---
layout: default
title: mesheryctl-system-provider-reset
permalink: reference/mesheryctl/system/provider/reset
redirect_from: reference/mesheryctl/system/provider/reset/
type: reference
display-title: "false"
language: en
command: system
subcommand: provider
---

# mesheryctl system provider reset

Clear the configured provider

## Synopsis

Clear the configured provider for the current context. This allows users to select a provider on the next Meshery start. This clears the enforced provider so that users are presented with the provider selection UI on next start.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system provider reset [flags]

</div>
</pre> 

## Examples

Clear the configured provider
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system provider reset

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for reset

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

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
