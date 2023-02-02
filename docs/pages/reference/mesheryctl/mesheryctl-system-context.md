---
layout: default
title: mesheryctl-system-context
permalink: reference/mesheryctl/system/context
redirect_from: reference/mesheryctl/system/context/
type: reference
display-title: "false"
language: en
command: system
subcommand: context
---

# mesheryctl system context

Configure your Meshery deployment(s)

## Synopsis

Configure and switch between different named Meshery server and component versions and deployments.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context [command] [flags]

</div>
</pre> 

## Examples

Base command
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --context string   (optional) temporarily change the current context.
  -h, --help             help for context

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

Go back to [command reference index](/reference/mesheryctl/) 
