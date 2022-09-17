---
layout: page
title: Documentación en Español 🇲🇽
permalink: es
display-title: "false"
language: es
---

{% assign sorted_pages = site.pages | where: "language", "es" | sort: "name" | alphabetical %}

<div style="display: block; text-align: center; margin-bottom: 30px;">
    <a href="https://layer5.io/meshery">
    <img style="width: calc(100% / 3.2); margin-bottom: 20px;"
         src="/assets/img/meshery/meshery-logo-light-text.svg" />
    </a>
    <p>
      <h1>Documentación en Español 🇲🇽</h1>
    </p>
</div>

<div class = "wrapper">

  <!-- Contribuir Inicio-->
  <!-- QUICK START -->

  <div>
    <a href="{{ site.baseurl }}/es/installation">
      <div class="overview">Inicio&nbsp;Rápido</div>
    </a>
    <b><a href="{{ site.baseurl }}/es/installation" class="text-black">Inicio Rápido</a></b>
    <ul>
      <li><a href="{{ site.baseurl }}/es/overview">Introducción a Meshery</a></li>
      <li><a href="{{ site.baseurl }}/es/project">Proyecto y Comunidad</a></li>
      <li><a href="{{ site.baseurl }}/es/project/contributing">Contribuir</a></li>
    </ul>
    <b><a href="{{ site.baseurl }}/es/installation/platforms" class="text-black">Plataformas Soportadas</a></b>
    <ul>
      {% for item in sorted_pages %}
        {% if item.type=="installation" %}
          <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
        {% endif %}
      {% endfor %}
    </ul>
  </div>

  <!-- Conceptos -->
  <!-- Concepts -->
  <div>
    <a href="{{ site.baseurl }}/concepts">
        <div class="overview">Conceptos</div>
    </a>
    <h6><a href="{{ site.baseurl }}/concepts" class="text-black section-title">Conceptos</a></h6>
    <ul>
      {% for item in sorted_pages %}
        {% if item.type=="concepts" and item.list!="exclude" -%}
          <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
        {% endif %}
      {% endfor %}
    </ul>
  </div>
</div>

<div class="wrapper" style="text-align: left;">
  
  <!-- Funcionalidad -->
  <!-- FUNCTIONALITY -->

  <div>
    <a href="{{ site.baseurl }}/es/functionality">
      <div class="overview">Funcionalidad</div>
    </a>
    <h6><a href="{{ site.baseurl }}/es/functionality" class="text-black section-title">Gestión de la malla de servicios</a></h6>
    <ul>
      {% for item in sorted_pages %}
        {% if item.type=="functionality" and item.list!="exclude" %}
          <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
        {% endif %}
      {% endfor %}
    </ul>
    <h6><a href="{{ site.baseurl }}/service-meshes" class="text-black section-title">Gestión específica de la malla de servicios</a></h6>
    <ul>
      {% for item in sorted_pages %}
        {% if item.type=="service-mesh" and item.list!="exclude" %}
          <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a> hello </li>
        {% endif %}
      {% endfor %}
      {% for adapter in site.adapters %}
        {% if adapter.project_status %}
          <li><img src="{{ adapter.image }}" style="width:20px" /> <a href="{{ site.baseurl }}{{ adapter.url }}">{{ adapter.name }}</a></li>
        {% endif -%}
      {% endfor %}
    </ul>
  </div>

  <!-- Guías -->
  <!-- GUIDES -->

  <div>
    <a href="{{ site.baseurl }}/guides">
        <div class="overview">Guías y referencia</div>
    </a>
    <h6><a href="{{ site.baseurl }}/guides" class="text-black section-title">Guías</a></h6>
    <ul>
      {% for item in sorted_pages %}
        {% if item.type=="Guides" and item.list!="exclude" %}
          <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
        {% endif %}
      {% endfor %}
    </ul>
    <h6><a href="{{ site.baseurl }}/reference" class="text-black section-title">Referencia</a></h6>
    <ul>
      {% for item in sorted_pages %}
        {% if item.type=="Reference" and item.list!="exclude" %}
          <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
        {% endif %}
      {% endfor %}
    </ul>
  </div>

</div>
