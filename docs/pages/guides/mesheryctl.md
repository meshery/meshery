---
layout: guide
title: Using mesheryctl
description: How to use mesheryctl
permalink: guides/mesheryctl
type: guide
---

### Meshery Lifecycle Management
Installation, troubleshooting and debugging of Meshery and its adapaters.

| command   | flag  | function                  | Usage                     |
|:----------|:-----:|:--------------------------|:--------------------------|
|cleanup    |       |Pulls current meshery.yaml from meshery repo. *Warning: Any local changes will be overwritten.* |  |
|log        |       |Starts tailing Meshery debug logs.              |   |
|start      |       |Start all Meshery containers.   |   |
|status     |       |Displays the status of Meshery's containers.       |   |
|stop       |       |Stop all Meshery containers.    |   |
|update     |       |Pull new Meshery images from Docker Hub. Pulls new `mesheryctl` client. This command may be executed while Meshery is running. |   |
|help       |       |Displays help about any command.     |   |


### Performance Management

| command   | flag          | function                  | Usage                     |
|:----------|:-------------:|:--------------------------|:--------------------------|
|test       |               |Performance test and benchmarking| |
|           | --name        |(optional) A memorable name for the test.<br> (default) a random string|   |
|           | --mesh optional)| Name of the service mesh.<br>(default) empty string|    |
|           | --file (optional)| URI to the service mesh performance test configuration file.<br>(default) empty string| |
|           | --url (required)| URL of the endpoint to use for the test| |
|           | --qps (optional)| Queries per second<br>(default) 1|   |
|           | --parallel-requests (optional)| Number of concurrent requests<br>(default) 1|  |
|           | --duration (optional) | Duration of the test like 10s, 5m, 2h. We are following the convention )|   |
|           | --load-generator (optional)| choice of load generator: fortio (OR) wrk2<br>(default) fortio|   |

### Service Mesh Lifecycle Management

| command   | flag          | function                  | Usage                     |
|:----------|:-------------:|:--------------------------|:--------------------------|
|mesh       |               | Lifecycle management of service meshes| |