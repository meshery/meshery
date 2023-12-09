---
layout: default
title: mesheryctl-app-import
permalink: reference/mesheryctl/app/import
redirect_from: reference/mesheryctl/app/import/
type: reference
display-title: "false"
language: en
command: app
subcommand: import
---

# mesheryctl app import

Import app manifests

## Synopsis

Import the app manifest into Meshery

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl app import [flags]

</div>
</pre> 

## Examples

Import app manifest
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl app import -f [file/URL] -s [source-type]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string          Path/URL to app file
  -h, --help                 help for import
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

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
