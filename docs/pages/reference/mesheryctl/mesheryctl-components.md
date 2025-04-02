---
layout: default
title: mesheryctl-components
permalink: reference/mesheryctl/components
redirect_from: reference/mesheryctl/components/
type: reference
display-title: "false"
language: en
command: components
subcommand: nil
---

# mesheryctl components

Manage components

## Synopsis

List, search and view component(s) and detailed informations
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl components [flags]

</div>
</pre> 

## Examples

Display number of available components in Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl components --count

</div>
</pre> 

List available component(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl components list

</div>
</pre> 

Search for component(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl component search [component-name]

</div>
</pre> 

View a specific component
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl component view [component-name]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --count   (optional) Get the number of components in total
  -h, --help    help for components

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

* [mesheryctl components list](/reference/mesheryctl/components/list)
* [mesheryctl components search](/reference/mesheryctl/components/search)
* [mesheryctl components view](/reference/mesheryctl/components/view)

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
