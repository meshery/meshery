---
layout: default
title: mesheryctl-exp-environment
permalink: reference/mesheryctl/exp/environment
redirect_from: reference/mesheryctl/exp/environment/
type: reference
display-title: "false"
language: en
command: exp
subcommand: environment
---

# mesheryctl exp environment

View list of environments and detail of environments

## Synopsis

View list of environments and detailed information of a specific environments
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp environment [flags]

</div>
</pre> 

## Examples

To view a list environments
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp environment list --orgID [orgId]

</div>
</pre> 

To create a environment
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp environment create --orgID [orgId] --name [name] --description [description]

</div>
</pre> 

Documentation for environment can be found at:
<pre class='codeblock-pre'>
<div class='codeblock'>
https://docs.layer5.io/cloud/spaces/environments/

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

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
