---
layout: default
title: mesheryctl-adapter-validate
permalink: reference/mesheryctl/adapter/validate
redirect_from: reference/mesheryctl/adapter/validate/
type: reference
display-title: "false"
language: en
command: adapter
subcommand: validate
---

# mesheryctl adapter validate

Validate conformance to predefined standards

## Synopsis

Validate predefined conformance to different standard specifications
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl adapter validate [flags]

</div>
</pre> 

## Examples

Validate conformance to predefined standards
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl adapter validate [mesh name] --adapter [name of the adapter] --tokenPath [path to token for authentication] --spec [specification to be used for conformance test] --namespace [namespace to be used]

</div>
</pre> 

Validate Istio to predefined standards
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl adapter validate istio --adapter meshery-istio --spec smi

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -a, --adapter string   (Required) Adapter to use for validation (default "meshery-nsm")
  -h, --help             help for validate
  -s, --spec string      (Required) specification to be used for conformance test (smi/istio-vet) (default "smi")
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

## Screenshots

Usage of mesheryctl adapter validate
![mesh-validate-usage](/assets/img/mesheryctl/mesh-validate.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
