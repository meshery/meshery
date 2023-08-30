---
layout: default
title: mesheryctl-system-login
permalink: reference/mesheryctl/system/login
redirect_from: reference/mesheryctl/system/login/
type: reference
display-title: "false"
language: en
command: system
subcommand: login
---

# mesheryctl system login

Authenticate to a Meshery Server

## Synopsis


Authenticate to the Local or a Remote Provider of a Meshery Server

The authentication mode is web-based browser flow

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system login [flags]

</div>
</pre> 

## Examples

Login with the Meshery Provider of your choice: the Local Provider or a Remote Provider.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system login

</div>
</pre> 

Login with the Meshery Provider by specifying it via -p or --provider flag.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system login -p Meshery

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help              help for login
  -p, --provider string   login Meshery with specified provider

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
