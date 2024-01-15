---
layout: default
title: Tutorials
display_title: true
permalink: guides/tutorials
redirect_from: guides/tutorials/
---

{% assign tutorials = site.pages | where: "type","tutorials" %}

<ul>
  {% for tutorial in tutorials %}
    <li><a href="{{ site.baseurl }}{{ tutorial.url }}">{{ tutorial.title }}</a>
    </li>

  {% endfor %}
</ul>