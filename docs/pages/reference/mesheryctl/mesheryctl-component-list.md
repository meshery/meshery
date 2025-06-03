---
layout: default
title: mesheryctl-component-list
permalink: reference/mesheryctl/component/list
redirect_from: reference/mesheryctl/component/list/
type: reference
display-title: "false"
language: en
command: component
subcommand: list
---

# mesheryctl component list

List registered components

## Synopsis

List all components registered in Meshery Server
Documentation for components can be found at https://docs.meshery.io/reference/mesheryctl/component/list
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl component list [flags]

</div>
</pre> 

## Examples

View list of components
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl component list

</div>
</pre> 

View list of components with specified page number (25 components per page)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl component list --page [page-number]

</div>
</pre> 

Display the number of components present in Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl component list --count

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --count          (optional) Display count only
  -h, --help           help for list
  -p, --page int       (optional) List next set of components with --page (default = 1) (default 1)
  -s, --pagesize int   (optional) List next set of components with --pagesize (default = 0)

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
