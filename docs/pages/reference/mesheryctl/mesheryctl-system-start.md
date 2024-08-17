---
layout: default
title: mesheryctl-system-start
permalink: reference/mesheryctl/system/start
redirect_from: reference/mesheryctl/system/start/
type: reference
display-title: "false"
language: en
command: system
subcommand: start
---

# mesheryctl system start

Start Meshery

## Synopsis

Start Meshery and each of its cloud native components.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system start [flags]

</div>
</pre> 

## Examples

Start meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system start

</div>
</pre> 

(optional) skip opening of MesheryUI in browser.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system start --skip-browser

</div>
</pre> 

(optional) skip checking for new updates available in Meshery.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system start --skip-update

</div>
</pre> 

Reset Meshery's configuration file to default settings.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system start --reset

</div>
</pre> 

Specify Platform to deploy Meshery to.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system start -p docker

</div>
</pre> 

Specify Provider to use.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system start --provider Meshery

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help              help for start
  -p, --platform string   platform to deploy Meshery to.
      --provider string   (optional) Defaults to the provider specified in the current context
      --reset             (optional) reset Meshery's configuration file to default settings.
      --skip-browser      (optional) skip opening of MesheryUI in browser.
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
