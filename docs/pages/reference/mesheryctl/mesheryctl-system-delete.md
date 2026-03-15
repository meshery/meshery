---
layout: default
title: mesheryctl-system-delete
permalink: reference/mesheryctl/system/delete
redirect_from: reference/mesheryctl/system/delete/
type: reference
display-title: "false"
language: en
command: system
subcommand: delete
---

# mesheryctl system delete

Delete Meshery containers

## Synopsis

Delete Meshery containers. This command removes all Meshery containers created by docker-compose.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system delete [flags]

</div>
</pre> 

## Examples

Delete Meshery containers
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system delete

</div>
</pre> 

Delete Meshery containers without confirmation
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system delete -y

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for delete

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
