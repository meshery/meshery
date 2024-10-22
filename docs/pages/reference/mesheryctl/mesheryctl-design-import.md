---
layout: default
title: mesheryctl-design-import
permalink: reference/mesheryctl/design/import
redirect_from: reference/mesheryctl/design/import/
type: reference
display-title: "false"
language: en
command: design
subcommand: import
---

# mesheryctl design import

Import a Meshery design

## Synopsis


		Import Helm Charts, Kubernetes Manifest, Docker Compose or Meshery designs by passing
		remote URL or local file system path to the file. Source type must be provided.

		YAML and TGZ (with helm only) format of file is accepted, if you are importing Meshery Design OCI file format is also supported

		If you are providing remote URL, it should be a direct URL to a downloadable file.
		For example, if the file is stored on GitHub, the URL should be 'https://raw.githubusercontent.com/path-to-file'.
	
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design import [flags]

</div>
</pre> 

## Examples

Import design manifest
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl design import -f [file/URL] -s [source-type] -n [name]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string          Path/URL to design file
  -h, --help                 help for import
  -n, --name string          Name for the design file
  -s, --source-type string   Type of source file (ex. manifest / compose / helm / design)

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

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
