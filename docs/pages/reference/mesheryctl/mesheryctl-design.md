---
layout: default
title: mesheryctl-design
permalink: reference/mesheryctl/design
redirect_from: reference/mesheryctl/design/
type: reference
display-title: "false"
language: en
command: design
subcommand: nil
---

# mesheryctl design

Cloud Native Designs Management

## Synopsis

Manage cloud and cloud native infrastructure using predefined designs.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design [flags]

</div>
</pre> 

## Examples

Apply design file:
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design apply --file [path to design file | URL of the file]

</div>
</pre> 

Delete design file:
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design delete --file [path to design file]

</div>
</pre> 

View design file:
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design view [design name | ID]

</div>
</pre> 

List all designs:
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design list

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help           help for design
  -t, --token string   Path to token file default from current context

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

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
