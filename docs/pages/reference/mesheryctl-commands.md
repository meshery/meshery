---
layout: default
title: mesheryctl Command Reference
abstract: "A guide to Meshery's CLI: mesheryctl"
permalink: reference/mesheryctl
#redirect_from: guides/mesheryctl
type: Reference
---

## Global Commands and Flags

| command    |        flag         | function                                                                                                                               | Usage                                                                   |
| :--------- | :-----------------: | :------------------------------------------------------------------------------------------------------------------------------------- | :---------------------------------------------------------------------- |
| mesheryctl |                     | Displays help about any command.                                                                                                       | `mesheryctl`                                                            |
|            |       version       | Displays the version of the Meshery Client (`mesheryctl`) and the SHA of the release binary.                                           | `mesheryctl system version`                                             |
|            |        help         | Displays help about any command.                                                                                                       | `mesheryctl --help`                                                     |
|            | --mesheryctl-config | (optional) path to Meshery Client (`mesheryctl`) configuration file (`~/.meshery/mesheryctl.yaml`) to overrides defaults.              | `mesheryctl perf <args> --mesheryctl-config=~/.meshery/mesheryctl.yaml` |
|            |      --config       | configures Meshery with the kubeconfig, generated with the help of user details, to provide cluster access for public clouds(GKE/EKS). | `mesheryctl system config gke --token "PATH TO TOKEN"`                  |

## Meshery Lifecycle Management

Installation, troubleshooting and debugging of Meshery and its adapters.

| command | arg    |     flag      | function                                                                                                                                | Usage                                                        |
| :------ | :----- | :-----------: | :-------------------------------------------------------------------------------------------------------------------------------------- | :----------------------------------------------------------- |
| system  | config |               | Configures Meshery to use a Kubernetes cluster.                                                                                         | `mesheryctl system config gke --token ~/Downloads/auth.json` |
|         | reset  |               | Resets `meshery.yaml` with copy from Meshery repo. _Warning: Any local changes will be overwritten._                                    | `mesheryctl system reset`                                    |
|         | log    |               | Starts tailing Meshery server debug logs.                                                                                               | `mesheryctl system log`                                      |
|         | start  |               | Start all Meshery containers.                                                                                                           | `mesheryctl system start`                                    |
|         |        |    --reset    | (optional) reset Meshery's configuration file to default settings.                                                                      | `mesheryctl system start --reset`                            |
|         |        | --skip-update | (optional) skip updates available in Meshery.                                                                                           | `mesheryctl system start --skip-update`                      |
|         | status |               | Displays the status of Meshery's containers (server and adapters).                                                                      | `mesheryctl system status`                                   |
|         | stop   |               | Stop all Meshery containers.                                                                                                            | `mesheryctl system stop`                                     |
|         |        |    --reset    | (optional) reset Meshery's configuration file to default settings.                                                                      | `mesheryctl system stop --reset`                             |
|         | update |               | Pull new Meshery images from Docker Hub. Does not pulls new `mesheryctl` client. This command may be executed while Meshery is running. | `mesheryctl system update`                                   |
|         | completion |               | Generates completion script.                                                                                               | `mesheryctl system completion [bash\|zsh\|fish]`              |
|         | help   |               | Displays help about any Meshery lifecycle management command.                                                                           | `mesheryctl system --help`                                   |

## Performance Management

| command |               flag               | function                                                                                                     | Usage                                                                                                                                                       |
| :------ | :------------------------------: | :----------------------------------------------------------------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| perf    |                                  | Performance management: baselining and testing.                                                              | `mesheryctl perf --name "a quick stress test" --url http://192.168.1.15/productpage --qps 300 --concurrent-requests 2 --duration 30s --load-generator wrk2` |
|         |        --name (optional)         | A memorable name for the test. (default) a random string                                                     |                                                                                                                                                             |
|         |         --mesh optional)         | Name of the service mesh. (default) empty string                                                             |                                                                                                                                                             |
|         |        --file (optional)         | URI of the service mesh performance test configuration file. (default) empty string                          | `--file soak-test-clusterA.yaml`                                                                                                                            |
|         |         --url (required)         | URL of the endpoint send load to during testing.                                                             | `http://my-service/api/v1/test`                                                                                                                             |
|         |         --qps (optional)         | Queries per second (default) 0 (0 - means to use the CPU unbounded to generate as many requests as possible. | `--qps 30`                                                                                                                                                  |
|         | --concurrent-requests (optional) | Number of concurrent requests (default) 1                                                                    | `--concurrent-requests 10`                                                                                                                                  |
|         |      --duration (optional)       | Duration of the test.                                                                                        | `10s`, `5m`, `2h`                                                                                                                                           |
|         |   --load-generator (optional)    | Choice of load generator: fortio (OR) wrk2 (default) fortio                                                  | `--load-generator=fortio`                                                                                                                                   |
|         |            --help, -h            | Displays help the performance management command.                                                            | `mesheryctl perf --help`                                                                                                                                    |

## Service Mesh Lifecycle Management

| command | arg  | flag       | function                                                           | Usage               |
| :------ | :--- | :--------- | :----------------------------------------------------------------- | :------------------ |
| mesh    |      |            | Lifecycle management of service meshes                             |                     |
|         | init |            | Provision a service mesh                                           |                     |
|         |      | --platform | Identify platform to provision service mesh on (e.g. Docker, K8s)  | `--platform docker` |
|         |      | --profile  | Use specific configuration profile                                 | `--profile mTLS`    |
|         |      | --help, -h | Displays help about any service mesh lifecycle management command. | `mesheryctl help`   |
