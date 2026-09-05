---
title: mesheryctl-system-check
display_title: false
command: system
subcommand: check
---

# mesheryctl system check

Pre-deployment and post-deployment healthchecks for Meshery

## Synopsis

Verify environment pre/post-deployment of Meshery.

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system check [flags]

</div>
</div>
</pre> 

## Examples

Run all system checks for both pre and post-deployment scenarios
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system check

</div>
</div>
</pre> 

Run pre-deployment checks (Docker and Kubernetes)
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system check --preflight

</div>
</div>
</pre> 

Run pre-deployment checks (Docker and Kubernetes)
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system check --pre

</div>
</div>
</pre> 

Run checks for all Meshery adapters
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system check --adapters

</div>
</div>
</pre> 

Run checks on a specific Meshery adapter
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system check --adapter meshery-istio:10000

</div>
</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system check --adapter meshery-istio

</div>
</div>
</pre> 

Verify the health of Meshery Operator's deployment with MeshSync and Broker
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system check --operator

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --adapter string   Check status of specified meshery adapter
      --adapters         Check status of meshery adapters
      --components       Check status of Meshery components
  -h, --help             help for check
      --operator         Verify the health of Meshery Operator's deployment with MeshSync and Broker
      --pre              Verify environment readiness to deploy Meshery
      --preflight        Verify environment readiness to deploy Meshery

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string    path to config file (default "/home/runner/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## Screenshots

Usage of mesheryctl system check
![check-usage](../../../images/check.png)

## See Also

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
