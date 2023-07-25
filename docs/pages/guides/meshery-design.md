---
layout: default
title: Meshery Designs 
description: Get to know what is Design & how its helps in configuring components
permalink: guides/meshery-design
type: Guides
language: en
---



We know in the world of microservice architecture, **service meshes** are crucial. And they comes with a **price** (complexity, error-prone, steep learning curve).

### We are here to help you out, with **Designs**.

## WHAT IS IT ?
**Here**, we have pre-build configuration yaml files that serve as a blueprint for the configuration of various service-mesh components. So that you dont need to bother for writing configuration files for your mesh & make things aggravated. <br>
when you apply the configuration files on kubernetes cluster. The service-mesh control plane reads the instructions (from the configuration file) & takes the necessary steps to setup the desired service-mesh environment. <br><br>
**Designs** help you to use your service-mesh's potential without any prior experience of using service-meshes. Like 

* Configuring control plane
* Defining traffic routing rules
* Implementing security policies
* Configuring observability


Along all these we also provide an opportuinity to create or import you own configuration files.


## WHERE TO FIND IT
_for those who has meshery-extention installed on docker_

**0. Login to meshery** 

**1. Launch meshery**

**2. Go to configuation**

**3. Click on designs**

**4. Get set go**

<img src="{{site.baseurl}}/assets/img/meshery-design/design-location.png" />

## IMPORT YOUR DESIGNS

**1.  Click on "import design"**

**2.  Provide the URL** :<br>
where you upload your config.yaml file (github/cloud storage).

**3.  Provide the file name**


<img src="{{site.baseurl}}/assets/img/meshery-design/import-design.png" />

## CREATE YOUR DESIGN 

**1. Click on "create design"**

**2. Select category**

**3. Select model**

**4. get set go**


<img src="{{site.baseurl}}/assets/img/meshery-design/create-design.png" />