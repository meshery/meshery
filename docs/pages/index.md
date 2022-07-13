---
layout: page
title: Meshery Documentation
permalink: /
display-title: "false"
---

{% assign sorted_pages = site.pages | sort: "name" | alphabetical %}

<div style="display: block; text-align: center; margin-bottom: 30px;">
    <a href="https://layer5.io/meshery">
    <img style="width: calc(100% / 2.5); " 
         src="/assets/img/meshery/meshery-logo-light-text-side.svg" />
    </a>
</div>
<p style="margin:auto;padding:1rem;font-size: 1.25rem;">Meshery is the open source, cloud native management plane that enables the adoption, operation, and management of Kubernetes, any service mesh, and their workloads.</p>

<p style="padding:1rem;">Follow on <a href="https://twitter.com/mesheryio">Twitter</a> or subscribe to our <a href="https://meshery.io/subscribe">newsletter</a> for the latest updates. Get support on our <a href="https://discuss.layer5.io">forum</a>. Join our <a href="https://slack.layer5.io">Slack</a> to interact directly with other users and contributors.</p>

<div class="wrapper" style="text-align: left;">

  <!-- QUICK START -->
  <div class="section">
    <a href="{{ site.baseurl }}/installation/quick-start">
        <div class="overview">Quick&nbsp;Start</div>
    </a>
    <h6>Getting Started</h6>
    <ul>
        <li><a href="{{ site.baseurl }}/getting-started/overview">Meshery Overview</a></li>
        <li><a href="{{ site.baseurl }}/project/community">Community</a></li>
        <li><a href="{{ site.baseurl }}/project/contributing">Contributing</a></li> 
    </ul>
    <h6><a href="{{ site.baseurl }}/installation/platforms" class="text-black section-title">Supported Platforms</a></h6>
    <ul>
        {% for item in sorted_pages %}
        {% if item.type=="installation" and item.list=="include" and item.language!="es"  -%}
          <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
          </li>
          {% endif %}
        {% endfor %}
      </ul>
  </div>
  
  <!-- CONCEPTS -->
  <div class="section">
    <a href="{{ site.baseurl }}/concepts">
        <div class="overview">Concepts</div>
    </a>
    <h6><a href="{{ site.baseurl }}/concepts" class="text-black section-title">Concepts</a></h6>
    <ul>
      {% for item in sorted_pages %}
      {% if item.type=="concepts" and item.list!="exclude" and item.language!="es" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
  </div>

</div>
<div class="wrapper" style="text-align: left;">

  <!-- FUNCTIONALITY -->
  <div class="section">
    <a href="{{ site.baseurl }}/functionality">
        <div class="overview">Functionality</div>
    </a>
    <h6><a href="{{ site.baseurl }}/functionality" class="text-black section-title">Cloud Native Management</a></h6>
    <ul>
      {% for item in sorted_pages %}
      {% if item.type=="functionality" and item.list!="exclude" and item.language !="es" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
    <h6><a href="{{ site.baseurl }}/service-meshes" class="text-black section-title">Service Mesh Specific Management</a></h6>
    <ul>
      {% for item in sorted_pages %}
      {% if item.type=="service-mesh" and item.list!="exclude" and item.language!="es"  -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
      {% for adapter in site.adapters -%}
      {% if adapter.project_status -%}
        <li><img src="{{ adapter.image }}" style="width:20px" /> <a href="{{ site.baseurl }}{{ adapter.url }}">{{ adapter.name }}</a></li>
      {% endif -%}
      {% endfor %}
    </ul>
  </div>

  <!-- GUIDES -->
  <div class="section">
    <a href="{{ site.baseurl }}/guides">
        <div class="overview">Guides</div>
    </a>
    <h6><a href="{{ site.baseurl }}/guides" class="text-black section-title">Guides</a></h6>
    <ul>
      {% for item in sorted_pages %}
      {% if item.type=="Guides" and item.list!="exclude"  and item.language!="es" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
    <h6><a href="{{ site.baseurl }}/reference" class="text-black section-title">Reference</a></h6>
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

<!-- <div style="text-align:center;padding:0;margin:0;">
<img src="https://layer5.io/assets/images/meshery/meshery-logo-shadow-light-white-text-side.svg" width="60%" />
<h1>Documentation</h1>
</div> -->

