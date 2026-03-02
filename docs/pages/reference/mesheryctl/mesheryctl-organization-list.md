---
layout: default
title: mesheryctl-organization-list
permalink: reference/mesheryctl/organization/list
redirect_from: reference/mesheryctl/organization/list/
type: reference
display-title: "false"
language: en
command: organization
subcommand: list
---

# mesheryctl organization list

List registered organizations

## Synopsis

List all registered organizations with their id, name and date of creation. Organization will be returned based on provider you logged in.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl organization list [flags]

</div>
</pre> 

## Examples

list all organizations
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl organization list

</div>
</pre> 

list organizations for a specified page
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl organization list --page [page-number]

</div>
</pre> 

Display number of available organizations
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl organization list --count

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --count          total number of registered orgs
  -h, --help           help for list
  -p, --page int       (optional) Page number of paginated results (default 1)
  -s, --pagesize int   (optional) Number of organizations per page (default 10)

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
