---
layout: default
title: mesheryctl-system-context-delete
permalink: reference/mesheryctl/system/context/delete
redirect_from: reference/mesheryctl/system/context/delete/
type: reference
display-title: "false"
language: en
command: system
subcommand: context
---

# mesheryctl system context delete

Delete context

## Synopsis

Delete an existing context (a named Meshery deployment) from Meshery config file
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context delete [context-name] [flags]

</div>
</pre> 

## Examples

### Delete context
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context delete [context name]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help         help for delete
  -s, --set string   New context to deploy Meshery

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
