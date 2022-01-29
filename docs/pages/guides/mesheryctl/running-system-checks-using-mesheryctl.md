---
layout: default
title: Running system checks using Meshery CLI
permalink: guides/mesheryctl/running-system-checks-using-mesheryctl
language: en
type: Guides
category: mesheryctl
---

Meshery's CLI, `mesheryctl`, includes commands for verifying system readiness for a Meshery deployment and health checks to confirm the health of an existing Meshery deployment. Whether you have yet to deploy Meshery or have already deployed Meshery, `mesheryctl system check` is a useful utility to ensure that your Meshery deployments are healthy. 

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">
mesheryctl system check

</div></div>
</pre>
<br/>
<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">
Verify environment pre/post-deployment of Meshery.

Usage:
  mesheryctl system check [flags]

Flags:
      --adapter     Check status of Meshery adapters
  -h, --help        help for check
      --operator    Check status of Meshery operators
      --pre         Verify environment readiness to deploy Meshery
      --preflight   Verify environment readiness to deploy Meshery

Global Flags:
      --config string    path to config file (default "/Users/navendu/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div></div>
</pre>

## Deployment checks

`mesheryctl system check` command can run two types of system checks. A pre-deployment check which verifies the environment to deploy Meshery and a post-deployment check which runs validation checks on a running Meshery deployment.

### Pre-deployment checks

Pre-deployment checks runs checks on the environment and verifies whether it is ready for deploying Meshery.

The following checks are done here:

 - Docker health checks: Checks for the availability of Docker and docker-compose in the user's machine
 - Kubernetes health checks: Checks for the availability of a Kubernetes cluster and checks if Meshery can initialize a Kubernetes client
 - Kubernetes version checks: Checks if kubectl and the Kubernetes version are higher than the minimum supported versions

Pre-deployment checks are run with the `--preflight` flag as shown below:

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">
mesheryctl system check --preflight

</div></div>
</pre>

### Post-deployment checks

Post-deployment checks are run after deploying Meshery in the user's environment. These checks ensure that the running deployment of Meshery and Meshery adapters are working as expected.

In addition to the pre-flight checks, the following checks are also run in this check:

 - Meshery version checks: Checks the version of Meshery server and CLI and shows if a new version is available
 - Meshery Adapter health checks: Checks if all the specified adapters are deployed and reachable

Post-deployment checks are run as shown below:

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">
mesheryctl system check

</div></div>
</pre>

## Additional checks

To check the status of the deployed adapters only, users can leverage the `--adapter` flag as shown below:

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">
mesheryctl system check --adapter

</div></div>
</pre>

Users can also narrow down the tests to just check the status of the Meshery operator deployed on their Kubernetes cluster:

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">
mesheryctl system check --operator

</div></div>
</pre>

## FAQ

##### Question: While running `mesheryctl system check --preflight` it says I didn't install Kubernetes, but I have Docker installed and the test returned "Meshery prerequisites met". Is that all good?
**Answer**: _Yes, as long as you've Docker installed, it's fine to run Meshery. But you will need a Kubernetes cluster to handle tasks such as deploying service mesh and so on, if you want to do them via Meshery._

##### Question: I ran a preflight check to see if I satisfy all requirements for Meshery in my system. It returned postive results but I couldn't start Meshery. What to do?
**Answer**: _Make sure if you've configured your system to run Meshery in smooth manner. For configuration, do check out the docs site and [this page](https://docs.meshery.io/installation/platforms) to see instructions related to the platform you use._

##### Question: Do I need a Kubernetes cluster or will a Docker host suffice for Meshery deployments?
**Answer**: _Meshery's [performance management](functionality/performance-management) functionality does not require a Kubernetes cluster. The rest of Meshery's functionality (e.g. service mesh management) does require a Kubernetes cluster._

##### Question: What are Meshery's production deployment requirements?
**Answer**: _One or more Kubernetes clusters. A stateful set for Meshery Server in order to persist performance test results. See [#2451](https://github.com/meshery/meshery/issues/2451)._

##### Question: For system checks, do I need any add-ons to pass the check?
**Answer**: _Not necessary. Basic requirements are enough to pass the check._

##### Question: The Adapter check is failing, it returns "Auth token not found". 
**Answer**: _You can log in to Meshery using `mesheryctl system login` which would generate an OAuth token. Once the OAuth token is generated, the check will start to function_

##### Question: I have a Kubernetes cluster enabled but Meshery couldn't reach the cluster and the checks are failing! What to do?
**Answer**: _To resolve this error, you can upload your kubeconfig file in the Meshery UI under settings and Meshery will reconfigure to use your Kubernetes cluster._

##### Question: Under Meshery Version test, I'm getting an error like "CLI is not up-to-date". Should I update mesheryctl often?
**Answer**: _Yes! You should update the mesheryctl often in order to run Meshery smoothly. The reason behind it is because not only the CLI is updated, but also the Meshery app. So it is advisable to update Meshery often._

##### Question: Is it advisable to keep Meshery in sleep mode while running system checks?
**Answer**: _Not necessary. It is good to keep Meshery up and running, else the system checks will fail to detect the Meshery version._

##### Question: What is the minimum version of k8s cluster and kubectl required to run Meshery?
**Answer**: _For Kubernetes, version >=1.12.0 is recommended. For kubectl version >=1.12 is recommended._

##### Question: In the "Meshery Adapter" section of check, I could see only some service mesh adapters up and running and not all. Is this fine?
**Answer**: _Not a problem, if you feel you need to have all mesh adapters to be up running, you can do so by creating a new context `mesheryctl system context create [context-name] --set` (if you voluntarily deleted mesh adapters in your current context)_

##### Question: I started Meshery fresh, didn't change any of the details in the context I have. But I see that all adapter checks are failing. What to do?
**Answer**: _Configure Meshery to use on your Kubernetes cluster, then upload the kubeconfig file via Meshery UI to notify Meshery to use that cluster. If that didn't work, feel free to [open up an issue](https://github.com/meshery/meshery/issues) in GitHub._

# Suggested Reading

For an exhaustive list of `mesheryctl` commands and syntax:

- See [`mesheryctl` Command Reference]({{ site.baseurl }}/reference/mesheryctl).

Guides to using Meshery's various features and components.

{% capture tag %}

<li><a href="{{ site.baseurl }}/guides/upgrade#upgrading-meshery-cli">Upgrading mesheryctl</a></li>

{% endcapture %}

{% include suggested-reading.html diffName="true" isDiffTag="true" diffTag=tag %}