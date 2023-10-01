---
layout: default
title: Error Code Reference
abstract: "Meshery Error Code Reference"
permalink: reference/error-codes
redirect_from: reference/error-codes/
type: Reference
language: en
---
<style>

.title {
  text-transform: capitalize;
}

.tbl-head-row{
  background-color:#F2F2F2;
}

.tbl-head-row .error-name-code{
  display:flex;
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
  display:flex;
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

Meshery and its components use a common framework (defined within MeshKit) to generate and document an error with a unique error code identifier: the combination of Meshery component moniker and numberic code - `[component-moniker]-[numeric code]`. Each error code identifies the source component for the error and a standard set of information to describe the error and provide helpful details for troubleshooting the situation surrounding the specific error.

{% include alert.html type="info" title="Error codes are combination of component moniker and numberic code" content="
Note: The numeric portion of error codes are component-scoped. The numeric portion of error codes are allowed to overlap between Meshery components. The combination of the `[component-moniker]-[numeric code]` is what makes a given error code globally unique." %}

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
              {% capture link %}meshery-adapter-for-{{component[1].component_name}}{% endcapture %}
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
      <th class="error-name-code"><span>Error Name - Code</span></th>
      <th style="width:15%">Severity</th>
      <th style="width:85%">Short Description</th>
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
        <td class="error-name-code">
          <code>{{ err_code[1]["name"] | xml_escape }}-{{ err_code[1]["code"] }}</code>
        </td>
        <td style="{{ severity }}">{{ err_code[1]["severity"] }}</td>
        <td>{{ err_code[1]["short_description"] | xml_escape }}</td>
      </tr>
      <tr id="{{ component[1].component_name }}-{{ err_code[1]["name"] }}-more-info" class="tbl-hidden-row">
        <td style="word-break:break-all;">
          <div><i><b>Probable Cause:</b></i></div>{{ err_code[1]["probable_cause"] | xml_escape }}
          </td>
          <td>
          <div><i><b>Suggested Remediation:</b></i></div>
          {{ err_code[1]["suggested_remediation"] | xml_escape }}
          </td>
          <td>
          <div><i><b>Long Description:</b></i></div>
          {{ err_code[1]["long_description"] | xml_escape }}
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


    
