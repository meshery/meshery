---
layout: default
title: Meshery Designs 
description: Get to know what is Design & how its helps in configuring components
permalink: guides/meshery-design
type: Guides
language: en
---


We know in the world of microservice architecture, service meshes are crucial, and they come with a price (complexity, error-prone, steep learning curve).
<br>
We are here to help you out with Designs.

## WHAT IS IT ?
Here, we have pre-built configuration YAML files that serve as a blueprint for the configuration of various service mesh components. So that you don't need to bother with writing configuration files for your mesh & make things aggravated. <br>
when you apply the configuration files on a Kubernetes cluster. The service mesh's control-plane reads the instructions (from the configuration file) & takes the necessary steps to setup the desired service mesh environment. <br>
Designs help you to use your service mesh's potential without any prior experience of using service meshes. Like 

* Configuring the control plane
* Defining traffic routing rules
* Implementing security policies
* Configuring observability


Along all these, we also provide an opportuinity to create or import your own configuration files.


## WHERE TO FIND IT
If your Meshery setup is done then, you can start meshery (eg : `$ mesheryctl system start`) & open your browser & search for http://< hostname >:9081/provider or http://localhost:9081/provider

* Select the Provider.

* Go to the "Configuration".

* Click on "Designs".

**For users who installed meshery extension on docker**

0. Login to Meshery.

1. Launch Meshery.

2. Select the Provider.

3. Go to the "Configuation".

4. Click on "Designs".

<img src="{{site.baseurl}}/assets/img/meshery-design/location.png" />


## IMPORT YOUR DESIGNS

1.  Click on "Import Design"

2.  Provide the URL :<br>
where you had uploaded your Configuration file (github/gitlab).

3.  Provide the file name


<img src="{{site.baseurl}}/assets/img/meshery-design/import.png" />

## CREATE YOUR DESIGN 

1. Click on "Create Design".

2. Select the Category.

3. Select the Model.

4. Write your Configuration file.


<img src="{{site.baseurl}}/assets/img/meshery-design/create.png" />