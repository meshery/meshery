---
title: "Compatibility Matrix"
description: "Meshery Server and Adapter compatibility with infrastructure platforms"
weight: 55
aliases:
  - /installation/compatibility-matrix/
display_title: "false"
---

<a name="compatibility-matrix"></a>

# Compatibility Matrix

Meshery Server and Meshery Adapters are tested daily for their compatibility with the infrastructure they manage and the platforms Meshery deploys on (Kubernetes and Docker). Integration test results are automatically posted to the following compatibility matrix.

{{< tabpane text=true >}}
  {{< tab header="Kubernetes" lang="en" >}}
    <h3>Tested Kubernetes Versions</h3>
    <p>Meshery is compatible with all versions of Kubernetes. The following table shows the tested versions:</p>
    <table class="table table-striped table-bordered">
      <tr style="text-align:center">
        <th>Kubernetes Version</th>
        <th>Status</th>
      </tr>
      <tr>
        <td>1.29.x</td>
        <td>✅ Tested</td>
      </tr>
      <tr>
        <td>1.28.x</td>
        <td>✅ Tested</td>
      </tr>
      <tr>
        <td>1.27.x</td>
        <td>✅ Tested</td>
      </tr>
      <tr>
        <td>1.26.x</td>
        <td>✅ Tested</td>
      </tr>
      <tr>
        <td>1.25.x</td>
        <td>✅ Tested</td>
      </tr>
    </table>
  {{< /tab >}}
  {{< tab header="Docker" lang="en" >}}
    <h3>The following minimum Docker build versions are required:</h3>
    <table class="table table-striped table-bordered">
      <tr style="text-align:center">
        <th>Name</th>
        <th>Version</th> 
      </tr>
      <tr>
        <td><a href="/installation/docker">Docker Engine</a></td>
        <td><b>19.x</b> and above</td>
      </tr>
      <tr>
        <td><a href="/installation/docker/docker-extension">Docker Desktop<br></a><span style="color:#999999;"><em>Used through <b>Docker Extension</b></em></span></td>
        <td><b>2.0.x</b> and above</td>
      </tr>
    </table>
  {{< /tab >}}
{{< /tabpane >}}

{{< alert type="info" title="Overview of the Integration Tests" >}}
For a complete overview of the latest integration tests and their status please visit Meshery <a href='/installation/compatibility-matrix#integration-tests'>Integration Tests</a>
{{< /alert >}}

{{< related-discussions tag="meshery" >}}
