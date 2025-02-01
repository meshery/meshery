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

### Prerequisites

1. Fork the [meshery/meshery](https://github.com/meshery/meshery) repository.
2. Install the Meshery CLI by following the [installation instructions](https://docs.meshery.io/installation/).

### Create a Model Definition

There are two approaches to creating Meshery Model definitions.

### Using Meshery Integration Spreadsheet

The Meshery project maintains a central Integration Spreadsheet that serves as the source of truth for all models. While this spreadsheet is publicly viewable at [Meshery Integration Sheet](https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw), direct editing access is restricted to maintain data integrity.

To work with the spreadsheet approach, you'll need to set up your way to access the sheet. Here is how:

1.  **Creating Your Spreadsheet Copy**: First, create your own copy of the [Meshery Integration Sheet](https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw)
2.  **Setting Up Google Cloud Project**: To access your spreadsheet programmatically, you'll need to set up a Google Cloud Project. [How to?](https://developers.google.com/workspace/guides/create-project)
3.  **Enabling the Google Sheets API**: Enable API access for your project. [How to?](https://support.google.com/googleapi/answer/6158841)
4.  **Creating Service Account Credentials**: [How to?](https://developers.google.com/workspace/guides/create-credentials#create_credentials_for_a_service_account)
5. **Preparing Credentials**: Convert your credentials to the required format:

```bash
# Convert the downloaded JSON credentials to base64
base64 -w 0 /path/to/your-service-account-creds.json

# Add the output to your shell configuration file
echo 'export SHEET_CRED="<paste-base64-output-here>"' >> ~/.bashrc  # or ~/.zshrc for Zsh users

# Reload your shell configuration
source ~/.bashrc  # or source ~/.zshrc for Zsh users
```

6. **Configuring Spreadsheet Access**: Grant access to your service account:

    -   Copy the service account email address (it ends with @developer.gserviceaccount.com)
    -   Open your copy of the Meshery Integration Sheet
    -   Click the "Share" button
    -   Paste the service account email
    -   Set the role to "Editor" and "Share"
7. **Publishing Your Spreadsheet**: Make your spreadsheet accessible:

    -   In your spreadsheet, click `File > Share > Publish` to web
    -   Choose "Entire Document" and "Comma-separated values (.csv)"
    -   Click "Publish"
    -   Keep note of your spreadsheet ID (it's the long string in your spreadsheet's URL between /d/ and /edit)
	    - Example: `https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw/edit`
	    - spreadsheet ID: `1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw`

### Using Local CSV Files

For a simpler approach that requires no additional setup, you can use local CSV files. This method is particularly useful for initial development and testing.

In your Meshery fork, you can find the CSV templates at `mesheryctl/templates/templates-csvs`. The directory contains three essential files:

-   `models.csv`: Defines model metadata and basic properties
-   `components.csv`: Specifies component definitions and their attributes
-   `relationships.csv`: Describes connections between components

Each file follows the same structure as the spreadsheet, but can be edited locally.

## Generating Models

### Using Meshery CLI

#### For Spreadsheet-Based Definitions:

If you set up the Google Sheets integration, use the following command:

```bash
mesheryctl registry generate --spreadsheet-id "<your-spreadsheet-id>" --spreadsheet-cred $SHEET_CRED
```

for specific model:

```bash
mesheryctl registry generate --spreadsheet-id "<your-spreadsheet-id>" --spreadsheet-cred $SHEET_CRED --model "<model-name-in-csv"
```

#### For Local CSV-Based Definitions:

If you're working with local CSV files, use:

```bash
mesheryctl registry generate --directory templates-csvs
```

The generated model will be placed in `server/meshmodel` directory under your model's name

### Using Meshery UI

The Meshery UI provides a visual interface for model generation. Here's how to use it:

-   Navigate to `Settings → Registry → Models`
-   Click "Generate" button
- In the generation dialog, select "Upload CSV Files"
-   You can now upload the three template CSV files from `mesheryctl/templates/templates-csvs`

Verify:

-   Navigate to `Settings → Registry → Models`
-   Look for your model in the list
-   Check that all components and relationships are present
-   Verify that properties and metadata are correctly displayed

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
