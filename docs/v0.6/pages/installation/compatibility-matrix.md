<!-- ---
layout: default
title: Compatibility Matrix
# permalink: installation/compatibility-matrix
type: installation
display-title: "false"
# redirect_from: installation/compatibility-matrix/
language: en
list: exclude
--- -->

<a name="compatibility-matrix"></a>

# Compatibility Matrix

Meshery Server and Meshery Adapters are tested daily for their compatibility with the infrastructure they manage and the platforms Meshery deploys on (Kubernetes and Docker). Integration test results are automatically posted to the following compatibility matrix.

<style>
.nav-link.active{
    border-bottom: #00B39F solid 5px!important;
    border-top: none !important;
     color:#00B39F!important

}
.nav-link:hover, .nav-link:focus{
    color:#00B39F!important
}
.nav-link{
    border: 0 !important;
    color: black
}
</style>

{% assign k8s_tests_group = site.compatibility | group_by: "k8s-version" | sort: "name" | reverse %}

<div>
    <ul class="nav nav-tabs nav-fill mb-3" id="myTab" role="tablist">
        <li class="nav-item" role="presentation">
          <a class="nav-link active" id="kubernetes-tab" data-toggle="tab" href="/v0.6/v0.6#kubernetes" role="tab" aria-controls="kubernetes" aria-selected="true">Kubernetes</a>
        </li>
        <li class="nav-item">
          <a class="nav-link" id="docker-tab" data-toggle="tab" href="/v0.6/v0.6#docker" role="tab" aria-controls="docker" aria-selected="false">Docker</a>
        </li>
         <!-- <li class="nav-item">
          <a class="nav-link" id="mac-tab" data-toggle="tab" href="/v0.6/v0.6#mac" role="tab" aria-controls="mac" aria-selected="false">Mac</a>
        </li>
         <li class="nav-item">
          <a class="nav-link" id="windows-tab" data-toggle="tab" href="/v0.6/v0.6#windows" role="tab" aria-controls="windows" aria-selected="false">Windows</a>
        </li> -->
      </ul>
      <div class="tab-content" id="myTabContent">
        <div class="tab-pane fade show active" id="kubernetes" role="tabpanel" aria-labelledby="kubernetes-tab">
            {%include compatibilityMatrix.md k8s_tests_group=k8s_tests_group %}
        </div>
        <div style="text-align:center" class="tab-pane fade" id="docker" role="tabpanel" aria-labelledby="docker-tab">
        <h3 style="text-align:left;">The following minimum Docker build versions are required:</h3>
        <table class="table table-striped table-bordered">
            <tr style="text-align:center">
              <th>Name</th>
              <th>Version</th> 
            </tr>
            <tr>
              <td><a href="/v0.6/v0.6{{site.baseurl}}/installation/docker">Docker Engine</a></td>
              <td><b>19.x</b> and above</td>
            </tr>
            <tr>
              <td><a href="/v0.6/v0.6{{site.baseurl}}/installation/docker-extension">Docker Desktop<br></a><span style="color:#999999; text-decororation:none;"><em>Used through <b>Docker Extension</b></em></span></td>
              <td><b>2.0.x</b> and above </td>
            </tr>
          </table>
        </div>
      </div>
        <!-- <div style="background-color:#E6E6E6; text-align:center;font-size:30px;padding:210px; color:#999999;" class="tab-pane fade" id="mac" role="tabpanel" aria-labelledby="mac-tab">Compatibility matrix not yet available<br> <a style="font-size:20px" href="/v0.6/v0.6{{ site.baseurl }}/installation/">Install Instruction</a></div>
        <div style="text-align:center;" class="tab-pane fade" id="windows" role="tabpanel" aria-labelledby="windows-tab">
          <h3 style="text-align:left;">The following minimum Windows build versions are required:</h3>
          <table class="table table-striped table-bordered">
            <tr style="text-align:center">
              <th>Name</th>
              <th>Version</th> 
            </tr>
            <tr>
              <td><a href="/v0.6/v0.6{{site.baseurl}}/installation/windows#wsl1">WSL1</a></td>
              <td><b>x64</b> - Windows 7 </td>
            </tr>
            <tr>
              <td><a href="/v0.6/v0.6{{site.baseurl}}/installation/windows#wsl1">WSL2</a></td>
              <td><b>x64</b> - Version 1903, Build 18362; <b>ARM 64</b> - Version 2004, Build 19041</td>
            </tr>
            <tr>
              <td><a href="/v0.6/v0.6https://docs.microsoft.com/en-us/windows/wsl/release-notes#build-18945">Custom Kernel</a></td>
              <td>Build 18945</td>
            </tr>
            <tr>
              <td><a href="/v0.6/v0.6https://docs.microsoft.com/en-us/windows/wsl/release-notes#build-19013">Kernel with K8s required modules</a></td>
              <td>Build 19013</td>
            </tr>
          </table>
        <br> <a style="font-size:20px" href="/v0.6/v0.6{{ site.baseurl }}/installation/windows">Install Instruction</a></div> -->

<div style="z-index:0">
  {% include alert.html type="info" title="<span style='margin:0;'>Overview of the Integration Tests</span>" content="For a complete overview of the latest integration tests and their status please visit Meshery <a href='/installation/compatibility-matrix#integration-tests'>Integration Tests</a>" %}
</div>
