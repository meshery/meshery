---
layout: default
title: Compatibility Matrix
permalink: project/compatibility-matrix
redirect_from: project/compatibility-matrix/
language: en
display-title: "false"
list: exclude
---

# Meshery Compatibility Matrix

{% assign sorted_tests = site.compatibility | sort: "type" | reverse %}

<table>
  <th>Status</th>
  <th>Meshery Component</th>
  <th>Meshery Component Version</th>
  <th>Meshery Server Version</th>
  <th>Service Mesh</th>
  <th>Service Mesh Version</th>

    {% for item in site.compatibility %}
    {% if item.overall-status == "passing" %}
      {% assign overall-status = "background-color: #83B71E;" %}
    {% elsif item.overall-status == "partial" %}
      {% assign overall-status = "background-color: #EBC017;" %}
    {% elsif item.overall-status == "failing" %}
      {% assign overall-status = "background-color: #B32700;" %}
    {% else %}
      {% assign overall-status = "" %}
    {% endif %}
    <tr>
      <td style="{{ overall-status }}"></td>
      <td><a href="{{ site.repo }}-{{ item.service-mesh }}">{{ item.meshery-component }}</a></td>
      <td><a href="{{ site.repo }}-{{ item.service-mesh }}/releases/tag/{{ item.meshery-component-version }}">{{ item.meshery-component-version }}</a></td>
      <td><a href="{{ site.repo }}-{{ item.service-mesh }}">{{ item.meshery-server-version }}</a></td>
      <td><a href="{{ site.repo }}-{{ item.service-mesh }}">{{ item.service-mesh }}</a></td>
      <td><a href="{{ site.repo }}-{{ item.service-mesh }}">{{ item.service-mesh-version }}</a></td>
    </tr>
    {% endfor %}

</table>
