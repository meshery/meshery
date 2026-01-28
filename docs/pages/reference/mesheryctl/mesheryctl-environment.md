---
layout: default
title: mesheryctl-environment
permalink: reference/mesheryctl/environment
redirect_from: reference/mesheryctl/environment/
type: reference
display-title: "false"
language: en
command: environment
subcommand: nil
---

# mesheryctl environment

Manage environments

## Synopsis

Create, delete, list of view details of environment(s) of a specific organization
Documentation for environment can be found at https://docs.meshery.io/concepts/logical/environments
	
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment [flags]

</div>
</pre> 

## Examples

Create an environment in an organization
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment create --orgID [orgID] --name [name] --description [description]

</div>
</pre> 

Delete an environment in an organization
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment delete environment-id

</div>
</pre> 

List of registered environments in an organization
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment list --orgID [orgID]

</div>
</pre> 

View a particular environment
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment view --orgID [orgID]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for environment

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

* [mesheryctl environment create](/reference/mesheryctl/environment/create)
* [mesheryctl environment delete](/reference/mesheryctl/environment/delete)
* [mesheryctl environment list](/reference/mesheryctl/environment/list)
* [mesheryctl environment view](/reference/mesheryctl/environment/view)

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
