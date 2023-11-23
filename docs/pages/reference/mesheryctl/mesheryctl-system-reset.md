---
layout: default
title: mesheryctl-system-reset
permalink: reference/mesheryctl/system/reset
redirect_from: reference/mesheryctl/system/reset/
type: reference
display-title: "false"
language: en
command: system
subcommand: reset
---

# mesheryctl system reset

Reset Meshery's configuration

## Synopsis

Reset Meshery to it's default configuration.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system reset [flags]

</div>
</pre> 

## Examples

Resets meshery.yaml file with a copy from Meshery repo
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system reset

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for reset

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

## Screenshots

Usage of mesheryctl system reset
![reset-usage](/assets/img/mesheryctl/reset.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
