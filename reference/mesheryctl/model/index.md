# mesheryctl-model

LLMS index: [llms.txt](/llms.txt)

---

# mesheryctl model

Manage models in the registry

## Synopsis

Export, generate, import, list, search and view model(s) and detailed informations

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl model [flags]

</div>
</div>
</pre> 

## Examples

Display number of available models in Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl model --count

</div>
</div>
</pre> 

Export registered models
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl model export [model-name]

</div>
</div>
</pre> 

Generate a model from a CSV directory
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl model generate [path-to-csv-directory]

</div>
</div>
</pre> 

Generate a model from a URL based on a JSON template
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl model generate --file [URL] --template [path-to-template.json]

</div>
</div>
</pre> 

Import model(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl model import -f [Uri]

</div>
</div>
</pre> 

List available model(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl model list

</div>
</div>
</pre> 

Delete available model(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl model delete [model-id]

</div>
</div>
</pre> 

Search for a specific model
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl model search [model-name]

</div>
</div>
</pre> 

View a specific model
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl model view [model-name]

</div>
</div>
</pre> 

Scaffold a folder structure for model creation
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl model init [model-name]

</div>
</div>
</pre> 

Create an OCI-compliant package from the model files
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl model build [model-name]

</div>
</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl model build [model-name]/[model-version]

</div>
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

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/cli/cli/#preserving-manually-added-documentation) for guidance.

---

Section pages:

- [mesheryctl-model-build](/reference/mesheryctl/model/build/)
- [mesheryctl-model-delete](/reference/mesheryctl/model/delete/)
- [mesheryctl-model-export](/reference/mesheryctl/model/export/)
- [mesheryctl-model-generate](/reference/mesheryctl/model/generate/)
- [mesheryctl-model-import](/reference/mesheryctl/model/import/)
- [mesheryctl-model-init](/reference/mesheryctl/model/init/)
- [mesheryctl-model-list](/reference/mesheryctl/model/list/)
- [mesheryctl-model-search](/reference/mesheryctl/model/search/)
- [mesheryctl-model-view](/reference/mesheryctl/model/view/)
