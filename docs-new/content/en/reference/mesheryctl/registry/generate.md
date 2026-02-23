---
title: mesheryctl-registry-generate
display_title: false
command: registry
subcommand: generate
---

# mesheryctl registry generate

Generate Models

## Synopsis

Prerequisite: Excecute this command from the root of a meshery/meshery repo fork.\n\nGiven a Google Sheet with a list of model names and source locations, generate models and components any Registrant (e.g. GitHub, Artifact Hub) repositories.\n\nGenerated Model files are written to local filesystem under "/server/models/<model-name>".
Documentation for components can be found at https://docs.meshery.io/reference/mesheryctl/registry/generate
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl registry generate [flags]

</div>
</pre> 

## Examples

Generate Meshery Models from a Google Spreadsheet (i.e. "Meshery Integrations" spreadsheet).
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred "$CRED"

</div>
</pre> 

Directly generate models from one of the supported registrants by using Registrant Connection Definition and (optional) Registrant Credential Definition
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl registry generate --registrant-def [path to connection definition] --registrant-cred [path to credential definition]

</div>
</pre> 

Generate a specific Model from a Google Spreadsheet (i.e. "Meshery Integrations" spreadsheet).
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred --model "[model-name]"

</div>
</pre> 

Generate Meshery Models and Component from csv files in a local directory.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl registry generate --directory [DIRECTORY_PATH]

</div>
</pre> 

Generate Meshery Models from individual CSV files.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl registry generate --model-csv [path/to/models.csv] --component-csv [path/to/components.csv] --relationship-csv [path/to/relationships.csv]

</div>
</pre> 

Generate models with a custom per-model timeout (e.g., 10 minutes per model).
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred "$CRED" --timeout 10m

</div>
</pre> 

Generate only the latest version of each model.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred "$CRED" --latest-only

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --component-csv string      path to the component CSV file
  -d, --directory string          Directory containing the Model and Component CSV files
  -h, --help                      help for generate
      --latest-only               generate only the latest version of each model
  -m, --model string              specific model name to be generated
      --model-csv string          path to the model CSV file
  -o, --output string             location to output generated models, defaults to ../server/meshmodels (default "../server/meshmodel")
      --registrant-cred string    path pointing to the registrant credential definition
      --registrant-def string     path pointing to the registrant connection definition
      --relationship-csv string   path to the relationship CSV file (optional)
      --spreadsheet-cred string   base64 encoded credential to download the spreadsheet
      --spreadsheet-id string     spreadsheet ID for the integration spreadsheet
      --timeout duration          timeout duration for generating each model (e.g., 5m, 10m, 1h), default: 5m (default 5m0s)

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
