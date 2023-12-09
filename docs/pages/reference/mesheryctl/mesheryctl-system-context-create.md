---
layout: default
title: mesheryctl-system-context-create
permalink: reference/mesheryctl/system/context/create
redirect_from: reference/mesheryctl/system/context/create/
type: reference
display-title: "false"
language: en
command: system
subcommand: context
---

# mesheryctl system context create

Create a new context (a named Meshery deployment)

## Synopsis

Add a new context to Meshery config.yaml file

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context create context-name [flags]

</div>
</pre> 

## Examples

Create new context
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context create [context-name]

</div>
</pre> 

Create new context and provide list of components, platform & URL
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context create context-name --components meshery-nsm --platform docker --url http://localhost:9081 --set --yes

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -a, --components stringArray   List of components
  -h, --help                     help for create
  -p, --platform string          Platform to deploy Meshery
  -s, --set                      Set as current context
  -u, --url string               Meshery Server URL with Port

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string    path to config file (default "/home/runner/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## Screenshots

Usage of mesheryctl context create
![context-create-usage](/assets/img/mesheryctl/newcontext.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
