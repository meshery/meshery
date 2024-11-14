---
layout: default
title: mesheryctl-design-list
permalink: reference/mesheryctl/design/list
redirect_from: reference/mesheryctl/design/list/
type: reference
display-title: "false"
language: en
command: design
subcommand: list
---

# mesheryctl design list

List designs

## Synopsis

Display list of all available design files.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design list [flags]

</div>
</pre> 

## Examples

list all available designs
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design list

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help       help for list
  -p, --page int   (optional) List next set of designs with --page (default = 1) (default 1)
  -v, --verbose    Display full length user and design file identifiers

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

Usage of mesheryctl design list
![pattern-list-usage](/assets/img/mesheryctl/patternList.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
