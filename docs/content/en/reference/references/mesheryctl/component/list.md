---
title: mesheryctl-component-list
display_title: false
command: component
subcommand: list
---

# mesheryctl component list

List registered components

## Synopsis

List all components registered in Meshery Server

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl component list [flags]

</div>
</div>
</pre> 

## Examples

View list of components
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl component list

</div>
</div>
</pre> 

View list of components with specified page number (10 components per page)
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl component list --page [page-number]

</div>
</div>
</pre> 

View list of components with specified page number with specified number of components per page
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl component list --page [page-number] --pagesize [page-size]

</div>
</div>
</pre> 

Display the number of components present in Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl component list --count

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --count          (optional) Display count only
  -h, --help           help for list
  -p, --page int       (optional) List next set of components with --page (default = 1) (default 1)
  -s, --pagesize int   (optional) List next set of components with --pagesize (default = 10) (default 10)

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
