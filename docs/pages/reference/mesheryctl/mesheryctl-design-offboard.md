---
layout: default
title: mesheryctl-design-offboard
permalink: reference/mesheryctl/design/offboard
redirect_from: reference/mesheryctl/design/offboard/
type: reference
display-title: "false"
language: en
command: design
subcommand: offboard
---

# mesheryctl design offboard

Offboard design

## Synopsis

Offboard design will trigger undeploy of design
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design offboard [flags]

</div>
</pre> 

## Examples

Offboard design by providing file path
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design offboard -f [filepath]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string   Path to design file
  -h, --help          help for offboard

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -t, --token string    Path to token file default from current context
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
