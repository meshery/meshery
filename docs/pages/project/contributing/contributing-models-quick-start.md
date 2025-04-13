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

[Meshery Models](/concepts/logical/models) help you visually represent the architecture of your system or application. Models—defined in JSON—enable you to see and manage how different components work together. This guide provides clear, step-by-step instructions to create your first model.

[Meshery Components](/concepts/logical/components) are the essential building blocks of a model. Each component can represent anything—from databases and microservices to entire servers. Using relationships, you can clearly indicate how these components depend on or interact with one another.

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
        <i class="fa fa-link"></i> URL Create
      </label>
      <section class="tabbed">
        <p>The URL Import feature allows you to create models in Meshery by providing URLs to source repositories or package registries. Here's how to use it:</p>

        <h4>Step 1: Open the Model Generation Interface</h4>
        <p>Go to the <a href="https://playground.meshery.io/settings?settingsCategory=Registry&tab=Models">Registry</a> section in the Meshery UI. Click the <strong>create</strong> button to start creating a new model.</p>

        <a href="/assets/img/registry/generate-model-from-UI.gif">
          <img src="/assets/img/registry/generate-model-from-UI.gif" alt="registry generate model" style="width: 50%; max-width: 400px;">
        </a>

        <h4>Step 2: Set Up Your Model Name</h4>
        <p>Enter the following names for your model:</p>
        <ul>
          <li><strong>Model Name:</strong> This is the unique, system identifier. Use lowercase letters and hyphens only (no spaces). <br>Example: <code>cert-manager</code></li>
          <li><strong>Display Name:</strong> This is the friendly name displayed in the UI. It can include spaces and uppercase letters. <br>Example: <code>Cert Manager</code></li>
        </ul>

        <a href="/assets/img/registry/model-create-name.png">
          <img src="/assets/img/registry/model-create-name.png" alt="registry generate model name" style="width: 50%; max-width: 400px;">
        </a>

        <h4>Step 3: Set Model Categorization</h4>
        <p>Select the appropriate <strong>Category</strong> and <strong>Subcategory</strong> from the dropdowns. If your model doesn't match any available category, choose Uncategorized to ensure it can be easily located later.</p>

        <a href="/assets/img/registry/model-create-category.png">
          <img src="/assets/img/registry/model-create-category.png" alt="registry generate model categorization" style="width: 50%; max-width: 400px;">
        </a>

        <h4>Step 4: Configure Model Styling</h4>
        <p>Customize your model's appearance:</p>
        <ul>
          <li><strong>Logos:</strong> Upload separate images for dark and light themes.</li>
          <li><strong>Colors:</strong> Define the primary and secondary colors for your model's visual elements.</li>
          <li><strong>Icon Shape:</strong> Choose a shape for the model’s icon as displayed in the UI.</li>
        </ul>
        <p>Note: If you don't provide custom styling, Meshery's default values will be used. You can change these later in the model definition.</p>

        <a href="/assets/img/registry/model-create-styling.png">
          <img src="/assets/img/registry/model-create-styling.png" alt="registry generate model styling" style="width: 50%; max-width: 400px;">
        </a>

        <h4>Step 5: Provide the Source Location</h4>
        <p>Specify the source for your model by entering one of the following URL types:</p>

        <p><strong>GitHub Repository:</strong><br>
        Enter a URL using the following format: <code>git://github.com/[organization]/[repository]/[branch]/path/to/crds</code>.<br>
        Example: <code>git://github.com/cert-manager/cert-manager/master/deploy/crds</code></p>

        <p><strong>ArtifactHub Package:</strong><br>
        Enter a URL in this format: <code>https://artifacthub.io/packages/search?ts_query_web={model-name}</code></p>

        {% include alert.html type="light" title="Need a Source URL?" content="Check the <code>sourceURL</code> column in the <a href='https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw' target='_blank' rel='noopener'>Meshery Integration Sheet</a> and try one of the listed sources." %}

        <a href="/assets/img/registry/model-create-source.png">
          <img src="/assets/img/registry/model-create-source.png" alt="registry generate model source selection" style="width: 50%; max-width: 400px;">
        </a>

        <h4>Step 6: Additional Options</h4>
        <ul>
        <li><strong>Immediate Registration:</strong> Enable this option to register your model instantly, making it available in your Meshery instance right away.</li>
        <li><strong>Visual Annotation Only:</strong> Mark this option if your model is intended solely for visual diagramming rather than infrastructure management.</li>
        </ul>

        {% include alert.html type="light" title="Visual Annotation Models" content="When a model is marked for visual annotation only, it means the model will be used purely for visualization and diagramming purposes within Meshery's interface, rather than for actual infrastructure management." %}

        <a href="/assets/img/registry/model-create-additional.png">
          <img src="/assets/img/registry/model-create-additional.png" alt="registry generate model additional settings" style="width: 50%; max-width: 400px;">
        </a>

        <h4>Step 7: Finalize and Create Your Model</h4>
        <p>Review your model details. If everything looks good, click "Generate" to create your model.</p>
        <p>Meshery will process the provided source URL, extract the necessary information, and create the model based on the specified details.</p>
        <p>If you want to make any changes before generation, you can go back to the previous steps using the navigation buttons and edit the details as needed.</p>

        <a href="/assets/img/registry/model-create-confirmation.png">
          <img src="/assets/img/registry/model-create-confirmation.png" alt="registry generate model generate" style="width: 50%; max-width: 400px;">
        </a>

        <h4>Step 8: Check Model Generation Status</h4>
        <p>After generation, a confirmation message will appear if your model is created successfully. You will be able to view all the components and their relationships.</p>
        <p>If an error occurs, an error message will display with details on what went wrong. You can then return to previous steps to adjust your model details or update the source URL.</p>
        <p>If model generation is successful, you will see a full model and component details in registry page. You can also view the generated model in the Kanvas.</p>

        <a href="/assets/img/registry/model-create-result.png">
          <img src="/assets/img/registry/model-create-result.png" alt="registry generate model success" style="width: 50%; max-width: 400px;">
        </a>

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
