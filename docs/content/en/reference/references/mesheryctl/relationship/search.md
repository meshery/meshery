---
title: mesheryctl-relationship-search
display_title: false
command: relationship
subcommand: search
---

# mesheryctl relationship search

Search registered relationship(s)

## Synopsis

Search registered relationship(s) used by different models.
	
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl relationship search [flags]

</div>
</div>
</pre> 

## Examples

Search for a specific relationship
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl relationship search [--kind &lt;kind&gt;] [--type &lt;type&gt;] [--subtype &lt;subtype&gt;] [--model &lt;model&gt;]

</div>
</div>
</pre> 

Search a relationship for a specified page
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl relationship search [--kind &lt;kind&gt;] [--page &lt;int&gt;]

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help             help for search
  -k, --kind string      (optional) Search relationships of a particular kind
  -m, --model string     (optional) Search relationships of a particular model name
  -p, --page int         (optional) Page number of results to fetch (default = 1) (default 1)
  -s, --subtype string   (optional) Search relationships of a particular subtype
  -t, --type string      (optional) Search relationships of a particular type

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

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
