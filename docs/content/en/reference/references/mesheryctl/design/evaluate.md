---
title: mesheryctl-design-evaluate
display_title: false
command: design
subcommand: evaluate
---

# mesheryctl design evaluate

Evaluate a design

## Synopsis

Evaluate a design by running relationship evaluation policies.
The evaluated design is saved to the specified output file while an overview
of evaluation actions is printed to the terminal.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design evaluate [ID] [flags]

</div>
</pre> 

## Examples

Evaluate a design from a file and save the result
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design evaluate -f design.yaml -o evaluated-design.yaml

</div>
</pre> 

Evaluate a design by ID
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design evaluate 12345678-abcd-efgh-ijkl-123456789012 -o result.yaml

</div>
</pre> 

Evaluate and save as JSON
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design evaluate -f design.yaml --output-format json -o evaluated-design.json

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string            Path to design file
  -h, --help                   help for evaluate
  -o, --output string          Path to save the evaluated design
      --output-format string   Output format for the evaluated design [json|yaml] (default "yaml")

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -t, --token string    Path to token file default from current context
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
