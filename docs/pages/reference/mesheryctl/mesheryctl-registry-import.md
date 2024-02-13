---
layout: default
title: mesheryctl-registry-import
permalink: reference/mesheryctl/registry/import
redirect_from: reference/mesheryctl/registry/import/
type: reference
display-title: "false"
language: en
command: registry
subcommand: import
---

# mesheryctl registry import

Import Models

## Synopsis

Import models from spreadsheet, GitHub or ArtifactHub repositories
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl registry import [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>
    // Import models from Meshery Integration Spreadsheet

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
    mesheryctl registry import --spreadsheet_url <url> --spreadsheet_cred <base64 encoded spreadsheet credential>

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
    

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
    // Directly import models from one of the supported registrants by using Registrant Connection Definition and (optional) Registrant Credential Definition

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
    mesheryctl registry import --registrant_def <path to connection definition> --registrant_cred <path to credential definition>

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
    

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help                      help for import
  -o, --output string             location to output generated models, defaults to ../server/meshmodels (default "../server/meshmodel")
      --registrant_cred string    path pointing to the registrant credetial definition
      --registrant_def string     path pointing to the registrant connection definition
      --spreadsheet_cred string   base64 encoded credential to download the spreadsheet
      --spreadsheet_id string     spreadsheet it for the integration spreadsheet

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
