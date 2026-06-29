---
title: mesheryctl-model-generate
display_title: false
command: model
subcommand: generate
---

# mesheryctl model generate

Generate models from a file

## Synopsis

Generate models by specifying the directory, file, or URL. You can also provide a template JSON file and registrant name

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model generate [flags]

</div>
</pre> 

## Examples

Generate a model from a CSV directory
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model generate -f [path-to-csv-directory]

</div>
</pre> 

Generate a model from a URL based on a JSON template
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model generate -f [URL] -t [path-to-template.json]

</div>
</pre> 

Generate a model from a URL based on a JSON template skipping registration
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model generate --file [URL] --template [path-to-template.json] --skip-registration

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string         Specify path to the file or directory
  -h, --help                help for generate
      --skip-registration   Skip registration of the model (default is false)
  -t, --template string     Specify path to the template JSON file

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
