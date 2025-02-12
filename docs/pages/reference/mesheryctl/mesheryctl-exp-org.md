---
layout: default
title: mesheryctl-exp-org
permalink: reference/mesheryctl/exp/org
redirect_from: reference/mesheryctl/exp/org/
type: reference
display-title: "false"
language: en
command: exp
subcommand: org
---

# mesheryctl exp org

view list of registered orgs

## Synopsis

view list of registered orgs with their name,id and date of creation
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp org [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>
	// Number of  registered orgs

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesheryctl org --count 

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	// List registerd orgs

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesheryctl org list

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --count   total number of registered orgs
  -h, --help    help for org

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

* [mesheryctl exp org list](/reference/mesheryctl/exp/org/list)

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
