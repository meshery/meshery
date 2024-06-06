---
layout: default
title: mesheryctl-system-context-list
permalink: reference/mesheryctl/system/context/list
redirect_from: reference/mesheryctl/system/context/list/
type: reference
display-title: "false"
language: en
command: system
subcommand: context
---

# mesheryctl system context list

list contexts

## Synopsis

List current context and available contexts
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context list [flags]

</div>
</pre> 

## Examples

List all contexts present
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context list

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for list

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
