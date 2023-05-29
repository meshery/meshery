---
layout: default
title: Catalog
permalink: concepts/catalog
redirect_from: catalog
type: concepts
abstract: Getting to know about catalog
language: en
list: include
---

# Meshery Catalog

Meshery Catalog functions much like a cloud marketplace, providing a user-friendly interface for browsing, discovering, and sharing configurations and patterns for cloud native infrastructure. With Meshery Catalog, you can easily find and deploy Kubernetes-based infrastructure and tools, making it easy to streamline your cloud native development and operations processes.

### Simplify Your Cloud Native Infrastructure Deployment and Management

Meshery Catalog functions much like a cloud marketplace, providing a user-friendly interface for browsing, discovering, and sharing configurations and patterns for cloud native infrastructure. With Meshery Catalog, you can easily find and deploy Kubernetes-based infrastructure and tools, making it easy to streamline your cloud native development and operations processes.

It also supports a collaborative environment, where DevOps engineers can share their experiences, feedback, and best practices with others in the community. Import cloud native patterns published by others into your Meshery Server. Benefit from and build upon each pattern by incorporating your own tips and tricks, then publish and share with the community at-large. This facilitates knowledge-sharing and helps to build a strong ecosystem of cloud native infrastructure experts.

## Add your content to the catalog

To add your content to the Cloud Native Catalog -

1. Head over [meshery](https://github.com/meshery/meshery.io/blob/master/README.md#add-your-content-to-the-cloud-native-catalog) repository and set up your local enviroment.
1. Create a new file in the catalog collection, under the respective type.
2. Follow the template in the [_defaults.md](https://github.com/meshery/meshery.io/blob/master/collections/_catalog/_defaults.md) and fill out the details.
3. Add your pattern file to the [catalog](https://github.com/meshery/meshery.io/tree/master/catalog) folder so that others can download the file.
4. That's it! You're good to go. Create a PR for your change and sit tight till a maintainer gets a chance to review it.


## To create catalog using Meshery UI

1. Open the Meshery UI in your web browser.
2. Navigate to the configuration section, usually located in the main navigation menu.
3. Head over to Designs and click on import or create design.
4. Select the category and Model as per your need and configure the application.
5. Voilà, You can publish or deploy you design.


## To Create catalog uisng Meshery CLI

1. Ensure that you have Meshery CLI installed on your machine and it is configured to connect to your desired Meshery instance.
2. Open a terminal or command prompt.
3. Use the Meshery CLI commands to interact with the catalog. `mesheryctl pattern`
4. Follow the prompts or instructions provided by the Meshery CLI help.
* Apply pattern file:  `mesheryctl pattern apply --file [path to pattern file | URL of the file]`
* Delete pattern file:  `mesheryctl pattern delete --file [path to pattern file]`
* View pattern file:  `mesheryctl pattern view [pattern name | ID]`
* List all patterns: `mesheryctl pattern list`
5. Onboarding an application. `mesheryctl app onboard -f [file-path]`
6. Applying WASM Filter. `mesheryctl exp filter apply --file [GitHub Link]`


If you have any questions or need assistance, please refer to the [Meshery Documentation](https://docs.meshery.io/) or reach out to our discussion form [layer5.io](https://discuss.layer5.io/).

