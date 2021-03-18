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
    <td rowspan="7" align="center">mesheryctl</td>
    <td>version</td>
    <td></td>
    <td>displays the version of mesheryctl and the SHA of the release binary</td>
    <td>`mesheryctl version`</td>
  </tr>
  <tr>
    <td rowspan="3"></td>
    <td rowspan="3">help, h</td>
    <td rowspan="3">displays help for any command</td>
    <td>`mesheryctl --help`</td>
  </tr>
  <tr>
    <td>`mesheryctl system --help`</td>
  </tr>
  <tr>
    <td>`mesheryctl system start --help`</td>
  </tr>
  <tr>
    <td></td>
    <td>config</td>
    <td>path to mesheryctl config file(~/.meshery/config.yaml)</td>
    <td>`mesheryctl system reset --config=&lt;path to config file&gt;`</td>
  </tr>
  <tr>
    <td rowspan="2"></td>
    <td rowspan="2">verbose, v</td>
    <td rowspan="2">displays verbose/debug logs</td>
    <td rowspan="2">`mesheryctl system update --verbose`</td>
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
    <td rowspan="14" align="center">system</td>
    <td></td>
    <td>context, c</td>
    <td></td>
    <td>mesheryctl system reset -c &lt;temporary context name&gt;</td>
  </tr>
  <tr>
    <td>completion</td>
    <td></td>
    <td>generates completion script</td>
    <td>mesheryctl system completion [bash\|zsh\|fish]</td>
  </tr>
  <tr>
    <td rowspan="4">start</td>
    <td></td>
    <td>start Meshery</td>
    <td>mesheryctl system start</td>
  </tr>
  <tr>
    <td>skip-update</td>
    <td>start Meshery- skip checking for new Meshery container images</td>
    <td>mesheryctl system start --skip-update</td>
  </tr>
  <tr>
    <td>reset</td>
    <td>start Meshery- reset Meshery's configuration file to default settings</td>
    <td>mesheryctl system start --reset</td>
  </tr>
  <tr>
    <td>silent</td>
    <td>start Meshery- silently create Meshery's configuration file with default settings</td>
    <td>mesheryctl system start --silent</td>
  </tr>
  <tr>
    <td rowspan="2">stop</td>
    <td></td>
    <td>stop Meshery</td>
    <td>mesheryctl system stop</td>
  </tr>
  <tr>
    <td>reset</td>
    <td>stop Meshery- reset Meshery's configuration file to default settings</td>
    <td>mesheryctl system stop --reset</td>
  </tr>
  <tr>
    <td>restart</td>
    <td></td>
    <td>restart all Meshery containers, their instances and their connected volumes</td>
    <td>mesheryctl system restart</td>
  </tr>
  <tr>
    <td>reset</td>
    <td></td>
    <td>resets meshery.yaml file with a copy from Meshery repo</td>
    <td>mesheryctl system reset</td>
  </tr>
  <tr>
    <td>update</td>
    <td></td>
    <td>pulls new Meshery images from Docker Hub. Does not update mesheryctl. This command can be run while Meshery is running. </td>
    <td>mesheryctl system update</td>
  </tr>
  <tr>
    <td>config</td>
    <td></td>
    <td>configures Meshery with kubeconfig. Generated with the help of user details to provide access for public clouds(GKE/EKS)</td>
    <td>mesheryctl system config gke --token &lt;path to token&gt;</td>
  </tr>
  <tr>
    <td>log</td>
    <td></td>
    <td>starts tailing Meshery server debug logs</td>
    <td>mesheryctl system log</td>
  </tr>
  <tr>
    <td>status</td>
    <td></td>
    <td>display the status of the Meshery containers</td>
    <td>mesheryctl system status</td>
  </tr>
  <tr>
    <td rowspan="4" align="center">system channel</td>
    <td>set</td>
    <td></td>
    <td>sets release channel and version of context in focus</td>
    <td>mesheryctl system channel set</td>
  </tr>
  <tr>
    <td rowspan="2">view</td>
    <td></td>
    <td>view release channel and version of context in focus</td>
    <td>mesheryctl system channel view</td>
  </tr>
  <tr>
    <td>all, a</td>
    <td>view release channel and version of all contexts</td>
    <td>mesheryctl system channel view --all</td>
  </tr>
  <tr>
    <td>switch</td>
    <td></td>
    <td>switch release channel and version of context in focus</td>
    <td>mesheryctl system channel switch</td>
  </tr>
  <tr>
    <td rowspan="10" align="center">system context</td>
    <td></td>
    <td></td>
    <td>display the current context</td>
    <td>mesheryctl system context</td>
  </tr>
  <tr>
    <td rowspan="4">create</td>
    <td></td>
    <td>create a new context in config.yaml file</td>
    <td>mesheryctl system context create &lt;context name&gt;</td>
  </tr>
  <tr>
    <td>url, u</td>
    <td>create a new context in config.yaml file- set Meshery server URL. Defaults to "https://localhost:9081"</td>
    <td>mesheryctl system context create &lt;context name&gt; --url &lt;URL&gt;</td>
  </tr>
  <tr>
    <td>set, s</td>
    <td>create a new context in config.yaml file- set as current context</td>
    <td>mesheryctl system context create &lt;context name&gt; --set</td>
  </tr>
  <tr>
    <td>adapters, a</td>
    <td>create a new context in config.yaml file- specify the list of adapters to be added</td>
    <td>mesheryctl system context create &lt;context name&gt; --adapters &lt;list of adapters&gt;</td>
  </tr>
  <tr>
    <td>delete</td>
    <td></td>
    <td>delete the specified context from the config.yaml file</td>
    <td>mesheryctl system context delete &lt;context name&gt;</td>
  </tr>
  <tr>
    <td>switch</td>
    <td></td>
    <td>switch between contexts</td>
    <td>mesheryctl system context switch &lt;context name&gt;</td>
  </tr>
  <tr>
    <td rowspan="3">view</td>
    <td></td>
    <td>view the configurations of the current context</td>
    <td>mesheryctl system context view</td>
  </tr>
  <tr>
    <td>--context</td>
    <td>view the configurations of the specified context</td>
    <td>mesheryctl system context view --context &lt;context name&gt;</td>
  </tr>
  <tr>
    <td>--all</td>
    <td>if set, shows the configurations of all the contexts</td>
    <td>mesheryctl system context view --all</td>
  </tr>
</thead>
</table>

## Service Mesh Performance Management

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
    <td rowspan="10" align="center">perf</td>
    <td></td>
    <td></td>
    <td>performance management- baseline and testing</td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td>url</td>
    <td>(required) endpoint URL to test</td>
    <td>mesheryctl perf --url &lt;URL&gt;</td>
  </tr>
  <tr>
    <td></td>
    <td>name</td>
    <td>name of the test</td>
    <td>mesheryctl perf --name "&lt;name&gt;" --url &lt;URL&gt;</td>
  </tr>
  <tr>
    <td></td>
    <td>mesh</td>
    <td>name of the service mesh</td>
    <td>mesheryctl perf --mesh &lt;name&gt; --url &lt;URL&gt;</td>
  </tr>
  <tr>
    <td></td>
    <td>qps</td>
    <td>queries per second</td>
    <td>mesheryctl perf --qps &lt;queries&gt; --url &lt;URL&gt;</td>
  </tr>
  <tr>
    <td></td>
    <td>concurrent-requests</td>
    <td>number of parallel requests</td>
    <td>mesheryctl perf --concurrent-requests &lt;number of requests&gt; --url &lt;URL&gt;</td>
  </tr>
  <tr>
    <td></td>
    <td>duration</td>
    <td>length of the test (e.g. 10s, 5m, 2h). For more, see https://golang.org/pkg/time/#ParseDuration</td>
    <td>mesheryctl perf --duration &lt;time&gt; --url &lt;URL&gt;</td>
  </tr>
  <tr>
    <td></td>
    <td>token</td>
    <td>path to Meshery auth token</td>
    <td>mesheryctl perf --token &lt;path to token&gt; --url &lt;URL&gt;</td>
  </tr>
  <tr>
    <td></td>
    <td>load-generator</td>
    <td>load generator to be used (fortio/wrk2)</td>
    <td>mesheryctl perf --load-generator [fortio/wrk2] --url &lt;URL&gt;</td>
  </tr>
  <tr>
    <td></td>
    <td>file</td>
    <td>file containing SMP compatible test configuration. For more, see https://github.com/layer5io/service-mesh-performance-specification</td>
    <td>mesheryctl perf --file &lt;path to file&gt; --url &lt;URL&gt;</td>
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
    <td rowspan="5" align="center">mesh</td>
    <td rowspan="5">validate</td>
    <td></td>
    <td>validate service mesh conformance to different standard specifications </td>
    <td></td>
  </tr>
  <tr>
    <td>adapter, a</td>
    <td>(required) adapter to use for validation. Defaults to "meshery-osm:10010"</td>
    <td rowspan="3">mesheryctl mesh validate --adapter &lt;name of the adapter&gt; --tokenPath &lt;path to token for authentication&gt; --spec &lt;specification to be used for conformance test&gt;</td>
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
    <td>mesheryctl mesh validate --adapter &lt;name of the adapter&gt; --tokenPath &lt;path to token for authentication&gt; --spec &lt;specification to be used for conformance test&gt; --namespace &lt;namespace to be used&gt;</td>
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
    <td rowspan="3" align="center">pattern</td>
    <td></td>
    <td>file, f</td>
    <td>(required) path to pattern file</td>
    <td></td>
  </tr>
  <tr>
    <td>apply</td>
    <td></td>
    <td>apply pattern file</td>
    <td>mesheryctl exp pattern apply --file &lt;path to pattern file&gt;</td>
  </tr>
  <tr>
    <td>delete</td>
    <td></td>
    <td>delete pattern file</td>
    <td>mesheryctl exp pattern delete --file &lt;path to pattern file&gt;</td>
  </tr>
</thead>
</table>