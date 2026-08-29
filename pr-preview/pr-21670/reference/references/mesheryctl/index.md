# Command Line Reference

> A guide to Meshery's CLI: mesheryctl

Source: /pr-preview/pr-21670/reference/references/mesheryctl/

## Categories and Command Structure

Meshery CLI commands are categorized by function, which are:

- [`mesheryctl`](#global-commands-and-flags) - Global flags and CLI configuration
- [`mesheryctl system`](#meshery-lifecycle-management-and-troubleshooting) - Meshery Lifecycle and Troubleshooting
- [`mesheryctl adapter`](#infrastructure-lifecycle-and-configuration-management) - Lifecycle & Configuration Management: provisioning and configuration best practices
- [`mesheryctl perf`](#service-performance-management) - Performance Management: Workload and cloud native performance characterization
- [`mesheryctl design`](#infrastructure-design-configuration-and-management) - Design Patterns: Cloud native patterns and best practices
- [`mesheryctl filter`](#data-plane-intelligence) - Data Plane Intelligence: Registry and configuration of WebAssembly filters for Envoy
- [`mesheryctl model`](#meshery-models) - A unit of packaging to define managed infrastructure and their relationships, and details specifics of how to manage them.
- [`mesheryctl component`](#meshery-components) - Fundamental building block used to represent and define the infrastructure under management
- [`mesheryctl registry`](#meshery-registry-management) - Model Database: Manage the state and contents of Meshery's internal registry of capabilities.
- [`mesheryctl environment`](#meshery-environment) - Logical group of connections and their associated credentials.
- [`mesheryctl connection`](#meshery-connection) - Managed or unmanaged resources that either through discovery or manual entry are tracked by Meshery.
- [`mesheryctl organization`](#meshery-organization) - Manage and interact with registered organizations.
- [`mesheryctl relationship`](#meshery-relationship) - View and manage relationships registered in Meshery.
- [`mesheryctl workspace`](#meshery-workspace) - View and manage workspaces under an organization.
- [`mesheryctl exp`](#experimental-featuresexp) - Experimental features

## Global Commands and Flags

<table class="mesheryctl-reference-table">
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
</thead>
<tbody>
      
      <tr><td rowspan="100"><a href="/pr-preview/pr-21670/reference/references/mesheryctl/">mesheryctl</a></td>
        <td></td>
        <td></td>
        <td>A global command that displays an overview of all commands.</td>
      </tr>

      
      <tr>
        <td></td>
        <td>--config</td>
        <td>configures Meshery with the meshconfig, generated with the help of user details to provide cluster access for public clouds(GKE/EKS).</td>
      </tr>
      <tr>
        <td></td>
        <td>--help, -h</td>
        <td>Displays helpful information about any command.</td>
      </tr>
      <tr>
        <td></td>
        <td>--verbose, -v</td>
        <td>Sets the log level to debug for verbose output and displays verbose/debug logs.</td>
      </tr>

      
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/completion/">completion</a></td>
        <td></td>
        <td>Output shell completion code for the specified shell.</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/version/">version</a></td>
        <td></td>
        <td>Displays the version of the Meshery Client (mesheryctl) and the SHA of the release binary.</td>
      </tr>
</tbody>
</table>


## Meshery Lifecycle Management and Troubleshooting

Installation, troubleshooting and debugging of Meshery and its adapters

<table class="mesheryctl-reference-table">
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
</thead>
<tbody>
      
      <tr><td rowspan="100"><a href="/pr-preview/pr-21670/reference/references/mesheryctl/system/">system</a></td>
        <td></td>
        <td></td>
        <td>Lifecycle management of Meshery deployments</td>
      </tr>

      
      <tr>
        <td></td>
        <td>--context, -c</td>
        <td>Temporarily override your current context by specifying an alternative context as an argument.</td>
      </tr>

      
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/system/check/">check</a></td>
        <td></td>
        <td>Run system checks for both pre and post mesh deployment scenarios on Meshery</td>
      </tr>
      <tr>
        <td></td>
        <td>--adapter</td>
        <td>Run checks on specific mesh adapter</td>
      </tr>
      <tr>
        <td></td>
        <td>--adapters</td>
        <td>Check status of meshery adapters</td>
      </tr>
      <tr>
        <td></td>
        <td>--components</td>
        <td>Check status of Meshery components</td>
      </tr>
      <tr>
        <td></td>
        <td>--operator</td>
        <td>Verify the health of Meshery Operator&#39;s deployment with MeshSync and Broker</td>
      </tr>
      <tr>
        <td></td>
        <td>--preflight, --pre</td>
        <td>Run Pre-mesh deployment checks (Docker and Kubernetes)</td>
      </tr>
      <tr><td>config</td>
        <td></td>
        <td>Configures Meshery to use a Kubernetes cluster.</td>
      </tr>
      <tr>
        <td></td>
        <td>--token</td>
        <td>Path to token for authenticating to Meshery API.</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/system/dashboard/">dashboard</a></td>
        <td></td>
        <td>Open Meshery UI in browser.</td>
      </tr>
      <tr>
        <td></td>
        <td>--port-forward</td>
        <td>(optional) Use port forwarding to access Meshery UI</td>
      </tr>
      <tr>
        <td></td>
        <td>--skip-browser</td>
        <td>(optional) skip opening of MesheryUI in browser.</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/system/login/">login</a></td>
        <td></td>
        <td>Authenticate with the Meshery Provider of your choice: the Local Provider or a Remote Provider.</td>
      </tr>
      <tr>
        <td></td>
        <td>--provider, -p</td>
        <td>Login Meshery with specified provider</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/system/logout/">logout</a></td>
        <td></td>
        <td>Invalidate current session with your Meshery Provider.</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/system/logs/">logs</a></td>
        <td></td>
        <td>Starts tailing Meshery server debug logs</td>
      </tr>
      <tr>
        <td></td>
        <td>--follow, -f</td>
        <td>(Optional) Follow the stream of the Meshery&#39;s logs. Defaults to false.</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/system/reset/">reset</a></td>
        <td></td>
        <td>Resets meshery.yaml file with a copy from Meshery repo</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/system/restart/">restart</a></td>
        <td></td>
        <td>restart all Meshery containers, their instances and their connected volumes</td>
      </tr>
      <tr>
        <td></td>
        <td>--skip-update</td>
        <td>(optional) skip checking for new updates available in Meshery.</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/system/start/">start</a></td>
        <td></td>
        <td>Start Meshery</td>
      </tr>
      <tr>
        <td></td>
        <td>--platform, -p</td>
        <td>platform to deploy Meshery to.</td>
      </tr>
      <tr>
        <td></td>
        <td>--reset</td>
        <td>Reset Meshery’s configuration file to default settings.</td>
      </tr>
      <tr>
        <td></td>
        <td>--skip-browser</td>
        <td>(optional) skip opening of MesheryUI in browser.</td>
      </tr>
      <tr>
        <td></td>
        <td>--skip-update</td>
        <td>(optional) skip checking for new updates available in Meshery.</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/system/status/">status</a></td>
        <td></td>
        <td>Check status of Meshery, Meshery adapters, Meshery Operator and its controllers.</td>
      </tr>
      <tr>
        <td></td>
        <td>--verbose, -v</td>
        <td>(optional) Extra data in status table</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/system/stop/">stop</a></td>
        <td></td>
        <td>Stop Meshery</td>
      </tr>
      <tr>
        <td></td>
        <td>--force</td>
        <td>(optional) uninstall Meshery resources forcefully</td>
      </tr>
      <tr>
        <td></td>
        <td>--keep-namespace</td>
        <td>(optional) keep the Meshery namespace during uninstallation</td>
      </tr>
      <tr>
        <td></td>
        <td>--reset</td>
        <td>Reset Meshery’s configuration file to default settings.</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/system/update/">update</a></td>
        <td></td>
        <td>Pull new Meshery images from Docker Hub. Does not update `mesheryctl`. This command may be executed while Meshery is running.</td>
      </tr>
      <tr>
        <td></td>
        <td>--skip-reset</td>
        <td>Pull the latest manifest files</td>
      </tr>
</tbody>
</table>


## Service Performance Management

<table class="mesheryctl-reference-table">
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
</thead>
<tbody>
      
      <tr><td rowspan="100"><a href="/pr-preview/pr-21670/reference/references/mesheryctl/perf/">perf</a></td>
        <td></td>
        <td></td>
        <td>Performance Management and Benchmarking using Meshery CLI</td>
      </tr>

      
      <tr>
        <td></td>
        <td>--output-format, -o</td>
        <td>format to display in json or yaml.</td>
      </tr>
      <tr>
        <td></td>
        <td>--token</td>
        <td>(required) Path to Meshery user&#39;s access token.</td>
      </tr>

      
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/perf/apply/">apply</a></td>
        <td></td>
        <td>Runs Performance test using existing profiles or using flags.</td>
      </tr>
      <tr>
        <td></td>
        <td>--cert-path</td>
        <td>(optional) Path to the certificate to be used for the load test</td>
      </tr>
      <tr>
        <td></td>
        <td>--concurrent-requests</td>
        <td>Number of Parallel Requests (default: 1).</td>
      </tr>
      <tr>
        <td></td>
        <td>--disable-cert</td>
        <td>(optional) Do not use certificate present in the profile</td>
      </tr>
      <tr>
        <td></td>
        <td>--duration</td>
        <td>Length of the test.</td>
      </tr>
      <tr>
        <td></td>
        <td>--file</td>
        <td>Path to cloud native performance test configuration file (default: empty string).</td>
      </tr>
      <tr>
        <td></td>
        <td>--load-generator</td>
        <td>Choice of load generator - fortio (default: fortio).</td>
      </tr>
      <tr>
        <td></td>
        <td>--mesh</td>
        <td>Name of the system under test.</td>
      </tr>
      <tr>
        <td></td>
        <td>--name</td>
        <td>A memorable name for the test (default: random string).</td>
      </tr>
      <tr>
        <td></td>
        <td>--qps</td>
        <td>Queries per second (default: 0) (0 - means to use the CPU unbounded to generate as many requests as possible).</td>
      </tr>
      <tr>
        <td></td>
        <td>--url</td>
        <td>(required/optional) URL of the endpoint send load to during testing.</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/perf/profile/">profile</a></td>
        <td></td>
        <td>List all the available performance profiles</td>
      </tr>
      <tr>
        <td></td>
        <td>--page, -p</td>
        <td>(optional) List next set of performance profiles with --page (default = 1) (default 1).</td>
      </tr>
      <tr>
        <td></td>
        <td>--view</td>
        <td>(optional) View more information of the performance profile.</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/perf/result/">result</a></td>
        <td></td>
        <td>View the results of a performance profile.</td>
      </tr>
      <tr>
        <td></td>
        <td>--page, -p</td>
        <td>(optional) List next set of performance results with --page (default = 1) (default 1).</td>
      </tr>
      <tr>
        <td></td>
        <td>--view</td>
        <td>(optional) View more information of the performance test results.</td>
      </tr>
</tbody>
</table>


## Infrastructure Lifecycle and Configuration Management

<table class="mesheryctl-reference-table">
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
</thead>
<tbody>
      
      <tr><td rowspan="100"><a href="/pr-preview/pr-21670/reference/references/mesheryctl/adapter/">adapter</a></td>
        <td></td>
        <td></td>
        <td>Lifecycle management of cloud native infrastructure</td>
      </tr>

      

      
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/adapter/deploy/">deploy</a></td>
        <td></td>
        <td>Deploy infrastructure into the cluster</td>
      </tr>
      <tr>
        <td></td>
        <td>--adapter, -a</td>
        <td>(required) adapter to use for validation. Defaults to &#34;meshery-osm:10010&#34;</td>
      </tr>
      <tr>
        <td></td>
        <td>--namespace, -n</td>
        <td>Kubernetes namespace to be used for deploying the validation tests and sample workload</td>
      </tr>
      <tr>
        <td></td>
        <td>--tokenpath, -t</td>
        <td>(required) path to token for authenticating to Meshery API</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/adapter/remove/">remove</a></td>
        <td></td>
        <td>Remove infrastructure in the connected Kubernetes cluster</td>
      </tr>
      <tr>
        <td></td>
        <td>--namespace, -n</td>
        <td>Kubernetes namespace from which to remove infrastructure</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/adapter/validate/">validate</a></td>
        <td></td>
        <td>validate configuration conformance to different standard specifications</td>
      </tr>
      <tr>
        <td></td>
        <td>--adapter, -a</td>
        <td>(required) adapter to use for validation. Defaults to &#34;meshery-osm:10010&#34;</td>
      </tr>
      <tr>
        <td></td>
        <td>--namespace, -n</td>
        <td>Kubernetes namespace to be used for deploying the validation tests and sample workload</td>
      </tr>
      <tr>
        <td></td>
        <td>--spec, -s</td>
        <td>(required) specification to be used for conformance test. Defaults to &#34;smi&#34;</td>
      </tr>
      <tr>
        <td></td>
        <td>--tokenpath, -t</td>
        <td>(required) path to token for authenticating to Meshery API</td>
      </tr>
</tbody>
</table>


## Infrastructure Design Configuration and Management

<table class="mesheryctl-reference-table">
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
</thead>
<tbody>
      
      <tr><td rowspan="100"><a href="/pr-preview/pr-21670/reference/references/mesheryctl/design/">design</a></td>
        <td></td>
        <td></td>
        <td>Manage cloud native infrastructure using predefined application designs</td>
      </tr>

      

      
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/design/apply/">apply</a></td>
        <td></td>
        <td>apply design file will trigger deploy of the design file</td>
      </tr>
      <tr>
        <td></td>
        <td>--file, -f</td>
        <td>apply design file will trigger deploy of the design file and also supports file retrieval from GitHub</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/design/delete/">delete</a></td>
        <td></td>
        <td>Deprovision cloud native infrastructure using a design /design file</td>
      </tr>
      <tr>
        <td></td>
        <td>--file, -f</td>
        <td>Deletes the resources that were created, using design file</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/design/deploy/">deploy</a></td>
        <td></td>
        <td>Deploy application</td>
      </tr>
      <tr>
        <td></td>
        <td>--file, -f</td>
        <td>Deploy application with application file</td>
      </tr>
      <tr>
        <td></td>
        <td>--skip-save</td>
        <td>Skip saving an application</td>
      </tr>
      <tr>
        <td></td>
        <td>--source-type, -s</td>
        <td>Type of source file (ex. manifest / compose / helm)</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/design/export/">export</a></td>
        <td></td>
        <td>Export a design from Meshery</td>
      </tr>
      <tr>
        <td></td>
        <td>--output, -o</td>
        <td>Specify the output directory to save the design</td>
      </tr>
      <tr>
        <td></td>
        <td>--type</td>
        <td>Specify the design type to export (default current). Supported types are oci,current,original</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/design/import/">import</a></td>
        <td></td>
        <td>Import app manifests</td>
      </tr>
      <tr>
        <td></td>
        <td>--file, -f</td>
        <td>Path/URL to app file</td>
      </tr>
      <tr>
        <td></td>
        <td>--source-type, -s</td>
        <td>Type of source file (ex. manifest / compose / helm)</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/design/list/">list</a></td>
        <td></td>
        <td>Displays a list of available design files</td>
      </tr>
      <tr>
        <td></td>
        <td>--count, -c</td>
        <td>display the total number of available designs</td>
      </tr>
      <tr>
        <td></td>
        <td>--page, -p</td>
        <td>list designs by specific page number (10 designs per page)</td>
      </tr>
      <tr>
        <td></td>
        <td>--pagesize</td>
        <td>number of designs to be displayed per page (default 10)</td>
      </tr>
      <tr>
        <td></td>
        <td>--verbose, -v</td>
        <td>show all design file metadata</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/design/undeploy/">undeploy</a></td>
        <td></td>
        <td>Undeploy application</td>
      </tr>
      <tr>
        <td></td>
        <td>--file, -f</td>
        <td>Undeploy application with an application file</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/design/view/">view</a></td>
        <td></td>
        <td>Displays the contents of a specific design file</td>
      </tr>
      <tr>
        <td></td>
        <td>--all, -a</td>
        <td>Show all design file content</td>
      </tr>
      <tr>
        <td></td>
        <td>-o json</td>
        <td>Display the content of a design in JSON format</td>
      </tr>
</tbody>
</table>


## Data Plane Intelligence

<table class="mesheryctl-reference-table">
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
</thead>
<tbody>
      
      <tr><td rowspan="100"><a href="/pr-preview/pr-21670/reference/references/mesheryctl/filter/">filter</a></td>
        <td></td>
        <td></td>
        <td>Envoy Filter Management</td>
      </tr>

      

      
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/filter/delete/">delete</a></td>
        <td></td>
        <td>Delete filter from Meshery Server</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/filter/import/">import</a></td>
        <td></td>
        <td>Upload WASM filter file to Meshery Server</td>
      </tr>
      <tr>
        <td></td>
        <td>--name, -n</td>
        <td>(optional) filter name</td>
      </tr>
      <tr>
        <td></td>
        <td>--wasm-config, -w</td>
        <td>(optional) WASM configuration filepath/string</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/filter/list/">list</a></td>
        <td></td>
        <td>List all WASM filters</td>
      </tr>
      <tr>
        <td></td>
        <td>--verbose, -v</td>
        <td>Display full length user and filter file identifiers</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/filter/view/">view</a></td>
        <td></td>
        <td>View the specified WASM filter</td>
      </tr>
      <tr>
        <td></td>
        <td>--all, -a</td>
        <td>(optional) view all filters available</td>
      </tr>
      <tr>
        <td></td>
        <td>--output-format, -o</td>
        <td>(optional) format to display in [json|yaml] (default &#34;yaml&#34;)</td>
      </tr>
</tbody>
</table>


## Meshery Registry Management

<table class="mesheryctl-reference-table">
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
</thead>
<tbody>
      
      <tr><td rowspan="100"><a href="/pr-preview/pr-21670/reference/references/mesheryctl/registry/">registry</a></td>
        <td></td>
        <td></td>
        <td>Manage the state and contents of Meshery&#39;s internal registry of capabilities.</td>
      </tr>

      
      <tr>
        <td></td>
        <td>--help, -h</td>
        <td>help for registry</td>
      </tr>

      
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/registry/generate/">generate</a></td>
        <td></td>
        <td>Generate Models to the registry</td>
      </tr>
      <tr>
        <td></td>
        <td>--output-format, -o</td>
        <td>(optional) format to display in [md | mdx | js]</td>
      </tr>
      <tr>
        <td></td>
        <td>--registrant-cred</td>
        <td>path pointing to the registrant credential definition</td>
      </tr>
      <tr>
        <td></td>
        <td>--registrant-def</td>
        <td>path pointing to the registrant connection definition</td>
      </tr>
      <tr>
        <td></td>
        <td>--spreadsheet-cred</td>
        <td>base64 encoded credential to download the spreadsheet</td>
      </tr>
      <tr>
        <td></td>
        <td>--spreadsheet-id</td>
        <td>id of the spreadsheet to download</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/registry/publish/">publish</a></td>
        <td></td>
        <td>Publish Meshery Models to Websites, Remote Provider, Meshery Server</td>
      </tr>
      <tr>
        <td></td>
        <td>--output-format, -o</td>
        <td>(optional) format to display in [md | mdx | js]</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/registry/update/">update</a></td>
        <td></td>
        <td>Update the registry with latest data.</td>
      </tr>
      <tr>
        <td></td>
        <td>--input, -i</td>
        <td>(optional) path to capability file</td>
      </tr>
      <tr>
        <td></td>
        <td>--spreadsheet-cred</td>
        <td>base64 encoded credential to download the spreadsheet</td>
      </tr>
      <tr>
        <td></td>
        <td>--spreadsheet-id</td>
        <td>id of the spreadsheet to download</td>
      </tr>
</tbody>
</table>


## Meshery Models

<table class="mesheryctl-reference-table">
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
</thead>
<tbody>
      
      <tr><td rowspan="100"><a href="/pr-preview/pr-21670/reference/references/mesheryctl/model/">model</a></td>
        <td></td>
        <td></td>
        <td>A unit of packaging to define managed infrastructure and their relationships, and details specifics of how to manage them.</td>
      </tr>

      
      <tr>
        <td></td>
        <td>--count</td>
        <td>(optional) Get the number of models in total</td>
      </tr>
      <tr>
        <td></td>
        <td>--help, -h</td>
        <td>help for model</td>
      </tr>

      
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/model/build/">build</a></td>
        <td></td>
        <td>Create an OCI-compliant package from the model files</td>
      </tr>
      <tr>
        <td></td>
        <td>--path, -p</td>
        <td>(optional) looks for a model folder under specified path, defaults to current dir.</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/model/delete/">delete</a></td>
        <td></td>
        <td>Delete a registered model by ID</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/model/export/">export</a></td>
        <td></td>
        <td>Export a model by name to the given output type (default is oci).</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/model/generate/">generate</a></td>
        <td></td>
        <td>Generate models by specifying the directory, file, or URL.</td>
      </tr>
      <tr>
        <td></td>
        <td>--file, -f</td>
        <td>Specify the path to the file, directory, or URL.</td>
      </tr>
      <tr>
        <td></td>
        <td>--skip-registration</td>
        <td>Skip registration of the model (default is false).</td>
      </tr>
      <tr>
        <td></td>
        <td>--template, -t</td>
        <td>Specify the path to the template JSON file (only required for URLs).</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/model/import/">import</a></td>
        <td></td>
        <td>Import models using a file or filepath directly to the meshery registry</td>
      </tr>
      <tr>
        <td></td>
        <td>--file, -f</td>
        <td>(optional) path to a single model file to import from local filesystem or URL to import from web location</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/model/init/">init</a></td>
        <td></td>
        <td>Scaffold a folder structure for model creation</td>
      </tr>
      <tr>
        <td></td>
        <td>--output-format, -o</td>
        <td>(optional) format to scaffold files to, supports json or yaml, planned to support csv, defaults to json.</td>
      </tr>
      <tr>
        <td></td>
        <td>--path, -p</td>
        <td>(optional) generates a folder structure under specified path, defaults to current dir.</td>
      </tr>
      <tr>
        <td></td>
        <td>--version</td>
        <td>(optional) model version, defaults to v0.1.0</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/model/list/">list</a></td>
        <td></td>
        <td>List models present in the registry</td>
      </tr>
      <tr>
        <td></td>
        <td>--page, -p</td>
        <td>list models by specific page number (25 models per page)</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/model/search/">search</a></td>
        <td></td>
        <td>Using a keyword(s) search the registry for matching model(s)</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/model/view/">view</a></td>
        <td></td>
        <td>View details of a specific model</td>
      </tr>
      <tr>
        <td></td>
        <td>-o, --output-format</td>
        <td>(optional) format to display in [json|yaml] (default &#34;yaml&#34;)</td>
      </tr>
</tbody>
</table>


## Meshery Components

<table class="mesheryctl-reference-table">
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
</thead>
<tbody>
      
      <tr><td rowspan="100"><a href="/pr-preview/pr-21670/reference/references/mesheryctl/component/">component</a></td>
        <td></td>
        <td></td>
        <td>Fundamental building block used to represent and define the infrastructure under management</td>
      </tr>

      
      <tr>
        <td></td>
        <td>--count</td>
        <td>(optional) Get the number of components in total</td>
      </tr>
      <tr>
        <td></td>
        <td>--help, -h</td>
        <td>help for components</td>
      </tr>

      
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/component/list/">list</a></td>
        <td></td>
        <td>List all components registered in Meshery Server</td>
      </tr>
      <tr>
        <td></td>
        <td>--page, -p</td>
        <td>list components by specific page number (25 models per page)</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/component/search/">search</a></td>
        <td></td>
        <td>Using a keyword(s) search the registry for matching component(s)</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/component/view/">view</a></td>
        <td></td>
        <td>View details of a specific component</td>
      </tr>
      <tr>
        <td></td>
        <td>-o, --output-format</td>
        <td>(optional) format to display in [json|yaml] (default &#34;yaml&#34;)</td>
      </tr>
      <tr>
        <td></td>
        <td>--save</td>
        <td>Save output as a JSON/YAML file</td>
      </tr>
</tbody>
</table>


## Meshery Environment

<table class="mesheryctl-reference-table">
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
</thead>
<tbody>
      
      <tr><td rowspan="100"><a href="/pr-preview/pr-21670/reference/references/mesheryctl/environment/">environment</a></td>
        <td></td>
        <td></td>
        <td>Logical group related to connections and their credentials</td>
      </tr>

      

      
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/environment/create/">create</a></td>
        <td></td>
        <td>create a new environments by providing the name and description of the environment</td>
      </tr>
      <tr>
        <td></td>
        <td>-d, --description</td>
        <td>description of the specified environment</td>
      </tr>
      <tr>
        <td></td>
        <td>--name, -n</td>
        <td>(optional) name of the environment</td>
      </tr>
      <tr>
        <td></td>
        <td>-o, --orgId</td>
        <td>organization id</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/environment/delete/">delete</a></td>
        <td></td>
        <td>delete a new environments by providing the name and description of the environment</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/environment/list/">list</a></td>
        <td></td>
        <td>list name of all registered environments</td>
      </tr>
      <tr>
        <td></td>
        <td>-o, --orgId</td>
        <td>organization id</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/environment/view/">view</a></td>
        <td></td>
        <td>view a environments registered in Meshery Server</td>
      </tr>
      <tr>
        <td></td>
        <td>-o, --output-format</td>
        <td>(optional) format to display in [json|yaml] (default &#34;yaml&#34;)</td>
      </tr>
      <tr>
        <td></td>
        <td>-s, --save</td>
        <td>(optional) save output as a JSON/YAML file</td>
      </tr>
</tbody>
</table>


## Meshery Connection

<table class="mesheryctl-reference-table">
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
</thead>
<tbody>
      
      <tr><td rowspan="100"><a href="/pr-preview/pr-21670/reference/references/mesheryctl/connection/">connection</a></td>
        <td></td>
        <td></td>
        <td>Manage Meshery connections</td>
      </tr>

      
      <tr>
        <td></td>
        <td>--count</td>
        <td>(optional) Get the total number of connections</td>
      </tr>
      <tr>
        <td></td>
        <td>--help, -h</td>
        <td>help for connection</td>
      </tr>

      
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/connection/create/">create</a></td>
        <td></td>
        <td>Create a new connection to a Kubernetes cluster or other supported platform</td>
      </tr>
      <tr>
        <td></td>
        <td>--token</td>
        <td>(optional) Path to token for authenticating to Meshery API</td>
      </tr>
      <tr>
        <td></td>
        <td>-t, --type</td>
        <td>(required) type of connection to create (aks, eks, gke, minikube)</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/connection/delete/">delete</a></td>
        <td></td>
        <td>delete a connection</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/connection/list/">list</a></td>
        <td></td>
        <td>list all the connections</td>
      </tr>
      <tr>
        <td></td>
        <td>--count</td>
        <td>(optional) Get the total number of connections</td>
      </tr>
      <tr>
        <td></td>
        <td>-k, --kind</td>
        <td>(optional) Filter by kind (repeatable)</td>
      </tr>
      <tr>
        <td></td>
        <td>-p, --page</td>
        <td>(optional) List a set of connection in specific page (default=1).</td>
      </tr>
      <tr>
        <td></td>
        <td>--pagesize</td>
        <td>(optional) Number of connections per page (default=10).</td>
      </tr>
      <tr>
        <td></td>
        <td>-s, --status</td>
        <td>(optional) Filter by status (repeatable)</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/connection/view/">view</a></td>
        <td></td>
        <td>view a connection queried by the connection id</td>
      </tr>
      <tr>
        <td></td>
        <td>-o, --output-format</td>
        <td>(optional) format to display in [json|yaml] (default &#34;yaml&#34;)</td>
      </tr>
      <tr>
        <td></td>
        <td>-s, --save</td>
        <td>(optional) save output as a JSON/YAML file</td>
      </tr>
</tbody>
</table>


## Meshery Organization

<table class="mesheryctl-reference-table">
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
</thead>
<tbody>
      
      <tr><td rowspan="100"><a href="/pr-preview/pr-21670/reference/references/mesheryctl/organization/">organization</a></td>
        <td></td>
        <td></td>
        <td>Interact with existing organizations</td>
      </tr>

      
      <tr>
        <td></td>
        <td>--count</td>
        <td>(optional) Get the total of registered organizations</td>
      </tr>
      <tr>
        <td></td>
        <td>--help, -h</td>
        <td>help for organization command</td>
      </tr>

      
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/organization/list/">list</a></td>
        <td></td>
        <td>List all organizations registered in Meshery Server</td>
      </tr>
      <tr>
        <td></td>
        <td>--count</td>
        <td>(optional) Get the total of registered organizations</td>
      </tr>
      <tr>
        <td></td>
        <td>--page, -p</td>
        <td>list organizations by specific page number (25 models per page)</td>
      </tr>
</tbody>
</table>


## Meshery Relationship

<table class="mesheryctl-reference-table">
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
</thead>
<tbody>
      
      <tr><td rowspan="100"><a href="/pr-preview/pr-21670/reference/references/mesheryctl/relationship/">relationship</a></td>
        <td></td>
        <td></td>
        <td>View list of relationships and details of relationship</td>
      </tr>

      

      
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/relationship/generate/">generate</a></td>
        <td></td>
        <td>generate relationships docs from the google spreadsheets</td>
      </tr>
      <tr>
        <td></td>
        <td>-s, --sheetId</td>
        <td>google sheet id</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/relationship/list/">list</a></td>
        <td></td>
        <td>list all relationships registered in Meshery Server</td>
      </tr>
      <tr>
        <td></td>
        <td>-p, --page</td>
        <td>(optional) List next set of relationship list with --page (default = 1) (default 1).</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/relationship/search/">search</a></td>
        <td></td>
        <td>search for relationship using a query</td>
      </tr>
      <tr>
        <td></td>
        <td>-k --kind</td>
        <td>kind of relationships</td>
      </tr>
      <tr>
        <td></td>
        <td>-m, --model</td>
        <td>model name</td>
      </tr>
      <tr>
        <td></td>
        <td>-s, --subtype</td>
        <td>subtype of relationships</td>
      </tr>
      <tr>
        <td></td>
        <td>-t, --type</td>
        <td>type of relationships</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/relationship/view/">view</a></td>
        <td></td>
        <td>view a relationship queried by the model name</td>
      </tr>
      <tr>
        <td></td>
        <td>-o, --output-format</td>
        <td>(optional) format to display in [json|yaml] (default &#34;yaml&#34;)</td>
      </tr>
</tbody>
</table>


## Meshery Workspace

<table class="mesheryctl-reference-table">
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
</thead>
<tbody>
      
      <tr><td rowspan="100"><a href="/pr-preview/pr-21670/reference/references/mesheryctl/workspace/">workspace</a></td>
        <td></td>
        <td></td>
        <td>View list of workspaces and detail of workspaces</td>
      </tr>

      

      
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/workspace/create/">create</a></td>
        <td></td>
        <td>create a new workspace by providing the name, description, and organization ID</td>
      </tr>
      <tr>
        <td></td>
        <td>-d, --description</td>
        <td>description of the workspace</td>
      </tr>
      <tr>
        <td></td>
        <td>-n, --name</td>
        <td>name of the workspace</td>
      </tr>
      <tr>
        <td></td>
        <td>-o, --orgId</td>
        <td>organization id</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/workspace/list/">list</a></td>
        <td></td>
        <td>list name of all registered workspaces</td>
      </tr>
      <tr>
        <td></td>
        <td>-o, --orgId</td>
        <td>organization id</td>
      </tr>
      <tr><td><a href="/pr-preview/pr-21670/reference/references/mesheryctl/workspace/view/">view</a></td>
        <td></td>
        <td>view a workspace by its ID or name</td>
      </tr>
      <tr>
        <td></td>
        <td>--orgId</td>
        <td>(optional) organization ID to search workspace by name</td>
      </tr>
      <tr>
        <td></td>
        <td>-o, --output-format</td>
        <td>(optional) format to display in [json|yaml] (default &#34;yaml&#34;)</td>
      </tr>
      <tr>
        <td></td>
        <td>-s, --save</td>
        <td>(optional) save output as a JSON/YAML file</td>
      </tr>
</tbody>
</table>


## Experimental Features(exp)

<table class="mesheryctl-reference-table">
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
</thead>
<tbody>
    <tr>
      <td colspan="4">Command "exp" not found in data</td>
    </tr>
</tbody>
</table>


## Frequently Asked Questions for Meshery CLI

Refer to the following [frequently asked questions](/pr-preview/pr-21670/guides/mesheryctl/working-with-mesheryctl/) related to Meshery CLI.

## See Also
