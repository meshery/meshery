---
layout: default
title: Permissions
permalink: reference/permissions
redirect_from: reference/permissions
language: en
list: exclude
abstract: List of default permissions.
---
<style>
      .scrollable-container {
        overflow-x: auto; 
      }
      table {
        border-collapse: collapse;
        width: auto; 
      }
    </style>

<table>
  {% for row in site.data.keys %}
    {% if forloop.first %}
    <tr>
      {% for pair in row %}
        <th>{{ pair[0] }}</th>
      {% endfor %}
    </tr>
    {% endif %}

    {% tablerow pair in row %}
      {{ pair[1] }}
    {% endtablerow %}
  {% endfor %}
</table>