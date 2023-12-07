<!-- ---
layout: default
title: Running SMI Conformance Tests
description: This guide is to help users get a better understanding of sample apps
permalink: guides/smi-conformance
type: guides
language: en

---

This guide will help you run SMI Conformance Tests with Meshery through the [UI](#running-smi-conformance-tests-through-meshery-ui), CLI and a GitHub action for your CI/CD pipelines. To learn more about Meshery and SMI Conformance, see [Meshery and Service Mesh Interface (SMI) Conformance]({{ site.baseurl }}/tasks/service-mesh-interface)

## Setup Meshery and Install a Service Mesh

Install and login to Meshery to start running SMI conformance tests. See [Installation]({{ site.baseurl }}/installation) documentation for detailed steps on how to install Meshery.

_Meshery dashboard_

<a href="{{ site.baseurl }}/assets/img/smi/dashboard.png"><img alt="Meshery Dashboard" src="{{ site.baseurl }}/assets/img/smi/dashboard.png" /></a>

Next, install the service mesh from Meshery. See [Service Meshes]({{ site.baseurl }}/service-meshes) for a list of supported service meshes and guides on how to install them.

_Installing Istio_

<a href="{{ site.baseurl }}/assets/img/smi/istio-dashboard.png"><img alt="Istio Dashboard" src="{{ site.baseurl }}/assets/img/smi/istio-dashboard.png" /></a>

**Alternatively**, you can use mesheryctl, Meshery's CLI to deploy a service mesh. See [mesheryctl mesh]({{ site.baseurl }}/reference/mesheryctl/mesh/) documentation for details.

## Running SMI Conformance Tests Through Meshery UI

Now that we have deployed the service mesh to validate (see [Setup Meshery and Install a Service Mesh](#setup-meshery-and-install-a-service-mesh)), we can run an SMI conformance test through the UI.

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

Once we have deployed the service mesh to validate  (see [Setup Meshery and Install a Service Mesh](#setup-meshery-and-install-a-service-mesh)), we can run an SMI conformance test through the CLI (mesheryctl).

You can also use mesheryctl, Meshery's CLI to deploy a service mesh. See [mesheryctl mesh]({{ site.baseurl }}/reference/mesheryctl/mesh/) documentation for details.

Download the token from the Meshery Dashboard by clicking on the profile icon on the top-right corner.

_Downloading the token_

<a href="{{ site.baseurl }}/assets/img/smi/download-token.png"><img alt="SMI Conformance Test Results" src="{{ site.baseurl }}/assets/img/smi/download-token.png" /></a>

Open a terminal and run the following command.

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl mesh validate -a [name of the adapter] -t [path to token for authentication] -s smi
 </div></div>
 </pre>

 For example to run an SMI Conformance test on Open Service Mesh, we can run:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl mesh validate -a meshery-osm:10009 -t ~/Downloads/auth.json -s smi
 </div></div>
 </pre>

 To view the results of the test, you can open up Meshery in the browser and navigate to the "Conformance" tab.

_SMI Conformance Test Results_

<a href="{{ site.baseurl }}/assets/img/smi/smi-conformance-page.png"><img alt="SMI Conformance Test Results" src="{{ site.baseurl }}/assets/img/smi/smi-conformance-page.png" /></a>

Click the dropdown button to view the results of a specific test.

_Viewing the Results_

<a href="{{ site.baseurl }}/assets/img/smi/smi-conformance-result.png"><img alt="SMI Conformance Test Results" src="{{ site.baseurl }}/assets/img/smi/smi-conformance-result.png" /></a>

## Running SMI Conformance Tests in CI/CD Pipelines

You can use [Meshery SMI Conformance GitHub action](https://github.com/layer5io/meshery-smi-conformance-action) to run SMI Conformance tests in your CI/CD pipelines.

The conformance test would be run in your pipeline and the results will be published on your Meshery Dashboard (see [Setup Meshery and Install a Service Mesh](#setup-meshery-and-install-a-service-mesh)).

Download the token from the Meshery Dashboard by clicking on the profile icon on the top-right corner.

_Downloading the token_

<a href="{{ site.baseurl }}/assets/img/smi/download-token.png"><img alt="SMI Conformance Test Results" src="{{ site.baseurl }}/assets/img/smi/download-token.png" /></a>

You can use this token to authenticate the instance of Meshery running in your CI/CD workflow.

{% include alert.html type="info" title="Using the token in GitHub workflows" content="You can use the <a href='https://docs.github.com/en/actions/reference/encrypted-secrets'>secrets feature in GitHub</a> to store the token." %}

A sample workflow that could be added to your `.github/workflows` folder is given below. This example uses Open Service Mesh.

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 name: SMI Conformance with Meshery
 on:
   push:
     tags:
       - 'v*'
 
 jobs:
   smi-conformance:
     name: SMI Conformance
     runs-on: ubuntu-latest
     steps:
 
       - name: Validate SMI Conformance
         uses: layer5io/mesheryctl-smi-conformance-action@master
         with:
           service_mesh: open_service_mesh
           provider_token: ${{ secrets.MESHERY_PROVIDER_TOKEN }}
           mesh_deployed: false
 </div></div>
 </pre>

 Once the tests are run, you can log into Meshery to view the results in the "Conformance" tab.

 _SMI Conformance Test Results_

<a href="{{ site.baseurl }}/assets/img/smi/smi-conformance-page.png"><img alt="SMI Conformance Test Results" src="{{ site.baseurl }}/assets/img/smi/smi-conformance-page.png" /></a>

Click the dropdown button to view the results of a specific test.

_Viewing the Results_

<a href="{{ site.baseurl }}/assets/img/smi/smi-conformance-result.png"><img alt="SMI Conformance Test Results" src="{{ site.baseurl }}/assets/img/smi/smi-conformance-result.png" /></a>

##### Suggested Reading

- Functionality: [Service Mesh Interface (SMI) Conformance]({{ site.baseurl }}/tasks/service-mesh-interface)
 -->
