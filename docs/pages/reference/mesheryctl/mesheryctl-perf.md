---
layout: default
title: mesheryctl-perf
permalink: /reference/mesheryctl/perf/
redirect_from: /reference/mesheryctl/perf/
type: reference
display-title: "false"
language: en
command: perf
---

# mesheryctl perf

Performance Management

## Synopsis

Performance Management & Benchmarking

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>


	// Run performance test
	mesheryctl perf apply test-3 --name "a quick stress test" --url http://192.168.1.15/productpage --qps 300 --concurrent-requests 2 --duration 30s
		
	// List performance profiles
	mesheryctl perf profile sam-test

	// List performance results
	mesheryctl perf result sam-test

	// Display Perf profile in JSON or YAML
	mesheryctl perf result -o json
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
      --config string   path to config file (default "/home/admin-pc/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

* [mesheryctl](/reference/mesheryctl/main)	 - Meshery Command Line tool
* [mesheryctl perf apply](apply/)	 - Run a Performance test
* [mesheryctl perf profile](profile/)	 - List performance profiles
* [mesheryctl perf result](result/)	 - List performance test results

