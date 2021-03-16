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

| Main command | Command | Flag | Function | Usage |
|-	|-	|-	|-	|-	|
| mesheryctl 	| version 	|  	| displays the version of mesheryctl and the SHA of the release binary 	| mesheryctl version 	|
|  	|  	| help, h 	| displays help for any command 	| mesheryctl --help 	|
|  	|  	|  	|  	| mesheryctl system --help 	|
|  	|  	|  	|  	| mesheryctl system start --help 	|
|  	|  	| config 	| path to mesheryctl config file(~/.meshery/config.yaml) 	| mesheryctl system reset --config=<path to config file> 	|
|  	|  	| verbose, v 	| displays verbose/debug logs 	| mesheryctl system update --verbose 	|
|  	|  	|  	|  	|  	|

## Meshery Lifecycle Management and Troubleshooting

Installation, troubleshooting and debugging of Meshery and its adapters.

| Main command | Arguments | Flag | Function | Usage |
|-	|-	|-	|-	|-	|
| system 	|  	| context, c 	|  	| mesheryctl system reset -c <temporary context name> 	|
|  	| completion 	|  	| generates completion script 	| mesheryctl system completion [bash\|zsh\|fish] 	|
|  	| start 	|  	| start Meshery 	| mesheryctl system start 	|
|  	|  	| skip-update 	| start Meshery- skip checking for new Meshery container images 	| mesheryctl system start --skip-update 	|
|  	|  	| reset 	| start Meshery- reset Meshery's configuration file to default settings 	| mesheryctl system start --reset 	|
|  	|  	| silent 	| start Meshery- silently create Meshery's configuration file with default settings 	| mesheryctl system start --silent 	|
|  	| stop 	|  	| stop Meshery 	| mesheryctl system stop 	|
|  	|  	| reset 	| stop Meshery- reset Meshery's configuration file to default settings 	| mesheryctl system stop --reset 	|
|  	| restart 	|  	| restart all Meshery containers, their instances and their connected volumes 	| mesheryctl system restart 	|
|  	| reset 	|  	| resets meshery.yaml file with a copy from Meshery repo 	| mesheryctl system reset 	|
|  	| update 	|  	| pulls new Meshery images from Docker Hub. Does not update mesheryctl. This command can be run while Meshery is running.  	| mesheryctl system update 	|
|  	| config 	|  	| configures Meshery with kubeconfig. Generated with the help of user details to provide access for public clouds(GKE/EKS) 	| mesheryctl system config gke --token <path to token> 	|
|  	| log 	|  	| starts tailing Meshery server debug logs 	| mesheryctl system log 	|
|  	| status 	|  	| display the status of the Meshery containers 	| mesheryctl system status 	|
|-	|-	|-	|-	|-	|
| system channel 	| set 	|  	| sets release channel and version of context in focus 	| mesheryctl system channel set 	|
|  	| view 	|  	| view release channel and version of context in focus 	| mesheryctl system channel view 	|
|  	|  	| all, a 	| view release channel and version of all contexts 	| mesheryctl system channel view --all 	|
|  	| switch 	|  	| switch release channel and version of context in focus 	| mesheryctl system channel switch 	|
|-	|-	|-	|-	|-	|
| system context 	|  	|  	| display the current context 	| mesheryctl system context 	|
|  	| create 	|  	| create a new context in config.yaml file 	| mesheryctl system context create <context name> 	|
|  	|  	| url, u 	| create a new context in config.yaml file- set Meshery server URL. Defaults to "https://localhost:9081" 	| mesheryctl system context create <context name> --url <URL> 	|
|  	|  	| set, s 	| create a new context in config.yaml file- set as current context 	| mesheryctl system context create <context name> --set 	|
|  	|  	| adapters, a 	| create a new context in config.yaml file- specify the list of adapters to be added 	| mesheryctl system context create <context name> --adapters <list of adapters> 	|
|  	| delete 	|  	| delete the specified context from the config.yaml file 	| mesheryctl system context delete <context name> 	|
|  	| switch 	|  	| switch between contexts 	| mesheryctl system context switch <context name> 	|
|  	| view 	|  	| view the configurations of the current context 	| mesheryctl system context view 	|
|  	|  	| --context 	| view the configurations of the specified context 	| mesheryctl system context view --context <context name> 	|
|  	|  	| --all 	| if set, shows the configurations of all the contexts 	| mesheryctl system context view --all 	|

## Service Mesh Performance Management

| Main command | Command | Flag | Function | Usage |
|-	|-	|-	|-	|-	|
| perf 	|  	|  	| performance management- baseline and testing 	|  	|
|  	|  	| url 	| (required) endpoint URL to test 	| mesheryctl perf --url <URL> 	|
|  	|  	| name 	| name of the test 	| mesheryctl perf --name "<name>" --url <URL> 	|
|  	|  	| mesh 	| name of the service mesh 	| mesheryctl perf --mesh <name> --url <URL> 	|
|  	|  	| qps 	| queries per second 	| mesheryctl perf --qps <queries> --url <URL> 	|
|  	|  	| concurrent-requests 	| number of parallel requests 	| mesheryctl perf --concurrent-requests <number of requests> --url <URL> 	|
|  	|  	| duration 	| length of the test (e.g. 10s, 5m, 2h). For more, see https://golang.org/pkg/time/#ParseDuration 	| mesheryctl perf --duration <time> --url <URL> 	|
|  	|  	| token 	| path to Meshery auth token 	| mesheryctl perf --token <path to token> --url <URL> 	|
|  	|  	| load-generator 	| load generator to be used (fortio/wrk2) 	| mesheryctl perf --load-generator [fortio/wrk2] --url <URL> 	|
|  	|  	| file 	| file containing SMP compatible test configuration. For more, see https://github.com/layer5io/service-mesh-performance-specification 	| mesheryctl perf --file <path to file> --url <URL> 	|

## Service Mesh Lifecycle and Configuration Management

| Main command | Command | Flag | Function | Usage |
|-	|-	|-	|-	|-	|
| mesh 	| validate 	|  	| validate service mesh conformance to different standard specifications  	|  	|
|  	|  	| adapter, a 	| (required) adapter to use for validation. Defaults to "meshery-osm:10010" 	| mesheryctl mesh validate --adapter <name of the adapter> --tokenPath <path to token for authentication> --spec <specification to be used for conformance test> 	|
|  	|  	| spec, s 	| (required) specification to be used for conformance test. Defaults to "smi" 	|  	|
|  	|  	| tokenPath, t 	| (required) path to token for authenticating to Meshery API 	|  	|
|  	|  	| namespace, n 	| Kubernetes namespace to be used for deploying the validation tests and sample workload 	| mesheryctl mesh validate --adapter <name of the adapter> --tokenPath <path to token for authentication> --spec <specification to be used for conformance test> --namespace <namespace to be used> 	|

## Service Mesh Pattern Configuration and Management

| Main command | Command | Flag | Function | Usage |
|-	|-	|-	|-	|-	|
| pattern 	|  	| file, f 	| (required) path to pattern file 	|  	|
|  	| apply 	|  	| apply pattern file 	| mesheryctl exp pattern apply --file <path to pattern file> 	|
|  	| delete 	|  	| delete pattern file 	| mesheryctl exp pattern delete --file <path to pattern file> 	|