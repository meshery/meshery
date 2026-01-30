---
layout: default
title: mesheryctl-environment-list
permalink: reference/mesheryctl/environment/list
redirect_from: reference/mesheryctl/environment/list/
type: reference
display-title: "false"
language: en
command: environment
subcommand: list
---

# mesheryctl environment list

List registered environments

## Synopsis

List detailed information of all registered environments
Documentation for environment can be found at https://docs.meshery.io/reference/mesheryctl/environment/list
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment list [flags]

</div>
</pre> 

## Examples

List all registered environment
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment list --orgID [orgID]

</div>
</pre> 

List count of all registered environment
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment list --orgID [orgID] --count

</div>
</pre> 

List all registered environment at a specific page
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment list --orgID [orgID] --page [page]

</div>
</pre> 

List all registered environment with a specific page size
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment list --orgID [orgID] --pagesize [pagesize]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --count          (optional) Display count only
  -h, --help           help for list
      --orgID string   Organization ID
      --page int       (optional) Page number of paginated results (default 1)
      --pagesize int   (optional) Number of results per page (default 10)

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
