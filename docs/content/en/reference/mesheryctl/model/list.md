---
title: mesheryctl-model-list
display_title: false
command: model
subcommand: list
---

# mesheryctl model list

List registered models

## Synopsis

List all registered models by pagingation (10 models per page)

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model list [flags]

</div>
</pre> 

## Examples

List of models
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model list

</div>
</pre> 

List of models for a specified page
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model list --page [page-number] --pagesize [pagesize]

</div>
</pre> 

Display number of available models in Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model list --count

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
    

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --count          (optional) Get the number of models in total
  -h, --help           help for list
  -p, --page int       (optional) List next set of models with --page (default = 1) (default 1)
  -s, --pagesize int   (optional) List next set of models with --pagesize (default = 10) (default 10)

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
