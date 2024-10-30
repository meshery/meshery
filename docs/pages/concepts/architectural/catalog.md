---
layout: default
title: Catalog
permalink: concepts/catalog
redirect_from: catalog
type: components
abstract: Browsing and using cloud native patterns
language: en
list: include
---

[Meshery Catalog](https://meshery.io/catalog) functions much like a cloud marketplace, providing a user-friendly interface for browsing, discovering, and sharing configurations and patterns for cloud native infrastructure. With Meshery Catalog, you can easily find and deploy Kubernetes-based infrastructure and tools, making it easy to streamline your cloud native development and operations processes. A Catolog is based on the Meshery's [Catalog Schema](https://github.com/meshery/schemas/blob/98560345814e4be036d9f0020759faf3202ec2e4/schemas/constructs/v1alpha1/catalog_data.json) with defined attributes.

### Simplify Your Cloud Native Infrastructure Deployment and Management

Meshery Catalog functions much like a cloud marketplace, providing a user-friendly interface for browsing, discovering, and sharing configurations and patterns for cloud native infrastructure. With Meshery Catalog, you can easily find and deploy Kubernetes-based infrastructure and tools, making it easy to streamline your cloud native development and operations processes.

It also supports a collaborative environment, where DevOps engineers can share their experiences, feedback, and best practices with others in the community. Import cloud native patterns published by others into your Meshery Server. Benefit from and build upon each pattern by incorporating your own tips and tricks, then publish and share with the community at-large. This facilitates knowledge-sharing and helps to build a strong ecosystem of cloud native infrastructure experts.


### To create a design pattern using Meshery UI

1. Navigate to the Meshery UI in your web browser.
2. Navigate to the configuration section, usually located in the main navigation menu.
3. Head over to Designs and click on import or create design.
4. Select the category and Model as per your need and configure the application.
5. Voilà, You can publish or deploy you design.

### To create design pattern using Meshery CLI

1. Ensure that you have [Meshery CLI](https://docs.meshery.io/installation/mesheryctl) installed on your machine and it is configured to connect to your desired Meshery instance.
2. Open a terminal or command prompt.
3. Use the Meshery CLI commands to interact with the catalog. `mesheryctl design`
4. Follow the prompts or instructions provided by the Meshery CLI help.
* Apply [design file](https://docs.meshery.io/guides/configuration-management):  `mesheryctl design apply --file [path to design file | URL of the file]`
* Delete design file:  `mesheryctl design delete --file [path to design file]`
* View design file:  `mesheryctl design view [design name | ID]`
* List all designs: `mesheryctl design list`
5. [Importing](https://docs.meshery.io/reference/mesheryctl#cloud-native-pattern-configuration-and-management) a design. `mesheryctl design import -f [file-path] -s [manifest | compose | helm]`
6. Applying [WASM Filter](https://docs.meshery.io/guides/configuration-management#wasm-filters). `mesheryctl filter import [file | URL] --wasm-config [filepath|string]`


### Publishing a Design to Meshery Catalog

1. **Request to Publish**: The author submits a request to publish their design to the Meshery Catalog, including a description and any relevant considerations.
2. **Review by Admin**: The workspace owner or admin reviews the design. They have the option to approve, deny, or request changes by commenting on the design.
3. **Approval Process**:
    - If the admin or workspace admin **approves** the design, a validation is performed to ensure the design data is accurate. Once validated, the design is published to the catalog.
    - If the admin **denies** the design, feedback is provided for necessary changes. After the changes are made, the design is automatically published with appropriate versioning.
4. **Ongoing Management**:
    - The author or workspace owner retains permission to edit, delete, or unpublish their designs from the catalog at any time.
5. **Notification for Changes**: If the design no longer adds value to the Meshery Catalog, a prior notification is sent to the author, and the design may be unpublished.
6. **GitHub Workflow Integration**: Once approved, a GitHub workflow is triggered to publish the design to the Meshery.io Catalog.

<a href="{{ site.baseurl }}/assets/img/architecture/Catalog-Publishing-Workflow.svg" class="lightbox-image">
<img src="{{ site.baseurl }}/assets/img/architecture/Catalog-Publishing-Workflow.svg" width="70%" /></a>
<figure>
  <figcaption>Figure: Workflow to publish a design in catalog</figcaption>
</figure>

### FAQ
<details>
    <summary>
<h6>Question: Why are images invisible for some designs in the Meshery Catalog?</h6>
</summary>
<p><strong>Answer:</strong> In certain instances, the images of published designs in <a href="https://meshery.io/catalog">Meshery Catalog</a> may not be visible due to bandwidth issues. This can occur when there are network constraints affecting the retrieval of image data. However, rest assured that the design information and other relevant details are still accessible.</p>
</details>

{% include alert.html
    type="info"
    title="Help with Meshery Catalog"
    content="If you have any questions or need assistance, reach out on the <a href='http://discuss.meshery.io/'>discussion forum</a>." %}
