---
layout: page
title: Documentación en Español 🇲🇽
permalink: es
language: es
lang: es
categories: es
display-title: "false"
---

{% assign sorted_pages = site.pages | where: "language", "es" | sort: "name" | alphabetical %}

<div style="display: block; text-align: center; margin-bottom: 30px;">
    <a href="https://layer5.io/meshery">
    <img style="width: calc(100% / 3.2); margin-bottom: 20px;"
         src="/assets/img/meshery/meshery-logo-light-text.svg" />
    </a>
    <p>
      <h1>Documentación en idioma Español 🇲🇽</h1>
    </p>
</div>

<!-- Contribuir Inicio-->
<p style="margin:auto;padding:1rem;font-size: 1.25rem;">Meshery es el plano de administración de malla de servicios (Service Mesh) de código abierto para permitir la adopción, operación y administración de cualquier malla de servicios y sus cargas de trabajo.</p>
<div class="wrapper" style="text-align: left;">

  <!-- QUICK START -->
  <div>
    <a href="{{ site.baseurl }}/es/installation/quick-start">
        <div class="overview">Inicio&nbsp;Rapido</div>
    </a>
    <ul><b>Empezando</b>
        <li><a href="{{ site.baseurl }}/es/overview">Presentando Meshery</a></li>
        <li><a href="{{ site.baseurl }}/es/project">Proyecto y Comunidad</a></li>        
    </ul>
    <ul><b><a href="{{ site.baseurl }}/es/installation/platforms" class="text-black">Plataformas compatibles</a></b>
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
    <a href="{{ site.baseurl }}/es/concepts">
        <div class="overview">Conceptos</div>
    </a>
    <ul><b><a href="{{ site.baseurl }}/es/concepts" class="text-black">Conceptos</a></b>
      {% for item in sorted_pages %}
      {% if item.type=="concepts" and item.list!="exclude" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
    <ul><b><a href="{{ site.baseurl }}/es/functionality" class="text-black">Funcionalidad</a></b>
      {% for item in sorted_pages %}
        <!-- {{ item.title }}|{{ item.type }}|{{ item.list }}|{{ item.language }}<br> -->
      {% if item.type=="functionality" and item.list!="exclude" and item.language =="es" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
    <ul><b><a href="{{ site.baseurl }}/es/service-meshes" class="text-black">Gestión de la malla de servicios</a></b>
      {% for item in sorted_pages %}
      {% if item.type=="service-mesh" and item.list!="exclude" and item.language =="es" -%}
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
    <a href="{{ site.baseurl }}/es/guides">
        <div class="overview">Guías</div>
    </a>
    <ul><b><a href="{{ site.baseurl }}/es/guides" class="text-black">Guías</a></b>
      {% for item in sorted_pages %}
      {% if item.type=="Guides" and item.list!="exclude" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
    <ul><b><a href="{{ site.baseurl }}/es/reference" class="text-black">Referencia</a></b>
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
{% include toc.html page=espanol %}
