---
layout: default
title: mesheryctl-system-token-view
permalink: reference/mesheryctl/system/token/view
redirect_from: reference/mesheryctl/system/token/view/
type: reference
display-title: "false"
language: en
command: system
subcommand: token
---

# mesheryctl system token view

View token

## Synopsis

View a specific token in meshery config

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system token view [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system token view [token-name]

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system token view (show token of current context)

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
      --config string    path to config file (default "/home/runner/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/) 
