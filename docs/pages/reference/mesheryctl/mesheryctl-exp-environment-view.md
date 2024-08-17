---
layout: default
title: mesheryctl-exp-environment-view
permalink: reference/mesheryctl/exp/environment/view
redirect_from: reference/mesheryctl/exp/environment/view/
type: reference
display-title: "false"
language: en
command: exp
subcommand: environment
---

# mesheryctl exp environment view

view registered environmnents

## Synopsis

view a environments registered in Meshery Server
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp environment view [flags]

</div>
</pre> 

## Examples

View details of a specific environment
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp environment view --orgID [orgId]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help                   help for view
  -o, --output-format string   (optional) format to display in [json|yaml] (default "yaml")
  -s, --save                   (optional) save output as a JSON/YAML file

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
