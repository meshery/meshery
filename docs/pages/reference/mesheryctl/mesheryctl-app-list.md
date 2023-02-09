---
layout: default
title: mesheryctl-app-list
permalink: reference/mesheryctl/app/list
redirect_from: reference/mesheryctl/app/list/
type: reference
display-title: "false"
language: en
command: app
subcommand: list
---

# mesheryctl app list

List applications

## Synopsis

Display list of all available applications.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl app list [flags]

</div>
</pre> 

## Examples

List all the applications
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl app list

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help      help for list
  -v, --verbose   Display full length user and app file identifiers

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -t, --token string    Path to token file default from current context

</div>
</pre>

## Screenshots

Usage of mesheryctl app list
![app-list-usage](/assets/img/mesheryctl/app-list.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/) 
