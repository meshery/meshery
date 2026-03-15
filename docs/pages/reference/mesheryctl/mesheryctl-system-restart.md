---
layout: default
title: mesheryctl-system-restart
permalink: reference/mesheryctl/system/restart
redirect_from: reference/mesheryctl/system/restart/
type: reference
display-title: "false"
language: en
command: system
subcommand: restart
---

# mesheryctl system restart

Stop, then start Meshery

## Synopsis

Restart all Meshery containers / pods.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system restart [flags]

</div>
</pre> 

## Examples

Restart all Meshery containers, their instances and their connected volumes
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system restart

</div>
</pre> 

(optional) skip checking for new updates available in Meshery.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system restart --skip-update

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help              help for restart
      --provider string   Provider to use with the Meshery server
      --skip-update       (optional) skip checking for new Meshery's container images.

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

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
