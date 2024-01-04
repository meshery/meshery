---
layout: default
title: mesheryctl-perf
permalink: reference/mesheryctl/perf
redirect_from: reference/mesheryctl/perf/
type: reference
display-title: "false"
language: en
command: perf
subcommand: nil
---

# mesheryctl perf

Performance Management

## Synopsis

Performance Management & Benchmarking.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf [flags]

</div>
</pre> 

## Examples

Run performance test:
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply test-3 --name "a quick stress test" --url http://192.168.1.15/productpage --qps 300 --concurrent-requests 2 --duration 30s

</div>
</pre> 

List performance profiles:
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf profile sam-test

</div>
</pre> 

List performance results:
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf result sam-test

</div>
</pre> 

Display Perf profile in JSON or YAML:
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf result -o json

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf result -o yaml

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help                   help for perf
  -o, --output-format string   (optional) format to display in [json|yaml]
  -t, --token string           (required) Path to meshery auth config
  -y, --yes                    (optional) assume yes for user interactive prompts.

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
