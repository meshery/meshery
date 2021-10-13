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

# Suggested Reading

For an exhaustive list of `mesheryctl` commands and syntax:

- See [`mesheryctl` Command Reference]({{ site.baseurl }}/reference/mesheryctl).

Guides to using Meshery's various features and components.

{% capture tag %}

<li><a href="{{ site.baseurl }}/guides/upgrade#upgrading-meshery-cli">Upgrading mesheryctl</a></li>

{% endcapture %}

{% include suggested-reading.html diffName="true" isDiffTag="true" diffTag=tag %}
