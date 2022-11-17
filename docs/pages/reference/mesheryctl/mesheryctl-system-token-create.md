---
layout: default
title: mesheryctl-system-token-create
permalink: reference/mesheryctl/system/token/create
redirect_from: reference/mesheryctl/system/token/create/
type: reference
display-title: "false"
language: en
command: system
subcommand: token
---

# mesheryctl system token create

Create a token in your meshconfig

## Synopsis

Create the token with provided token name (optionally token path) to your meshconfig tokens.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system token create [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system token create [token-name] -f [token-path]

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system token create [token-name] (default path is auth.json)

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system token create [token-name] -f [token-path] --set

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --filepath string   Add the token location
  -h, --help              help for create
  -s, --set               Set as current token

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
