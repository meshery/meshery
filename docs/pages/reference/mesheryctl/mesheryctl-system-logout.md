---
layout: default
title: mesheryctl-system-logout
permalink: reference/mesheryctl/system/logout
redirect_from: reference/mesheryctl/system/logout/
type: reference
display-title: "false"
language: en
command: system
subcommand: logout
---

# mesheryctl system logout

Remove authentication for Meshery Server

## Synopsis


Remove authentication for Meshery Server

This command removes the authentication token from the user's filesystem

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system logout [flags]

</div>
</pre> 

## Examples

Logout current session with your Meshery Provider.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system logout

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for logout

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
