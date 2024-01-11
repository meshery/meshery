---
layout: default
title: Articles
permalink: /test
---
  <div class="section">
    <a href="{{ site.baseurl }}/extensions">
        <div class="btn-primary">Extensions</div>
    </a>
    <!-- <h6><a href="{{ site.baseurl }}/extensions" class="text-black section-title">Extensions</a></h6> -->
    <ul>
      {% assign sorted_items = site.pages | category: "name" | alphabetical %}
      <pre>{{ sorted_items }}</pre>
      {% for item in sorted_items %}
      {% if item.type=="extensions" and item.list!="exclude" and item.language!="es" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        </li>
        {% endif %}
      {% endfor %}
    </ul>
  </div>