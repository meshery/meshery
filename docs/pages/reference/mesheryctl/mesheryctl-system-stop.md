---
layout: default
title: mesheryctl-system-stop
permalink: reference/mesheryctl/system/stop
redirect_from: reference/mesheryctl/system/stop/
type: reference
display-title: "false"
language: en
command: system
subcommand: stop
---

# mesheryctl system stop

Stop Meshery

## Synopsis

Stop all Meshery containers / remove all Meshery resources.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system stop [flags]

</div>
</pre> 

## Examples

Stop Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system stop

</div>
</pre> 

Reset Meshery's configuration file to default settings.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system stop --reset

</div>
</pre> 

(optional) keep the Meshery namespace during uninstallation
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system stop --keep-namespace

</div>
</pre> 

Stop Meshery forcefully (use it when system stop doesn't work)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system stop --force

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --force            (optional) uninstall Meshery resources forcefully
  -h, --help             help for stop
      --keep-namespace   (optional) keep the Meshery namespace during uninstallation
      --reset            (optional) reset Meshery's configuration file to default settings.

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
