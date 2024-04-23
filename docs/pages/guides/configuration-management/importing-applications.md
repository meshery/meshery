---
layout: default
title: Importing Designs
abstract: Learn how to import your existing designs and your existing infrastructure configurations into Meshery as you to manage, operate, and observe your cloud native infrastructure more effectively.
permalink: guides/configuration-management/importing-designs
category: configuration
type: guides
language: en
---

Import your existing designs and existing infrastructure configurations into Meshery. The platform supports a variety of application definition formats, and you can import designs using either the Meshery CLI or the Meshery UI.

## Supported Design Definition Formats

Meshery supports the following design definition formats:

- [Kubernetes Manifests](https://kubernetes.io/docs/concepts/overview/working-with-objects/kubernetes-objects/)
- [Helm Charts](https://helm.sh/docs/topics/charts/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Meshery Designs](/concepts/logical/designs)

## Import Designs Using Meshery CLI

**Step 1: Install Meshery CLI**

Before you can use the Meshery CLI to import a Docker Compose app, you must first install it. You can install Meshery CLI by [following the instructions]({{site.baseurl}}/installation#install-mesheryctl).

**Step 2: Import the Design Manifest**

Once you have created your Design Definition file, you can use the Meshery CLI to import your Docker Compose app into Meshery. To do this, run the following command:

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl pattern import -f [file/url] -s [source-type]</div></div>
</pre>

This command enable users to import their existing designs from sources as

- Helm Charts
- Kubernetes Manifests
- Docker Compose

**Example :**

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl pattern import -f ./SampleDesign.yml -s "Kubernetes Manifest"</div></div>
</pre>

## Import Designs Using Meshery UI

**Step 1: Access the Meshery UI**

To import a Docker Compose app into Meshery using the Meshery UI, you must first [install Meshery](/installation/quick-start)

**Step 2: Navigate to the Designs section in the Configuration**

Once you have accessed the Meshery UI, navigate to the Design Import page. This page can be accessed by clicking on the "Designs" menu item and then selecting "Import Design".

<a href="{{ site.baseurl }}/assets/img/applications/Menu.png"><img alt="Application-Navigation" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/applications/Menu.png" /></a>

**Step 3: Upload the Design**

On the Design Import page, you can upload your design by select File Type from the options and clicking on the "Browse" button and selecting the file from your local machine or uploading in through URL. Once you have selected the file, click on the "Import" button to import design into Meshery.

When you import an design into Meshery, it will create a Meshery design based on definition. This Meshery design will include all of the services, ports, and other parameters defined in the File.

<a href="{{ site.baseurl }}/assets/img/applications/ImportDesign.png"><img alt="Import-Application" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/applications/ImportDesign.png" /></a>

Once the Meshery design has been created, you can use Meshery to manage, operate and observe your cloud native infrastructure. You can also use Meshery to deploy your Meshery design to any of your connected kubernetes clusters. For more information, see [connections](/installation/kubernetes)

