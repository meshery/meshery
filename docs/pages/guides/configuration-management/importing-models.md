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

{% include alert.html type="warning" title="Limitation on Importing Connections" content="The `mesheryctl model import` command currently supports importing Models containing Components, Relationships, and Policies. Importing Models with `Connection` definitions is not yet supported. This functionality may be added in a future release." %}

**Step 1: Install Meshery CLI**

Before you can use the Meshery CLI to import a [Model](/concepts/logical/models), you must first install it. You can install Meshery CLI by [following the instructions]({{site.baseurl}}/installation#install-mesheryctl).


**Step 2: Import the Model**

A model can be imported using `Directory, URL` or `File`. This command is used for specifically importing models and not for creating a new model. For Meshery models, you can refer to the exported models at [Meshery catalog](https://meshery.io/catalog/models).

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl model import -f [file/url] </div></div>
</pre>

The extensions supported for this operation are `.yml`, `.tar.gz`, `.gz`, `.tgz`, `.yaml`, `.json`, `.tar`, `.tar.tgz` and the supported registrants are `Github`,`Meshery` and `Artifacthub`. The URL format must be in this order.

https://github.com/{owner}/{repo}/raw/refs/heads/main/filename

**Example :**

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl model import -f istio-base.tar</div></div>
</pre>

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl model import -f "https://github.com/{owner}/{repo}/raw/refs/heads/main/filename"</div></div>
</pre>

For importing Meshery models, go to [Meshery catalog](https://meshery.io/catalog/models), click on the model you want to import, download it and use:

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl model import -f path/to/your/file</div></div>
</pre>



## Import Models Using Meshery UI

**Step 1: Access the Meshery UI**

To import a model into Meshery using the Meshery UI, you must first [install Meshery](/installation/quick-start)

**Step 2: Navigate to Registry under Settings Page**

Once you have accessed the Meshery UI, navigate to the Registry under Settings. This page can be accessed by clicking on the Settings on the top right on setting icon and then selecting "Registry" and then choose model.

<a href="{{ site.baseurl }}/assets/img/export/Registry.png"><img alt="Registry-Navigator" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/export/Registry.png" /></a>

**Step 3: Upload the Model**

On the Registry page, you can import your model by clicking the import button in registry page. Select URL or File and then hit Import. Extensions supported for this operation are `.yml`, `.tar.gz`, `.gz`, `.tgz`, `.yaml`, `.json`, `.tar`, `.tar.tgz`.

This operation is only used for importing an existing model and not for creating one. For importing meshery models head over to [Meshery catalog](https://meshery.io/catalog/models). 

This Meshery model will include components, relationships.

<a href="{{ site.baseurl }}/assets/img/import/ImportModel.gif"><img alt="Import-Model" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/import/ImportModel.gif" /></a>


