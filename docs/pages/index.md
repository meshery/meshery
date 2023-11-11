---
layout: page
title: Meshery Documentation
permalink: /
display-title: "false"
display-toc: "false"
language: en
---

{% assign sorted_pages = site.pages | sort: "name" | alphabetical %}

<div class="flex flex-col--1">
  <div style="align-self:center; margin-bottom:0px; margin-top:0px;padding-top:0px; padding-bottom:0px;width:clamp(170px, 50%, 800px);">
    {% include svg/meshery-logo.html %}
  </div>
  <h3>As a self-service engineering platform, Meshery enables collaborative design and operation of cloud native infrastructure.</h3>
</div>
<div class="flex flex-col--2"  style="text-align: left; padding:1.6rem ;--col-gap:1rem">
  <!-- OVERVIEW -->
  <div class="section">
    <a href="{{ site.baseurl }}/project/overview">
        <div class="btn-primary">Overview</div>
    </a>
    <!-- <h6>Getting Started</h6> -->
    <ul>
        <!-- <li><a href="{{ site.baseurl }}/project">Project Overview</a></li> -->
        <li><a href="{{ site.baseurl }}/installation/quick-start">Quick Start</a></li>
        <!-- <li><a href="{{ site.baseurl }}/project">Essential Features</a></li>  -->
    </ul>
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{ site.baseurl }}/installation/" class="text-black">Supported Platforms</a>
        </p>
      </summary>
      <ul class="section-title">
      {% for item in sorted_pages %}
      {% if item.type=="installation" and item.category!="integrations" and item.list=="include" and item.language == "en" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        {% if item.abstract %}
          -  {{ item.abstract }}
        {% endif %}
        </li>
        {% endif %}
      {% endfor %}
    </ul>
    </details>
  </div>

    <!-- CONCEPTS -->
  <div class="section">
    <a href="{{ site.baseurl }}/concepts">
        <div class="btn-primary">Concepts</div>
    </a>
   <!-- <h6><a href="{{ site.baseurl }}/concepts/logical" class="text-black section-title">Conceptual</a></h6>
     <ul>
      {% for item in sorted_pages %}
      {% if item.type=="concepts" and item.list!="exclude" and item.language!="es" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul> -->
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{ site.baseurl }}/concepts/logical" class="text-black">Logical</a>
        </p>
      </summary>
      <ul class="section-title">
          {% for item in sorted_pages %}
          {% if item.type=="concepts" and item.language=="en" -%}
            <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
            {% if item.abstract != " " %}
              -  {{ item.abstract }}
            {% endif %}
            </li>
            {% endif %}
          {% endfor %}
      </ul>
    </details>
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{ site.baseurl }}/concepts/architecture" class="text-black section-title">Architectural</a>
        </p>
      </summary>
      <ul>
          {% for item in sorted_pages %}
          {% if item.type=="components" and item.language=="en" -%}
            <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
            {% if item.abstract != " " %}
              -  {{ item.abstract }}
            {% endif %}
            </li>
            {% endif %}
          {% endfor %}
      </ul>
    </details>
  </div>
</div>

<div class="flex flex-col--2"  style="text-align: left; padding:1.6rem ;--col-gap:1rem">

  <!-- TASKS -->
  <div class="section">
    <a href="{{ site.baseurl }}/tasks">
        <div class="btn-primary">Tasks</div>
    </a>
    <!-- <h6><a href="{{ site.baseurl }}/tasks" class="text-black section-title">Cloud Native Management</a></h6> -->
    <ul>
      {% for item in sorted_pages %}
      {% if item.type=="tasks" and item.list!="exclude" and item.language !="es" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
    <!-- <h6><a href="{{ site.baseurl }}/service-meshes" class="text-black section-title">Service Mesh Specific Management</a></h6> -->
  </div>

 <!-- Extensions -->
  <div class="section">
    <a href="{{ site.baseurl }}/extensions">
        <div class="btn-primary">Extensions</div>
    </a>
    <!-- <h6><a href="{{ site.baseurl }}/extensions" class="text-black section-title">Extensions</a></h6> -->
    <ul>
      {% for item in sorted_pages %}
      {% if item.type=="extensions" and item.list!="exclude" and item.language!="es" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
  </div>

</div> 

<div class="flex flex-col--2"  style="text-align: left; padding:1.6rem ;--col-gap:1rem">

<!-- GUIDES -->
  <div class="section">
    <a href="{{ site.baseurl }}/guides">
        <div class="btn-primary">Guides</div>
    </a>
    <!-- <h6><a href="{{ site.baseurl }}/guides" class="text-black section-title">Guides</a></h6> -->
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{ site.baseurl }}/guides/mesheryctl/" class="text-black">Meshery CLI</a>
        </p>
      </summary>
      <ul class="section-title">
          {% for item in sorted_pages %}
          {% if item.type=="guides" and item.category=="mesheryctl" and item.language=="en" -%}
            <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
            {% if item.abstract != " " %}
              -  {{ item.abstract }}
            {% endif %}
            </li>
            {% endif %}
          {% endfor %}
      </ul>
    </details>
    <ul>
      {% for item in sorted_pages %}
      {% if item.type=="guides" and item.category!="mesheryctl" and item.language=="en" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
    <!-- <h6><a href="{{ site.baseurl }}/service-meshes" class="text-black section-title">Service Mesh Specific Management</a></h6> -->
    <!-- <ul>
      {% for item in sorted_pages %}
      {% if item.type=="service-mesh" and item.list!="exclude" and item.language!="es"  -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
      {% for adapter in site.adapters -%}
      {% if adapter.project_status -%}
        <li><img src="{{ adapter.image }}" style="width:20px;height:20px; transform:translateY(5px)"/> <a href="{{ site.baseurl }}{{ adapter.url }}">{{ adapter.name }}</a></li>
      {% endif -%}
      {% endfor %}
    </ul> -->
  </div> 


    <!-- CONTRIBUTING -->
  <div class="section">
    <a href="{{ site.baseurl }}/project">
        <div class="btn-primary">Contributing and Community</div>
    </a>
    <!-- <h6><a href="{{ site.baseurl }}/tasks" class="text-black section-title">Cloud Native Management</a></h6> -->
    <ul>
      <li><a href="{{ site.baseurl }}/project/community" class="text-black">Community</a></li>
    </ul>
    <details>
      <summary>
        <p style="display:inline">
          <a href="{{ site.baseurl }}/project/contributing" class="text-black">Contributing</a>
        </p>
      </summary>
      <ul class="section-title">
          {% for item in sorted_pages %}
          {% if item.category=="contributing" and item.language=="en" -%}
            <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
            {% if item.abstract != " " %}
              -  {{ item.abstract }}
            {% endif %}
            </li>
            {% endif %}
          {% endfor %}
      </ul>
    </details>
  </div>
    
</div>
<div class="flex flex-col--2"  style="text-align: left; padding:1.6rem ;--col-gap:1rem">

    <!-- PROJECT -->
  <div class="section">
    <a href="{{ site.baseurl }}/project/overview">
        <div class="btn-primary">Project</div>
    </a>
    <!-- <h6><a href="{{ site.baseurl }}/tasks" class="text-black section-title">Cloud Native Management</a></h6> -->
    {% assign sorted_pages = site.pages | sort: "name" | alphabetical %}

    <ul>
      {% for item in sorted_pages %}
      {% if item.type=="project" and item.category!="contributing" and item.list=="include" and item.language =="en" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
  </div>

    <!-- REFERENCE -->
  <div class="section">
  <a href="{{ site.baseurl }}/installation/quick-start">
        <div class="btn-primary">Reference</div>
    </a>
    <!-- <h6><a href="{{ site.baseurl }}/reference" class="text-black section-title">Reference</a></h6> -->
    <ul>
        {% for item in sorted_pages %}
        {% if item.type=="Reference" and item.list!="exclude"  and item.language!="es"  -%}
          <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
          </li>
          {% endif %}
        {% endfor %}
      </ul>
    </div>

</div>

<p width="100%">Follow on <a href="https://twitter.com/mesheryio">Twitter</a> or subscribe to our <a href="https://meshery.io/subscribe">newsletter</a> for the latest updates. Get support on our <a href="http://discuss.meshery.io">forum</a>. Join our <a href="https://slack.meshery.io">Slack</a> to interact directly with other users and contributors.</p>





<!-- <div style="text-align:center;padding:0;margin:0;">
<img src="https://layer5.io/assets/images/meshery/meshery-logo-shadow-light-white-text-side.svg" width="60%" />
<h1>Documentation</h1>
</div> -->
