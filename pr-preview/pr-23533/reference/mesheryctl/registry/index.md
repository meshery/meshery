# mesheryctl-registry

LLMS index: [llms.txt](/pr-preview/pr-23533/llms.txt)

---

# mesheryctl registry

Manage the capability registry

## Synopsis

Generate, publish, update, and purge models and components in Meshery's capability registry.

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry [flags]

</div>
</div>
</pre> 

## Examples

Generate models from a spreadsheet
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry generate --spreadsheet-id [spreadsheet-id] --spreadsheet-cred [base64-encoded-google-service-account-credential]

</div>
</div>
</pre> 

Publish models to Meshery Cloud / Remote Provider
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path] -o [output-format]

</div>
</div>
</pre> 

Update models from spreadsheets or website
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry update --spreadsheet-id [spreadsheet-id] --spreadsheet-cred [base64-encoded-google-service-account-credential]

</div>
</div>
</pre> 

Dry run purging models directory, keeping only the latest version of each model
Run from the root of a meshery/meshery fork; purge operates on ./models
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry purge --dry-run

</div>
</div>
</pre> 

Purge older model versions, retaining the latest 2 versions per model
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry purge --retain 2

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for registry

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "~/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/pr-preview/pr-23533/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/pr-preview/pr-23533/project/contributing/cli/cli/#documenting) for guidance.

---

Section pages:

- [mesheryctl-registry-generate](/pr-preview/pr-23533/reference/mesheryctl/registry/generate/)
- [mesheryctl-registry-publish](/pr-preview/pr-23533/reference/mesheryctl/registry/publish/)
- [mesheryctl-registry-purge](/pr-preview/pr-23533/reference/mesheryctl/registry/purge/)
- [mesheryctl-registry-update](/pr-preview/pr-23533/reference/mesheryctl/registry/update/)
