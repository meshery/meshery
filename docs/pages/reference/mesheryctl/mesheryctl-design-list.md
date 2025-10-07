---
layout: default
title: mesheryctl-design-list
permalink: reference/mesheryctl/design/list
redirect_from: reference/mesheryctl/design/list/
type: reference
display-title: "false"
language: en
command: design
subcommand: list
---

# mesheryctl design list

List designs

## Synopsis

Display list of all available designs.
Documentation for design can be found at https://docs.meshery.io/reference/mesheryctl/design/list

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design list [flags]

</div>
</pre> 

## Examples

Display a list of all available designs
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design list

</div>
</pre> 

Display a list of all available designs with verbose output
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design list --verbose

</div>
</pre> 

Display a list of all available designs with specified page number (10 designs per page by default)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design list --page [pange-number]

</div>
</pre> 

Display a list of all available designs with custom page size (10 designs per page by default)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design list --pagesize [page-size]

</div>
</pre> 

Display only the count of all available designs
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design list --count

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --count          (optional) Display count only
  -h, --help           help for list
  -p, --page int       (optional) List next set of designs with --page (default 1)
      --pagesize int   (optional) Number of designs to be displayed per page (default 10)
  -v, --verbose        (optional) Display full length user and design file identifiers

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
![pattern-list-usage](/assets/img/mesheryctl/patternList.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
