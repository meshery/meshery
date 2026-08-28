---
title: mesheryctl-component
display_title: false
command: component
subcommand: nil
---

# mesheryctl component

Manage Meshery components

## Synopsis

List, search and view component(s) and detailed informations

<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl component [flags]</code>
</div>
</pre> 

## Examples

Display number of available components in Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl component --count</code>
</div>
</pre> 

List available component(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl component list</code>
</div>
</pre> 

Search for component(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl component search [component-name]</code>
</div>
</pre> 

View a specific component
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl component view [component-name | component-id]</code>
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
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
