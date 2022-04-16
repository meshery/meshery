---
layout: default
title: mesheryctl-system-channel
permalink: /reference/mesheryctl/system/channel/
redirect_from: /reference/mesheryctl/system/channel/
type: reference
display-title: "false"
language: en
command: system
---

# mesheryctl system channel

Switch between release channels

## Synopsis

Subscribe to a release channel. Choose between either 'stable' or 'edge' channels.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system channel [flags]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for channel

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
* [mesheryctl system channel set](channel/set/)	 - set release channel and version
* [mesheryctl system channel switch](channel/switch/)	 - switch release channel and version
* [mesheryctl system channel view](channel/view/)	 - view release channel and version

