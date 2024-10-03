---
layout: default
title: Generating Models
abstract: Generating Model 
permalink: guides/configuration-management/generating-models
category: configuration
type: guides
language: en
---

Generating Models annotation and existing custom resource definition (CRD) into Meshery. The platform supports generation from URL, csv you can generate models using either the Meshery CLI or the Meshery UI.

## Generate Models Using Meshery CLI

**Step 1: Install Meshery CLI**

Before you can use the Meshery CLI to generate a [Model](/concepts/logical/models), you must first install it. You can install Meshery CLI by [following the instructions]({{site.baseurl}}/installation#install-mesheryctl).


**Step 2: Generate the Model**

Model can generated in 2 different format ```URL, CSV```. A template file is required that contains some required properties: Registrant, Model name, Model DisplayName, Category. The template file is only required when you use URL for generation of model from a crd. This command enable users to import their new models from CRD and existing Meshery Model

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl model generate -f [file/url] </div></div>
</pre>

The supported registrant for generating from URL is `github` and `artifacthub`.The URL format must be in this order.


Registrant `Artifacthub`:
- https://artifacthub.io/packages/search?ts_query_web={ model-name } 
- https://istio-release.storage.googleapis.com/charts/base-1.19.0-alpha.1.tgz&sa=D&source=editors&ust=1726839249773905&usg=AOvVaw0j88gkt6FOS1LLSRCYq95X 

Registrant `Github`:
- git:://github.com/cert-manager/cert-manager/master/deploy/crds
- https://github.com/UffizziCloud/uffizzi_controller/releases/download/uffizzi-controller-2.0.1/uffizzi-controller-2.0.1.tgz&sa=D&source=editors&ust=1726839320133140&usg=AOvVaw2AryFwXIPKFnWRjRRTApzp


**Example :**

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl model generate -f "git:://github.com/cert-manager/cert-manager/master/deploy/crds" -t template.json</div></div>
</pre>


**Note:** A `-r` flag is present to skip registaration. If the flag is used then no new model would be registered though they would be generated and stored inside `.meshery/models` directory.
The `template.json` can be viewed [here](#).In template.json if the field `isAnnotaion` is true then we would only consider that component as an annotation (svg icon) rather than a normal component.



**Example :**


<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl model generate -f "git:://github.com/cert-manager/cert-manager/master/deploy/crds" -t template.json -r</div></div>
</pre>



## Generate Models Using Meshery UI

**Step 1: Access the Meshery UI**

To generate a model into Meshery using the Meshery UI, you must first [install Meshery](/installation/quick-start)

**Step 2: Navigate to Registry under Settings Page**

Once you have accessed the Meshery UI, navigate to the Registry under Settings. This page can be accessed by clicking on the Settings on the top right on setting icon and then selecting "Registry" and then choose model.

<a href="{{ site.baseurl }}/assets/img/export/Registry.png"><img alt="Registry-Navigator" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/export/Registry.png" /></a>

### Step 3: Generate the Model

On the **Registry** page, you can generate your model by clicking the **Generate** button. You have two options for input: URL or CSV. Below are the steps for both methods:

---

#### **From CSV**
1. **Upload CSV Files**: You can generate your model by uploading a **components model CSV**. Uploading a **relationships CSV** is optional.
   
2. **Template CSV**: If you don’t have a CSV file ready, you can use our [Spreasheet template](https://docs.google.com/spreadsheets/d/19JEpqvHrG8UL-Bc-An9UIcubf1NVhlfnQSN1TD7JOZ4/) to create one. Simply  fill in your details, download the template as csv, and upload it to generate your model.

---
<a href="{{ site.baseurl }}/assets/img/generate/CsvTemplate.gif"><img alt="Import-Model" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/generate/CsvTemplate.gif" /></a>


<a href="{{ site.baseurl }}/assets/img/generate/GenerateFromCsv.gif"><img alt="Import-Model" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/generate/GenerateFromCsv.gif" /></a>


#### **From URL**
1. **Paste URL**: Enter the URL of your model. 

2. **Fill in Values**: After entering the URL, you will be prompted to fill in the values corresponding to your model's details.

3. **Specify Options**: Lastly, specify if the model is an annotation or if you wish to register it.

Once you complete these steps, click **Next** to continue with the model generation process.

---
<a href="{{ site.baseurl }}/assets/img/generate/GenerateFromURl.gif"><img alt="Import-Model" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/generate/GenerateFromURl.gif" /></a>


