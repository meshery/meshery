# mesheryctl-design-undeploy

LLMS index: [llms.txt](/pr-preview/pr-24333/llms.txt)

---

# mesheryctl design undeploy

Undeploy design

## Synopsis

Undeploy design will trigger undeploy of design.
	
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl design undeploy [flags]

</div>
</div>
</pre> 

## Examples

Undeploy design by providing file path
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl design undeploy -f [filepath]

</div>
</div>
</pre> 

Undeploy design without listing the target clusters and asking for confirmation
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl design undeploy -f [filepath] --yes

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string   Path to design file
  -h, --help          help for undeploy
  -y, --yes           (optional) assume yes for user interactive prompts.

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

Go back to [command reference index](/pr-preview/pr-24333/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/pr-preview/pr-24333/project/contributing/cli/cli/#documenting) for guidance.
