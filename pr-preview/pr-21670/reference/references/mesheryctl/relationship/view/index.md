# mesheryctl-relationship-view

Source: /pr-preview/pr-21670/reference/references/mesheryctl/relationship/view/

# mesheryctl relationship view

view relationships of a model by its name

## Synopsis

view a relationship queried by the model name.
	
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl relationship view [flags]

</div>
</div>
</pre> 

## Examples

View relationships of a model in default format yaml
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl relationship view [model-name]

</div>
</div>
</pre> 

View relationships of a model in JSON format
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl relationship view [model-name] --output-format json

</div>
</div>
</pre> 

View relationships of a model in json format and save it to a file
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl relationship view [model-name] --output-format json --save

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help                   help for view
  -o, --output-format string   (optional) format to display in [json|yaml] (default "yaml")
  -s, --save                   (optional) save output as a JSON/YAML file

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

Go back to [command reference index](/pr-preview/pr-21670/reference/references/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/pr-preview/pr-21670/project/contributing/cli/cli/#preserving-manually-added-documentation) for guidance.
