---
title: mesheryctl-perf-delete
display_title: false
command: perf
subcommand: delete
categories: [mesheryctl-perf]
aliases:
- /reference/references/mesheryctl/perf/delete/
---

# mesheryctl perf delete

Delete a performance profile

## Synopsis

Delete a performance profile by its name.
	
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl perf delete [profile-name] [flags]

</div>
</div>
</pre> 

## Examples

Delete a performance profile (asks for confirmation)
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl perf delete test-3

</div>
</div>
</pre> 

Delete a performance profile without confirmation
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl perf delete test-3 --yes

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
      --config string          path to config file (default "~/.meshery/config.yaml")
  -o, --output-format string   (optional) format to display in [json|yaml]
  -t, --token string           (required) Path to meshery auth config
  -v, --verbose                verbose output
  -y, --yes                    (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

Go back to [command reference index]({{< ref "reference/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#documenting" >}}) for guidance.
