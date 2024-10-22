---
layout: default
title: mesheryctl-design-view
permalink: reference/mesheryctl/design/view
redirect_from: reference/mesheryctl/design/view/
type: reference
display-title: "false"
language: en
command: design
subcommand: view
---

# mesheryctl design view

Display a design content

## Synopsis

Display the content of a specific design based on name or id
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design view design name [flags]

</div>
</pre> 

## Examples

view a design
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design view [design-name | ID]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -a, --all                    (optional) view all designs available
  -h, --help                   help for view
  -o, --output-format string   (optional) format to display in [json|yaml] (default "yaml")

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

Usage of mesheryctl design view
![pattern-view-usage](/assets/img/mesheryctl/patternView.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
