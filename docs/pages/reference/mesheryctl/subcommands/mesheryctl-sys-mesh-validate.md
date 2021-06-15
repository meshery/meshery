---
layout: default
title: mesheryctl system mesh validate
permalink: reference/mesheryctl/mesh/validate
type: reference
display-title: "false"
language: en
categories: en
list: exclude
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl system mesh validate

<!-- Description of the command. Preferably a paragraph -->
## Description

{% for subcommand_hash in site.data.mesheryctlcommands.meshes.validate.subcommand %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl mesh validate [flags]
  </div>
</pre>

## Examples

<h6>(required) adapter to use for validation. Defaults to "meshery-osm:10010"<h6>
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl mesh validate --adapter [name of the adapter]
  </div>
</pre>
<h6>Kubernetes namespace to be used for deploying the validation tests and sample workload<h6>
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl mesh validate --namespace [namespace to be used]
  </div>
</pre>
<h6>(required) specification to be used for conformance test. Defaults to "smi"<h6>
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl mesh validate --spec [specification to be used for conformance test]
  </div>
</pre>
<h6>(required) path to token for authenticating to Meshery API<h6>
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl mesh validate --tokenPath [path to token for authentication]
  </div>
</pre>
<br/>


<!-- Options/Flags available in this command -->
<h2> Options & Flags </h2>

<h6>(required) adapter to use for validation. Defaults to "meshery-osm:10010"<h6>
<pre class="codeblock-pre">
  <div class="codeblock">
   --adapter, -a
  </div>
</pre>
<h6>Kubernetes namespace to be used for deploying the validation tests and sample workload<h6>
<pre class="codeblock-pre">
  <div class="codeblock">
  --namespace, -n
  </div>
</pre>
<h6>(required) specification to be used for conformance test. Defaults to "smi"<h6>
<pre class="codeblock-pre">
  <div class="codeblock">
  --spec, -s
  </div>
</pre>
<h6>(required) path to token for authenticating to Meshery API<h6>
<pre class="codeblock-pre">
  <div class="codeblock">
  --tokenpath, -t
  </div>
</pre>
<br/>

<h2> Options inherited from parent commands </h2>
<pre class="codeblock-pre">
  <div class="codeblock">
  --help, -h # Shows help for the command
  </div>
</pre>
