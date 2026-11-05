# mesheryctl-design-delete

LLMS index: [llms.txt](/llms.txt)

---

# mesheryctl design delete

Delete design file

## Synopsis

delete design file will trigger deletion of the design file.
	When a positional target is passed, it is resolved with the following precedence:
	1. Existing local file path
	2. Remote URL
	3. Design ID or name lookup
	
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl design delete [flags]

</div>
</div>
</pre> 

## Examples

delete a design file
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl design delete [file | URL]

</div>
</div>
</pre> 

delete a design by ID
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl design delete [ID]

</div>
</div>
</pre> 

delete a design using --file
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl design delete -f [file | URL]

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string   Path to design file
  -h, --help          help for delete

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "~/.meshery/config.yaml")
  -t, --token string    Path to token file default from current context
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/cli/cli/#documenting) for guidance.
