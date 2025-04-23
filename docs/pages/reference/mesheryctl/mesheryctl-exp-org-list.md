---
layout: default
title: mesheryctl-exp-org-list
permalink: reference/mesheryctl/exp/org/list
redirect_from: reference/mesheryctl/exp/org/list/
type: reference
display-title: "false"
language: en
command: exp
subcommand: org
---

# mesheryctl exp org list

List registered orgs

## Synopsis

Print all registered organizations with their id, name and date of creation. Organization will be returned based on provider you logged in.
	
	Documentation for organizations can be found at 
	https://docs.layer5.io/cloud/identity/organizations/
	
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp org list [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>
	// list all organizations

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesheryctl exp org list

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	// list organizations (using flags)

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesheryctl exp org --page [page_no] --size [size]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for list

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
