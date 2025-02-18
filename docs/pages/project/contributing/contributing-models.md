---
layout: page
title: Contributing to Models
permalink: project/contributing/contributing-models
redirect_from: project/contributing/contributing-models/
abstract: How to contribute to Meshery Models, Components, Relationships, Policies...
language: en
type: project
category: contributing
list: include
---

<!-- Concepts for which docs needs to be updated: -->
<!-- Scopes - What they mean to contributors/expand on which takes precedence?
1. Which policies get loaded?
   1. What policies are loaded by default?
   2. What happens in conflict?
2. What controls are exposed to model contributors?
3. Are there any global Meshery defaults (can user change them?)
4. Instructions for Creating a New Connection
5. Instructions for Creating a New Component -->

# Understanding the internals of Meshery's logical object model

Meshery's internal object model is designed to provide a consistent and extensible way of capturing and characterizing the resources under Meshery's management and the capabilities Meshery has at its disposal. Meshery Models serve as the unit of packaging for the object models that define a registered capability or a type of managed infrastructure and their relationships, and details specifics of how to manage them.Models often represent infrastructure and application technologies, however, they are also capable of defining other types of entities like annotations, like shapes (infrastructure ambiguous components). Models are used to define the capabilities of Meshery. _See [Models]({{site.baseurl}}/concepts/logical/models) to learn more about models as a logical concept._

Each model includes a set of entities (in the form of definitions) that Meshery can manage. Models are defined and versioned using on the [Model Schema](https://github.com/meshery/schemas/blob/master/schemas/constructs/openapi/meshmodels.yml). The schema defines the structure of the model, including the entities it contains, their relationships, and the properties they have. The schema also defines the version of the model and the version of the schema itself. _See [Registry]({{site.baseurl}}/concepts/logical/registry) to learn more about Meshery's internal registry and how to use it._

[![Model Entity Classification]({{ site.baseurl }}/assets/img/meshmodel/meshmodel-architecture.svg)]({{ site.baseurl }}/assets/img/concepts/meshery-models.png)
_Figure: Model Entity Classification_

## Meshery Entities and their Lifecycle

This section aids in your understanding of the vernacular of Meshery's internal object model and discusses the difference beteween schemas, definitions, declarations, and instances. The lifecycle of Meshery entities (components, relationships, policies) is represented by the following terms, which are used to describe the various stages of their lifecycle.

### Schema

**Schema** _(static)_ **: the skeletal structure representing a logical view of the size, shape, characteristics of a construct.**

The schema represents the skeletal structure of an entity and provides a logical view of its size, shape, and characteristics. It defines the expected properties and attributes of the entity. The schema serves as a blueprint or template for creating instances of the entity. It is a static representation that defines the structure and properties but does not contain specific configuration values.

{% include alert.html type="info" title="Schema example" content='<details><summary>Component schema excerpt</summary><pre> {
"$id": "https://schemas.meshery.io/component.json",
  "$schema": "<http://json-schema.org/draft-07/schema#>",
"description": "Components are the atomic units for designing infrastructure. Learn more at <https://docs.meshery.io/concepts/components>",
"required": [
"apiVersion",
"kind",
"schema",
"model"
],
"additionalProperties": false,
"type": "object",
"properties": {
"apiVersion": {
"type": "string",
"description": "API Version of the component."
},
"kind": {
"type": "string",
"description": "Kind of the component."
.
.
.

</pre></details> See <a href="https://github.com/meshery/schemas">github.com/meshery/schemas</a> for more details.' %}

### Definition

**Definition** _(static)_ **: An implementation of the Schema containing an outline of the specific attributes of a given, unconfigured entity.**

A definition is an implementation of the schema. It contains specific configurations and values for the entity at hand. The definition provides the actual configuration details for a specific instance of the entity. It is static because it is created based on the schema but does not change once created. The definition is used to instantiate declarations of the entity.

{% include alert.html type="info" title="Definition example" content="a generic, unconfigured Kubernetes Pod." %}

### Declaration

**Declaration** _(static)_ **: - A configured entity with detailed intentions of a given Definition.**

{% include alert.html type="info" title="Declaration example" content="NGINX container as a Kubernetes Pod with port 443 and SSL termination." %}

### Instance

**Instance** _(dynamic)_ **: A realized entity (deployed/discovered); An instantiation of the declaration.**

An _instance_ represents a realized entity. An _instance_ is a dynamic representation that corresponds to a deployed or discovered instantiation of a _declaration_. An _instance_ is created based on its corresponding _definition_ and represents an actual running or deployed version of the entity within the environment.

{% include alert.html type="info" title="Instance example" content="NGINX-as234z2 pod running in a cluster as a Kubernetes Pod with port 443 and SSL termination." %}

## Instructions for Creating a New Model

{% include alert.html type="info" title="Creating Models Quick Start" content="See the <a href='/project/contributing/contributing-models-quick-start'>quick start</a> for a no fluff guide to creating your first Meshery model." %}

All of Meshery's Models can be found in the [Meshery Integrations spreadsheet](https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw/edit#). This spreadsheet serves as the source of truth for the definition of Meshery's models and is refreshed daily.

{% include alert.html type="light" title="Model Source Code" content="See examples of <a href='https://github.com/meshery/meshery/tree/master/server/meshmodel'>Models defined in JSON in meshery/meshery</a>." %}

<div class="tab-container">
  <!-- First Level Tab: mesheryctl -->
  <input type="radio" id="tab1" name="tabs" checked>
  <label for="tab1">
    <i class="fa fa-terminal"></i> mesheryctl
  </label>
  <section class="tabbed">
    <h3>Prerequisites:</h3>
    <ul>
      <li>Fork the <a href="https://github.com/meshery/meshery" target="_blank" rel="noopener">meshery/meshery repository.</a></li>
      <li>Install the Meshery CLI by following the <a href="https://docs.meshery.io/installation/" target="_blank" rel="noopener">installation instructions.</a></li>
    </ul>

    {% include alert.html type="info" title="Generating Models does not require Meshery Server" content="Meshery Server is not required to generate models. The Meshery CLI can be used to generate models. Model and Component generation logic is MeshKit. `mesheryctl` and Meshery Server both utilize MeshKit’s libraries for ongoing programmatic generation of models and components." %}

    <br />
    <!-- Second Level Tabs under mesheryctl -->
    <div class="tab-container">
      <!-- CSV Method -->
      <input type="radio" id="csv-tab" name="mesheryctl-tabs" checked>
      <label for="csv-tab">
        <i class="fa fa-list"></i> Using CSV
      </label>
      <section class="tabbed">
        <h4>1. Understanding the Template Directory</h4>
        <p>Inside your forked Meshery repository, you'll find the templates-csvs directory containing three essential CSV files:</p>
        <code>
            mesheryctl/templates/templates-csvs/
                <br />
                ├── models.csv       # Define model metadata and core properties
                <br />
                ├── components.csv   # Specify individual components and their characteristics
                <br />
                └── relationships.csv # Define how components interact and connect
                <br />
        </code>

        <h4>2. Customizing Your Model</h4>
        <p>Creating your model involves modifying these CSV files to match your specific requirements. When making changes, you have two valuable references at your disposal: the existing entries in the CSV files serve as practical examples, while the <a href="https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw">integration spreadsheet</a> provides comprehensive documentation of all possible fields and their purposes.</p>

        <h4>3. Generating Your Model</h4>
        <p>Once you've customized your CSV files, you can generate your model using a single command. Ensure you're in the root directory of your forked Meshery repository, as this maintains the correct file path relationships:</p>
        <code>mesheryctl registry generate --directory templates-csvs --model "YOUR_MODEL_NAME"</code>

        <h4>4. Locating Generated Files</h4>
        <p>After successful generation, your model's files will be created in the Meshery server's model directory. You can find these files at <code>meshery/server/meshmodel/[YOUR_MODEL_NAME]/</code>. Take time to review these generated files to ensure they accurately reflect your intended model structure.</p>

        <h4>5. Troubleshooting</h4>
        <p>If you encounter issues during the generation process, you can use these diagnostic approaches to identify and resolve problems:</p>
        <ul>
          <li>Examine the detailed error logs at <code>~/.meshery/logs/registry/</code> to understand specific generation issues.</li>
          <li>Review your CSV files for proper formatting, ensuring all required columns are present and correctly populated.</li>
          <li>Confirm you're executing the command from the root of your forked Meshery repository.</li>
        </ul>
      </section>

      <!-- Spreadsheet Method -->
      <input type="radio" id="spreadsheet-tab" name="mesheryctl-tabs">
      <label for="spreadsheet-tab">
        <i class="fa fa-table"></i> Using Integration Spreadsheet
      </label>
      <section class="tabbed">
        <h3>Setting Up Your Environment</h3>
        <p>Before you begin working with the Integration Spreadsheet, you'll need to complete several important setup steps:</p>

        <h4>1. Spreadsheet Preparation</h4>
        <p>Start by creating your own copy of the Meshery Integration Sheet:</p>
        <ol>
          <li>Visit the <a href="https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" target="_blank" rel="noopener">Meshery Integration Sheet</a></li>
          <li>Make a copy using File > Make a copy</li>
          <li>Look at the URL of your new spreadsheet and note the ID (the long string between /d/ and /edit)</li>
        </ol>

        <h4>2. Google Cloud Configuration</h4>
        <p>Set up your Google Cloud environment with these steps:</p>
        <ol>
          <li><a href="https://developers.google.com/workspace/guides/create-project" target="_blank" rel="noopener">Create a new Google Cloud Project or select an existing one</a></li>
          <li><a href="https://support.google.com/googleapi/answer/6158841" target="_blank" rel="noopener">Enable the Google Sheets API for your project</a></li>
          <li><a href="https://developers.google.com/workspace/guides/create-credentials#create_credentials_for_a_service_account" target="_blank" rel="noopener">Create service account credentials</a></li>
        </ol>

        <h4>3. Credential Configuration</h4>
        <p>Set up your credentials in your local environment:</p>
        <code>
          base64 -w 0 /path/to/your-service-account-creds.json
        </code>
        <br />
        <code>
          echo 'export SHEET_CRED="[paste-base64-output-here]"' >> ~/.bashrc
        </code>
        <br />
        <code>
          source ~/.bashrc
        </code>

        <h4>4. Spreadsheet Access Configuration</h4>
        <ol>
          <li>Open your copied spreadsheet</li>
          <li>Click "Share" in the top right</li>
          <li>Add your service account email (ends with @developer.gserviceaccount.com)</li>
          <li>Grant "Editor" permissions</li>
          <li>Publish the spreadsheet:
            <ul>
              <li>File > Share > Publish to web</li>
              <li>Select "Comma-separated values (.csv)"</li>
              <li>Click "Publish"</li>
            </ul>
          </li>
        </ol>

        <h3>Working with the Integration Spreadsheet</h3>
        <p>Once your environment is set up, you can begin working with the spreadsheet:</p>

        <h4>1. Adding Your Model</h4>
        <p>The integration spreadsheet contains existing model definitions that serve as practical examples. You can either create a new entry following the patterns in existing rows, or practice by generating an existing model first to understand the process. Each row represents a complete model definition, use them as reference for creating a new row.</p>

        <h4>2. Generating the Model</h4>
        <p>Use mesheryctl to generate your models. Make sure to run the command inside your forked <code>meshery/meshery</code> repo</p>
        <code>mesheryctl registry generate --spreadsheet-id "YOUR_SPREADSHEET_ID" --spreadsheet-cred "$SHEET_CRED" --model "YOUR_MODEL_NAME"</code>

        <p>The command will:</p>
        <ul>
          <li>Read your spreadsheet data</li>
          <li>Validate the model definition</li>
          <li>Generate the model files</li>
          <li>For error logs your can checkout <code>~/.meshery/logs/registry/</code></li>
        </ul>

        <h4>3. Verification</h4>
        <p>The model will be generated in <code>meshery/server/meshmodels/[YOUR_MODEL_NAME]</code></p>
      </section>
    </div>


    {% include alert.html type="info" title="Using Meshery CLI with the Meshery Registry and Meshery Models" content="Meshery CLI has a set of commands that pertain to the lifecycle management of models:
    <br />

    - <code>mesheryctl registry</code> - interact with and update spreadsheets
    <br />
    - <code>mesheryctl models</code> - interact with and update Meshery Server
    <br />
    - <code>mesheryctl components</code> - interact with and update Meshery Server
    <br />
    - <code>mesheryctl relationships</code> - interact with and update Meshery Server" %}

  </section>

  <!-- First Level Tab: Meshery UI -->
  <input type="radio" id="tab2" name="tabs">
  <label for="tab2">
    <i class="fa fa-desktop"></i> Meshery UI
  </label>
  <section class="tabbed">
    <!-- Second Level Tabs under Meshery UI -->
    <div class="tab-container">
      <!-- URL Import Method -->
      <input type="radio" id="url-tab" name="ui-tabs" checked>
      <label for="url-tab">
        <i class="fa fa-link"></i> URL Import
      </label>
      <section class="tabbed">
        <p>The URL Import feature allows you to generate models in Meshery by providing URLs to source repositories or package registries. Here's how to use it:</p>

        <h4>1. Access the Model Generation Interface</h4>
        <p>Navigate to <a href="https://playground.meshery.io/settings?settingsCategory=Registry&tab=Models">Registry</a> in the Meshery UI. Click the "Generate" button to begin creating a new model. In the Upload Method dialog, select "URL Import" and click Next.</p>

        <a href="/assets/img/registry/generate-model-from-UI.png">
          <img src="/assets/img/registry/generate-model-from-UI.png" alt="registry generate model" style="width: 50%; max-width: 400px;">
        </a>

        <h4>2. Provide Source Location</h4>
        <p>You can provide either a GitHub repository URL or an ArtifactHub package URL as your source:</p>

        <p><strong>For GitHub repositories:</strong><br>
        Enter a URL in the format: <code>git://github.com/[organization or username]/[repository]/[branch]/path/to/crds</code><br>
        For example: <code>git://github.com/cert-manager/cert-manager/master/deploy/crds</code></p>

        <p><strong>For ArtifactHub packages:</strong><br>
        Enter a URL in the format: <code>https://artifacthub.io/packages/search?ts_query_web={model-name}</code></p>

        <p><strong>Pro tip:</strong> Check the <code>sourceURL</code> column in the <a href="https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" target="_blank" rel="noopener">Meshery Integration Sheet</a> and try one of the listed sources.</p>

        <a href="/assets/img/registry/url-import-github-url.png">
          <img src="/assets/img/registry/url-import-github-url.png" alt="registry generate model source selection" style="width: 50%; max-width: 400px;">
        </a>

        <h4>3. Configure Model Details</h4>
        <p>Enter the required information for your model:</p>
        <ul>
          <li><strong>Model Name:</strong> Should be in lowercase with hyphens. For example, cert-manager.</li>
          <li><strong>Display Name:</strong> How you want your model to be named. For example, Cert Manager.</li>
        </ul>

        <a href="/assets/img/registry/url-import-model-name.png">
          <img src="/assets/img/registry/url-import-model-name.png" alt="registry generate model name" style="width: 50%; max-width: 400px;">
        </a>

        <h4>4. Set Model Categorization</h4>
        <p>Choose appropriate category and subcategory for your model from the dropdown menus. If your model doesn't fit existing categories, select "Uncategorized". This helps organize models in the registry and makes them easier to find.</p>

        <a href="/assets/img/registry/url-import-category.png">
          <img src="/assets/img/registry/url-import-category.png" alt="registry generate model categorization" style="width: 50%; max-width: 400px;">
        </a>

        <h4>5. Configure Model Styling</h4>
        <p>Customize your model's appearance:</p>
        <ul>
          <li>Upload logos for both dark and light themes</li>
          <li>Set primary and secondary colors for visual elements</li>
          <li>Select a shape for the model's icon in the UI</li>
        </ul>
        <p>Note: If you don't provide custom styling, Meshery's default values will be used. You can change these later in the model definition.</p>

        <a href="/assets/img/registry/url-import-model-styling.png">
          <img src="/assets/img/registry/url-import-model-styling.png" alt="registry generate model styling" style="width: 50%; max-width: 400px;">
        </a>

        <h4>6. Additional Settings</h4>
        <p>Before finishing, you can:</p>
        <ul>
          <li>Choose to register the model immediately for instant availability in Meshery instance.</li>
          <li>Specify if the model is for visual annotation only</li>
        </ul>

        {% include alert.html type="light" title="Visual Annotation Models" content="When a model is marked for visual annotation only, it means the model will be used purely for visualization and diagramming purposes within Meshery's interface, rather than for actual infrastructure management." %}

        <a href="/assets/img/registry/url-import-model-additional-settings.png">
          <img src="/assets/img/registry/url-import-model-additional-settings.png" alt="registry generate model additional settings" style="width: 50%; max-width: 400px;">
        </a>

        <p>After completing these steps, click "Finish" to generate your model. Once generated, you can find your model in the Registry section (if you checked "Register Model Immediately") else it'll download the generated model in an archive, ready for use in your Meshery environment.</p>
      </section>

      <!-- CSV Import Method -->
      <input type="radio" id="ui-csv-tab" name="ui-tabs">
      <label for="ui-csv-tab">
        <i class="fa fa-list"></i> Using CSV
      </label>
      <section class="tabbed">
        <p>The CSV Import feature allows you to generate models in Meshery by providing template CSV files that define your model structure, components, and relationships. Here's a comprehensive guide on how to use this feature:</p>

        <h4>1. Access the Model Generation Interface</h4>
        <p>Navigate to <a href="https://playground.meshery.io/settings?settingsCategory=Registry&tab=Models">Registry</a> in the Meshery UI. Click the "Generate" button to begin creating a new model. In the Upload Method dialog, select "URL Import" and click Next.</p>

        <a href="/assets/img/registry/generate-model-from-UI.png">
          <img src="/assets/img/registry/generate-model-from-UI.png" alt="registry generate model" style="width: 50%; max-width: 400px;">
        </a>

        <a href="/assets/img/registry/csv-import.png">
          <img src="/assets/img/registry/csv-import.png" alt="CSV Import Initial Screen" style="width: 50%; max-width: 400px;">
        </a>

        <h4>2. Prepare Your CSV Files</h4>
        <p>You'll need three essential CSV files to define your model. You can find templates for these files in the <a href="https://github.com/meshery/meshery/tree/a514f8689260791077bde8171646933cff15dd08/mesheryctl/templates/template-csvs" target="_blank" rel="noopener noreferrer">Meshery repository</a>. Each file serves a specific purpose:</p>

        <ul>
          <li><strong>models.csv:</strong> Defines your model's core metadata, including name, version, and general properties</li>
          <li><strong>components.csv:</strong> Describes the individual components that make up your model</li>
          <li><strong>relationships.csv:</strong> Specifies how different components interact and connect with each other</li>
        </ul>

        <p><strong>Pro tip:</strong> Look at existing models in the <a href="https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw/edit?gid=0#gid=0" target="_blank" rel="noopener noreferrer">Meshery Integration Sheet</a> to understand how to structure your CSV files effectively.</p>

        <h4>3. Upload Models, Components, and Relationships CSV</h4>

        <a href="/assets/img/registry/csv-import-upload-model-csv.png">
          <img src="/assets/img/registry/csv-import-upload-model-csv.png" alt="Model CSV Upload" style="width: 50%; max-width: 400px;">
        </a>

        <a href="/assets/img/registry/csv-import-upload-components-csv.png">
          <img src="/assets/img/registry/csv-import-upload-components-csv.png" alt="Component CSV Upload" style="width: 50%; max-width: 400px;">
        </a>

        <a href="/assets/img/registry/csv-import-upload-relationship-csv.png">
          <img src="/assets/img/registry/csv-import-upload-relationship-csv.png" alt="Relationship CSV Upload" style="width: 50%; max-width: 400px;">
        </a>

        <h4>6. Model Registration</h4>
        <p>In the final step, you can choose to register your model immediately in your Meshery instance. This makes the model available for immediate use after generation.</p>

        <a href="/assets/img/registry/csv-import-register-model.png">
          <img src="/assets/img/registry/csv-import-register-model.png" alt="Model Registration Options" style="width: 50%; max-width: 400px;">
        </a>

        <p>After completing these steps and successfully generating your model, you can find it in the Registry section if you chose to register it immediately. Otherwise, you'll receive a downloaded archive containing your generated model files.</p>
      </section>
    </div>
  </section>
</div>

### Importing Generated Models

The generated model can be importing using both Mesheryctl and Meshery UI. Read [Importing Models]({{site.baseurl}}/guides/configuration-management/importing-models) for detailed instructions on how to import models.

### Post Model Generation

During model generation, corresponding components are created. Next step is to enrich these component details and define their capabilities and relationships.

1. **Enrich Component Details**
   When a Component is initially generated, a new Component definition is created with default properties (e.g. colors, icons, capabilities, etc.), some of which are inherited from their respective Model.

   - **3.1. Customize Shapes and Colors**

     - Default shape for new components is a circle
     - Consider enriching components' details based on what they represent
     - Reference Cytoscape [node types](https://js.cytoscape.org/demos/node-types/) for possible shapes
     - Example: Use a pentagon shape to represent a Deployment
     - Know more about [components shapes and colors](https://docs.meshery.io/extensions/component-shape-guide)

   - **3.2. Customize Icons**

     - Components inherit the icon (colored and white SVGs) of their respective Model by default
     - Propose specific icons best suited to visually represent each component
     - Example: Use a skull icon for a DaemonSet

   - **3.3. Review Capabilities**
     - Review and confirm assigned capabilities
     - Modify capabilities as needed

    See the [Contributing to Components]({{site.baseurl}}/project/contributing/contributing-components) for detailed instructions.

4. **Identify Relationships**

   - **4.1. Review Available Types**
     Review and familiarize yourself with the predefined relationship kinds, types, and subtypes. See ["Relationships logical concepts"]({{ site.baseurl }}/concepts/logical/relationships)

   - **4.2. Map Component Relationships**

     - Identify appropriate relationships for your new components
     - Consider how components relate to others within the same model
     - Consider relationships with components in other models

   - **4.3. Create Definitions**
     Codify the relationships you have identified into a Relationship Definition

    See the [Contributing to Relationships]({{site.baseurl}}/project/contributing/contributing-relationships) for detailed instructions.

<!-- ### Instructions for Creating a New Connection

### Managed and Unmanaged Connections

Each Meshery Model can contain one more ConnectionDefinitions (files), each Definition representing one Connection, and also, (as a matter of convenience multiple Connections can be described in the same ConnectionDefinition file).

Connections can be:

1. a ConnectionDefinition based Meshery's [Connection Schema](https://github.com/meshery/schemas/) with hand-curated Connection attributes.
2. a custom ConnectionDefinition based Meshery's Connection Schema that references an existing Component within the same Model. -->

## Next Steps

The Meshery team is currently working on the following:

- Extending the model to support additional entities
- Improving the tooling for working with models
- Defining relationships between components and embedding those policies within models

We encourage you to get involved in the development of Meshery Models and to share your feedback.
{% include alert.html type="info" title="Meshery Models are extensible" content="Meshery Models are designed to be extensible, allowing you to define new components as needed. If you have an idea for a new component, please create one and share it with the Meshery community." %}
