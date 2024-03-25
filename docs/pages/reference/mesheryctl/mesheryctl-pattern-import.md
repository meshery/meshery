---
layout: default
title: mesheryctl-pattern-import
permalink: reference/mesheryctl/pattern/import
redirect_from: reference/mesheryctl/pattern/import/
type: reference
display-title: "false"
language: en
command: pattern
subcommand: import
---

# mesheryctl pattern import

Import pattern manifests

## Synopsis

Import the pattern manifest into Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl pattern import [flags]

</div>
</pre> 

## Examples

Import pattern manifest
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl pattern import -f [file/URL] -s [source-type]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string          Path/URL to pattern file
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
