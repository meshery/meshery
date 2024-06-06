---
layout: default
title: mesheryctl-pattern-onboard
permalink: reference/mesheryctl/pattern/onboard
redirect_from: reference/mesheryctl/pattern/onboard/
type: reference
display-title: "false"
language: en
command: pattern
subcommand: onboard
---

# mesheryctl pattern onboard

Onboard pattern

## Synopsis

Command will trigger deploy of pattern
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl pattern onboard [flags]

</div>
</pre> 

## Examples

Onboard pattern by providing file path
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl pattern onboard -f [filepath] -s [source type]

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl pattern onboard -f ./pattern.yml -s "Kubernetes Manifest"

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string          Path to pattern file
  -h, --help                 help for onboard
      --skip-save            Skip saving a pattern
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

Usage of mesheryctl pattern onboard
![pattern-onboard-usage](/assets/img/mesheryctl/pattern-onboard.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
