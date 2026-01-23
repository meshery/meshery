---
layout: default
title: mesheryctl-connection-view
permalink: reference/mesheryctl/connection/view
redirect_from: reference/mesheryctl/connection/view/
type: reference
display-title: "false"
language: en
command: connection
subcommand: view
---

# mesheryctl connection view

View a connection

## Synopsis

View a connection by its ID or name.
Documentation for viewing connection can be found at https://docs.meshery.io/reference/mesheryctl/connection/view
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection view [flags]

</div>
</pre> 

## Examples

View details of a specific connection in default format (yaml)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection view [connection-name|connection-id]

</div>
</pre> 

View details of a specific connection in JSON format
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection view [connection-name|connection-id] --output-format json

</div>
</pre> 

View details of a specific connection in json format and save it to a file
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection view [connection-name|connection-id] --output-format json --save

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
