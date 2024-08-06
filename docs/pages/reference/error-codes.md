---
layout: default
title: Error Code Reference
abstract: "Meshery Error Code Reference"
permalink: reference/error-codes
redirect_from: reference/error-codes/
type: Reference
language: en
abstract: "Meshery Error Code Reference for all Meshery components so that you can troubleshoot issues."
---
<style>

.title {
  text-transform: capitalize;
}
div.error-heading {
  text-transform: uppercase;
}
p.error-details {
    margin-left: 1.5rem;
    font-size: 1rem;
    text-wrap: wrap;
    width:85%
}
td {
  vertical-align: middle;
}
.tbl-head-row{
  background-color:#F2F2F2;
  text-align: left
}
.tbl-head-row .error-name-code{
  /* display:flex; */
  justify-content:space-between;
  align-items:flex-end;
  height:5rem;
}

.tbl .tbl-body .tbl-body-row{
  background-color:#FFFFFF;
}

.tbl .tbl-body .tbl-body-row.hover-effect:hover{
  background-color:#ccfff9;
  cursor:pointer;
}

.tbl-body-row .error-name-code{
  /* display:flex; */
  justify-content:flex-start;
}

.tbl .tbl-body .tbl-hidden-row{
  visibility:hidden;
  display:none;
  background-color:#FAFAFA;
  width:100%
}

</style>

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

## Error Codes and Troubleshooting

Meshery and its components use a common framework (defined within MeshKit) to generate and document an event with a unique error code identifier as the combination of `[component type]-[component name]-[event moniker]-[numeric code]`. Each error code identifies the source component for the error and a standard set of information to describe the error and provide helpful details for troubleshooting the situation surrounding the specific error.

{% include alert.html type="info" title="Error codes are combination of component type, component name, event moniker and numberic code" content="Error codes are a hyphenated collection of details that include:
<ul>
<li><b>Component Type</b> (string): The type of the component that emits this error event; e.g. <code>adapter</code></li>
<li><b>Component Name</b> (string): The name of the component that emits this error event; e.g. <code>ameshery-istio</code></li>
<li><b>Error Moniker</b> (string): A semi-human readable short key used in descriptive reference to the specific event at-hand; e.g. <code>ErrClosingDatabaseInstanceCode</code></li>
<li><b>Numberic Code</b> (number): Unique number identifying a specific error as scoped by a specific component; e.g. <code>a1000</code></li>
</ul>
The numeric portion of error codes are component-scoped. The numeric portion of error codes are allowed to overlap between Meshery components. The combination of the <code>[component type]-[component name]-[event moniker]-[numeric code]</code> is what makes a given error code globally unique." %}

### See Also

Troubleshooting guides to using Meshery's various features and components.

{% assign sorted_guides = site.pages | sort: "type" | reverse %}

<ul>
    {% for item in sorted_guides %}
    {% if item.type=="guides" and item.category=="troubleshooting" and item.list!="exclude"  -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
      {% endif %}
    {% endfor %}
</ul>

## Error Code Categories by Component

<table style="margin:auto;padding-right:25%; padding-left:20%;">
<thead>
  <tr>
    <th style="text-align:left">Component Type</th>
    <th style="text-align:left">Component Name</th>
  </tr>
</thead>
<tbody>
  {% for files in site.data.errorref %}
    {% for eachFile in files %}
      {% for component in eachFile  %}
          {% comment %} <tr><td colspan="2">{{component}}</td></tr> {% endcomment %}
           {% capture thecycle %}{% cycle 'odd', 'even' %}{% endcapture %}
            {% if thecycle == 'even' %} 
            {% if component[1].component_type == 'adapter' %}
              {% capture link %}meshery-adapter-for-{{component[1].component_name | lowercase}}{% endcapture %}
            {% elsif component[1].component_type == 'component' %}
               {% capture link %}meshery-server{% endcapture %}
            {% else %}
              {% capture link %}{{ component[1].component_name  | camelcase }}-{{ component[1].component_type }}{% endcapture %}      
            {% endif %}
            <tr>
              <td style="text-align:left">{{ component[1].component_type }}</td>
              <td class="title"><a href="#{{ link}}">{{ component[1].component_name }}</a></td>
            </tr>
        {% endif %}
      {% endfor %}
    {% endfor %}
  {% endfor %}
</tbody>
</table>
 <a href="#error-code-reference">Top</a>
  <hr>
  <br>

  {% for files in site.data.errorref %}    
  {% for eachFile in files %}
    {% for component in eachFile %}
      {% capture thecycle %}{% cycle 'odd', 'even' %}{% endcapture %}
      {% if thecycle == 'even' %}
        {% if component[1].component_type == 'adapter' %}
          {% capture heading %}
            Meshery Adapter for {{ component[1].component_name }}
          {% endcapture %}
        {% elsif component[1].component_type == 'client' %}
          {% capture heading %}
            {{ component[1].component_name }} client
          {% endcapture %}
        {% elsif component[1].component_type == 'library' %}
          {% capture heading %}
            {{ component[1].component_name }} {{ component[1].component_type | camelcase }}
          {% endcapture %}
        {% elsif component[1].component_name == 'meshery-server' %}
          {% capture heading %}
            Meshery Server
          {% endcapture %}
        {% endif %}


<h2 class="title">{{ heading }}</h2>
<table class="tbl">
  <thead>
    <tr class="tbl-head-row">
      <th style="width:15%">Severity</th>
      <th class="error-name-code"><span>Error Name - Code</span></th>
      <th style="width:85%">Short Description</th>
      <th>Discussion</th>
    </tr>
  </thead>
  <tbody class="tbl-body">
    {% for err_code in component[1].errors %}
      {% if err_code[1]["severity"] == "Fatal" %}
        {% assign severity = "background-color: #FF0101; color: white;" %}
      {% elsif err_code[1]["severity"] == "Alert" %}
        {% assign severity = "background-color: #FEA400; color: white;" %}
      {% else %}
        {% assign severity = "background-color: transparent; color: black;" %}
      {% endif %}
      <tr class="tbl-body-row hover-effect" onclick="toggle_visibility('{{ component[1].component_name }}-{{ err_code[1]["name"] }}-more-info');">
        <td style="{{ severity }}">{{ err_code[1]["severity"] }}</td>
        <td id="{{ heading | slugify }}-{{err_code[1]["code"] }}" class="error-name-code">
          <code>{{ err_code[1]["name"] | xml_escape }}-{{ err_code[1]["code"] }}</code>
        </td>
        <td>{{ err_code[1]["short_description"] | xml_escape }}</td>
        <td><a href="https://discuss.layer5.io/search?q={{ err_code[1]['name'] | xml_escape }}-{{ err_code[1]['code'] }}" target="_blank">search forum</a></td>
      </tr>
      <tr id="{{ component[1].component_name }}-{{ err_code[1]["name"] }}-more-info" class="tbl-hidden-row">
        <td style="word-break:break-all;" colspan="3">
          <div class="error-heading">Long Description</div>
          <p class="error-details">{{ err_code[1]["long_description"] | xml_escape }}</p>
          <div class="error-heading">Probable Cause</div>
          <p class="error-details">{{ err_code[1]["probable_cause"] | xml_escape }}</p>
          <div class="error-heading">Suggested Remediation</div>
          <p class="error-details">{{ err_code[1]["suggested_remediation"] | xml_escape }}</p>
        </td>
      </tr>
    {% endfor %}
  </tbody>
</table>
<a href="#error-code-reference">Top</a>
<hr>
<br>
{% endif %}
{% endfor %}
{% endfor %}
{% endfor %}


    
