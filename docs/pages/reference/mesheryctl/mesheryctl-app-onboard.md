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

Command will trigger deploy of application
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl app onboard [flags]

</div>
</pre> 

## Examples

Onboard application by providing file path
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl app onboard -f [filepath] -s [source type]

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl app onboard -f ./application.yml -s "Kubernetes Manifest"

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string          Path to app file
  -h, --help                 help for onboard
      --skip-save            Skip saving a app
  -s, --source-type string   Type of source file (ex. manifest / compose / helm)

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

Usage of mesheryctl app onboard
![app-onboard-usage](/assets/img/mesheryctl/app-onboard.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
