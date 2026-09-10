---
title: mesheryctl-version
display_title: false
command: version
subcommand: nil
---

# mesheryctl version

Show Meshery CLI and Server versions

## Synopsis

Version of Meshery command line client - mesheryctl.
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl version [flags]

</div>
</div>
</pre> 

## Examples

To view the current version and SHA of release binary of mesheryctl client 
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl version

</div>
</div>
</pre> 

To view the version in JSON or YAML format
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl version -o json

</div>
</div>
</pre>

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for version
  -o, --output-format string   (optional) format to display in [json|yaml|table|string]

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## Screenshots

Usage of mesheryctl version
![version-usage](../../images/version.png)

## See Also

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
