---
layout: default
title: mesheryctl-exp-organization-list
permalink: reference/mesheryctl/exp/organization/list
redirect_from: reference/mesheryctl/exp/organization/list/
type: reference
display-title: "false"
language: en
command: exp
subcommand: organization
---

# mesheryctl exp organization list

List registered organizations

## Synopsis

List all registered organizations with their id, name and date of creation. Organization will be returned based on provider you logged in.
Documentation for organizations can be found at https://docs.meshery.io/reference/mesheryctl/exp/organizations/list
	
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp organization list [flags]

</div>
</pre> 

## Examples

list all organizations
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp organizations list

</div>
</pre> 

list organizations for a specified page
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp organizations list --page [page-number]

</div>
</pre> 

list organizations for a specified page
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp organizations list --count

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --count      total number of registered orgs
  -h, --help       help for list
  -p, --page int   page number

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/n2/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
