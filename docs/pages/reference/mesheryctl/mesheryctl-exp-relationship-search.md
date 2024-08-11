---
layout: default
title: mesheryctl-exp-relationship-search
permalink: reference/mesheryctl/exp/relationship/search
redirect_from: reference/mesheryctl/exp/relationship/search/
type: reference
display-title: "false"
language: en
command: exp
subcommand: relationship
---

# mesheryctl exp relationship search

Searches registered relationships

## Synopsis

Searches and finds the realtionship used by different models based on the query-text.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp relationship search [flags]

</div>
</pre> 

## Examples

Search for relationship using a query
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp relationship search --[flag] [query-text]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help             help for search
  -k, --kind string      search particular kind of relationships
  -m, --model string     search relationships of particular model name
  -s, --subtype string   search particular subtype of relationships
  -t, --type string      search particular type of relationships

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
