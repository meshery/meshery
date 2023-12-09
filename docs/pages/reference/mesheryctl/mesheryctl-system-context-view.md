---
layout: default
title: mesheryctl-system-context-view
permalink: reference/mesheryctl/system/context/view
redirect_from: reference/mesheryctl/system/context/view/
type: reference
display-title: "false"
language: en
command: system
subcommand: context
---

# mesheryctl system context view

view current context

## Synopsis

Display active Meshery context

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context view [context-name | --context context-name| --all] --flags [flags]

</div>
</pre> 

## Examples

View default context
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context view

</div>
</pre> 

View specified context
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context view context-name

</div>
</pre> 

View specified context with context flag
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context view --context context-name

</div>
</pre> 

View config of all contexts
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context view --all

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --all              Show configs for all of the context
  -c, --context string   Show config for the context
  -h, --help             help for view

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -v, --verbose         verbose output
  -y, --yes             (optional) assume yes for user interactive prompts.

</div>
</pre>

## Screenshots

Usage of mesheryctl context view
![context-view-usage](/assets/img/mesheryctl/context-view.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
