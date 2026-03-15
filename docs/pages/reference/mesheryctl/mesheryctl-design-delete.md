---
layout: default
title: mesheryctl-design-delete
permalink: reference/mesheryctl/design/delete
redirect_from: reference/mesheryctl/design/delete/
type: reference
display-title: "false"
language: en
command: design
subcommand: delete
---

# mesheryctl design delete

Delete design file

## Synopsis

delete design file will trigger deletion of the design file
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design delete [flags]

</div>
</pre> 

## Examples

delete a design file
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design delete [file | URL]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string   Path to design file
  -h, --help          help for delete

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
