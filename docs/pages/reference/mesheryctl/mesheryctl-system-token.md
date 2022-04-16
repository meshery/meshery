---
layout: default
title: mesheryctl-system-token
permalink: /reference/mesheryctl/system/token/
redirect_from: /reference/mesheryctl/system/token/
type: reference
display-title: "false"
language: en
command: system
---

# mesheryctl system token

Manage Meshery user tokens

## Synopsis


	Manipulate user tokens and their context assignments in your meshconfig

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system token [flags]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for token

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string    path to config file (default "/home/admin-pc/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

* [mesheryctl system](system/)	 - Meshery Lifecycle Management
* [mesheryctl system token create](token/create/)	 - Create a token in your meshconfig
* [mesheryctl system token delete](token/delete/)	 - Delete a token from your meshconfig
* [mesheryctl system token list](token/list/)	 - List tokens
* [mesheryctl system token set](token/set/)	 - Set token for context
* [mesheryctl system token view](token/view/)	 - View token

