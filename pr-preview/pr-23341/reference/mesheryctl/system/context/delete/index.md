# mesheryctl-system-context-delete

LLMS index: [llms.txt](/docs-preview/pr-preview/pr-23341/llms.txt)

---

# mesheryctl system context delete

Delete context

## Synopsis

Delete an existing context (a named Meshery deployment) from Meshery config file.

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system context delete [context-name] [flags]

</div>
</div>
</pre> 

## Examples

### Delete context
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system context delete [context name]

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help         help for delete
  -s, --set string   New context to deploy Meshery

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string    path to config file (default "~/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

Go back to [command reference index](/docs-preview/pr-preview/pr-23341/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/docs-preview/pr-preview/pr-23341/project/contributing/cli/cli/#documenting) for guidance.
