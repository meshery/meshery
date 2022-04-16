---
layout: default
title: mesheryctl-perf-result
permalink: /reference/mesheryctl/perf/result/
redirect_from: /reference/mesheryctl/perf/result/
type: reference
display-title: "false"
language: en
command: perf
---

# mesheryctl perf result

List performance test results

## Synopsis

List all the available test results of a performance profile

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf result profile-name [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>

	// List Test results (maximum 25 results)
	mesheryctl perf result saturday-profile

	// View other set of performance results with --page (maximum 25 results)
	mesheryctl perf result saturday-profile --page 2

	// View single performance result with detailed information
	mesheryctl perf result saturday-profile --view


</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help       help for result
  -p, --page int   (optional) List next set of performance results with --page (default = 1) (default 1)
      --view       (optional) View single performance results with more info

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string          path to config file (default "/home/admin-pc/.meshery/config.yaml")
  -o, --output-format string   (optional) format to display in [json|yaml]
  -t, --token string           (required) Path to meshery auth config
  -v, --verbose                verbose output
  -y, --yes                    (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

* [mesheryctl perf](perf/)	 - Performance Management

