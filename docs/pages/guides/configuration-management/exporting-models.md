---
layout: default
title: Exporting Models
abstract: Exporting Meshery Model as an OCI or Compress file
permalink: guides/configuration-management/exporting-models
category: configuration
type: guides
language: en
---

Export existing Meshery [Models](/concepts/logical/models) or your custom infrastructure crd imported as [Models](/concepts/logical/models) in Meshery. The platform supports a variety of application definition formats, and you can export designs using either the Meshery CLI or the Meshery UI.

## Import Models Using Meshery CLI

**Step 1: Install Meshery CLI**

Before you can use the Meshery CLI to export a [Model]((/concepts/logical/models)), you must first install it. You can install Meshery CLI by [following the instructions]({{site.baseurl}}/installation#install-mesheryctl).


**Step 2: Import the Model**

Model can be exported in 2 format and OCI or Compress file. These files contain different [components](/concepts/logical/components), [relationships](/concepts/logical/relationships) files and a model file. The files can be exported in either yaml or json by default it is exported as yaml.

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl model export [model-name] </div></div>
</pre>

Different file type and output type can be set using flags `-o` for oci and tar and `-t` for json or yaml. A location can also be specified using the -l flag where user needs to export the model


**Example :**

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl model export istio-base -o oci</div></div>
</pre>

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl model export istio-base -o oci -t json</div></div>
</pre>


**Note:** A `--discard-components --discard-relationships` flag is present to skip [component](/concepts/logical/components) and [relationships](/concepts/logical/relationships) generation in the export only the model is exported without any components and relationships.



## Import Models Using Meshery UI

**Step 1: Access the Meshery UI**

To import a model into Meshery using the Meshery UI, you must first [install Meshery](/installation/quick-start)

**Step 2: Navigate to Registry under Settings Page**

Once you have accessed the Meshery UI, navigate to the Registry under Settings. This page can be accessed by clicking on the Settings on the top right on setting icon and then selecting "Registry" and then choose model.

<a href="{{ site.baseurl }}/assets/img/export/Registry.png"><img alt="Registry-Navigator" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/export/Registry.png" /></a>

**Step 3: Upload the Models**

On the Registry page, you can upload your model by select File Type from the options and clicking on the "Browse" button and selecting the file (tar.gz or OCI) from your local machine or uploading in through URL. If you choose URL you need to provide in some fields like registrant, model name, model displayname and category. Once you have selected the file or entered the URL and filled the form, click on the "Import" button to import model into Meshery.

When you import an model into Meshery, it will create a Meshery model based on definition. This Meshery model will include components, relationships.

<a href="{{ site.baseurl }}/assets/img/export/ExportingModel.png"><img alt="Export-Model" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/export/ExportingModel.png" /></a>

Once the Meshery model has been exported, you can use [`Import`](/guides/configuration-management/importing-models) to import the model and then visualize,operate and observe your components on Kanvas. You can also use Meshery to deploy your Meshery Model in form of a design to any of your connected kubernetes clusters. For more information, see [connections](/installation/kubernetes)

