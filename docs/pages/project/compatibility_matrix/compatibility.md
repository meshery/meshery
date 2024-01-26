---
layout: default
title: Compatibility Matrix
permalink: installation/compatibility-matrix
redirect_from: 
 - project/compatibility-matrix/
abstract: An installation compatibility matrix and project test status dashboard.
language: en
display-title: "false"
list: include
type: "installation"
---

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
    function handleEdgeCheckboxChange(){
        let e = document.getElementsByClassName("edge")
        let stable = document.getElementsByClassName("stable")
        let stable_box = document.getElementById("checkbox_stable")
        for(let i = 0; i <e.length;i++){
            if(e[i].classList.contains("edge_visible") ){
                e[i].classList.remove("edge_visible")
                if(!stable_box.checked){
                    stable_box.checked=true;
                    handleStableCheckboxChange();
                }
        }
        else{
                e[i].classList.add("edge_visible")
        }
        }
    }
    function handleStableCheckboxChange(){
        let e = document.getElementsByClassName("stable")
        let edge_box = document.getElementById("checkbox_edge")
        for(let i = 0; i <e.length;i++){
            if(e[i].classList.contains("stable_visible")){
                e[i].classList.remove("stable_visible")
                if(!edge_box.checked){
                    edge_box.checked=true;
                    handleEdgeCheckboxChange();
                }
        }
        else{
                e[i].classList.add("stable_visible")
        }
        }
    }

    
</script>

<style>
td:hover,tr:hover {
    background-color: #ccfff9;
    cursor:pointer;
  }
  .edge_visible{
    display: table-row !important;
    visibility: visible !important;
  }
  .stable_visible{
    display: table-row !important;
    visibility: visible !important;
  }
  .checkbox{
    display: flex;
    justify-content: flex-end;
    align-items: center;
    text-align: left;
  }
  td.details {
    background-color: #fafafa;
    cursor:text;
  }
  .edge_test_text{
    margin-right: 20px;
  }

.tablesorter-header-inner{
    display: flex;
    justify-content: space-between;
}
th span.sort-by { 
	padding-right: 18px;
	position: relative;
    text-decoration: none;
    margin-right:0 ;
}
span.sort-by:before,
span.sort-by:after {
	border: 4px solid transparent;
	content: "";
	display: block;
	height: 0;
	right: 5px;
	top: 50%;
	position: absolute;
	width: 0;
}
span.sort-by:before {
	border-bottom-color: #666;
	margin-top: -9px;
}
span.sort-by:after {
	border-top-color: #666;
	margin-top: 1px;
}

table.tablesorter th.headerSortDown,
table.tablesorter th.headerSortUp {
background-color: #8dbdd8;
}
.tablesorter-headerRow{
    background-color: rgb(244, 244, 244);
}
.tablesorter-headerRow:hover{
    background-color:rgb(244, 244, 244) !important ;
}

</style>

{% assign sorted_tests_group = site.compatibility | group_by: "meshery-component" %}
{% assign k8s_tests_group = site.compatibility | group_by: "k8s-version" | sort: "name" | reverse %}
{% assign service_meshes = site.adapters  %}

 <div>
 </div>

# Compatibility Matrix

Meshery Server and Meshery Adapters are tested daily for their compatibility with the infrastructure they manage and the platforms Meshery deploys on (Kubernetes and Docker). End-to-end test results are automatically posted to the following compatibility matrix.


{%include compatibilityMatrix.md k8s_tests_group=k8s_tests_group %}

## Integration Tests

As a key aspect of Meshery, its integrations with other systems are routinely tested. Unit and integration tests before and after every pull request (before code is to be merged into the project and after code is merged into the project). End-to-end tests are run nightly and automatically posted to the following test matrix.

<div class="checkbox">
    <div>
    <input onchange="handleEdgeCheckboxChange();" type="checkbox" id="checkbox_edge" value="Edge Tests" checked>
    <label for="checkbox_edge" class="edge_test_text">Edge Channel</label>
    </div>
    <div>
    <input onchange="handleStableCheckboxChange();" type="checkbox" id="checkbox_stable" value="Stable Tests" checked>
    <label for="checkbox_stable">Stable Channel</label>
    </div>
</div>

<table id = "test-table" class="tablesorter" style="text-align: center;">
  <thead>
    <th style="text-align: center;" >Status <span class="sort-by">  </span></th>
        <th style="text-align: center;" ><span class="sort-by"> Meshery Component  </span> </th>
        <th style="text-align: center;" ><span class="sort-by"> Meshery Component Version </span> </th>
        <th style="text-align: center;" ><span class="sort-by"> Meshery Server Version </span> </th>
        <th style="text-align: center;" > <span class="sort-by"> Service Mesh </span> </th>
        <th style="text-align: center;" ><span class="sort-by"> Service Mesh Version </span></th>
  </thead>
  <tbody>
    {% for group in sorted_tests_group %}
      {% assign items = group.items | sort: "timestamp" | reverse %}
      {% for item in items limit: 1 %}
        {% if item.meshery-component-version == "edge" %}
          {% if item.overall-status == "passing" %}
            {% assign overall-status = "background-color: #56B257; color: white;" %}
            {% assign result-state = "/assets/img/passing.svg" %}
          {% elsif item.overall-status == "partial" %}
            {% assign overall-status = "background-color: #EBC017; color: white;" %}
            {% assign result-state = "/assets/img/YellowCheck.svg" %}
          {% elsif item.overall-status == "failing" %}
            {% assign overall-status = "background-color: #B32700; color: white;" %}
            {% assign result-state = "/assets/img/failing.svg" %}
          {% else %}
            {% assign overall-status = "" %}
          {% endif %}
          <tr style="visibility: hidden; display: none; background: white" class="test-details edge edge_visible" onclick="toggle_visibility('{{item.meshery-component}}');">
            <td style="{{ overall-status }}">{{ item.timestamp }}</td>
            <td style="white-space:nowrap;"><a href="{{ site.repo }}-{{ item.service-mesh }}">{{ item.meshery-component }}</a></td>
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
            <td style="white-space: nowrap;"><img style="height: 2rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/{{item.service-mesh | downcase }}.svg" />&nbsp;<a href="{{ site.baseurl }}/service-meshes/adapters/{{ item.service-mesh }}">{{ item.service-mesh }}</a></td>
            <td>{{ item.service-mesh-version }}</td>
          </tr>
          <tr class="hidden-details" id="{{item.meshery-component}}" style="visibility:hidden; display:none;">
            <td colspan="2" class="details">
              <i>Platform:</i>
              <li><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/kubernetes-icon-color.svg" />  {{item.k8s-distro}}  {{item.k8s-version}}</li>
            </td>
            <td colspan="3" class="details">
              <i>Test results:</i>
              <table style="border:0">
              {% for test in item.tests %}
                  <tr><td><img style="height: 24px; width: 24px" src="{{ result-state }}"></td><td>{{test[0] }}</td></tr>
              {% endfor %}
              </table>
            </td>
            <td>
              <a href = "{{site.baseurl}}/installation/compatibility-matrix/{{item.meshery-component}}-past-results">To see past results click here </a>
            </td>
          </tr>

        <!-- if the latest test is stable as we require edge test to show too and since sorted through timestamp second element will always be an edge tests. -->

        {% elsif items[1].meshery-component-version == "edge" %}
          {% if items[1].overall-status == "passing" %}
            {% assign overall-status = "background-color: #56B257; color: white;" %}
            {% assign result-state = "/assets/img/passing.svg" %}
          {% elsif items[1].overall-status == "partial" %}
            {% assign overall-status = "background-color: #EBC017; color: white;" %}
            {% assign result-state = "/assets/img/YellowCheck.svg" %} \
          {% elsif items[1].overall-status == "failing" %}
            {% assign overall-status = "background-color: #B32700; color: white;" %}
            {% assign result-state = "/assets/img/failing.svg" %}
          {% else %}
            {% assign overall-status = "" %}
          {% endif %}
          <tr style="visibility: hidden; display: none; background:white" class="test-details edge edge_visible" onclick="toggle_visibility('{{items[1].meshery-component}}');">
            <td style="{{ overall-status }}">{{ items[1].timestamp }}</td>
            <td style="white-space:nowrap;"><a href="{{ site.repo }}-{{ items[1].service-mesh }}">{{ items[1].meshery-component }}</a></td>
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
            <td style="white-space:nowrap;"><img style="height: 2rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/{{items[1].service-mesh | downcase }}.svg" />&nbsp;<a href="{{ site.baseurl }}/service-meshes/adapters/{{ items[1].service-mesh }}">{{ items[1].service-mesh }}</a></td>
            <td>{{ items[1].service-mesh-version }}</td>
          </tr>
          <tr class="hidden-details" id="{{items[1].meshery-component}}" style="visibility:hidden; display:none;">
            <td colspan="2" class="details">
              <i>Platform:</i>
              <li><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/kubernetes-icon-color.svg" />  {{items[1].k8s-distro}}  {{items[1].k8s-version}}</li>
            </td>
            <td colspan="3" class="details">
              <i>Test results:</i>
              <table style="border:0">
              {% for test in item.tests %}
                  <tr><td><img style="height:24px; width: 24px" src="{{result-state}}"></td><td>{{test[0] }}</td></tr>
              {% endfor %}
              </table>
            </td>
            <td>
              <a href = "{{site.baseurl}}/installation/compatibility-matrix/{{item.meshery-component}}-past-results">To see past results click here </a>
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
            {% assign overall-status = "background-color: #56B257; color: white;" %}
            {% assign result-state = "/assets/img/passing.svg" %}
          {% elsif item.overall-status == "partial" %}
            {% assign overall-status = "background-color: #EBC017; color: white;" %}
            {% assign result-state = "/assets/img/YellowCheck.svg" %}
          {% elsif item.overall-status == "failing" %}
            {% assign overall-status = "background-color: #B32700; color: white;" %}
            {% assign result-state = "/assets/img/failing.svg" %}
          {% else %}
            {% assign overall-status = "" %}
          {% endif %}
          <tr style="visibility: hidden; display: none; background:white" class="test-details stable stable_visible" onclick="toggle_visibility('{{item.meshery-component}}-stable');">
            <td style="{{ overall-status }}">{{ item.timestamp }}</td>
            <td style = "white-space:nowrap;"><a href="{{ site.repo }}-{{ item.service-mesh }}">{{ item.meshery-component }}</a></td>
            <td><a href="{{ site.repo }}-{{ item.service-mesh }}/releases/tag/{{ item.meshery-component-version }}">{{ item.meshery-component-version }}</a></td>
            <td><a href="{{ site.repo }}/releases/tag/{{ item.meshery-server-version }}">{{ item.meshery-server-version }}</a></td>
            <td style="white-space:nowrap;"><img style="height: 2rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/{{item.service-mesh | downcase }}.svg" />&nbsp;<a href="{{ site.baseurl }}/service-meshes/adapters/{{ item.service-mesh }}">{{ item.service-mesh }}</a></td>
            <td>{{ item.service-mesh-version }}</td>
          </tr>
          <tr class="hidden-details" id="{{item.meshery-component}}-stable" style="visibility:hidden; display:none;">
            <td colspan="2" class="details">
              <i>Platform:</i>
              <li><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/kubernetes-icon-color.svg" />  {{item.k8s-distro}}  {{item.k8s-version}}</li>
            </td>
            <td colspan="3" class="details">
              <i>Test results:</i>
              <table style="border:0">
              {% for test in item.tests %}
                  <tr><td><img style="height:24px; width: 24px" src="{{result-state}}"></td><td>{{test[0] }}</td></tr>
              {% endfor %}
              </table>
            </td>
            <td>
              <a href = "{{site.baseurl}}/installation/compatibility-matrix/{{item.meshery-component}}-past-results">To see past results click here </a>
            </td>
          </tr>
        {% endif %}
      {% endfor %}
    {% endfor %}

  </tbody>
</table>

<script type="text/javascript" >

    $(function($) {
            console.log("sorting table");
            $("#test-table").tablesorter({
                cssChildRow: "hidden-details",
            });
    });
</script>

