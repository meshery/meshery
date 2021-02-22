---
layout: page
title: Türkçe dokümantasyon 🇹🇷
permalink: /tr/index
language: tr
lang: tr
categories: tr
display-title: "false"
---
{% assign sorted_pages = site.pages | where: "language", "tr" | sort: "name" | alphabetical %}

<div style="display: block; text-align: center; margin-bottom: 30px;">
    <a href="https://layer5.io/meshery">
    <img style="width: calc(100% / 3.2); margin-bottom: 20px;" 
         src="/assets/img/meshery/meshery-logo-light-text.svg" />
    </a>
    <p>
      <h1>Türkçe dokümantasyon 🇹🇷</h1>
    </p>
</div>

<p style="margin:auto;padding:1rem;font-size: 1.25rem;">Meshery, herhangi bir hizmet ağının ve iş yüklerinin benimsenmesini, çalıştırılmasını ve yönetilmesini sağlayan açık kaynaklı, hizmet ağı yönetim düzlemidir.</p>
<div class="wrapper" style="text-align: left;">

  <!-- QUICK START -->
  <div>
    <a href="{{ site.baseurl }}/tr/installation/quick-start">
        <div class="overview">Hızlı&nbsp;başlangıç</div>
    </a>
    <ul><b>Başlarken</b>
        <li><a href="{{ site.baseurl }}/tr/overview">Meshery ile tanışın</a></li>
        <li><a href="{{ site.baseurl }}/tr/project">Proje ve Topluluk</a></li>        
    </ul>
    <ul><b><a href="{{ site.baseurl }}/tr/installation/platforms" class="text-black">Desteklenen Platformlar</a></b>
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
    <a href="{{ site.baseurl }}/tr/concepts">
        <div class="overview">Kavramlar</div>
    </a>
    <ul><b><a href="{{ site.baseurl }}/tr/concepts" class="text-black">Kavramlar</a></b>
      {% for item in sorted_pages %}
      {% if item.type=="concepts" and item.list!="exclude" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
    <ul><b><a href="{{ site.baseurl }}/functionality" class="text-black">İşlevsellik</a></b>
      {% for item in sorted_pages %}
        <!-- {{ item.title }}|{{ item.type }}|{{ item.list }}|{{ item.language }}<br> -->
      {% if item.type=="functionality" and item.list!="exclude" and item.language =="tr" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
    <ul><b><a href="{{ site.baseurl }}/service-meshes" class="text-black">Hizmet Ağı Yönetimi</a></b>
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
    <a href="{{ site.baseurl }}/tr/guides">
        <div class="overview">Kılavuzlar</div>
    </a>
    <ul><b><a href="{{ site.baseurl }}/tr/guides" class="text-black">Kılavuzlar</a></b>
      {% for item in sorted_pages %}
      {% if item.type=="Guides" and item.list!="exclude" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
    <ul><b><a href="{{ site.baseurl }}/tr/reference" class="text-black">Referans</a></b>
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

<!-- TODO: check about this value page=espanol or page=turkish? -->
{% include toc.html page=turk %}
