---
layout: default
title: mesheryctl-organization
permalink: reference/mesheryctl/organization
redirect_from: reference/mesheryctl/organization/
type: reference
display-title: "false"
language: en
command: organization
subcommand: nil
---

# mesheryctl organization

Interact with registered organizations

## Synopsis

Interact with registered organizations to display detailed information
Documentation for organizations can be found at https://docs.meshery.io/reference/mesheryctl/organizations
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl organization [flags]

</div>
</pre> 

## Examples

Number of  registered orgs
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl organization --count

</div>
</pre> 

List registered orgs
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl organization list

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
      --config string   path to config file (default "/home/vydeh/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

* [mesheryctl organization list](/reference/mesheryctl/organization/list)

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
