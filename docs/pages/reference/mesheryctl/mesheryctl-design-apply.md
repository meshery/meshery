---
layout: default
title: mesheryctl-design-apply
permalink: reference/mesheryctl/design/apply
redirect_from: reference/mesheryctl/design/apply/
type: reference
display-title: "false"
language: en
command: design
subcommand: apply
---

# mesheryctl design apply

Apply design file

## Synopsis

Apply design will trigger deploy of the design file
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design apply [flags]

</div>
</pre> 

## Examples

apply a design file
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design apply -f [file | URL]

</div>
</pre> 

deploy a saved design
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design apply [design-name]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string   Path to design file
  -h, --help          help for apply
      --skip-save     Skip saving a design

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

## Screenshots

Usage of mesheryctl design apply
![pattern-apply-usage](/assets/img/mesheryctl/patternApply.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
