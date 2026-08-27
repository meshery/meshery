---
title: mesheryctl-model
display_title: false
command: model
subcommand: nil
---

# mesheryctl model

Manage models in the registry

## Synopsis

Export, generate, import, list, search and view model(s) and detailed informations

<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model [flags]</code>
</div>
</pre> 

## Examples

Display number of available models in Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model --count</code>
</div>
</pre> 

Export registered models
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model export [model-name]</code>
</div>
</pre> 

Generate a model from a CSV directory
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model generate [path-to-csv-directory]</code>
</div>
</pre> 

Generate a model from a URL based on a JSON template
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model generate --file [URL] --template [path-to-template.json]</code>
</div>
</pre> 

Import model(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model import -f [Uri]</code>
</div>
</pre> 

List available model(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model list</code>
</div>
</pre> 

Delete available model(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model delete [model-id]</code>
</div>
</pre> 

Search for a specific model
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model search [model-name]</code>
</div>
</pre> 

View a specific model
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model view [model-name]</code>
</div>
</pre> 

Scaffold a folder structure for model creation
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model init [model-name]</code>
</div>
</pre> 

Create an OCI-compliant package from the model files
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model build [model-name]</code>
</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl model build [model-name]/[model-version]</code>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --count   (optional) Get the number of models in total
  -h, --help    help for model

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
