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
    <a href="{{< ref "project/_index.md" >}}">
        <div class="btn-primary">Overview & Installation</div>
    </a>
    <ul>
        <li>🚀 <a href="{{< ref "installation/quick-start/index.md" >}}">Quick Start</a> , <a href="{{< ref "project/_index.md" >}}">Overview</a>, and <a href="{{< ref "project/FAQ.md" >}}">FAQs</a></li>
    </ul>
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{< ref "installation/_index.md" >}}" class="text-black">Installation</a>
        </p>
      </summary>
      <ul class="section-title">
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
      <summary>
        <p style="display:inline">
          <a href="{{< ref "concepts/logical/_index.md" >}}" class="text-black">Logical</a>
        </p>
      </summary>
      <ul>
        {{< section-pages section="concepts/logical" >}}
      </ul>
    </details>
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{< ref "concepts/architecture/_index.md" >}}" class="text-black section-title">Architectural</a>
        </p>
      </summary>
      <ul>
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
      <summary>
        <p style="display:inline">
          <a href="{{< ref "guides/mesheryctl/_index.md" >}}" class="text-black">Using Meshery CLI Guides</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="guides/mesheryctl" >}}
      </ul>
    </details>
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{< ref "guides/tutorials/_index.md" >}}" class="text-black">🧑‍🔬 Tutorials</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="guides/tutorials" >}}
      </ul>
    </details>
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{< ref "guides/infrastructure-management/_index.md" >}}" class="text-black">Infrastructure Management</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="guides/infrastructure-management" >}}
      </ul>
    </details>
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{< ref "guides/performance-management/_index.md" >}}" class="text-black">Performance Management</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="guides/performance-management" >}}
      </ul>
    </details>
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{< ref "guides/configuration-management/_index.md" >}}" class="text-black">Configuration Management</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="guides/configuration-management" >}}
      </ul>
    </details>  
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{< ref "guides/troubleshooting/_index.md" >}}" class="text-black">Troubleshooting Guides</a>
        </p>
      </summary>
      <ul class="section-title">
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
      <summary>
        <p style="display:inline">
          <a href="{{< ref "extensions/_index.md" >}}" class="text-black">Extensions</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="extensions" >}}
      </ul>
    </details>
    <details>
      <summary>
        <p style="display:inline">
<a href="{{< ref "extensions/models/_index.md" >}}" class="text-black">Models</a>
        </p>
      </summary>
      <ul class="section-title">
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
      <summary>
        <p style="display:inline">
          <a href="{{< ref "project/_index.md" >}}" class="text-black">Community</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="project" >}}
      </ul>
    </details>
    <!-- CONTRIBUTING -->
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{< ref "project/contributing/_index.md" >}}" class="text-black">Contributing</a>
        </p>
      </summary>
      <ul class="section-title">
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
      <summary>
        <p style="display:inline">
          <a href="{{< ref "reference/_index.md" >}}" class="text-black">Reference</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="reference/references" >}}
      </ul>
    </details>
    <!-- Extensibility -->
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{< ref "reference/extensibility/_index.md" >}}" class="text-black">Extensibility</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="reference/extensibility" >}}
      </ul>
    </details>
  </div>

</div>

<p width="100%">Follow on <a href="https://x.com/mesheryio">X</a> or subscribe to our <a href="https://meshery.io/subscribe">newsletter</a> for the latest updates. Get support on our <a href="https://discuss.meshery.io/">forum</a>. Join our <a href="https://slack.meshery.io">Slack</a> to interact directly with other users and contributors.</p>
