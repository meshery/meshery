---
layout: default
title: Error Code Reference
abstract: Meshery Error Code Reference for all Meshery components to help troubleshoot issues.
permalink: reference/error-codes
redirect_from: reference/error-codes/
type: Reference
language: en
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
    width: 85%;
  }
  td {
    vertical-align: middle;
  }
  .tbl-head-row {
    background-color: #F2F2F2;
    text-align: left;
  }
  .tbl-head-row .error-name-code {
    justify-content: space-between;
    align-items: flex-end;
    height: 5rem;
  }
  .tbl .tbl-body .tbl-body-row {
    background-color: #FFFFFF;
  }
  .tbl .tbl-body .tbl-body-row.hover-effect:hover {
    background-color: #ccfff9;
    cursor: pointer;
  }
  .tbl-body-row .error-name-code {
    justify-content: flex-start;
  }
  .tbl .tbl-body .tbl-hidden-row {
    visibility: hidden;
    display: none;
    background-color: #FAFAFA;
    width: 100%;
  }
  table.tbl, table {
    width: 100%;
    border-collapse: collapse;
    display: block;
    overflow-x: auto;
    white-space: nowrap;
  }
  th, td {
    padding: 0.75rem;
    text-align: left;
    white-space: normal;
    word-break: break-word;
  }
  @media (max-width: 768px) {
    .tbl-head-row, .tbl-body-row {
      font-size: 0.9rem;
    }
    td, th {
      padding: 0.5rem;
    }
    /* Consider removing or adjusting these lines to enable toggling on mobile
    .tbl .tbl-body .tbl-hidden-row {
      display: block;
      visibility: visible;
    }
    */
  }
</style>

<script type="text/javascript">
  function toggle_visibility(id) {
    var e = document.getElementById(id);
    if (e.style.display === 'table-row') {
      e.style.display = 'none';
      e.style.visibility = 'hidden';
    } else {
      e.style.display = 'table-row';
      e.style.visibility = 'visible';
    }
  }
</script>

# Error Codes and Troubleshooting

Meshery and its components use a common framework (defined within MeshKit) to generate and document an event with a unique error code identifier as the combination of `[component type]-[component name]-[event moniker]-[numeric code]`. Each error code identifies the source component for the error and provides a standard set of information to describe the error and offer helpful details for troubleshooting.

{% include alert.html type="info" title="Error Code Structure" content="Error codes are a hyphenated collection of details that include:

- **Component Type** (string): The type of the component that emits this error event; e.g., adapter.
- **Component Name** (string): The name of the component that emits this error event; e.g., meshery-istio.
- **Error Moniker** (string): A semi-human-readable short key used to describe the specific event; e.g., ErrClosingDatabaseInstanceCode.
- **Numeric Code** (number): A unique number identifying a specific error within a component; e.g., a1000.

The numeric portion of error codes is component-scoped and may overlap between components. The combination of `[component type]-[component name]-[event moniker]-[numeric code]` ensures a globally unique error code." %}

## See Also
Troubleshooting guides for using Meshery's various features and components:

{% assign sorted_guides = site.pages | sort: "type" | reverse %}
{% for item in sorted_guides %}
  {% if item.type == "guides" and item.category == "troubleshooting" and item.list != "exclude" %}
  - [{{ item.title }}]({{ item.url }})
  {% endif %}
{% endfor %}

## Error Code Categories by Component

{% if site.data.errorref %}
{% for files in site.data.errorref %}
  {% for eachFile in files %}
    {% for component in eachFile %}
      {% if component[1].component_type == 'adapter' %}
        {% capture heading %}Meshery Adapter for {{ component[1].component_name | capitalize }}{% endcapture %}
      {% elsif component[1].component_type == 'client' %}
        {% capture heading %}{{ component[1].component_name | capitalize }} Client{% endcapture %}
      {% elsif component[1].component_type == 'library' %}
        {% capture heading %}{{ component[1].component_name | capitalize }} {{ component[1].component_type | capitalize }}{% endcapture %}
      {% elsif component[1].component_name == 'meshery-server' %}
        {% capture heading %}Meshery Server{% endcapture %}
      {% else %}
        {% capture heading %}{{ component[1].component_name | capitalize }}{% endcapture %}
      {% endif %}

      <h2 class="title">{{ heading }}</h2>
      <table class="tbl">
        <thead>
          <tr class="tbl-head-row">
            <th style="width:5%">Severity</th>
            <th class="error-name-code">Error Name - Code</th>
            <th style="width:85%">Short Description</th>
            <th>Discussion</th>
          </tr>
        </thead>
        <tbody class="tbl-body">
          {% for err_code in component[1].errors %}
            {% if err_code[1].severity == "Fatal" %}
              {% assign severity = "background-color: #FF0101; color: white; writing-mode: vertical-rl; text-orientation: mixed;" %}
            {% elsif err_code[1].severity == "Alert" %}
              {% assign severity = "background-color: #FEA400; color: white; writing-mode: vertical-rl; text-orientation: mixed;" %}
            {% else %}
              {% assign severity = "background-color: transparent; color: black; writing-mode: vertical-rl; text-orientation: mixed;" %}
            {% endif %}
            <tr class="tbl-body-row hover-effect" onclick="toggle_visibility('{{ component[1].component_name }}-{{ err_code[1].name }}-more-info');">
              <td style="{{ severity }}">{{ err_code[1].severity }}</td>
              <td id="{{ heading | slugify }}-{{ err_code[1].code }}" class="error-name-code">
                <code>{{ err_code[1].name | xml_escape }}-{{ err_code[1].code }}</code>
              </td>
              <td>{{ err_code[1].short_description | xml_escape }}</td>
              <td><a href="https://meshery.io/community#discussion-forums/search?q={{ err_code[1].name | xml_escape }}-{{ err_code[1].code }}" target="_blank">Search Forum</a></td>
            </tr>
            <tr id="{{ component[1].component_name }}-{{ err_code[1].name }}-more-info" class="tbl-hidden-row">
              <td colspan="4">
                <div class="error-heading">Long Description</div>
                <p class="error-details">{{ err_code[1].long_description | xml_escape }}</p>
                <div class="error-heading">Probable Cause</div>
                <p class="error-details">{{ err_code[1].probable_cause | xml_escape }}</p>
                <div class="error-heading">Suggested Remediation</div>
                <p class="error-details">{{ err_code[1].suggested_remediation | xml_escape }}</p>
              </td>
            </tr>
          {% endfor %}
        </tbody>
      </table>
      <a href="#error-codes-and-troubleshooting">Back to Top</a>
      <hr><br>
    {% endfor %}
  {% endfor %}
{% endfor %}
{% else %}
  <p>No error codes are currently available. Please check the data source or contact the Meshery team for assistance.</p>
{% endif %}
