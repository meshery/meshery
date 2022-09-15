---
layout: default
title: mesheryctl-perf-apply
permalink: reference/mesheryctl/perf/apply
redirect_from: reference/mesheryctl/perf/apply/
type: reference
display-title: "false"
language: en
command: perf
subcommand: apply
---

# mesheryctl perf apply

Run a Performance test

## Synopsis

Run Performance test using existing profiles or using flags

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply [profile-name | --file] --flags [flags]

</div>
</pre> 

## Examples

Execute a Performance test with the specified performance profile
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile --flags

</div>
</pre> 

Execute a Performance test with creating a new performance profile
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile-new --url "https://google.com"

</div>
</pre> 

Run Performance test using SMP compatible test configuration
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply -f perf-config.yaml

</div>
</pre> 

Run performance test using SMP compatible test configuration and override values with flags
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply -f [filepath] --flags

</div>
</pre> 

Choice of load generator - fortio or wrk2 (default: fortio)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-test --load-generator wrk2

</div>
</pre> 

Execute a Performance test with specified queries per second
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply local-perf --url https://192.168.1.15/productpage --qps 30

</div>
</pre> 

Execute a Performance test with specified service mesh
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply local-perf --url https://192.168.1.15/productpage --mesh istio

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --concurrent-requests string   (optional) Number of Parallel Requests
      --duration string              (optional) Length of test (e.g. 10s, 5m, 2h). For more, see https://golang.org/pkg/time/#ParseDuration
  -f, --file string                  (optional) file containing SMP-compatible test configuration. For more, see https://github.com/layer5io/service-mesh-performance-specification
  -h, --help                         help for apply
      --load-generator string        (optional) Load-Generator to be used (fortio/wrk2)
      --mesh string                  (optional) Name of the Service Mesh
      --name string                  (optional) Name of the Test
      --qps string                   (optional) Queries per second
      --url string                   (optional) Endpoint URL to test (required with --profile)

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string          path to config file (default "/home/runner/.meshery/config.yaml")
  -o, --output-format string   (optional) format to display in [json|yaml]
  -t, --token string           (required) Path to meshery auth config
  -v, --verbose                verbose output
  -y, --yes                    (optional) assume yes for user interactive prompts.

</div>
</pre>

## Screenshots

Usage of mesheryctl perf apply
![perf-apply-usage](/assets/img/mesheryctl/perf-apply.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/) 
