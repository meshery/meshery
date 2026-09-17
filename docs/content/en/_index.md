---
title: Meshery Documentation
display_title: false
toc_hide: true
display_toc: false
linkTitle: Documentation
---

<div style="display:grid; justify-items:center">
  <div class="home-hero-logo" style="align-self: center; margin: 0; padding: 0; width: clamp(170px, 50%, 800px);">
    {{< svg/meshery-logo >}}
  </div>
  <h3 style="font-size:1.6rem">As a self-service engineering platform, Meshery enables collaborative design and operation of cloud and cloud native infrastructure.</h3>
</div>

<div class="flex container">
  <!-- OVERVIEW -->
  <div class="section">
    <a href="{{< ref "installation/_index.md" >}}">
        <div class="btn-primary">Overview & Installation</div>
    </a>
    <ul>
        <li>🚀 <a href="{{< ref "installation/quick-start/index.md" >}}">Quick Start</a> , <a href="{{< ref "project/_index.md" >}}">Overview</a>, and <a href="{{< ref "project/FAQ.md" >}}">FAQs</a></li>
    </ul>
    <details>
      <summary>Installation</summary>
      <ul class="section-title">
        <li><a href="{{< ref "installation/_index.md" >}}">Getting Started</a></li>
        {{< section-pages section="installation" >}}
      </ul>
    </details>
  </div>

  <!-- CONCEPTS -->
  <div class="section">
    <a href="{{< ref "concepts/_index.md" >}}">
        <div class="btn-primary">Concepts</div>
    </a>
    <details>
      <summary>Logical</summary>
      <ul>
        <li><a href="{{< ref "concepts/logical/_index.md" >}}">Getting Started</a></li>
        {{< section-pages section="concepts/logical" >}}
      </ul>
    </details>
    <details>
      <summary>Architectural</summary>
      <ul>
        <li><a href="{{< ref "concepts/architecture/_index.md" >}}">Getting Started</a></li>
        {{< section-pages section="concepts/architecture" >}}
      </ul>
    </details>
  </div>
</div>

<div class="flex container">

<!-- GUIDES -->
  <div class="section">
    <a href="{{< ref "guides/_index.md" >}}">
        <div class="btn-primary">Guides & Tutorials</div>
    </a>
    <details>
      <summary>Using Meshery CLI Guides</summary>
      <ul class="section-title">
        <li><a href="{{< ref "guides/mesheryctl/_index.md" >}}">Getting Started</a></li>
        {{< section-pages section="guides/mesheryctl" >}}
      </ul>
    </details>
    <details>
      <summary>🧑‍🔬 Tutorials</summary>
      <ul class="section-title">
        <li><a href="{{< ref "guides/tutorials/_index.md" >}}">Getting Started</a></li>
        {{< section-pages section="guides/tutorials" >}}
      </ul>
    </details>
    <details>
      <summary>Infrastructure Management</summary>
      <ul class="section-title">
        <li><a href="{{< ref "guides/infrastructure-management/_index.md" >}}">Getting Started</a></li>
        {{< section-pages section="guides/infrastructure-management" >}}
      </ul>
    </details>
    <details>
      <summary>Performance Management</summary>
      <ul class="section-title">
        <li><a href="{{< ref "guides/performance-management/_index.md" >}}">Getting Started</a></li>
        {{< section-pages section="guides/performance-management" >}}
      </ul>
    </details>
    <details>
      <summary>Configuration Management</summary>
      <ul class="section-title">
        <li><a href="{{< ref "guides/configuration-management/_index.md" >}}">Getting Started</a></li>
        {{< section-pages section="guides/configuration-management" >}}
      </ul>
    </details>
    <details>
      <summary>Troubleshooting Guides</summary>
      <ul class="section-title">
        <li><a href="{{< ref "guides/troubleshooting/_index.md" >}}">Getting Started</a></li>
        {{< section-pages section="guides/troubleshooting" >}}
      </ul>
    </details>
  </div>

  <!-- Extensions -->
  <div class="section">
    <a href="{{< ref "extensions/_index.md" >}}">
        <div class="btn-primary">Integrations & Extensions</div>
    </a>
    <details>
      <summary>Extensions</summary>
      <ul class="section-title">
        <li><a href="{{< ref "extensions/_index.md" >}}">Getting Started</a></li>
        {{< section-pages section="extensions" >}}
      </ul>
    </details>
    <details>
      <summary>Models</summary>
      <ul class="section-title">
        <li><a href="{{< ref "extensions/models/_index.md" >}}">Getting Started</a></li>
        {{< model-categories >}}
        <li>See all <a href="{{< ref "extensions/models/_index.md" >}}">{{< model-count >}} models</a></li>
      </ul>
    </details>
  </div>
   
</div>

<div class="flex container">

<!-- Contributing & Community -->
  <div class="section">
    <a href="{{< ref "project/_index.md" >}}">
        <div class="btn-primary">Contributing & Community</div>
    </a>
    <details>
      <summary>Community</summary>
      <ul class="section-title">
        <li><a href="{{< ref "project/_index.md" >}}">Getting Started</a></li>
        {{< section-pages section="project" >}}
      </ul>
    </details>
    <!-- CONTRIBUTING -->
    <details>
      <summary>Contributing</summary>
      <ul class="section-title">
        <li><a href="{{< ref "project/contributing/_index.md" >}}">Getting Started</a></li>
        {{< section-pages section="project/contributing" >}}
      </ul>
    </details>
  </div>

  <!-- REFERENCE -->
  <div class="section">
    <a href="{{< ref "reference/_index.md" >}}">
        <div class="btn-primary">Extensibility & Reference</div>
    </a>
    <!-- Reference -->
    <details>
      <summary>Reference</summary>
      <ul class="section-title">
        <li><a href="{{< ref "reference/_index.md" >}}">Getting Started</a></li>
        {{< section-pages section="reference/references" >}}
      </ul>
    </details>
    <!-- Extensibility -->
    <details>
      <summary>Extensibility</summary>
      <ul class="section-title">
        <li><a href="{{< ref "reference/extensibility/_index.md" >}}">Getting Started</a></li>
        {{< section-pages section="reference/extensibility" >}}
      </ul>
    </details>
  </div>

</div>

<p width="100%">Follow on <a href="https://x.com/mesheryio">X</a> or subscribe to our <a href="https://meshery.io/subscribe">newsletter</a> for the latest updates. Get support on our <a href="https://discuss.meshery.io/">forum</a>. Join our <a href="https://slack.meshery.io">Slack</a> to interact directly with other users and contributors.</p>
