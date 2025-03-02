---
layout: page
title: Contributing to Models Quick Start
permalink: project/contributing/contributing-models-quick-start
redirect_from: project/contributing/contributing-models-models-quick-start/
abstract: A no-fluff guide to creating your own Meshery Models quickly.
language: en
type: project
category: contributing
list: include
---

[Meshery Models](/concepts/logical/models) are a way to represent the architecture of a system or application. Models are defined in JSON and can be used to visualize the components and relationships between them. This guide will walk you through the process of creating a new model.

[Meshery Components](/concepts/logical/components) are the building blocks of a model. Each component represents a different part of the system or application. Components can be anything from a database to a microservice to a server. Relationships define how components interact with each other. For example, a database component might have a relationship with a microservice component that represents the microservice's dependency on the database.

## Creating your first Meshery Model

The following instructions are a no-fluff guide to creating your own Meshery Models quickly. For more detailed information, see the [Contributing to Models](/project/contributing/contributing-models) documentation.

<div class="tab-container">
  <!-- First Level Tab: Meshery UI -->
  <input type="radio" id="tab2" name="tabs" checked>
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
        <p>The CSV Import feature allows you to generate models in Meshery by providing <a href="https://github.com/meshery/meshery/tree/a514f8689260791077bde8171646933cff15dd08/mesheryctl/templates/template-csvs" target="_blank" rel="noopener noreferrer">template CSV files</a> that define your model structure, components, and relationships. Here's a comprehensive guide on how to use this feature:</p>

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

**Congratulations! You have successfully created a new model.**

### Contributing a Model Definition

1. Fork the [meshery/meshery.io](https://github.com/meshery/meshery.io) repository.
1. Create a new branch in your fork of the meshery/meshery.io repository.
1. Add your model definition to the `collections/_models` directory.
1. Create a pull request to the meshery/meshery.io repository.
1. Once your pull request is merged, your model will be available in the next Meshery release.

## Next Steps

{% include alert.html type="info" title="Contributing to Models" content="See the <a href='/project/contributing/contributing-models'>full Contributing to Models</a> documentation for a detailed understanding of models and the many ways in which you can customize them." %}

We encourage you to get involved in the development of Meshery Models and to share your feedback.

{% include alert.html type="info" title="Meshery Models are extensible" content="Meshery Models are designed to be extensible, allowing you to define new components as needed. If you have an idea for a new component, please create one and share it with the Meshery community." %}
