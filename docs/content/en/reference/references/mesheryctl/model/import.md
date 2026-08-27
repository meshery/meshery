---
title: mesheryctl-model-import
display_title: false
command: model
subcommand: import
---

# mesheryctl model import

Import models

## Synopsis

Import models by specifying the directory, file, or URL. You can also provide a template JSON file and registrant name

<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model import [flags]</code>
</div>
</pre> 

## Examples

Import model
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model import --file [URI]</code>
</div>
</pre> 

Import model from a URL to a meshery model
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model import --file [URL]</code>
</div>
</pre> 

Import model from an OCI artifact
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model import --file [OCI]</code>
</div>
</pre> 

Import model from a tar.gz file
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model import --file [path-to-model.tar.gz]</code>
</div>
</pre> 

Import model from a path
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model import --file [path-to-model]</code>
</div>
</pre> 

Import model using CSV files
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model import --file [path-to-csv-directory]</code>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string   Specify path to the file or directory
  -h, --help          help for import

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
