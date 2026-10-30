# mesheryctl-system-delete

LLMS index: [llms.txt](/pr-preview/pr-22685/llms.txt)

---

# mesheryctl system delete

Delete Meshery containers

## Synopsis

Delete Meshery containers. This command removes all Meshery containers created by docker-compose.
	
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system delete [flags]

</div>
</div>
</pre> 

## Examples

Delete Meshery containers
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system delete

</div>
</div>
</pre> 

Delete Meshery containers without confirmation
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system delete -y

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for delete

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string    path to config file (default "C:\\Users\\saisu\\.meshery\\config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

Go back to [command reference index](/pr-preview/pr-22685/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/pr-preview/pr-22685/project/contributing/cli/cli/#preserving-manually-added-documentation) for guidance.
