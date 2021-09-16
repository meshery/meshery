---
layout: default
title: Performance Management with Meshery
description: This guide is to help users get a better overview of running and managing performance tests in Meshery
permalink: guides/performance-management
type: Guides
command: perf
---

This guide is useful to learn how to manage performance tests with Meshery through the [UI](#performance-management-through-meshery-ui), [CLI](#performance-management-through-meshery-cli-(mesheryctl)) and a [GitHub Action](#performance-management-through-meshery-using-github-action) for your test results. To learn more about Meshery and Service Mesh Performance Management commands, see [Service Mesh Performance Management]({{ site.baseurl }}/reference/mesheryctl/commands/mesheryctl-perf)

## Setup Meshery and run a Performance Test on your Service Mesh

Install and login to Meshery to start running SMI conformance tests. See [Installation]({{ site.baseurl }}/installation) documentation for detailed steps on how to install Meshery.

_Meshery dashboard_

<a href="{{ site.baseurl }}/assets/img/smi/dashboard.png"><img alt="Meshery Dashboard" src="{{ site.baseurl }}/assets/img/smi/dashboard.png" /></a>

Next, we navigate to the main Performance Testing dashboard. See [Performance Management]({{ site.baseurl }}/functionality/performance) to learn more about performance profiles, load generators and Kubernetes Cluster and Service Mesh Metrics.

_Performance dashboard_

<a href="{{ site.baseurl }}/assets/img/performance-management/performance-dashboard.PNG"><img alt="Performance Dashboard" src="{{ site.baseurl }}/assets/img/performance-management/performance-dashboard.PNG" /></a>

The dashboard allows us to run our Performance Tests , it also shows the different tests run over the lifetime of the profile sorted by time and day. 

## Performance Management Through Meshery UI

<a href="{{ site.baseurl }}/assets/img/performance-management/performance-run-test.PNG"><img alt="Performance Test Options" src="{{ site.baseurl }}/assets/img/performance-management/performance-run-test.PNG" /></a>

Now that we have our setup, we can run a performance test with our required constraints. We can modify multiple testing scenarios by altering the Service Mesh used, adjusting the Concurrent Requests taking place or the number of queries run or select a particular Load Generator. For the sake of simplicity, we can also personalize the name of the test for future reference. 

Upon personalization, select **Run Test** to execute the current test or **Save Profile** to save the properties of the test and execute it at a convinient time. 

_SMI Conformance Test Results_

<a href="{{ site.baseurl }}/assets/img/performance-management/perf-test-compare-option.PNG"><img alt="Performance Test Options" src="{{ site.baseurl }}/assets/img/performance-management/perf-test-compare-option.PNG" /></a>

Once we run the tests and get the desired results, we can compare these test results against one another by selecting tests. For the best experience, you should compare between two or more tests of similar configuration. Tests that are configured with a high degree of variance (e.g., one test ran for a duration of 5 minutes, while another test ran for a duration of 1 hour) will produce comparisons that are more difficult to extrapolate insights from.

Clicking on the **Compare Selected** button takes us to the next step of visualizing the comparision of the results.

_Performance Test Results_

<a href="{{ site.baseurl }}/assets/img/performance-management/performance-test-comparison.PNG"><img alt="Performance Test Results" src="{{ site.baseurl }}/assets/img/performance-management/performance-test-comparison.PNG" /></a>

The results are depicted in the form of Histograms showing the count and Line Graphs showing the cumulative percentage.

Latency and throughput are the two most significantly examined signals.

Meshery will use different algorithms to calculate results depending upon which load generator was used to run the test.
  
_Viewing the Results_

<a href="{{ site.baseurl }}/assets/img/performance-management/chart.png"><img alt="Performance Test Results" src="{{ site.baseurl }}/assets/img/performance-management/chart.png" /></a>

## Performance Management Through Meshery CLI (mesheryctl)

_Description_

<!-- Description of the command. Preferably a paragraph -->
{% assign name = site.data.mesheryctlcommands.cmds[page.command] %}
{{ name.description }}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ name.usage }}
  </div>
</pre>

_Examples_

{% for flag_hash in name.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.usage }}
  </div>
</pre>
{% endfor %}
<br/>

_Options & Flags_

{% for flag_hash in name.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.name }}
  </div>
</pre>
{% endfor %}
<br/>

_Options inherited from parent commands_
<pre class="codeblock-pre">
  <div class="codeblock">
  --help, -h # Shows help for the command
  </div>
</pre>

## Performance Management Through Meshery using GitHub Action 

You can use [Performance Management GitHub action](https://github.com/layer5io/meshery-smp-action) to run Performance Management tests in your CI/CD pipelines.

### Meshery SMP Action

This repository is used for storing a GitHub action for performing SMP tests using Meshery

_Inputs_

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
    # token to connect with the remote provider
    provider_token:
      description: "Provider token to use. NOTE: value of the 'token' key in auth.json"
      required: true

    # platform to deploy meshery
    platform:
      description: "Platform to deploy meshery on. Possible values: docker, kubernetes"
      default: docker

    # SUPPLY EITHER "profile_filename" or profile_name

    # name of the file storing the performance profile (keep in .github)
    profile_filename:
      description: "Name of the file containing SMP profile"

    # name of the prformance profile to use
    profile_name:
      description: "Name of the performance profile" 
  </div></div>
  </pre>


_Sample Configuration_

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
    name: Testing SMP action
    on:
      push:
        branches:
          'perf'

    jobs:
      job1:
        name: Run Performance Test
        runs-on: ubuntu-latest
        steps:
          - name: checkout
            uses: actions/checkout@v2
            with:
              ref: 'perf'

          - name: Deploy k8s
            uses: manusa/actions-setup-minikube@v2.4.1
            with:
              minikube version: 'v1.21.0'
              kubernetes version: 'v1.20.7'
              driver: docker

          - name: Performance test
            uses: layer5io/meshery-smp-action@master
            with:
              provider_token: ${{ secrets.PROVIDER_TOKEN }}
              platform: docker
              profile_name: demo
  </div></div>
  </pre>