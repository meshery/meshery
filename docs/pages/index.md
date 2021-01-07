---
layout: page
title: Meshery Documentation
permalink: /
display-title: "false"
---

{% assign sorted_pages = site.pages | sort: "name" | alphabetical %}


<div style="display: block; text-align: center; margin-bottom: 30px;">
    <a href="https://layer5.io/meshery">
    <img style="width: calc(100% / 3.2); " 
         src="/assets/img/meshery/meshery-logo-light-text.svg" />
    </a>
</div>
<p style="margin:auto;padding:1rem;font-size: 1.25rem;">Meshery is the open source, service mesh management plane for enabling the adoption, operation, and management of any service mesh and their workloads.</p>
<div class="wrapper" style="text-align: left;">

  <!-- QUICK START -->
  <div>
    <a href="{{ site.baseurl }}/installation/quick-start">
        <div class="overview">Quick&nbsp;Start</div>
    </a>
    <ul><b>Getting Started</b>
        <li><a href="{{ site.baseurl }}/overview">Introducing Meshery</a></li>
        <li><a href="{{ site.baseurl }}/project">Project and Community</a></li>        
    </ul>
    <ul><b><a href="{{ site.baseurl }}/installation/platforms" class="text-black">Supported Platforms</a></b>
        {% for item in sorted_pages %}
        {% if item.type=="installation" and item.list=="include" -%}
          <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
          </li>
          {% endif %}
        {% endfor %}
      </ul>
  </div>
  
  <!-- CONCEPTS -->
  <div>
    <a href="{{ site.baseurl }}/concepts">
        <div class="overview">Concepts</div>
    </a>
    <ul><b><a href="{{ site.baseurl }}/concepts" class="text-black">Concepts</a></b>
      {% for item in sorted_pages %}
      {% if item.type=="concepts" and item.list!="exclude" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
    <ul><b><a href="{{ site.baseurl }}/functionality" class="text-black">Functionality</a></b>
      {% for item in sorted_pages %}
      {% if item.type=="functionality" and item.list!="exclude" and item.language !="es" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
    <ul><b><a href="{{ site.baseurl }}/service-meshes" class="text-black">Service Mesh Management</a></b>
      {% for item in sorted_pages %}
      {% if item.type=="service-mesh" and item.list!="exclude" -%}
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
  <div>
    <a href="{{ site.baseurl }}/guides">
        <div class="overview">Guides</div>
    </a>
    <ul><b><a href="{{ site.baseurl }}/guides" class="text-black">Guides</a></b>
      {% for item in sorted_pages %}
      {% if item.type=="Guides" and item.list!="exclude" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
    <ul><b><a href="{{ site.baseurl }}/reference" class="text-black">Reference</a></b>
        {% for item in sorted_pages %}
        {% if item.type=="Reference" and item.list!="exclude" -%}
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
