---
layout: default
title: mesheryctl-system-provider
permalink: reference/mesheryctl/system/provider
redirect_from: reference/mesheryctl/system/provider/
type: reference
display-title: "false"
language: en
command: system
subcommand: provider
---

# mesheryctl system provider

Switch between providers

## Synopsis

Enforce a provider. Choose between available Meshery providers

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system provider [flags]

</div>
</pre> 

## Examples

To view provider
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system provider view

</div>
</pre> 

To list all available providers
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system provider list

</div>
</pre> 

To set a provider
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system provider set [provider]

</div>
</pre> 

To switch provider and redeploy Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system provider switch [provider]

</div>
</pre> 

To reset provider to default
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system provider reset

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

Go back to [command reference index](/reference/mesheryctl/) 
