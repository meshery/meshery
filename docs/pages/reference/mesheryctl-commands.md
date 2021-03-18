---
layout: default
title: mesheryctl Command Reference
abstract: "A guide to Meshery's CLI: mesheryctl"
permalink: reference/mesheryctl
redirect_from: 
 - guides/mesheryctl
 - guides/mesheryctl-commands
type: Reference
---
## Categories and Command Structure

Meshery CLI commands are categorized by function, which are:

- `mesheryctl` - Global flags
- `mesheryctl system` - Meshery Lifecycle and Troubleshooting
- `mesheryctl mesh` - Service Mesh Lifecycle & Configuration Management
- `mesheryctl perf` -  Service Mesh Performance Management
- `mesheryctl pattern` - Service Mesh Pattern Configuration & Management

## Global Commands and Flags

<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Command</th>
    <th>Flag</th>
    <th>Function</th>
    <th>Usage</th>
  </tr>
  <tr>
    <td rowspan="7" align="center" bgcolor="gainsboro">mesheryctl</td>
    <td>version</td>
    <td></td>
    <td>displays the version of mesheryctl and the SHA of the release binary</td>
    <td><code>mesheryctl version</code></td>
  </tr>
  <tr>
    <td rowspan="3"></td>
    <td rowspan="3">help, h</td>
    <td rowspan="3">displays help for any command</td>
    <td><code>mesheryctl --help</code></td>
  </tr>
  <tr>
    <td><code>mesheryctl system --help</code></td>
  </tr>
  <tr>
    <td><code>mesheryctl system start --help</code></td>
  </tr>
  <tr>
    <td></td>
    <td>config</td>
    <td>path to mesheryctl config file(~/.meshery/config.yaml)</td>
    <td><code>mesheryctl system reset --config=&lt;path to config file&gt;</code></td>
  </tr>
  <tr>
    <td rowspan="2"></td>
    <td rowspan="2">verbose, v</td>
    <td rowspan="2">displays verbose/debug logs</td>
    <td rowspan="2"><code>mesheryctl system update --verbose</code></td>
  </tr>
  <tr>
  </tr>
</thead>
</table>

## Meshery Lifecycle Management and Troubleshooting

Installation, troubleshooting and debugging of Meshery and its adapters.

<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Arguments</th>
    <th>Flag</th>
    <th>Function</th>
    <th>Usage</th>
  </tr>
  <tr>
    <td rowspan="14" align="center" bgcolor="gainsboro">system</td>
    <td></td>
    <td>context, c</td>
    <td></td>
    <td><code>mesheryctl system reset -c &lt;temporary context name&gt;</code></td>
  </tr>
  <tr>
    <td>completion</td>
    <td></td>
    <td>generates completion script</td>
    <td><code>mesheryctl system completion [bash\|zsh\|fish]</code></td>
  </tr>
  <tr>
    <td rowspan="4" bgcolor="lightgray">start</td>
    <td></td>
    <td>start Meshery</td>
    <td><code>mesheryctl system start</code></td>
  </tr>
  <tr>
    <td>skip-update</td>
    <td>start Meshery- skip checking for new Meshery container images</td>
    <td><code>mesheryctl system start --skip-update</code></td>
  </tr>
  <tr>
    <td>reset</td>
    <td>start Meshery- reset Meshery's configuration file to default settings</td>
    <td><code>mesheryctl system start --reset</code></td>
  </tr>
  <tr>
    <td>silent</td>
    <td>start Meshery- silently create Meshery's configuration file with default settings</td>
    <td><code>mesheryctl system start --silent</code></td>
  </tr>
  <tr>
    <td rowspan="2" bgcolor="lightgray">stop</td>
    <td></td>
    <td>stop Meshery</td>
    <td><code>mesheryctl system stop</code></td>
  </tr>
  <tr>
    <td>reset</td>
    <td>stop Meshery- reset Meshery's configuration file to default settings</td>
    <td><code>mesheryctl system stop --reset</code></td>
  </tr>
  <tr>
    <td>restart</td>
    <td></td>
    <td>restart all Meshery containers, their instances and their connected volumes</td>
    <td><code>mesheryctl system restart</code></td>
  </tr>
  <tr>
    <td>reset</td>
    <td></td>
    <td>resets meshery.yaml file with a copy from Meshery repo</td>
    <td><code>mesheryctl system reset</code></td>
  </tr>
  <tr>
    <td>update</td>
    <td></td>
    <td>pulls new Meshery images from Docker Hub. Does not update mesheryctl. This command can be run while Meshery is running. </td>
    <td><code>mesheryctl system update</code></td>
  </tr>
  <tr>
    <td>config</td>
    <td></td>
    <td>configures Meshery with kubeconfig. Generated with the help of user details to provide access for public clouds(GKE/EKS)</td>
    <td><code>mesheryctl system config gke --token &lt;path to token&gt;</code></td>
  </tr>
  <tr>
    <td>log</td>
    <td></td>
    <td>starts tailing Meshery server debug logs</td>
    <td><code>mesheryctl system log</code></td>
  </tr>
  <tr>
    <td>status</td>
    <td></td>
    <td>display the status of the Meshery containers</td>
    <td><code>mesheryctl system status</code></td>
  </tr>
  <tr>
    <td rowspan="4" align="center" bgcolor="gainsboro">system channel</td>
    <td>set</td>
    <td></td>
    <td>sets release channel and version of context in focus</td>
    <td><code>mesheryctl system channel set</code></td>
  </tr>
  <tr>
    <td rowspan="2" bgcolor="lightgray">view</td>
    <td></td>
    <td>view release channel and version of context in focus</td>
    <td><code>mesheryctl system channel view</code></td>
  </tr>
  <tr>
    <td>all, a</td>
    <td>view release channel and version of all contexts</td>
    <td><code>mesheryctl system channel view --all</code></td>
  </tr>
  <tr>
    <td>switch</td>
    <td></td>
    <td>switch release channel and version of context in focus</td>
    <td><code>mesheryctl system channel switch</code></td>
  </tr>
  <tr>
    <td rowspan="10" align="center" bgcolor="gainsboro">system context</td>
    <td></td>
    <td></td>
    <td>display the current context</td>
    <td><code>mesheryctl system context</code></td>
  </tr>
  <tr>
    <td rowspan="4" bgcolor="lightgray">create</td>
    <td></td>
    <td>create a new context in config.yaml file</td>
    <td><code>mesheryctl system context create &lt;context name&gt;</code></td>
  </tr>
  <tr>
    <td>url, u</td>
    <td>create a new context in config.yaml file- set Meshery server URL. Defaults to "https://localhost:9081"</td>
    <td><code>mesheryctl system context create &lt;context name&gt; --url &lt;URL&gt;</code></td>
  </tr>
  <tr>
    <td>set, s</td>
    <td>create a new context in config.yaml file- set as current context</td>
    <td><code>mesheryctl system context create &lt;context name&gt; --set</code></td>
  </tr>
  <tr>
    <td>adapters, a</td>
    <td>create a new context in config.yaml file- specify the list of adapters to be added</td>
    <td><code>mesheryctl system context create &lt;context name&gt; --adapters &lt;list of adapters&gt;</code></td>
  </tr>
  <tr>
    <td>delete</td>
    <td></td>
    <td>delete the specified context from the config.yaml file</td>
    <td><code>mesheryctl system context delete &lt;context name&gt;</code></td>
  </tr>
  <tr>
    <td>switch</td>
    <td></td>
    <td>switch between contexts</td>
    <td><code>mesheryctl system context switch &lt;context name&gt;</code></td>
  </tr>
  <tr>
    <td rowspan="3" bgcolor="lightgray">view</td>
    <td></td>
    <td>view the configurations of the current context</td>
    <td><code>mesheryctl system context view</code></td>
  </tr>
  <tr>
    <td>--context</td>
    <td>view the configurations of the specified context</td>
    <td><code>mesheryctl system context view --context &lt;context name&gt;</code></td>
  </tr>
  <tr>
    <td>--all</td>
    <td>if set, shows the configurations of all the contexts</td>
    <td><code>mesheryctl system context view --all</code></td>
  </tr>
</thead>
</table>

## Service Mesh Performance Management

<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Flag</th>
    <th>Function</th>
    <th>Usage</th>
  </tr>
  <tr>
    <td rowspan="10" align="center" bgcolor="gainsboro">perf</td>
    <td></td>
    <td>performance management- baseline and testing</td>
    <td></td>
  </tr>
  <tr>
    <td>url</td>
    <td>(required) endpoint URL to test</td>
    <td><code>mesheryctl perf --url &lt;URL&gt;</code></td>
  </tr>
  <tr>
    <td>name</td>
    <td>name of the test</td>
    <td><code>mesheryctl perf --name "&lt;name&gt;" --url &lt;URL&gt;</code></td>
  </tr>
  <tr>
    <td>mesh</td>
    <td>name of the service mesh</td>
    <td><code>mesheryctl perf --mesh &lt;name&gt; --url &lt;URL&gt;</code></td>
  </tr>
  <tr>
    <td>qps</td>
    <td>queries per second</td>
    <td><code>mesheryctl perf --qps &lt;queries&gt; --url &lt;URL&gt;</code></td>
  </tr>
  <tr>
    <td>concurrent-requests</td>
    <td>number of parallel requests</td>
    <td><code>mesheryctl perf --concurrent-requests &lt;number of requests&gt; --url &lt;URL&gt;</code></td>
  </tr>
  <tr>
    <td>duration</td>
    <td>length of the test (e.g. 10s, 5m, 2h). For more, see https://golang.org/pkg/time/#ParseDuration</td>
    <td><code>mesheryctl perf --duration &lt;time&gt; --url &lt;URL&gt;</code></td>
  </tr>
  <tr>
    <td>token</td>
    <td>path to Meshery auth token</td>
    <td><code>mesheryctl perf --token &lt;path to token&gt; --url &lt;URL&gt;</code></td>
  </tr>
  <tr>
    <td>load-generator</td>
    <td>load generator to be used (fortio/wrk2)</td>
    <td><code>mesheryctl perf --load-generator [fortio/wrk2] --url &lt;URL&gt;</code></td>
  </tr>
  <tr>
    <td>file</td>
    <td>file containing SMP compatible test configuration. For more, see https://github.com/layer5io/service-mesh-performance-specification</td>
    <td><code>mesheryctl perf --file &lt;path to file&gt; --url &lt;URL&gt;</code></td>
  </tr>
</thead>
</table>


## Service Mesh Lifecycle and Configuration Management

<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Command</th>
    <th>Flag</th>
    <th>Function</th>
    <th>Usage</th>
  </tr>
  <tr>
    <td rowspan="5" align="center" bgcolor="gainsboro">mesh</td>
    <td rowspan="5">validate</td>
    <td></td>
    <td>validate service mesh conformance to different standard specifications </td>
    <td></td>
  </tr>
  <tr>
    <td>adapter, a</td>
    <td>(required) adapter to use for validation. Defaults to "meshery-osm:10010"</td>
    <td rowspan="3"><code>mesheryctl mesh validate --adapter &lt;name of the 
    adapter&gt; --tokenPath &lt;path to token for authentication&gt; 
    --spec &lt;specification to be used for conformance test&gt;</code></td>
  </tr>
  <tr>
    <td>spec, s</td>
    <td>(required) specification to be used for conformance test. Defaults to "smi"</td>
  </tr>
  <tr>
    <td>tokenPath, t</td>
    <td>(required) path to token for authenticating to Meshery API</td>
  </tr>
  <tr>
    <td>namespace, n</td>
    <td>Kubernetes namespace to be used for deploying the validation tests and sample workload</td>
    <td><code>mesheryctl mesh validate --adapter &lt;name of the adapter
    &gt; --tokenPath &lt;path to token for authentication&gt; --spec &lt;
    specification to be used for conformance test&gt; --namespace 
    &lt;namespace to be used&gt;</code></td>
  </tr>
</thead>
</table>

## Service Mesh Pattern Configuration and Management

<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Command</th>
    <th>Flag</th>
    <th>Function</th>
    <th>Usage</th>
  </tr>
  <tr>
    <td rowspan="3" align="center" bgcolor="gainsboro">pattern</td>
    <td></td>
    <td>file, f</td>
    <td>(required) path to pattern file</td>
    <td></td>
  </tr>
  <tr>
    <td>apply</td>
    <td></td>
    <td>apply pattern file</td>
    <td><code>mesheryctl exp pattern apply --file &lt;path to pattern file&gt;</code></td>
  </tr>
  <tr>
    <td>delete</td>
    <td></td>
    <td>delete pattern file</td>
    <td><code>mesheryctl exp pattern delete --file &lt;path to pattern file&gt;</code></td>
  </tr>
</thead>
</table>