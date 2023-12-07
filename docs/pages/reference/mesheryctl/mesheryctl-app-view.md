---
layout: default
title: mesheryctl-app-view
permalink: reference/mesheryctl/app/view
redirect_from: reference/mesheryctl/app/view/
type: reference
display-title: "false"
language: en
command: app
subcommand: view
---

# mesheryctl app view

Display application(s)

## Synopsis

Displays the contents of a specific application based on name or id

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl app view application name [flags]

</div>
</pre> 

## Examples

View applictaions with name
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl app view [app-name]

</div>
</pre> 

View applications with id
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl app view [app-id]

</div>
</pre> 

View all applications
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl app view --all

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -a, --all                    (optional) view all applications available
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

Usage of mesheryctl app view
![app-view-usage](/assets/img/mesheryctl/app-view.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
