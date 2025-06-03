---
layout: default
title: mesheryctl-component-view
permalink: reference/mesheryctl/component/view
redirect_from: reference/mesheryctl/component/view/
type: reference
display-title: "false"
language: en
command: component
subcommand: view
---

# mesheryctl component view

View registered components

## Synopsis

View a component registered in Meshery Server
Documentation for components can be found at https://docs.meshery.io/reference/mesheryctl/component/view
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl component view [flags]

</div>
</pre> 

## Examples

View details of a specific component
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl component view [component-name]

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
      --config string   path to config file (default "/home/n2/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
