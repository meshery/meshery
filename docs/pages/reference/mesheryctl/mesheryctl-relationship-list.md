---
layout: default
title: mesheryctl-relationship-list
permalink: reference/mesheryctl/relationship/list
redirect_from: reference/mesheryctl/relationship/list/
type: reference
display-title: "false"
language: en
command: relationship
subcommand: list
---

# mesheryctl relationship list

List registered relationships

## Synopsis

List all relationships registered in Meshery Server
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl relationship list [flags]

</div>
</pre> 

## Examples

List all relationships
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl relationship list

</div>
</pre> 

List relationships for a specified page
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl relationship list --page [page-number]

</div>
</pre> 

List relationships with a custom page size
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl relationship list --pagesize [page-size]

</div>
</pre> 

Display the total number of available relationships in Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl relationship list --count

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --count          (optional) Display the total count of relationships only
  -h, --help           help for list
  -p, --page int       (optional) List next set of relationships with --page (default = 1) (default 1)
      --pagesize int   (optional) Number of results per page (default = 10) (default 10)

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
