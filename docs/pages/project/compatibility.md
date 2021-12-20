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

<script type="text/javascript">
<!--
    function toggle_visibility(id) {
       var e = document.getElementById(id);
       if(e.style.visibility == 'visible') {
          e.style.display = 'none';
          e.style.visibility = 'hidden';
      }
       else {
         
          e.style.display = 'table-row';
          e.style.visibility = 'visible';
          }
    }
//-->
</script>

<style>
  tr:hover {
    background-color: #ffff99;
   }
</style>

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
      {% assign overall-status = "background-color: #83B71E; color: white;" %}
    {% elsif item.overall-status == "partial" %}
      {% assign overall-status = "background-color: #EBC017; color: white;" %}
    {% elsif item.overall-status == "failing" %}
      {% assign overall-status = "background-color: #B32700; color: white;" %}
    {% else %}
      {% assign overall-status = "" %}
    {% endif %}
    <tr onclick="toggle_visibility('{{forloop.index}}');"> 
      <td style="{{ overall-status }}">{{ item.timestamp }}</td>
      <td><a href="{{ site.repo }}-{{ item.service-mesh }}">{{ item.meshery-component }}</a></td>
      <td><a href="{{ site.repo }}-{{ item.service-mesh }}/releases/tag/{{ item.meshery-component-version }}">{{ item.meshery-component-version }}</a></td>
      <td><a href="{{ site.repo }}-{{ item.service-mesh }}">{{ item.meshery-server-version }}</a></td>
      <td><a href="{{ site.repo }}-{{ item.service-mesh }}">{{ item.service-mesh }}</a></td>
      <td><a href="{{ site.repo }}-{{ item.service-mesh }}">{{ item.service-mesh-version }}</a></td>
    </tr>
    <tr id="{{forloop.index}}" style="visibility:hidden; display:none;">
      <td colspan="6">
        Test results:
        {% for testa in item.tests %}
          <li>{{ testa[0] }}: {{testa[1] }}</li>
        {% endfor %}      
        <!-- {{ item.tests }} -->
        
      </td>
      
    </tr>
    {% endfor %}

</table>
