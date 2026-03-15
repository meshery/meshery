---
layout: default
title: Meshery CLI Guides
permalink: guides/mesheryctl
redirect_from: guides/mesheryctl/
language: en
type: guides
category: mesheryctl
list: exclude
abstract: Guides for common tasks while using Meshery's CLI, mesheryctl.
---

From the Command Line: Guides to using Meshery's various features and components.

{% assign sorted_guides = site.pages | sort: "name" %}

<ul>
  {% for item in sorted_guides %}
  {% if item.type=="guides" and item.category=="mesheryctl" and item.list!="exclude" and item.language=="en" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
    <li><a href="{{ site.baseurl }}/installation/upgrades#upgrading-meshery-cli">Upgrading Meshery CLI</a></li>
</ul>

{% include related-discussions.html tag="mesheryctl" %}

<!-- {% include toc.html page=Guides %} -->

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

