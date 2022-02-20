---
layout: default
title: Compatibility Matrix
permalink: project/compatibility-matrix
redirect_from: /
description: a complete compatibility matrix and project test status dashboard.
language: en
display-title: "false"
list: exclude
type: "project"
---

# Meshery Compatibility Matrix

## Integration Tests

As a key aspect of Meshery, its integrations with other systems are routinely tested. Unit, integration testing occurs before and after every pull request (before code is to be merged into the project and after code is merged into the project). Regression tests are run nightly.

<script type="text/javascript">
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
</script>

<style>
  
  td:hover, tr:hover {
    background-color: #ccfff9;
    cursor:pointer;
   }
  td.details {
    background-color: #fafafa;
    cursor:text;
  }
</style>

{% assign sorted_tests_group = site.compatibility | group_by: "meshery-component" %}

<table>
  <th>Status</th>
  <th>Meshery Component</th>
  <th>Meshery Component Version</th>
  <th>Meshery Server Version</th>
  <th>Service Mesh</th>
  <th>Service Mesh Version</th>


    {% for group in sorted_tests_group %}
      {% assign items = group.items | sort: "timestamp" | reverse %}
      {% for item in items limit: 1 %}
        {% if item.meshery-component-version == "edge" %}
          {% if item.overall-status == "passing" %}
            {% assign overall-status = "background-color: #83B71E; color: white;" %}
          {% elsif item.overall-status == "partial" %}
            {% assign overall-status = "background-color: #EBC017; color: white;" %}
          {% elsif item.overall-status == "failing" %}
            {% assign overall-status = "background-color: #B32700; color: white;" %}
          {% else %}
            {% assign overall-status = "" %}
          {% endif %}
          <tr onclick="toggle_visibility('{{item.meshery-component}}');"> 
            <td style="{{ overall-status }}">{{ item.timestamp }}</td>
            <td><a href="{{ site.repo }}-{{ item.service-mesh }}">{{ item.meshery-component }}</a></td>
            {% if item.meshery-component-version == "edge" %}
              <td><a href="{{ site.repo }}-{{ item.service-mesh }}/releases">{{ item.meshery-component-version }}</a></td>
            {% else %}
              <td><a href="{{ site.repo }}-{{ item.service-mesh }}/releases/tag/{{ item.meshery-component-version }}">{{ item.meshery-component-version }}</a></td>
            {% endif %}
            {% if item.meshery-server-version == "edge" %}
              <td><a href="{{ site.repo }}/releases{{ item.meshery-server-version }}">{{ item.meshery-server-version }}</a></td>
            {% else %}
              <td><a href="{{ site.repo }}/releases/tag/{{ item.meshery-server-version }}">{{ item.meshery-server-version }}</a></td>
            {% endif %}
            <td><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/{{item.service-mesh | downcase }}.svg" />&nbsp;<a href="{{ site.baseurl }}/service-meshes/adapters/{{ item.service-mesh }}">{{ item.service-mesh }}</a></td>
            <td>{{ item.service-mesh-version }}</td>
          </tr>
          <tr id="{{item.meshery-component}}" style="visibility:hidden; display:none;">
            <td colspan="2" class="details">
              <i>Platform:</i>
              <li><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/kubernetes-icon-color.svg" />  {{item.k8s-distro}}  {{item.k8s-version}}</li>
            </td>
            <td colspan="3" class="details">
              <i>Test results:</i>
              <ol>
              {% for test in item.tests %}
                <li>{{ test[0] }}: {{test[1] }}</li>
              {% endfor %}      
              </ol>      
            </td>
            <td>
              <a href = "{{site.baseurl}}/project/compatibility-matrix/{{item.meshery-component}}-past-results">To see past results click here </a>
            </td>
          </tr>

        <!-- if the latest test is stable as we require edge test to show too and since sorted through timestamp second element will always be an edge tests. -->

        {% else %} 
          {% if items[1].overall-status == "passing" %}
            {% assign overall-status = "background-color: #83B71E; color: white;" %}
          {% elsif items[1].overall-status == "partial" %}
            {% assign overall-status = "background-color: #EBC017; color: white;" %}
          {% elsif items[1].overall-status == "failing" %}
            {% assign overall-status = "background-color: #B32700; color: white;" %}
          {% else %}
            {% assign overall-status = "" %}
          {% endif %}
          <tr onclick="toggle_visibility('{{items[1].meshery-component}}');"> 
            <td style="{{ overall-status }}">{{ items[1].timestamp }}</td>
            <td><a href="{{ site.repo }}-{{ items[1].service-mesh }}">{{ items[1].meshery-component }}</a></td>
            {% if items[1].meshery-component-version == "edge" %}
              <td><a href="{{ site.repo }}-{{ items[1].service-mesh }}/releases">{{ items[1].meshery-component-version }}</a></td>
            {% else %}
              <td><a href="{{ site.repo }}-{{ items[1].service-mesh }}/releases/tag/{{ items[1].meshery-component-version }}">{{ items[1].meshery-component-version }}</a></td>
            {% endif %}
            {% if items[1].meshery-server-version == "edge" %}
              <td><a href="{{ site.repo }}/releases{{ items[1].meshery-server-version }}">{{ items[1].meshery-server-version }}</a></td>
            {% else %}
              <td><a href="{{ site.repo }}/releases/tag/{{ items[1].meshery-server-version }}">{{ items[1].meshery-server-version }}</a></td>
            {% endif %}
            <td><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/{{items[1].service-mesh | downcase }}.svg" />&nbsp;<a href="{{ site.baseurl }}/service-meshes/adapters/{{ items[1].service-mesh }}">{{ items[1].service-mesh }}</a></td>
            <td>{{ items[1].service-mesh-version }}</td>
          </tr>
          <tr id="{{items[1].meshery-component}}" style="visibility:hidden; display:none;">
            <td colspan="2" class="details">
              <i>Platform:</i>
              <li><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/kubernetes-icon-color.svg" />  {{items[1].k8s-distro}}  {{items[1].k8s-version}}</li>
            </td>
            <td colspan="3" class="details">
              <i>Test results:</i>
              <ol>
              {% for test in items[1].tests %}
                <li>{{ test[0] }}: {{test[1] }}</li>
              {% endfor %}      
              </ol>      
            </td>
            <td>
              <a href = "{{site.baseurl}}/project/compatibility-matrix/{{item.meshery-component}}-past-results">To see past results click here </a>
            </td>
          </tr>        
        {% endif %}  
      {% endfor %}
    {% endfor %}
    
    <!-- display tests from the stable channel -->

    {% for group in sorted_tests_group %}
      {% assign items = group.items | sort: "meshery-component-version" | reverse %}
      {% for item in items limit: 1 %}
        {% if item.meshery-component-version != "edge" %}
          {% if item.overall-status == "passing" %}
            {% assign overall-status = "background-color: #83B71E; color: white;" %}
          {% elsif item.overall-status == "partial" %}
            {% assign overall-status = "background-color: #EBC017; color: white;" %}
          {% elsif item.overall-status == "failing" %}
            {% assign overall-status = "background-color: #B32700; color: white;" %}
          {% else %}
            {% assign overall-status = "" %}
          {% endif %}
          <tr onclick="toggle_visibility('{{item.meshery-component}}-stable');"> 
            <td style="{{ overall-status }}">{{ item.timestamp }}</td>
            <td><a href="{{ site.repo }}-{{ item.service-mesh }}">{{ item.meshery-component }}</a></td>
            <td><a href="{{ site.repo }}-{{ item.service-mesh }}/releases/tag/{{ item.meshery-component-version }}">{{ item.meshery-component-version }}</a></td>
            <td><a href="{{ site.repo }}/releases/tag/{{ item.meshery-server-version }}">{{ item.meshery-server-version }}</a></td>
            <td><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/{{item.service-mesh | downcase }}.svg" />&nbsp;<a href="{{ site.baseurl }}/service-meshes/adapters/{{ item.service-mesh }}">{{ item.service-mesh }}</a></td>
            <td>{{ item.service-mesh-version }}</td>
          </tr>
          <tr id="{{item.meshery-component}}-stable" style="visibility:hidden; display:none;">
            <td colspan="2" class="details">
              <i>Platform:</i>
              <li><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/kubernetes-icon-color.svg" />  {{item.k8s-distro}}  {{item.k8s-version}}</li>
            </td>
            <td colspan="3" class="details">
              <i>Test results:</i>
              <ol>
              {% for test in item.tests %}
                <li>{{ test[0] }}: {{test[1] }}</li>
              {% endfor %}      
              </ol>      
            </td>
            <td>
              <a href = "{{site.baseurl}}/project/compatibility-matrix/{{item.meshery-component}}-past-results">To see past results click here </a>
            </td>
          </tr>
        {% endif %}  
      {% endfor %}
    {% endfor %}

</table>
