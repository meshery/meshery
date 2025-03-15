---
layout: default
title: Generating Models
abstract: A comprehensive guide on generating models in Meshery, covering both CLI and UI methods.
permalink: guides/configuration-management/generating-models
category: configuration
type: guides
language: en
---

Meshery lets you create models by processing Custom Resource Definitions (CRDs) or importing existing resources. Models can be generated from URLs or CSV files and are classified as either **Static Models** (pre-defined with each release) or **Dynamic Models** (created at runtime from external sources). This guide explains how to generate models using the Meshery CLI and UI, and clarifies the differences between the two primary commands:

- **`mesheryctl model generate`** – Generates models locally.
- **`mesheryctl registry generate`** – Generates models and registers them into the Meshery Registry.

## Generate Models Using Meshery CLI

Meshery provides two distinct CLI commands for model generation:

### `mesheryctl model generate`
  - **What It Does:** Creates models from a file, directory, or URL.
  - **Where It Stores Models:** Locally, typically in the `.meshery/models` directory.
  - **When to Use:** When you want to generate models without immediate registration.

### `mesheryctl registry generate`
  - **What It Does:** Generates models and automatically registers them into the Meshery Registry.
  - **Where It Stores Models:** Under `/server/models/<model-name>`, with logs saved in `~/.meshery/logs/registry`.
  - **When to Use:** When you need the models to be immediately available for Meshery’s runtime use.

> **Note:** In the examples below, we demonstrate the usage of `mesheryctl model generate` for simplicity. 

### Steps to Generate Models

### 1. Install Meshery CLI

Before using the CLI, ensure it is installed by following the [installation instructions](https://docs.meshery.io/installation).

### 2. Generate the Model

Models can be generated in two formats: **URL** and **CSV**. When generating from a URL, a template JSON file is required. This file must include fields such as *Registrant, Model Name, DisplayName,* and *Category*. The template also defines whether a component should be treated as a regular component or simply as an annotation (for example, an SVG icon).

#### 2.1 Using URL as Input

To generate a model from a CRD URL, run:

```sh
mesheryctl model generate -f "git:://github.com/cert-manager/cert-manager/master/deploy/crds" -t template.json
```

If you want to skip automatic registration (i.e. only generate and store the model locally), add the `-r` flag:

```sh
mesheryctl model generate -f "git:://github.com/cert-manager/cert-manager/master/deploy/crds" -t template.json -r
```

When generating models from a URL, Meshery supports the following sources:

- **ArtifactHub:** `mesheryctl model generate -f "https://artifacthub.io/packages/search?ts_query_web={model-name}"`
  
- **GitHub:** `mesheryctl model generate -f "git:://github.com/cert-manager/cert-manager/master/deploy/crds"`

> **Note:** The order and format of the URL are important. The `-t` flag points to the template file that maps required fields. If the template’s `isAnnotation` field is set to `true`, the component is treated as an annotation rather than a standard component.

#### 2.2 Using CSV Files as Input

When using CSV files, ensure your directory includes:
- A **model CSV** file (with model definitions)
- A **components CSV** file (with component details)
- Optionally, a **relationships CSV** file

The CLI will parse these files and generate models accordingly.

## Generate Models Using Meshery UI

### 1. Access the Meshery UI

Ensure Meshery is installed by following the [Quick Start instructions]({{site.baseurl}}/installation/quick-start).

### 2. Navigate to the Registry Section

- Click the **Settings** icon (top right).
- Select **Registry** to manage and generate models.

<a href="{{ site.baseurl }}/assets/img/export/Registry.png">
  <img alt="Registry Navigator" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/export/Registry.png" />
</a>

### 3. Generate the Model

On the **Registry** page, click the **Generate** button. You can choose one of two methods:

#### 3.1 From CSV

- Upload CSV Files: Upload your components CSV file and, optionally, a relationships CSV file.  

  [![CSV Template](/assets/img/generate/CsvTemplate.gif)](/assets/img/generate/CsvTemplate.gif)

- Use the Spreadsheet Template: If you don’t have a CSV file, use our [Spreadsheet template](https://docs.google.com/spreadsheets/d/19JEpqvHrG8UL-Bc-An9UIcubf1NVhlfnQSN1TD7JOZ4/).  Fill in your details, download it as a CSV, and upload it.

  [![Generate From CSV](/assets/img/generate/GenerateFromCsv.gif)](/assets/img/generate/GenerateFromCsv.gif)

#### 3.2 From URL

- **Paste the URL:**Enter the URL for your model.
- **Fill in the Details:**Provide the required model details when prompted.
- **Specify Options:**Indicate if the model should be treated as an annotation or if it should be registered immediately.

## Understanding How Meshery Generates Models

Meshery employs both manual and automated techniques to generate models for its Registry. Here’s how it works:

1. **Parsing Input:**The CLI reads data from a file, directory, or URL. A template file is used in URL mode to map required fields.
2. **Validating Data:**The process verifies that all necessary fields (Registrant, Model Name, DisplayName, Category) are present.
3. **Generating and Registering:**
   - **Local Generation (`model generate`):** Models are created and stored locally.
   - **Registry Generation (`registry generate`):** Models are created and then registered in the Registry, making them available to Meshery.
4. **Logging:**Detailed logs and error messages are recorded to help troubleshoot any issues during generation.

#### What Happens If Models Are Not Registered?

If you generate models with the `-r` flag (to skip registration) using `mesheryctl model generate`, the models are stored locally but Meshery will not recognize them. For Meshery to use the models, they must be registered in the Registry. You can later [import](https://docs.meshery.io/guides/configuration-management/importing-models) them with:

```sh
mesheryctl model import -f <path-to-model>
```