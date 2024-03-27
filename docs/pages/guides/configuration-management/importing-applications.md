---
layout: default
title: Importing Applications
abstract: Learn how to import your existing application definitions and your existing infrastructure configurations into Meshery as you to manage, operate, and observe your cloud native infrastructure more effectively.
permalink: guides/configuration-management/importing-apps
category: configuration
type: guides
language: en
---

Import your existing application definitions and existing infrastructure configurations into Meshery. The platform supports a variety of application definition formats, and you can import apps using either the Meshery CLI or the Meshery UI.

## Supported Application Definition Formats

Meshery supports the following application definition formats:

- [Kubernetes Manifests](https://kubernetes.io/docs/concepts/overview/working-with-objects/kubernetes-objects/)
- [Helm Charts](https://helm.sh/docs/topics/charts/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Meshery Designs](/concepts/logical/designs)

## Import Apps Using Meshery CLI

**Step 1: Install Meshery CLI**

Before you can use the Meshery CLI to import a Docker Compose app, you must first install it. You can install Meshery CLI by [following the instructions]({{site.baseurl}}/installation#install-mesheryctl).

**Step 2: Import the App Manifest**

Once you have created your App Definition file, you can use the Meshery CLI to import your Docker Compose app into Meshery. To do this, run the following command:

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl app import -f [file/url] -s [source-type]</div></div>
</pre>

This command enable users to import their existing applications from sources as

- Helm Charts
- Kubernetes Manifests
- Docker Compose

**Example :**

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl app import -f ./SampleApplication.yml -s "Kubernetes Manifest"</div></div>
</pre>

## Import Apps Using Meshery UI

**Step 1: Access the Meshery UI**

To import a Docker Compose app into Meshery using the Meshery UI, you must first [install Meshery](/installation/quick-start)

**Step 2: Navigate to the Application section in the Configuration**

Once you have accessed the Meshery UI, navigate to the App Import page. This page can be accessed by clicking on the "Applications" menu item and then selecting "Import Application".

<a href="{{ site.baseurl }}/assets/img/applications/Menu.png"><img alt="Application-Navigation" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/applications/Menu.png" /></a>

**Step 3: Upload the Application**

On the App Import page, you can upload your application by select File Type from the options and clicking on the "Browse" button and selecting the file from your local machine or uploading in through URL. Once you have selected the file, click on the "Import" button to import app into Meshery.

When you import an app into Meshery, it will create a Meshery App based on definition. This Meshery App will include all of the services, ports, and other parameters defined in the File.

<a href="{{ site.baseurl }}/assets/img/applications/ImportApp.png"><img alt="Import-Application" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/applications/ImportApp.png" /></a>

Once the Meshery App has been created, you can use Meshery to manage, operate and observe your cloud native infrastructure. You can also use Meshery to deploy your Meshery App to any of your connected kubernetes clusters. For more information, see [connections](/installation/kubernetes)

