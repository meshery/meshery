---
layout: default
title: Operation Guides
permalink: guides
language: en
list: exclude
---

Guides to using Meshery's various features and components. 

{% assign sorted_guides = site.pages | sort: "type" | reverse %}



<ul><b>General</b>
    {% for item in sorted_guides %}
    {% if item.type=="Guides" and item.category!="mesheryctl" and item.list!="exclude"  -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
      {% endif %}
    {% endfor %}
</ul>

<ul><b><a href="{{ site.baseurl }}/guides/mesheryctl" class="text-black">Meshery CLI</a></b>
  {% for item in sorted_guides %}
  {% if item.type=="Guides" and item.category=="mesheryctl" and item.list!="exclude" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
</ul>

{% include toc.html page=Guides %}

{:toc}

<!-- {% comment %}
#
#  Change date order by adding '| reversed'
#  To sort by title or other variables use {% assign sorted_posts = category[1] | sort: 'title' %}
#
{% endcomment %}

{% for guide in site.adapter %}
<h2 id="{{guide[0] | uri_escape | downcase }}">{{guide[0] | capitalize}}1</h2>

{% endfor %}

{% assign sorted_guides = site.guides | sort %}
{% for guide in sorted_guides %}
<h2 id="{{guide[0] | uri_escape | downcase }}">{{guide[0] | capitalize}}</h2>

{% endfor %} -->