---
layout: default
title: Importing Models
abstract: Importing Existing Model and CRD-based Infrastructure Configurations into Meshery as Model
permalink: guides/configuration-management/importing-models
category: configuration
type: guides
language: en
---

Import your existing Models and existing custom resource definition (CRD) into Meshery. The platform supports a variety of application definition formats, and you can import designs using either the Meshery CLI or the Meshery UI.

**Note:** A [Model](/concepts/logical/models) can be only imported if it contains atleast a valid [Component](/concepts/logical/components) or [Relationship](/concepts/logical/relationships).

## Import Models Using Meshery CLI

**Step 1: Install Meshery CLI**

Before you can use the Meshery CLI to import a [Model](/concepts/logical/models), you must first install it. You can install Meshery CLI by [following the instructions]({{site.baseurl}}/installation#install-mesheryctl).


**Step 2: Import the Model**

Model can imported in 2 different format ```URL, File, OCI, Compress(tar.gz)```. A template file is required that contains some required properties: Registrant, Model name, Model DisplayName, Category. The template file is only required when you use URL for generation of model from a crd. This command enable users to import their new models from CRD and existing Meshery Model

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl model import -f [file/url] </div></div>
</pre>

The supported registrant for importing from URL is `github` and `artifacthub`.The URL format must be in this order.


Registrant `Artifacthub`:
- https://artifacthub.io/packages/search?ts_query_web={ model-name } 
- https://istio-release.storage.googleapis.com/charts/base-1.19.0-alpha.1.tgz&sa=D&source=editors&ust=1726839249773905&usg=AOvVaw0j88gkt6FOS1LLSRCYq95X 

Registrant `Github`:
- git:://github.com/cert-manager/cert-manager/master/deploy/crds
- https://github.com/UffizziCloud/uffizzi_controller/releases/download/uffizzi-controller-2.0.1/uffizzi-controller-2.0.1.tgz&sa=D&source=editors&ust=1726839320133140&usg=AOvVaw2AryFwXIPKFnWRjRRTApzp


**Example :**

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl model import -f istio-base.tar</div></div>
</pre>

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl model import -f "git:://github.com/cert-manager/cert-manager/master/deploy/crds" -t template.json</div></div>
</pre>


**Note:** A `-r` flag is present to skip registaration. If the flag is used then no new model would be registered though they would be generated and stored inside `.meshery/models` directory.
The `template.json` can be viewed [here](#).In template.json if the field `isAnnotaion` is true then we would only consider that component as an annotation (svg icon) rather than a normal component.



## Import Models Using Meshery UI

**Step 1: Access the Meshery UI**

To import a model into Meshery using the Meshery UI, you must first [install Meshery](/installation/quick-start)

**Step 2: Navigate to Registry under Settings Page**

Once you have accessed the Meshery UI, navigate to the Registry under Settings. This page can be accessed by clicking on the Settings on the top right on setting icon and then selecting "Registry" and then choose model.

<a href="{{ site.baseurl }}/assets/img/export/Registry.png"><img alt="Registry-Navigator" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/export/Registry.png" /></a>

**Step 3: Upload the Model**

On the Registry page, you can export your model by select Model in registry page. If no specific version the latest version is downloaded and if the specific version of model is selected then on clicking the `Export` button the model is exported as an OCI.

This Meshery model will include components, relationships.

<a href="{{ site.baseurl }}/assets/img/import/ImportingModel.png"><img alt="Import-Model" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/import/ImportingModel.png" /></a>

Once the Meshery model has been exported, you can export your model anytime back using `Import` on UI and then visualize on Kanvas, operate and observe your components that are geneated from the crd. You can also use Meshery to deploy your Meshery Model in form of a design to any of your connected kubernetes clusters. For more information, see [connections](/installation/kubernetes)

