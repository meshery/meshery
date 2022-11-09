---
layout: default
title: Supported Platforms
permalink: installation/platforms
type: installation
display-title: "false"
redirect_from: installation/platforms/
language: en
list: exclude
---

# Supported Platforms<a name="compatibility-matrix"></a>

Meshery deploys as a set of Docker containers, which can be deployed to either a Docker host or Kubernetes cluster. See the complete list of supported platforms in the table below. With service meshes having sprung to life in the context of Kubernetes, so too, can Meshery’s deployment models be characterized in the context of Kubernetes. A given deployment of Meshery can be described as either an _in-cluster_ or an _out-of-cluster_ deployment. Meshery deploys as a stand-alone, management plane on a Docker host (_out-of-cluster_) or as a management plane in a Kubernetes cluster (_in-cluster_).

## Platform Compatibility Matrix



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

{% assign k8s_tests_group = site.compatibility | group_by: "k8s-version" %}

<div>
    <ul class="nav nav-tabs nav-fill mb-3" id="myTab" role="tablist">
        <li class="nav-item" role="presentation">
          <a class="nav-link active" id="kubernetes-tab" data-toggle="tab" href="#kubernetes" role="tab" aria-controls="kubernetes" aria-selected="true">Kubernetes</a>
        </li>
        <li class="nav-item">
          <a class="nav-link" id="docker-tab" data-toggle="tab" href="#docker" role="tab" aria-controls="docker" aria-selected="false">Docker</a>
        </li>
         <li class="nav-item">
          <a class="nav-link" id="mac-tab" data-toggle="tab" href="#mac" role="tab" aria-controls="mac" aria-selected="false">Mac</a>
        </li>
         <li class="nav-item">
          <a class="nav-link" id="windows-tab" data-toggle="tab" href="#windows" role="tab" aria-controls="windows" aria-selected="false">Windows</a>
        </li>
      </ul>
      <div class="tab-content" id="myTabContent">
        <div class="tab-pane fade show active" id="kubernetes" role="tabpanel" aria-labelledby="kubernetes-tab">
            {%include compatibilityMatrix.md k8s_tests_group=k8s_tests_group %}
            <p style="font-size:20px; text-align:center;"><a href="{{ site.baseurl }}/installation/platforms/kubernetes">Install Instruction</a></p>
        </div>
        <div style="width: 1054px; height:477px; background-color:#E6E6E6; text-align:center;font-size:30px;padding:210px; color:#999999;" class="tab-pane fade" id="docker" role="tabpanel" aria-labelledby="docker-tab">Compatibility matrix not yet available<br><a style="font-size:20px" href="{{ site.baseurl }}/installation/platforms/docker">Install Instruction</a></div>
        <div style="width: 1054px; height:477px; background-color:#E6E6E6; text-align:center;font-size:30px;padding:210px; color:#999999;" class="tab-pane fade" id="mac" role="tabpanel" aria-labelledby="mac-tab">Compatibility matrix not yet available<br> <a style="font-size:20px" href="{{ site.baseurl }}/installation/">Install Instruction</a></div>
        <div style="width: 1054px; height:477px; background-color:#E6E6E6; text-align:center;font-size:30px;padding:210px; color:#999999;" class="tab-pane fade" id="windows" role="tabpanel" aria-labelledby="windows-tab">Compatibility matrix not yet available<br> <a style="font-size:20px" href="{{ site.baseurl }}/installation/platforms/windows">Install Instruction</a></div>
      </div>
    </div>

