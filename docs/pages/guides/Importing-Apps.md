---
layout: default
title: Importing  Applications
description: This guide is to help users get a better understanding of sample apps
permalink: guides/importing-apps
type: Guides
language: en
---


Importing apps into Meshery as Meshery Apps is a powerful feature that enables users to manage, operate, and observe their service mesh more effectively. With this feature, users can easily import their app manifest and store it in the database.

You can import apps into Meshery using either the Meshery CLI or the Meshery UI. We will discuss both methods in detail below.


## Using Meshery CLI:


**Step 1: Install Meshery CLI**

Before you can use the Meshery CLI to import a Docker Compose app, you must first install it. You can install Meshery CLI by following the instructions on the [Meshery documentation site](../installation/mesheryctl.md).


**Step 2: Import the App Manifest**

Once you have created your App Definition file, you can use the Meshery CLI to import your Docker Compose app into Meshery. To do this, run the following command:



<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl app import -f [file/url] -s [source-type]</div></div>
</pre>

This command enable users to import their existing applications from 
- Helm Charts
- k8s manifests
- Docker Compose

Example : 

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">mesheryctl app import -f ./SampleApplication.yml -s "Kubernetes Manifest"</div></div>
</pre>


## Using Meshery UI:


**Step 1: Access the Meshery UI**


To import a Docker Compose app into Meshery using the Meshery UI, you must first Install the Meshery. You can install Meshery by following the instructions on the [Meshery documentation site](../installation/quick-start.md) 



**Step 2: Navigate to the Application section in the Configuration**


Once you have accessed the Meshery UI, navigate to the App Import page. This page can be accessed by clicking on the "Applications" menu item and then selecting "Import Application".


<a href="{{ site.baseurl }}/assets/img/applications/Menu.png"><img alt="Application-Navigation" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/applications/Menu.png" /></a>


**Step 3: Upload the Application**

On the App Import page, you can upload your application by select File Type from the options and clicking on the "Browse" button and selecting the file from your local machine or uploading in through URL Once you have selected the file, click on the "Import" button to import app into Meshery.
When you import a app into Meshery, it will create a Meshery App based on definition. This Meshery App will include all of the services, ports, and other parameters defined in the File.

<a href="{{ site.baseurl }}/assets/img/applications/ImportApp.png"><img alt="Import-Application" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/applications/ImportApp.png" /></a>

Once the Meshery App has been created, you can use Meshery to manage, operate and observe your service mesh. You can also use Meshery to deploy your Meshery App to your service mesh.
