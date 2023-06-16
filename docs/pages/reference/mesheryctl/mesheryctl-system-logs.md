---
layout: default
title: mesheryctl-system-logs
permalink: reference/mesheryctl/system/logs
redirect_from: reference/mesheryctl/system/logs/
type: reference
display-title: "false"
language: en
command: system
subcommand: logs
---

# mesheryctl system logs

Print logs

## Synopsis

Print history of Meshery's logs and begin tailing them.

It also shows the logs of a specific component.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system logs [flags]

</div>
</pre> 

## Examples

Starts tailing Meshery server debug logs (works with components also)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system logs --verbose

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system logs meshery-istio

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --follow   (Optional) Follow the stream of the Meshery's logs. Defaults to false.
  -h, --help     help for logs

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
