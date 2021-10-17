---
layout: default
title: Running system checks using Meshery CLI
permalink: guides/mesheryctl/running-system-checks-using-cli
language: en
type: Guides
category: mesheryctl
---

To check if Meshery is suitable on a system for both pre and post mesh deployment scenarios, `mesheryctl system check` helps to run several checks on the system.

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">
mesheryctl system check

</div></div>
</pre>

## Types of deployment checks

`mesheryctl system check` command has two variants to run different kinds of checks:
 - Pre-deployment
 - Post-deployment

### Pre-deployment check

In this kind of test, the `mesheryctl system check` runs the following checks to determine if a system is compatible with Meshery or not. The checks done here are:
 - Docker healthchecks
 - Kubernetes healthchecks
 - Kubernetes version healthchecks

To run pre-deployment checks, type the following command:
<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">
mesheryctl system check --preflight

</div></div>
</pre>

### Post deployment check

Here, the command runs several checks to determine if the system is good to deploy a service mesh over a kubernetes cluster using Meshery. The checks done here are:
 - Docker healthchecks
 - Kubernetes healthchecks
 - Kubernetes version healthchecks
 - Meshery version healthchecks
 - Meshery Adapter healthchecks

To run post-deployment checks, type the following command:
<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">
mesheryctl system check

</div></div>
</pre>

**Note: This deployment check is a full test along with pre-deployment checks**

## Additional checks

In addition to pre and post deployment checks, the `mesheryctl system check` also has flags that let's you allow to run checks on specific components like Mesh adapters, Meshery Operator and so on. Refer [`mesheryctl system check`]({{ site.baseurl }}/reference/mesheryctl/system/check) documentation page for more details.

## FAQ

##### Question: While running `mesheryctl system check --preflight` it says I didn't install Kubernetes, but I have Docker installed and the test returned "Meshery prerequisites met". Is that all good?
**Answer**: _Yes, as long as you've Docker installed, it's fine to run Meshery. But while handling tasks like deploying service mesh and so on, you need a Kubernetes cluster to do so if you want to do them via Meshery._

##### Question: I ran a preflight check to see if I satisfy all requirements for Meshery in my system. It returned postive results but I couldn't start Meshery. What to do?
**Answer**: _Make sure if you've configured your system to run Meshery in smooth manner. For configuration, do check out the docs site and [this page](https://docs.meshery.io/installation/platforms) to see instructions related to the platform you use._

##### Question: Is Docker alone enough to run meshery?
**Answer**: _For basic tasks, yes. But do note that, some functionality will not be available if you don't use a Kubernetes cluster._

##### Question: What should I have in order to run and use Meshery in full scale?
**Answer**: _You should have a Kubernetes cluster or Docker with Kubernetes in order to use Meshery in full scale_

##### Question: For system checks, do I need any add-ons to pass the check?
**Answer**: _Not necessary. Basic requirements are enough to pass the check._

# Suggested Reading

For an exhaustive list of `mesheryctl` commands and syntax:

- See [`mesheryctl` Command Reference]({{ site.baseurl }}/reference/mesheryctl).

Guides to using Meshery's various features and components.

{% capture tag %}

<li><a href="{{ site.baseurl }}/guides/upgrade#upgrading-meshery-cli">Upgrading mesheryctl</a></li>

{% endcapture %}

{% include suggested-reading.html diffName="true" isDiffTag="true" diffTag=tag %}
