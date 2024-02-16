---
layout: default
title: mesheryctl-exp-registry-publish
permalink: reference/mesheryctl/exp/registry/publish
redirect_from: reference/mesheryctl/exp/registry/publish/
type: reference
display-title: "false"
language: en
command: exp
subcommand: registry
---

# mesheryctl exp registry publish

Publish Meshery Models to Websites, Remote Provider, Meshery

## Synopsis

Publishes metadata about Meshery Models to Websites, Remote Provider and Meshery by reading from a Google Spreadsheet.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp registry publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path] [flags]

</div>
</pre> 

## Examples

Publish To System
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp registry publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path] -o [output-format]

</div>
</pre> 

Publish To Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp registry publish meshery GoogleCredential GoogleSheetID [repo]/server/meshmodel

</div>
</pre> 

Publish To Remote Provider
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp registry publish remote-provider GoogleCredential GoogleSheetID [repo]/meshmodels/models [repo]/ui/public/img/meshmodels

</div>
</pre> 

Publish To Website
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp registry publish website GoogleCredential GoogleSheetID [repo]/integrations [repo]/ui/public/img/meshmodels

</div>
</pre> 

Publishing to meshery docs
<pre class='codeblock-pre'>
<div class='codeblock'>
cd docs;

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp registry publish website $CRED 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw docs/pages/integrations docs/assets/img/integrations -o md

</div>
</pre> 

Publishing to mesheryio site
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp registry publish website $CRED 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw meshery.io/integrations meshery.io/assets/images/integration -o js

</div>
</pre> 

Publishing to layer5 site
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl/mesheryctl exp registry publish website $CRED 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw layer5/src/collections/integrations layer5/src/collections/integrations -o mdx

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
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
