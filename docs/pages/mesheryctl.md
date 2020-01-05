---
layout: guide
title: mesheryctl
description: How to use mesheryctl
permalink: guides/mesheryctl
type: guide
---

| command | flag  | function  | example usage |
| :------------ |:---------------:| -----:| -----:|
|cleanup|       |Pulls current meshery.yaml from meshery repo| Troubleshooting and Debugging|
|init||Print logs|Troubleshooting and Debugging|
|start||Start Meshery containers|Troubleshooting and Debugging|
|status||Check Meshery status|Troubleshooting and Debugging|
|stop||Stop Meshery containers|Troubleshooting and Debugging|
|update||Pull new Meshery images from Docker Hub|Troubleshooting and Debugging|
|help||Help about any command|Troubleshooting and Debugging|
|mesh||Lifecycle management of service meshes|Management|
|test||Performance test and benchmarking|Management|
|--name||(optional) A memorable name for the test.<br> (default) a random string|Test|
|--mesh||(optional) Name of the service mesh.<br>(default) empty string|Test|
|--file||(optional) URI to the service mesh performance test configuration file.<br>(default) empty string|Test |
|--url||(required) URL of the endpoint to use for the test|Test|
|--qps||(optional) Queries per second<br>(default) 1|Test|
|--parallel-requests||(optional) Number of concurrent requests<br>(default) 1|Test|
|--duration||(optional) Duration of the test like 10s, 5m, 2h. We are following the convention )|Test|
|--load-generator||(optional) choice of load generator: fortio (OR) wrk2<br>(default) fortio|Test|

