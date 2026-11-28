---
title: mesheryctl-registry-publish
display_title: false
command: registry
subcommand: publish
categories: [mesheryctl-registry]
aliases:
- /reference/references/mesheryctl/registry/publish/
---

# mesheryctl registry publish

Publish Meshery Models to Websites and Remote Provider

## Synopsis

Publishes metadata about Meshery Models to Websites or Remote Provider, including model and component icons by reading from a Google Spreadsheet and outputing to markdown or json format.
For the remote-provider system, [models-output-path] and [imgs-output-path] are resolved relative to the parent of the current working directory, not the current directory itself, and icons are written to a subdirectory named for the model: from /work/project, models writes to /work/models and icons to /work/imgs/my-model.

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path] [flags]

</div>
</div>
</pre> 

## Examples

Publish To System
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path] -o [md|mdx|js]

</div>
</div>
</pre> 

Publish To Remote Provider
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry publish remote-provider GoogleCredential GoogleSheetID [models-output-path] [imgs-output-path]

</div>
</div>
</pre> 

Publish To Website
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry publish website GoogleCredential GoogleSheetID [models-output-path] [imgs-output-path] -o [md|mdx|js]

</div>
</div>
</pre> 

Publishing to meshery docs
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
cd docs;

</div>
</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry publish website "$CRED" 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw docs/pages/integrations docs/assets/img/integrations -o md

</div>
</div>
</pre> 

Publishing to mesheryio site
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry publish website "$CRED" 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw meshery.io/integrations meshery.io/assets/images/integration -o js

</div>
</div>
</pre> 

Publishing to any website
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry publish website "$CRED" 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw path/to/models path/to/icons -o mdx

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help                   help for publish
  -o, --output-format string   output format [md | mdx | js]

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

Go back to [command reference index]({{< ref "reference/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#documenting" >}}) for guidance.
