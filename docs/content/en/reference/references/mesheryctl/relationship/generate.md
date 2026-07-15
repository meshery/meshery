---
title: mesheryctl-relationship-generate
display_title: false
command: relationship
subcommand: generate
---

# mesheryctl relationship generate

Generate relationships documents

## Synopsis

Generate relationships documents from a CSV file or Google Spreadsheet
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl relationship generate [flags]

</div>
</pre> 

## Examples

Generate relationships documents from a CSV file
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl relationship generate --file <path-to-relationships.csv>

</div>
</pre> 

Generate relationships documents with a custom output path
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl relationship generate --file <path-to-relationships.csv> --output <path-to-output.json>

</div>
</pre> 

Generate relationships documents from a Google Spreadsheet
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl relationship generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string               path to the relationships CSV file
  -h, --help                      help for generate
  -o, --output string             path to the output JSON file
      --spreadsheet-cred string   base64 encoded credential to download the spreadsheet
      --spreadsheet-id string     spreadsheet ID for the integration spreadsheet

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

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
