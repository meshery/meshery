---
layout: default
title: mesheryctl-system-dashboard
permalink: reference/mesheryctl/system/dashboard
redirect_from: reference/mesheryctl/system/dashboard/
type: reference
display-title: "false"
language: en
command: system
subcommand: dashboard
---

# mesheryctl system dashboard

Open Meshery UI in browser.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system dashboard [flags]

</div>
</pre> 

## Examples

Open Meshery UI in browser
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system dashboard

</div>
</pre> 

Open Meshery UI in browser and use port-forwarding (if default port is taken already)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system dashboard --port-forward

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help           help for dashboard
      --port-forward   (optional) Use port forwarding to access Meshery UI
      --skip-browser   (optional) skip opening of MesheryUI in browser.

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
