---
title: mesheryctl-relationship
display_title: false
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
<div class='clipboardjs'>
mesheryctl relationship [flags]

</div>
</div>
</pre> 

## Examples

Display number of available relationships in Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl relationship --count

</div>
</div>
</pre> 

Generate a relationship documentation 
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl relationship generate [flags]

</div>
</div>
</pre> 

List available relationship(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl relationship list [flags]

</div>
</div>
</pre> 

Search for a specific relationship
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl relationship search [--kind &lt;kind&gt;] [--type &lt;type&gt;] [--subtype &lt;subtype&gt;] [--model &lt;model&gt;]

</div>
</div>
</pre> 

View a specific relationship
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl relationship view [model-name]

</div>
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

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
