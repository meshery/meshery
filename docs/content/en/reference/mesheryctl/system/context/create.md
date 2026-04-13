---
title: mesheryctl-system-context-create
display_title: false
command: system
subcommand: context
---

# mesheryctl system context create

Create a new context (a named Meshery deployment)

## Synopsis

Add a new context to Meshery config.yaml file.

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

Create new context and provide list of components, platform & URL and set it as current context
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context create [context-name] --components [meshery-nsm] --platform [docker|kubernetes] --url [server-url] --set --yes

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -a, --components stringArray   List of components
  -h, --help                     help for create
  -p, --platform string          Platform to deploy Meshery (docker or kubernetes)
      --provider string          Provider to use with the Meshery server (Layer5 or None)
  -s, --set                      Set as current context
  -u, --url string               Meshery Server URL with Port (default: http://localhost:9081)

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
![context-create-usage](/reference/images/newcontext.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
