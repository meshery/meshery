---
layout: enhanced
title: Performance Management with Meshery
abstract: This guide is to help users get a better overview of running and managing performance tests in Meshery
permalink: guides/performance-management/performance-management
redirect_from: functionality/performance-management
type: guides
category: performance
command: perf
language: en
---

This guide walks through running performance benchmarks using Meshery. Users can either use the Meshery UI, the CLI, mesheryctl or run performance tests in their CI/CD pipelines using the Meshery GitHub Action.

## Installing Meshery

Install and login to Meshery to start running performance benchmarks. See [Installation]({{ site.baseurl }}/installation) documentation for detailed steps on how to install Meshery.

_Meshery dashboard_

<a href="{{ site.baseurl }}/assets/img/smi/dashboard.png"><img alt="Meshery Dashboard" src="{{ site.baseurl }}/assets/img/smi/dashboard.png" /></a>

If you are looking to run performance benchmarks on a service mesh, you can use Meshery's service mesh lifecycle management capabilities to deploy a service mesh and deploy your application on the mesh. 

With Meshery's performance benchmarking feature, you can also deploy you application off the mesh and compare the performance and determine the overhead when the app runs on the mesh.

To install a service mesh, see [this guide]({{ site.baseurl }}/service-meshes).

Meshery also comes with a set of [sample applications](h{{ site.baseurl }}/guides/sample-apps) that you can quickly deploy to test out the capabilities of your service mesh.

Next, we navigate to the main Performance Testing dashboard. See [Performance Management]({{ site.baseurl }}/tasks/performance/managing-performance) to learn more about performance profiles, load generators, Kubernetes cluster, and service mesh metrics.

## Running Performance Benchmarks Through Meshery UI

Meshery UI provides an easy-to-use interface in which users can create performance profiles to run repeated tests with similar configuration and can also even schedule performance tests to be run at particular times through the calendar.

On the navigation menu, click on performance.

This will open the performance management dashboard as shown below. 

<a href="{{ site.baseurl }}/assets/img/performance-management/performance-management-dashboard.png"><img alt="Performance Management Dashboard" src="{{ site.baseurl }}/assets/img/performance-management/performance-management-dashboard.png" /></a>

_Performance Management Dashboard_

To create a performance profile, click on "Manage Profiles" and then select "+ Add Performance Profile".

<a href="{{ site.baseurl }}/assets/img/performance-management/performance-profiles.png"><img alt="Performance Profiles" src="{{ site.baseurl }}/assets/img/performance-management/performance-profiles.png" /></a>

_Performance Profiles_

This will open up a popup window where you can fill out the test configuration for running your benchmarks.

<a href="{{ site.baseurl }}/assets/img/performance-management/running-tests.png"><img alt="Performance Test Configuration" src="{{ site.baseurl }}/assets/img/performance-management/running-tests.png" /></a>

_Configuring a Performance Test_

You can now save the profile and run the test. The test will continue to run in the background even if you close the popup.

<a href="{{ site.baseurl }}/assets/img/performance-management/running-tests-spinner.png"><img alt="Running a performance test" src="{{ site.baseurl }}/assets/img/performance-management/running-tests-spinner.png" /></a>

Once the test is done, you would be able to see the results of the test below.

<a href="{{ site.baseurl }}/assets/img/performance-management/result-chart.png"><img alt="Running a performance test" src="{{ site.baseurl }}/assets/img/performance-management/result-chart.png" /></a>

_Performance Test Results_

You can then go back to your performance profile and get these results anytime.

You can also compare different test results in Meshery and draw insights. For example, you can run this test with your application running on different service meshes and check which performs better.

<a href="{{ site.baseurl }}/assets/img/performance-management/comparison-table.png"><img alt="Comparing Test Results" src="{{ site.baseurl }}/assets/img/performance-management/comparison-table.png" /></a>

_Selecting multiple test results_

<a href="{{ site.baseurl }}/assets/img/performance-management/comparison.png"><img alt="Comparing Test Results" src="{{ site.baseurl }}/assets/img/performance-management/comparison.png" /></a>

_Comparing multiple test results_

## Running Performance Benchmarks Through mesheryctl

The `mesheryctl perf` subcommand provides the performance management features of Meshery in the CLI.

To run a performance test based on a performance profile, run:

```
mesheryctl perf apply Istio Performance Test
```

You can also use flags to configure your performance test. For example:

```
mesheryctl perf apply istio-soak-test --concurrent-requests 1 --duration 15s --load-generator nighthawk --mesh istio --url http://localhost:2323
```

mesheryctl also supports test configurations written in SMP compatible format as shown below:

```
test:
  smp_version: v0.0.1
  name: Istio Performance Test
  labels: {}
  clients:
    - internal: false
      load_generator: fortio
      protocol: 1
      connections: 2
      rps: 10
      headers: {}
      cookies: {}
      body: ''
      content_type: ''
      endpoint_urls:
        - 'http://localhost:2323/productpage'
  duration: '30m'
mesh:
  type: 3
```

And then you can pass this file to mesheryctl as:

```
mesheryctl perf apply -f perf-config.yaml
```

You can also override the configuration passed in the file with flags like shown below:

```
mesheryctl perf apply -f perf-config.yaml --url http://localhost:2323/productpage?u=test --load-generator nighthawk --qps 5
```

## Running Performance Benchmarks in your Pipelines

Meshery also has a [meshery-smp-action](https://github.com/layer5io/meshery-smp-action) which is a GitHub action that can be used to run performance tests in your CI/CD pipelines.

Download the token from the Meshery Dashboard by clicking on the profile icon on the top-right corner.

_Downloading the token_

<a href="{{ site.baseurl }}/assets/img/smi/download-token.png"><img alt="SMI Conformance Test Results" src="{{ site.baseurl }}/assets/img/smi/download-token.png" /></a>

You can use this token to authenticate the instance of Meshery running in your CI/CD workflow.

{% include alert.html type="info" title="Using the token in GitHub workflows" content="You can use the <a href='https://docs.github.com/en/actions/reference/encrypted-secrets'>secrets feature in GitHub</a> to store the token." %}

The action can be used by defining your test configuration in a performance profile in Meshery or by writing your test configuration in [SMP compatible format](https://github.com/layer5io/meshery-smp-action#smp-compatible-test-configuration-file).

The action can then be configured as shown below:

```
name: Meshery SMP Action
on:
  push:
    branches:
      'master'

jobs:
  performance-test:
    name: Performance Test
    runs-on: ubuntu-latest
    steps:
      - name: checkout
        uses: actions/checkout@v2
        with:
          ref: 'perf'

      - name: Deploy k8s-minikube
        uses: manusa/actions-setup-minikube@v2.4.1
        with:
          minikube version: 'v1.21.0'
          kubernetes version: 'v1.20.7'
          driver: docker

      - name: Run Performance Test
        uses: layer5io/meshery-smp-action@master
        with:
          provider_token: ${{ secrets.PROVIDER_TOKEN }}
          platform: docker
          profile_name: soak-test
```

More configuration details of the action can be found [here](https://github.com/layer5io/meshery-smp-action/blob/master/action.yml).

See [sample configurations](https://github.com/layer5io/meshery-smp-action#sample-configuration) for more workflow examples using this action.
