---
layout: default
title: Running SMI Conformance Tests
description: This guide is to help users get a better understanding of sample apps
permalink: guides/smi-conformance
type: Guides
---

This guide will help you run SMI Conformance Tests with Meshery through the [UI](#running-smi-conformance-tests-through-meshery-ui), CLI and a GitHub action for your CI/CD pipelines. To learn more about Meshery and SMI Conformance, see [Meshery and Service Mesh Interface (SMI) Conformance]({{ site.baseurl }}/functionality/service-mesh-interface)

## Setup Meshery and Install a Service Mesh

Install and login to Meshery to start running SMI conformance tests. See [Installation]({{ site.baseurl }}/installation) documentation for detailed steps on how to install Meshery.

_Meshery dashboard_

<a href="{{ site.baseurl }}/assets/img/smi/dashboard.png"><img alt="Meshery Dashboard" src="{{ site.baseurl }}/assets/img/smi/dashboard.png" /></a>

Next, install the service mesh from Meshery. See [Service Meshes]({{ site.baseurl }}/service-meshes) for a list of supported service meshes and guides on how to install them.

_Installing Istio_

<a href="{{ site.baseurl }}/assets/img/smi/istio-dashboard.png"><img alt="Istio Dashboard" src="{{ site.baseurl }}/assets/img/smi/istio-dashboard.png" /></a>

**Alternatively**, you can use mesheryctl, Meshery's CLI to deploy a service mesh. See [mesheryctl mesh]({{ site.baseurl }}/reference/mesheryctl/mesh/) documentation for details.

## Running SMI Conformance Tests Through Meshery UI

Now that we have deployed the service mesh to validate, we can run an SMI conformance test through the UI.

In the "Validate Service Mesh Configuration" click on the "+" button and select "SMI Conformance".

_Running SMI Conformance Test on Istio_

<a href="{{ site.baseurl }}/assets/img/smi/smi-conformance-run.png"><img alt="Running SMI Conformance Tests on Istio" src="{{ site.baseurl }}/assets/img/smi/smi-conformance-run.png" /></a>

This will start running the SMI Conformance tests.

Once the tests are done, you can navigate to SMI Conformance results page by clicking the conformance tab on the menu.

_SMI Conformance Test Results_

<a href="{{ site.baseurl }}/assets/img/smi/smi-conformance-page.png"><img alt="SMI Conformance Test Results" src="{{ site.baseurl }}/assets/img/smi/smi-conformance-page.png" /></a>

Click the dropdown button to view the results of a specific test.

_Viewing the Results_

<a href="{{ site.baseurl }}/assets/img/smi/smi-conformance-result.png"><img alt="SMI Conformance Test Results" src="{{ site.baseurl }}/assets/img/smi/smi-conformance-result.png" /></a>

## Running SMI Conformance Tests Through Meshery CLI (mesheryctl)

## Running SMI Conformance Tests in CI/CD Pipelines


##### Suggested Reading

- Functionality: [Service Mesh Interface (SMI) Conformance]({{ site.baseurl }}/functionality/service-mesh-interface)

