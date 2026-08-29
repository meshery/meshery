# mesheryctl-system-reset

Source: /pr-preview/pr-21670/reference/references/mesheryctl/system/reset/

# mesheryctl system reset

Reset Meshery's configuration

## Synopsis

Reset Meshery to it's default configuration.
	
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system reset [flags]

</div>
</div>
</pre> 

## Examples

Resets meshery.yaml file with a copy from Meshery repo
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system reset

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for reset

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string    path to config file (default "/home/runner/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## Screenshots

Usage of mesheryctl system reset
![reset-usage](../../../images/reset.png)

## See Also

Go back to [command reference index](/pr-preview/pr-21670/reference/references/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/pr-preview/pr-21670/project/contributing/cli/cli/#preserving-manually-added-documentation) for guidance.
