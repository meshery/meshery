---
layout: page
title: Frequently Asked Questions
permalink: project/faq
abstract: General commonly asked questions and answers about Meshery.
language: en
type: project
---

## General FAQs

<details>
    <summary>
    <h6>Question: What is Meshery?</h6>
</summary>

<p><strong>Answer:</strong> As a self-service engineering platform, Meshery enables collaborative design and operation of cloud native infrastructure.</p>
</details>

<details>
    <summary>
    <h6>Question: Why was Meshery created?</h6>
</summary>

<p><strong>Answer:</strong> As an open source, vendor neutral project, Meshery was created out of the necessity to enable platform engineers, site reliability engineers, devops engineers... engineers to collaborate in the management of their infrastucture and workloads. Meshery was created to enable you to expect more from your infrastructure and to do so with confidence.</p>
</details>

<details>
    <summary>
    <h6>Question: What does Meshery do?</h6>
</summary>

<p><strong>Answer:</strong>Collaborative infrastructure management. Meshery enables you to design and operate cloud native infrastructure visually, collaboratively, with confidence and in partnership with your teammates.</p>
</details>

<!-- - _offers a catalog of operational best practices._
- _offersompare apples-to-apples performance across different infrastructure configurations._
- _Understand behavioral differences between service deployments._
- _Track your application performance from version to version._ -->

<details>
    <summary>
    <h6>Question: Is Meshery open source project?</h6>
</summary>
<p><strong>Answer:</strong> Yes, Meshery is a Cloud Native Computing Foundation (CNCF) project and is licensed under Apache v2. As the cloud native management plane, Meshery is an extensible platform, offering multiple extension points within which users and partners can customize and extend Meshery's functionality.</p>
</details>

<details>
    <summary>
<h6>Question: Why should I use Meshery?</h6>
</summary>
<p><strong>Answer:</strong> Meshery is a powerful tool for managing ​Kubernetes infrastructure. It seamlessly integrates with different hundreds of tools and offers extensibility through many different <a href="{{site.baseurl}}/extensibility/#extension-points">extension points</a>. With Meshery, you can easily discover your environment, collaboratively manage multiple Kubernetes clusters, connect your Git and Helm repos, and analyze app and infra performance.</p>
</details>


## User FAQs

<details>
    <summary>
    <h6>Question: What is mesheryctl?</h6>
</summary>
<strong>Answer:</strong> A command line interface to manage Meshery. `mesheryctl` can manage any number of Meshery deployments.
</details>

<details>
<summary>
<h6>Question: How do I install Meshery?</h6>
</summary>
<p><strong>Answer:</strong> Meshery runs on a <a href="{site.baseurl}}/installation">number of platforms</a>. You are encouraged to use <code>mesheryctl</code> to configure and control Meshery deployments. Install `mesheryctl` using any of these options:</p>
<ul>
<li><a href="/installation/linux-mac/bash">Bash user</a></li>
<li><a href="/installation/linux-mac/brew">Brew user</a></li>
<li><a href="/installation/windows/scoop">Scoop user</a></li>
<li><a href="https://github.com/meshery/meshery/releases/latest">Direct download</a></li>
</ul>
</details>

<details>
<summary><h6>Question: What architecture does Meshery have?</h6></summary>
<p><strong>Answer:</strong> An extensible architecture. There are several components, languages and they have different purposes. See Meshery's <a href="/concepts/architecture">Architecture</a>.</p>
</details>

<details>
<summary>
<h6>Question: What is the difference between <code>make server</code> and <code>mesheryctl system start</code>? Do they both run Meshery on my local machine?</h6>
</summary>
<strong>Answer:</strong> Yes, both of them do run Meshery on your local machine. `make server` builds Meshery from source and runs it on your local OS, while `mesheryctl system start` runs Meshery as a set of containers in Docker or in Kubernetes on your local machine.
</details>

<details>
<summary>
<h6>Question: What systems can I deploy Meshery onto?</h6>
</summary>
<strong>Answer:</strong> Many. See Meshery's <a href="{{site.baseurl}}/installation">Compatibility Matrix</a>.
</details>

<details>
<summary><h6>Question: What systems does Meshery manage?</h6></summary>
<p><strong>Answer:</strong> Many. See Meshery's <a href="https://meshery.io/integrations">Integrations</a></p>
</details>

<details>
<summary><h6>Question: Why is Meshery Server only receiving MeshSync updates from one of my Kubernetes Clusters?</h6></summary>
<p><strong>Answer:</strong> In order to receive MeshSync updates, Meshery Server subscribes for updates Meshery Broker. In other words, Meshery Server connects to the `meshery-broker` service port in order to subscribe for streaming MeshSync updates. By default, the Meshery Broker service is deployed as type Kubernetes Service type <code>LoadBalancer</code>, which requires that your Kubernetes cluster provides an external IP address to the Meshery Broker service, exposing it external to the Kubernetes cluster.</p>
<p>If you're running Kubernetes in Docker Desktop, an external IP address of <code>localhost</code> is assigned. If you're running Minikube, and execute <code>minikube tunnel</code> to gain access to Meshery Broker's service, you will find that both Meshery Broker service endpoints (from two different clusters) are sharing the same <code>localhost:4222</code> address and port number. This port sharing causes conflict and Meshery Server is only able to connect to one of the Meshery Brokers.</p>

<p>Few ways to solve this problem:</p>

<ul>
<li>Use an external cloud provider which provides you with the LoadBalancer having an external IP address other than localhost.</li>
<li>Use <a href="https://kind.sigs.k8s.io">Kind</a> cluster with <a href="https://metallb.universe.tf">MetalLB</a> configuration</li>
</ul>
</details>

<details><summary>
<h6>Question: Why does the dashboard not show the infrastructure provisioned or discovered by Meshery?</h6></summary>
<strong>Answer:</strong> <p>This issue is typically caused by either lack of connectivity between Meshery Server and Meshery Broker or by database corruption. Use the following troubleshooting steps to resolve this issue:</p>

<p><strong>Lack of Connectivity</strong></p>

<ol>
<li>Confirm that the Meshery Broker service is exposed from your cluster using <code>kubectl get svc -n meshery</code> and that an hostname or IP address is displayed in the External Address column. Meshery Server should be able to reach this address.</li>
<li>It is possible that MeshSync is not healthy and not sending cluster updates, check for MeshSync status by navigating to Settings in Meshery UI and clicking on the MeshSync connection.</li>
<li>If MeshSync is healthy, check the status of Meshery Broker by clicking on the NATS connection.</li>
</ol>

<p>If either is the case, Meshery Operator will make sure MeshSync and Meshery Broker deployments are again healthy, wait for some time, otherwise try redeploying Meshery Operator.</p>

<p><strong>Database Corruption</strong></p>

<p>If MeshSync, Meshery Broker and Meshery Operator are healthy, then perhaps, there is corruption in the Meshery Database. Use the following troubleshooting steps to resolve this issue:</p>
<ul>
<li>Try clearing the database by clicking on the `Flush MeshSync` button associated with the corresponding cluster.</li>
<li>If you don't see the specific entities in Meshery UI, you may choose to reset Meshery's database. This option is in the <code>Reset System</code> Tab in <code>Settings</code> page.</li>
</ul>

<p>Note: You can also verify health of your system using <a href="{{site.baseurl}}/reference/mesheryctl/system/check">mesheryctl system check</a></p>

</details>

## Contributing FAQs

<details>
<summary>
<strong>Question: Getting an error while running <code>make server</code> on Windows?</strong>
</summary><strong>Answer:</strong> <p>On Windows, set up the project on Ubuntu WSL2 and you will be able to run the Meshery UI and the server. For more information please visit <a href="/project/contributing/meshery-windows">Setting up Meshery Development Environment on Windows</a>.</p>
</details>

{% include discuss.html %}

<!--Add other questions-->
