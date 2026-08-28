---
title: mesheryctl-design-list
display_title: false
command: design
subcommand: list
---

# mesheryctl design list

List designs

## Synopsis

Display list of all available designs.

<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl design list [flags]</code>
</div>
</pre> 

## Examples

Display a list of all available designs
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl design list</code>
</div>
</pre> 

Display a list of all available designs with verbose output
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl design list --verbose</code>
</div>
</pre> 

Display a list of all available designs with specified page number (10 designs per page by default)
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl design list --page [pange-number]</code>
</div>
</pre> 

Display a list of all available designs with custom page size (10 designs per page by default)
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl design list --pagesize [page-size]</code>
</div>
</pre> 

Display only the count of all available designs
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl design list --count</code>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --count          (optional) Display count only
  -h, --help           help for list
  -p, --page int       (optional) List next set of designs with --page (default 1)
      --pagesize int   (optional) Number of designs to be displayed per page (default 10)
  -v, --verbose        (optional) Display full length owner identifiers and detailed timestamps

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -t, --token string    Path to token file default from current context

</div>
</pre>

## Screenshots

Usage of mesheryctl design list
![pattern-list-usage](../../../images/patternList.png)

## See Also

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
