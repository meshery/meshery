---
layout: default
title: mesheryctl-pattern-export
permalink: reference/mesheryctl/pattern/export
redirect_from: reference/mesheryctl/pattern/export/
type: reference
display-title: "false"
language: en
command: pattern
subcommand: export
---

# mesheryctl pattern export

Export a design from Meshery

## Synopsis

The 'export' command allows you to export a specific design from your Meshery server.
You can specify the design by its name or ID and optionally define the type of design.
The command also supports specifying an output directory where the exported design will be saved.
By default, the exported design will be saved in the current directory. The different types of design
type allowed are oci, original, and current. The default design type is current.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl pattern export [pattern-name | ID] [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>
	# Export a design with a specific ID

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesheryctl pattern export [pattern-name | ID]

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	# Export a design with a specific ID and type

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesheryctl pattern export [pattern-name | ID] --type [design-type]

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	# Export a design and save it to a specific directory

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesheryctl pattern export [pattern-name | ID] --output ./designs

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	# Export a design with a specific type and save it to a directory

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesheryctl pattern export [pattern-name | ID] --type [design-type] --output ./exports

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help            help for export
  -o, --output string   Specify the output directory to save the design
      --type string     Specify the design type to export

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -t, --token string    Path to token file default from current context
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
