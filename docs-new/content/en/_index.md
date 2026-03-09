---
title: Meshery Documentation
display_title: false
display_toc: false
linkTitle: Documentation
---

<div style="display:grid; justify-items:center">
  <div style="align-self:center; margin-bottom:0px; margin-top:0px;padding-top:0px; padding-bottom:0px;width:clamp(170px, 50%, 800px);">
    {{< svg/meshery-logo >}}
  </div>
  <h3 style="font-size:1.6rem">As a self-service engineering platform, Meshery enables collaborative design and operation of cloud and cloud native infrastructure.</h3>
</div>

<div class="flex container">
  <!-- OVERVIEW -->
  <div class="section">
    <a href="/installation">
        <div class="btn-primary">Overview & Installation</div>
    </a>
    <ul>
        <li>🚀 <a href="/installation/quick-start">Quick Start</a> and <a href="/project/faq">FAQs</a></li>
    </ul>
    <details>
      <summary>
        <p style="display:inline">
          <a href="/installation/" class="text-black">Installation</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="installation" >}}
      </ul>
    </details>
  </div>

  <!-- CONCEPTS -->
  <div class="section">
    <a href="/concepts">
        <div class="btn-primary">Concepts</div>
    </a>
    <details>
      <summary>
        <p style="display:inline">
          <a href="/concepts/logical" class="text-black">Logical</a>
        </p>
      </summary>
      <ul>
        {{< section-pages section="concepts/logical" >}}
      </ul>
    </details>
    <details>
      <summary>
        <p style="display:inline">
          <a href="/concepts/architecture" class="text-black section-title">Architectural</a>
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
    <a href="/guides">
        <div class="btn-primary">Guides & Tutorials</div>
    </a>
    <details>
      <summary>
        <p style="display:inline">
          <a href="/guides/mesheryctl/" class="text-black">Using Meshery CLI Guides</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="guides/mesheryctl" >}}
      </ul>
    </details>
    <details>
      <summary>
        <p style="display:inline">
          <a href="/guides/tutorials/" class="text-black">🧑‍🔬 Tutorials</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="guides/tutorials" >}}
      </ul>
    </details>
    <details>
      <summary>
        <p style="display:inline">
          <a href="/guides/infrastructure-management" class="text-black">Infrastructure Management</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="guides/infrastructure-management" >}}
      </ul>
    </details>
    <details>
      <summary>
        <p style="display:inline">
          <a href="/guides/performance-management" class="text-black">Performance Management</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="guides/performance-management" >}}
      </ul>
    </details>
    <details>
      <summary>
        <p style="display:inline">
          <a href="/guides/configuration-management" class="text-black">Configuration Management</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="guides/configuration-management" >}}
      </ul>
    </details>  
    <details>
      <summary>
        <p style="display:inline">
          <a href="/guides/troubleshooting" class="text-black">Troubleshooting Guides</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="guides/troubleshooting" >}}
      </ul>
    </details>
  </div>

  <!-- Extensions -->
  <div class="section">
    <a href="/extensions">
        <div class="btn-primary">Integrations & Extensions</div>
    </a>
    <details>
      <summary>
        <p style="display:inline">
          <a href="/extensions" class="text-black">Extensions</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="extensions" >}}
      </ul>
    </details>
    <details>
      <summary>
        <p style="display:inline">
          <a href="/extensibility/integrations" class="text-black">Integrations</a>
        </p>
      </summary>
      <ul class="section-title">
        <li>See all <a href="/extensions/models">{{< model-count >}} integrations</a></li>
      </ul>
    </details>
  </div>
   
</div>

<div class="flex container">

<!-- Contributing & Community -->
  <div class="section">
    <a href="/project">
        <div class="btn-primary">Contributing & Community</div>
    </a>
    <details>
      <summary>
        <p style="display:inline">
          <a href="/project" class="text-black">Community</a>
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
          <a href="/project/contributing" class="text-black">Contributing</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="project/contributing" >}}
      </ul>
    </details>
  </div>

  <!-- REFERENCE -->
  <div class="section">
    <a href="/reference">
        <div class="btn-primary">Extensibility & Reference</div>
    </a>
    <!-- Reference -->
    <details>
      <summary>
        <p style="display:inline">
          <a href="/reference" class="text-black">Reference</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="reference" >}}
      </ul>
    </details>
    <!-- Extensibility -->
    <details>
      <summary>
        <p style="display:inline">
          <a href="/extensibility" class="text-black">Extensibility</a>
        </p>
      </summary>
      <ul class="section-title">
        {{< section-pages section="extensibility" >}}
      </ul>
    </details>
  </div>

</div>

<p width="100%">Follow on <a href="https://x.com/mesheryio">X</a> or subscribe to our <a href="https://meshery.io/subscribe">newsletter</a> for the latest updates. Get support on our <a href="https://meshery.io/community#discussion-forums">forum</a>. Join our <a href="https://slack.meshery.io">Slack</a> to interact directly with other users and contributors.</p>
