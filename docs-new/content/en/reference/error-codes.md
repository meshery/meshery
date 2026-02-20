---
title: Error Code Reference
description: Meshery Error Code Reference for all Meshery components so that you can troubleshoot issues.
display_title: false
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

{{% alert color="info" title="Error codes are combination of component type, component name, event moniker and numberic code" %}}Error codes are a hyphenated collection of details that include:

<ul>
<li><b>Component Type</b> (string): The type of the component that emits this error event; e.g. <code>adapter</code></li>
<li><b>Component Name</b> (string): The name of the component that emits this error event; e.g. <code>ameshery-istio</code></li>
<li><b>Error Moniker</b> (string): A semi-human readable short key used in descriptive reference to the specific event at-hand; e.g. <code>ErrClosingDatabaseInstanceCode</code></li>
<li><b>Numberic Code</b> (number): Unique number identifying a specific error as scoped by a specific component; e.g. <code>a1000</code></li>
</ul>
The numeric portion of error codes are component-scoped. The numeric portion of error codes are allowed to overlap between Meshery components. The combination of the <code>[component type]-[component name]-[event moniker]-[numeric code]</code> is what makes a given error code globally unique.{{% /alert %}}

### See Also

Troubleshooting guides to using Meshery's various features and components.

{{< troubleshooting-guides-list >}}

## Error Code Categories by Component

{{< error-codes-index >}}
 <a href="#error-code-reference">Top</a>
  <hr>
  <br>

{{< error-codes-detail >}}
