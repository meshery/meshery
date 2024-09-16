---
layout: default
title: Tutorials
display_title: true
type: guides
category: tutorials
language: en
list: exclude
abstract: "Explore the tutorials to learn how to use Meshery for collaboratively managing infrastructure."
permalink: guides/tutorials
redirect_from: guides/tutorials/
---
🧑‍🔬 Explore these tutorials to learn how to use Meshery for collaboratively managing infrastructure. Access the [Meshery Playground]({{site.baseurl}}/installation/playground) as a convenient resource for the labs in these tutorials.

{% assign tutorials = site.pages | where: "category", "tutorials" %}

{% assign items_grouped = tutorials | group_by: 'model' %}
{% for group in items_grouped %}
  {% if group.name != "" %}
  <h2>{{ group.name | capitalize }} </h2>
  
    {% for item in group.items %}
0. [{{ item.title }}]({{ site.baseurl }}{{ item.url }})
    {% endfor %}
  {% endif %}
{% endfor %}
