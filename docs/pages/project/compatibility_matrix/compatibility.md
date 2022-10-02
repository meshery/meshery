---
layout: default
title: Compatibility Matrix
permalink: project/compatibility-matrix
redirect_from: project/compatibility-matrix/
description: a complete compatibility matrix and project test status dashboard.
language: en
display-title: "false"
list: exclude
type: "project"
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
            console.log(e[i].classList)
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
            console.log(typeof(e[i].classList["1"]))
            if(e[i].classList.contains("stable_visible")){
                e[i].classList.remove("stable_visible")
                if(!edge_box.checked){
                    edge_box.checked=true;
                    handleEdgeCheckboxChange();
                }
        }
        else{
          console.log("stable")
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
  .yellowCheckbox{
    width:1.5rem
  }
  .tooltipss{
    position:relative;
    width:fit-content;
    cursor:pointer;
  }
  .tooltipss .tooltiptext {
  visibility: hidden;
  width: 120px;
  background-color: #555;
  color: #fff;
  text-align: center;
  border-radius: 6px;
  padding: 5px 0;
  position: absolute;
  z-index: 1;
  bottom: 125%;
  left: 50%;
  margin-left: -60px;
  opacity: 0;
  transition: opacity 0.3s;
}

.tooltipss .tooltiptext::after {
  content: "";
  position: absolute;
  top: 100%;
  left: 50%;
  margin-left: -5px;
  border-width: 5px;
  border-style: solid;
  border-color: #555 transparent transparent transparent;
}

.tooltipss:hover .tooltiptext {
  visibility: visible;
  opacity: 1;
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
{% assign k8s_tests_group = site.compatibility | group_by: "k8s-version" %}
{% assign service_meshes = site.adapters  %}

 <div>
 </div>

# Compatibility Matrix

Compatibility of Meshery with other integrated systems.

<table>
  <th>Kubernetes Version</th>

  <th><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/istio.svg" /><a href="{{ site.repo }}-istio">meshery-istio</a></th>
  <th><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/linkerd.svg" /><a href="{{ site.repo }}-linkerd">meshery-linkerd</a></th>
  <th><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/kuma.svg" /><a href="{{ site.repo }}-kuma">meshery-kuma</a></th>
  <th><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/osm.svg" /><a href="{{ site.repo }}-osm">meshery-osm</a></th>
  <th><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/nginx-sm.svg" /><a href="{{ site.repo }}-nginx-sm">meshery-nginx-sm</a></th>
  <th><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/traefik-mesh.svg" /><a href="{{ site.repo }}-traefik-mesh">meshery-traefik-mesh</a></th>
  <th><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/cilium.svg" /><a href="{{ site.repo }}-cilium">meshery-cilium</a></th>

{% for k8s in k8s_tests_group%}

<tr class = "first-row">
{% assign successfull_istio = 0 %}
{% assign successfull_linkerd = 0 %}
{% assign successfull_cilium = 0 %}
{% assign successfull_osm = 0 %}
{% assign successfull_kuma = 0 %}
{% assign successfull_traefik_mesh = 0 %}
{% assign successfull_nginx_sm = 0 %}
<td>{{k8s.name}}</td>
{% assign k8s_items = k8s.items | group_by: "meshery-component"  %}
{% for k8s_item in k8s_items %}
{% if k8s_item.name == "meshery-linkerd" %}
{% assign linkerd_size = k8s_item.size | times:1.0 %}
{% for single in k8s_item.items %}
{% if single.overall-status == "passing" %}
{% assign successfull_linkerd = successfull_linkerd | plus:1 %}
{% endif %}
{% endfor %}
{% elsif k8s_item.name == "meshery-istio" %}
{% assign istio_size = k8s_item.size | times:1.0 | times:1.0 %}
{% for single in k8s_item.items %}
{% if single.overall-status == "passing" %}
{% assign successfull_istio = successfull_istio | plus:1 %}
{% endif %}
{% endfor %}
{% elsif k8s_item.name == "meshery-kuma" %}
{% assign kuma_size = k8s_item.size | times:1.0 %}
{% for single in k8s_item.items %}
{% if single.overall-status == "passing" %}
{% assign successfull_kuma = successfull_kuma | plus:1 %}
{% endif %}
{% endfor %}
{% elsif k8s_item.name == "meshery-osm" %}
{% assign osm_size = k8s_item.size | times:1.0 %}
{% for single in k8s_item.items %}
{% if single.overall-status == "passing" %}
{% assign successfull_osm = successfull_osm | plus:1 %}
{% endif %}
{% endfor %}
{% elsif k8s_item.name == "meshery-cilium" %}
{% assign cilium_size = k8s_item.size | times:1.0 %}
{% for single in k8s_item.items %}
{% if single.overall-status == "passing" %}
{% assign successfull_cilium = successfull_cilium | plus:1 %}
{% endif %}
{% endfor %}
{% elsif k8s_item.name == "meshery-nginx-sm"%}
{% assign nginx_size = k8s_item.size | times:1.0 %}
{% for single in k8s_item.items %}
{% if single.overall-status == "passing" %}
{% assign successfull_nginx_sm= successfull_nginx_sm| plus:1 %}
{% endif %}
{% endfor %}
{% elsif k8s_item.name == "meshery-traefik-mesh" %}
{% assign traefik_size = k8s_item.size | times:1.0 %}
{% for single in k8s_item.items %}
{% if single.overall-status == "passing" %}
{% assign successfull_traefik_mesh = successfull_traefik_mesh | plus:1 %}
{% endif %}
{% endfor %}
{% endif %}
{% endfor %}
{% assign istio_percentage = successfull_istio | divided_by:istio_size | times:100 | round:2 %}
<td onclick = "clickIcon(`meshery-istio`)" class = "compatibility">{{istio_percentage}}
</td>
{% assign linkerd_percentage = successfull_linkerd | divided_by:linkerd_size | times:100 | round:2 %}
<td onclick = "clickIcon(`meshery-linkerd`)" class = "compatibility">{{linkerd_percentage}}</td>
{% assign kuma_percentage = successfull_kuma | divided_by:kuma_size | times:100 | round:2 %}
<td onclick = "clickIcon(`meshery-kuma`)" class = "compatibility">{{kuma_percentage}}</td>
{% assign osm_percentage = successfull_osm | divided_by:osm_size | times:100 | round:2 %}
<td onclick = "clickIcon(`meshery-osm`)" class = "compatibility">{{osm_percentage}}%</td>
{% assign nginx_percentage = successfull_nginx_sm | divided_by:nginx_size | times:100 | round:2 %}
<td onclick = "clickIcon(`meshery-nginx-sm`)" class = "compatibility">{{nginx_percentage}}% </td>
{% assign traefik_percentage = successfull_traefik_mesh | divided_by:traefik_size | times:100 | round:2 %}
<td onclick = "clickIcon(`meshery-traefik-mesh`)" class = "compatibility">{{traefik_percentage}}%</td>
{% assign cilium_percentage = successfull_cilium | divided_by:cilium_size | times:100 | round:2 %}
<td onclick = "clickIcon(`meshery-cilium`)" class = "compatibility">{{cilium_percentage}}%</td>
</tr>
{% endfor %}


</table>

<script>
  function showCompatability () {
      let percentContainer = document.querySelectorAll(".compatibility")
      console.log(percentContainer);
      for(let i = 0 ; i<percentContainer.length;i++){
        console.log(parseFloat(percentContainer[i].innerHTML));
        let percentage = parseFloat(percentContainer[i].innerHTML);
        if (percentage >= 90.00){
          percentContainer[i].innerHTML = `
            <div class = "tooltipss">
              <img src = "{{site.baseurl}}/assets/img/passing.svg" class = "yellowCheckbox" >
              <span class = "tooltiptext">${percentage}%</span>
            </div>
          `
        }
        else if(percentage >=1 && percentage<=89.99){
          percentContainer[i].innerHTML = `<div class = "tooltipss">
              <img src = "{{site.baseurl}}/assets/img/YellowCheck.svg" class = "yellowCheckbox" >
              <span class = "tooltiptext">${percentage}%</span>
            </div>`
        }
        else{
           percentContainer[i].innerHTML = `<div class = "tooltipss">
              <img src = "{{site.baseurl}}/assets/img/failing.svg" class = "yellowCheckbox" >
              <span class = "tooltiptext">${percentage}%</span>
            </div>`
        }
      }
    }
  function clickIcon(serviceMesh){
    console.log("clicked",serviceMesh);
    location.href = `{{site.baseurl}}/project/compatibility-matrix/${serviceMesh}-past-results`
  }

showCompatability()
</script>

## Integration Tests

As a key aspect of Meshery, its integrations with other systems are routinely tested. Unit, integration testing occurs before and after every pull request (before code is to be merged into the project and after code is merged into the project). Regression tests are run nightly.

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

<table id = "test-table" class="tablesorter">
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
            <td style="white-space: nowrap;"><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/{{item.service-mesh | downcase }}.svg" />&nbsp;<a href="{{ site.baseurl }}/service-meshes/adapters/{{ item.service-mesh }}">{{ item.service-mesh }}</a></td>
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
              <a href = "{{site.baseurl}}/project/compatibility-matrix/{{item.meshery-component}}-past-results">To see past results click here </a>
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
            <td style="white-space:nowrap;"><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/{{items[1].service-mesh | downcase }}.svg" />&nbsp;<a href="{{ site.baseurl }}/service-meshes/adapters/{{ items[1].service-mesh }}">{{ items[1].service-mesh }}</a></td>
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
            <td style="white-space:nowrap;"><img style="height: 1rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/{{item.service-mesh | downcase }}.svg" />&nbsp;<a href="{{ site.baseurl }}/service-meshes/adapters/{{ item.service-mesh }}">{{ item.service-mesh }}</a></td>
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
              <a href = "{{site.baseurl}}/project/compatibility-matrix/{{item.meshery-component}}-past-results">To see past results click here </a>
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


