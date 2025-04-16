---
layout: default
title: Create Models
abstract: A comprehensive guide on creating models in Meshery, covering both CLI and UI methods.
permalink: guides/configuration-management/creating-models
redirect_from: guides/configuration-management/generating-models
category: configuration
type: guides
language: en
---

Meshery offers two primary ways to bring models into your system: **Create** and **[Import](/guides/configuration-management/importing-models)**.

- **Import** is ideal when you already have a model definition file (e.g., JSON, CSV) and simply want to bring it into Meshery.
- **Create** is the recommended approach if you’re starting fresh — it guides you through building a model step by step using an intuitive UI wizard.

This section walks you through the **Create** workflow — perfect for users who want full control over naming, styling, and categorizing their models while sourcing CRDs from GitHub or ArtifactHub.

> Tip: If you don’t have any files yet, start with **Create**. You can always [export](/guides/configuration-management/exporting-models) or re-import your model later.

<div class="tab-container">

  <!-- First Level Tab: Meshery UI -->
  <input type="radio" id="tab1" name="tabs" checked>
  <label for="tab1">
    <i class="fa fa-desktop"></i> Meshery UI
  </label>
  <section class="tabbed">
    <p>The URL Create feature allows you to create models in Meshery by providing URLs to source repositories or package registries. Here's how to use it:</p>

    <h4>Step 1: Open the Model Generation Interface</h4>
    <p>Go to the <a href="https://playground.meshery.io/settings?settingsCategory=Registry&tab=Models">Registry</a> section in the Meshery UI. Click the <strong>create</strong> button to start creating a new model.</p>

    <a href="/assets/img/registry/create-ui-where.gif">
      <img src="/assets/img/registry/create-ui-where.gif" alt="registry generate model" style="width: 50%; max-width: 400px;">
    </a>

    <h4>Step 2: Set Up Your Model Name</h4>
    <p>Enter the following names for your model:</p>
    <ul>
      <li><strong>Model Name:</strong> This is the unique, system identifier. Use lowercase letters and hyphens only (no spaces). <br>Example: <code>cert-manager</code></li>
      <li><strong>Display Name:</strong> This is the friendly name displayed in the UI. It can include spaces and uppercase letters. <br>Example: <code>Cert Manager</code></li>
    </ul>

    <a href="/assets/img/registry/create-ui-name.png.png">
      <img src="/assets/img/registry/create-ui-name.png" alt="registry generate model name" style="width: 50%; max-width: 400px;">
    </a>

    <h4>Step 3: Set Model Categorization</h4>
    <p>Select the appropriate <strong>Category</strong> and <strong>Subcategory</strong> from the dropdowns. If your model doesn't match any available category, choose Uncategorized to ensure it can be easily located later.</p>

    <a href="/assets/img/registry/create-ui-category.png">
      <img src="/assets/img/registry/create-ui-category.png" alt="registry generate model categorization" style="width: 50%; max-width: 400px;">
    </a>

    <h4>Step 4: Configure Model Styling</h4>
    <p>Customize your model's appearance:</p>
    <ul>
      <li><strong>Logos:</strong> Upload separate images for dark and light themes.</li>
      <li><strong>Colors:</strong> Define the primary and secondary colors for your model's visual elements.</li>
      <li><strong>Icon Shape:</strong> Choose a shape for the model’s icon as displayed in the UI.</li>
    </ul>
    <p>Note: If you don't provide custom styling, Meshery's default values will be used.</p>

    <a href="/assets/img/registry/create-ui-styling.png">
      <img src="/assets/img/registry/create-ui-styling.png" alt="registry generate model styling" style="width: 50%; max-width: 400px;">
    </a>

    <h4>Step 5: Provide the Source Location</h4>
    <p>Specify the source for your model by entering one of the following URL types:</p>
    <ul>
      <li>
        <strong>GitHub Repository:</strong><br>
        Enter a URL using the following format:  
        <code>git://github.com/[organization]/[repository]/[branch]/path/to/crds</code><br>
        Example:  
        <code>git://github.com/cert-manager/cert-manager/master/deploy/crds</code>
      </li>
      <li>
        <strong>ArtifactHub Package:</strong><br>
        Enter a URL in this format:  
        <code>https://artifacthub.io/packages/search?ts_query_web={model-name}</code>
      </li>
    </ul>

    {% include alert.html type="light" title="Need a Source URL?" content="Check the <code>sourceURL</code> column in the <a href='https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw' target='_blank' rel='noopener'>Meshery Integration Sheet</a> and try one of the listed sources." %}

    <a href="/assets/img/registry/create-ui-source.png">
      <img src="/assets/img/registry/create-ui-source.png" alt="registry generate model source selection" style="width: 50%; max-width: 400px;">
    </a>

    <h4>Step 6: Additional Settings</h4>
    <ul>
    <li><strong>Immediate Registration:</strong> Enable this option to register your model instantly, making it available in your Meshery instance right away.</li>
    <li><strong>Visual Annotation Only:</strong> Mark this option if your model is intended solely for visual diagramming rather than infrastructure management.</li>
    </ul>

    {% include alert.html type="light" title="Visual Annotation Models" content="When a model is marked for visual annotation only, it means the model will be used purely for visualization and diagramming purposes within Meshery's interface, rather than for actual infrastructure management." %}
    <a href="/assets/img/registry/create-ui-details.png">
      <img src="/assets/img/registry/create-ui-details.png" alt="registry generate model additional settings" style="width: 50%; max-width: 400px;">
    </a>

    <h4>Step 7: Finalize and Create Your Model</h4>
    <p>Review your model details. If everything looks good, click "Generate" to create your model.</p>
    <p>Meshery will process the provided source URL, extract the necessary information, and create the model based on the specified details.</p>
    <p>If you want to make any changes before generation, you can go back to the previous steps using the navigation buttons and edit the details as needed.</p>

    <a href="/assets/img/registry/create-ui-finalize.png">
      <img src="/assets/img/registry/create-ui-finalize.png" alt="registry generate model generate" style="width: 50%; max-width: 400px;">
    </a>

    <h4>Step 8: Check Model Generation Status</h4>
    <p>After generation, a confirmation message will appear if your model is created successfully. You will be able to view all the components and their relationships.</p>
    <p>If an error occurs, an error message will display with details on what went wrong. You can then return to previous steps to adjust your model details or update the source URL.</p>
    <p>If model generation is successful, you will see a full model and component details in registry page. You can also view the created model in the Kanvas.</p>

    <a href="/assets/img/registry/create-ui-finish.png">
      <img src="/assets/img/registry/create-ui-finish.png" alt="registry generate model success" style="width: 50%; max-width: 400px;">
    </a>

    <p>Meshery also displays a notification center to confirm the operation. If there are any issues, an error message will appear with details to help you resolve them:</p>

    <a href="/assets/img/registry/create-ui-notification.png">
      <img src="/assets/img/registry/create-ui-notification.png" alt="registry generate model success notification" style="width: 50%; max-width: 400px;">
    </a>
    </section>

  <!-- Second Level Tab: mesheryctl -->
  <input type="radio" id="tab2" name="tabs">
  <label for="tab2">
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

    {% include alert.html type="info" title="Using Meshery CLI with the Meshery Registry and Meshery Models" content="Meshery CLI has a set of commands that pertain to the lifecycle management of models:
    <br />

    - <code>mesheryctl registry</code> - interact with and update spreadsheets
    <br />
    - <code>mesheryctl models</code> - interact with and update Meshery Server
    <br />
    - <code>mesheryctl component</code> - interact with and update Meshery Server
    <br />
    - <code>mesheryctl relationships</code> - interact with and update Meshery Server" %}

      </section>
  </div>