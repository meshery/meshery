---
layout: default
title: mesheryctl-environment-view
permalink: reference/mesheryctl/environment/view
redirect_from: reference/mesheryctl/environment/view/
type: reference
display-title: "false"
language: en
command: environment
subcommand: view
---

# mesheryctl environment view

View registered environmnents

## Synopsis

View details of an environment registered in Meshery Server for a specific organization
Documentation for environment can be found at https://docs.meshery.io/reference/mesheryctl/environment/view
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment view [flags]

</div>
</pre> 

## Examples

View details of a specific environment
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment view --orgID [orgID]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help                   help for view
      --orgID string           Organization ID
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
