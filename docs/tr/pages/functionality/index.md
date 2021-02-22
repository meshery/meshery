---
layout: page
title: İşlevsellik
permalink: tr/functionality
type: functionality
language: tr
lang: tr
categories: tr
list: exclude
---

# Meshery işlevleri

Bu bölümler, Meshery'nin sunduğu işlevler hakkında kullanıcı kılavuzları sağlar.

{% assign sorted_functionality = site.pages | sort: "functionality" %}

<ul>
    {% for item in sorted_functionality %}
    {% if item.type=="functionality" and item.language=="tr" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
      {% endif %}
    {% endfor %}
</ul>

{% include toc.html page=functionality %}

{:toc}
