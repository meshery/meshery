---
layout: default
title: mesheryctl-mesh-validate
permalink: reference/mesheryctl/mesh/validate
redirect_from: reference/mesheryctl/mesh/validate/
type: reference
display-title: "false"
language: en
command: mesh
subcommand: validate
---

# mesheryctl mesh validate

Validate conformance to service mesh standards

## Synopsis

Validate service mesh conformance to different standard specifications

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh validate [flags]

</div>
</pre> 

## Examples

Validate conformance to service mesh standards
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh validate --adapter [name of the adapter] --tokenPath [path to token for authentication] --spec [specification to be used for conformance test] --namespace [namespace to be used]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -a, --adapter string   Adapter to use for validation (default "meshery-osm")
  -h, --help             help for validate
  -s, --spec string      specification to be used for conformance test (default "smi")
  -t, --token string     Path to token for authenticating to Meshery API
  -w, --watch            Watch for events and verify operation (in beta testing)

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

Go back to [command reference index](/reference/mesheryctl/) 
