---
layout: default
title: Codespaces
permalink: installation/codespaces
type: installation
display-title: "false"
language: en
list: include
image: /assets/img/platforms/codespaces.png
---

<h1>Quick Start with {{ page.title }} <img src="{{ page.image }}" style="width:35px;height:35px;" /></h1>

Use Minikube in GitHub Codespace to setup your development environment for Meshery.

**To Setup and run Meshery on Minikube** :

- [Steps](#steps)
  - [1. Start Minikube](#1-start-minikube)
  - [2. Install Meshery](#2-install-meshery)
  - [3. Open Meshery UI](#3-open-meshery-ui)
- [Additional Guides](#additional-guides)

You can develop and run Meshery in a GitHub Codespace using your choice of tool:

- A command shell, via an SSH connection initiated using GitHub CLI.
- One of the JetBrains IDEs, via the JetBrains Gateway.
- The Visual Studio Code desktop application.
- A browser-based version of Visual Studio Code.

{% include alert.html type="dark" title="Choice of Codespace Tool" content="For the best experience, run Codespace in your locally <a href='https://docs.github.com/en/codespaces/developing-in-codespaces/developing-in-a-codespace'>installed IDE</a>. Alternatively, you can <br /><a href='https://github.com/codespaces/new?hide_repo_select=true&ref=master&repo=157554479&machine=premiumLinux'><img alt='Open in GitHub Codespaces' src='https://github.com/codespaces/badge.svg' /></a>" %}

## Steps

After starting the Codespace in your Meshery fork, perform the following steps in the `workspaces` folder in order to run Meshery on Minikube:

### 1. Start Minikube

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">minikube start --cpus 4 --memory 4096</div></div>
 </pre>

Confirm that Minikube is running:

<pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">minikube status</div></div>
 </pre>

### 2. Install Meshery

<pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">$ curl -L https://meshery.io/install | bash -</div></div>
 </pre>

### 3. Open Meshery UI
Open Meshery UI and use Meshery. Run `mesheryctl system dashboard`, if needed, to locate the port over which Meshery is exposed.

## Additional Guides

<div class="section">
    <ul>
        <li><a href="{{ site.baseurl }}/guides/troubleshooting/installation">Troubleshooting Meshery Installations</a></li>
        <li><a href="{{ site.baseurl }}/reference/error-codes">Meshery Error Code Reference</a></li>
        <li><a href="{{ site.baseurl }}/reference/mesheryctl/system/check">Mesheryctl system check</a></li> 
    </ul>
</div>
<script src="/assets/js/terminal.js" data-termynal-container="#termynal0|#termynal1|#termynal2"></script>