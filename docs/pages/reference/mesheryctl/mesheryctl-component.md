---
layout: default
title: mesheryctl-component
permalink: reference/mesheryctl/component
redirect_from: reference/mesheryctl/component/
type: reference
display-title: "false"
language: en
command: component
subcommand: nil
---

# mesheryctl component

Manage components

## Synopsis

List, search and view component(s) and detailed informations
Documentation for components can be found at https://docs.meshery.io/reference/mesheryctl/component
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl component [flags]

</div>
</pre> 

## Examples

Display number of available components in Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl component --count

</div>
</pre> 

List available component(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl component list

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
  -h, --help    help for component

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/n2/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

* [mesheryctl component list](/reference/mesheryctl/component/list)
* [mesheryctl component search](/reference/mesheryctl/component/search)
* [mesheryctl component view](/reference/mesheryctl/component/view)

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
