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
mesheryctl perf apply [profile-name] [flags]

</div>
</pre> 

## Examples

Execute a Performance test with the specified performance profile
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile [flags]

</div>
</pre> 

Execute a Performance test with creating a new performance profile
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile-new --url "https://google.com"

</div>
</pre> 

Execute a Performance test creating a new performance profile and pass certificate to be used 
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile-new --url "https://google.com" --cert-path path/to/cert.pem

</div>
</pre> 

Execute a performance profile without using the certificate present in the profile
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile --url "https://google.com" --disable-cert

</div>
</pre> 

Run Performance test using SMP compatible test configuration
If the profile already exists, the test will be run overriding the values with the ones provided in the configuration file
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile -f path/to/perf-config.yaml

</div>
</pre> 

Run performance test using SMP compatible test configuration and override values with flags
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile -f path/to/perf-config.yaml [flags]

</div>
</pre> 

Choice of load generator - fortio, wrk2 or nighthawk (default: fortio)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile --load-generator wrk2

</div>
</pre> 

Execute a Performance test with specified queries per second
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile --url https://192.168.1.15/productpage --qps 30

</div>
</pre> 

Execute a Performance test with specified service mesh
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile --url https://192.168.1.15/productpage --mesh istio

</div>
</pre> 

Execute a Performance test creating a new performance profile and pass options to the load generator used
If any options are already present in the profile or passed through flags, the --options flag will take precedence over the profile and flag options 
Options for nighthawk - https://github.com/layer5io/getnighthawk/blob/v1.0.5/pkg/proto/options.pb.go#L882-L1018
Options for fortio - https://github.com/fortio/fortio/blob/v1.57.0/fhttp/httprunner.go#L77-L84
Options for wrk2 - https://github.com/layer5io/gowrk2/blob/v0.6.1/api/gowrk2.go#L47-L53
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile-new --url "https://google.com" --options [filepath|json-string]

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile-new --url "https://google.com" --options path/to/options.json

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile-new --url "https://google.com" --load-generator nighthawk --options '{"requests_per_second": 10, "max_pending_requests": 5}'

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile-new --url "https://google.com" --load-generator fortio --options '{"MethodOverride": "POST"}'

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf apply meshery-profile-new --url "https://google.com" --load-generator wrk2 --options '{"DurationInSeconds": 15, "Thread": 3}'

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -b, --body string                  (optional) Load test body. Can be a filepath/string
      --cert-path string             (optional) Path to the certificate to be used for the load test
      --concurrent-requests string   (optional) Number of Parallel Requests
      --disable-cert                 (optional) Do not use certificate present in the profile
      --duration string              (optional) Length of test (e.g. 10s, 5m, 2h). For more, see https://golang.org/pkg/time/#ParseDuration
  -f, --file string                  (optional) File containing SMP-compatible test configuration. For more, see https://github.com/layer5io/service-mesh-performance-specification
  -h, --help                         help for apply
      --load-generator string        (optional) Load-Generator to be used (fortio/wrk2/nighthawk)
      --mesh string                  (optional) Name of the Service Mesh
      --name string                  (optional) Name of the Test
      --options string               (optional) Additional options to be passed to the load generator. Can be a json string or a filepath containing json
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
