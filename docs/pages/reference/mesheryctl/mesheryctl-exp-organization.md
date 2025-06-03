---
layout: default
title: mesheryctl-exp-organization
permalink: reference/mesheryctl/exp/organization
redirect_from: reference/mesheryctl/exp/organization/
type: reference
display-title: "false"
language: en
command: exp
subcommand: organization
---

# mesheryctl exp organization

Interact with registered orgnizations

## Synopsis

Interact with registered organizations to display detailled informations
Documentation for organizations can be found at https://docs.meshery.io/reference/mesheryctl/exp/organizations
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp organization [flags]

</div>
</pre> 

## Examples

Number of  registered orgs
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl organizations --count 

</div>
</pre> 

List registerd orgs
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl organizations list	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --count   total number of registered organizations
  -h, --help    help for organization

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

* [mesheryctl exp organization list](/reference/mesheryctl/exp/organization/list)

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
