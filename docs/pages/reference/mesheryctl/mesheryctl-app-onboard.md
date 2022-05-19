---
layout: default
title: mesheryctl-app-onboard
permalink: reference/mesheryctl/app/onboard
redirect_from: reference/mesheryctl/app/onboard/
type: reference
display-title: "false"
language: en
command: app
subcommand: onboard
---

# mesheryctl app onboard

Onboard application

## Synopsis

Command will trigger deploy of Application file

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl app onboard [flags]

</div>
</pre> 

## Examples

Onboard application by providing file path
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl app onboard -f [filepath]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string   Path to app file
  -h, --help          help for onboard
      --skip-save     Skip saving a app

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

Go back to [command reference index](/reference/mesheryctl/) 
