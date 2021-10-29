---
layout: page
title: Documentación en Español 🇲🇽
permalink: es
display-title: "false"
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

<!-- Contribuir Inicio-->
<!-- QUICK START -->
  <div>
    <a href="{{ site.baseurl }}/es/installation">
        <div class="overview">Inicio&nbsp;Rápido</div>
    </a>
    <ul>
        <b><a href="{{ site.baseurl }}/es/installation" class="text-black">Inicio Rápido</a></b>
        <li><a href="{{ site.baseurl }}/es/overview">Introducción a Meshery</a></li>
        <li><a href="{{ site.baseurl }}/es/project">Proyecto y Comunidad</a></li>
    </ul>
    <ul><b><a href="{{ site.baseurl }}/es/installation/platforms" class="text-black">Plataformas Soportadas</a></b>
        {% for item in sorted_pages %}
        {% if item.type=="installation" %}
          <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
          {% endif %}
        {% endfor %}
      </ul>
  </div>

{% include toc.html page=espanol %}
