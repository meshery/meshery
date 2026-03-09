---
layout: default
title: mesheryctl-relationship
permalink: reference/mesheryctl/relationship
redirect_from: reference/mesheryctl/relationship/
type: reference
display-title: "false"
language: en
command: relationship
subcommand: nil
---

# mesheryctl relationship

Manage relationships

## Synopsis

Generate, list, search and view relationship(s) and detailed information
Meshery uses relationships to define how interconnected components interact.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl relationship [flags]

</div>
</pre> 

## Examples

Display number of available relationships in Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl relationship --count

</div>
</pre> 

Generate a relationship documentation 
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl relationship generate [flags]

</div>
</pre> 

List available relationship(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl relationship list [flags]

</div>
</pre> 

Search for a specific relationship
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl relationship search [--kind <kind>] [--type <type>] [--subtype <subtype>] [--model <model>]

</div>
</pre> 

View a specific relationship
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl relationship view [model-name]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --count   (optional) Get the number of relationship(s) in total
  -h, --help    help for relationship

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

* [mesheryctl relationship generate](/reference/mesheryctl/relationship/generate)
* [mesheryctl relationship list](/reference/mesheryctl/relationship/list)
* [mesheryctl relationship search](/reference/mesheryctl/relationship/search)
* [mesheryctl relationship view](/reference/mesheryctl/relationship/view)

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
