---
layout: default
title: mesheryctl-environment-delete
permalink: reference/mesheryctl/environment/delete
redirect_from: reference/mesheryctl/environment/delete/
type: reference
display-title: "false"
language: en
command: environment
subcommand: delete
---

# mesheryctl environment delete

Delete an environment

## Synopsis

Delete an environment by providing the environment ID
Documentation for environment can be found at Documentation for environment can be found at https://docs.meshery.io/reference/mesheryctl/environment/delete
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment delete [flags]

</div>
</pre> 

## Examples

delete a new environment
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment delete [environmentId]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for delete

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
