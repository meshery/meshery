---
layout: default
title: mesheryctl-pattern-list
permalink: reference/mesheryctl/pattern/list
redirect_from: reference/mesheryctl/pattern/list/
type: reference
display-title: "false"
language: en
command: pattern
subcommand: list
---

# mesheryctl pattern list

List patterns

## Synopsis

Display list of all available pattern files.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl pattern list [flags]

</div>
</pre> 

## Examples

list all available patterns
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl pattern list

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help       help for list
  -p, --page int   (optional) List next set of patterns with --page (default = 1) (default 1)
  -v, --verbose    Display full length user and pattern file identifiers

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -t, --token string    Path to token file default from current context

</div>
</pre>

## Screenshots

Usage of mesheryctl pattern list
![pattern-list-usage](/assets/img/mesheryctl/patternList.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
