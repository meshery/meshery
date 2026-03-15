---
layout: default
title: mesheryctl-design-deploy
permalink: reference/mesheryctl/design/deploy
redirect_from: reference/mesheryctl/design/deploy/
type: reference
display-title: "false"
language: en
command: design
subcommand: deploy
---

# mesheryctl design deploy

Deploy design

## Synopsis

Command will trigger deploy of design
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design deploy [flags]

</div>
</pre> 

## Examples

Deploy design by providing file path
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design deploy -f [filepath] -s [source type]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string          Path to design file
  -h, --help                 help for deploy
      --skip-save            Skip saving a design
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

Usage of mesheryctl design deploy
![pattern-onboard-usage](/assets/img/mesheryctl/pattern-onboard.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
