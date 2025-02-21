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
1. Install the Meshery CLI by following the [installation instructions](https://docs.meshery.io/installation/).

### Create a Model Definition

<ol>

<li>In your meshery/meshery fork, open the template <a href="https://docs.google.com/spreadsheets/d/19JEpqvHrG8UL-Bc-An9UIcubf1NVhlfnQSN1TD7JOZ4/edit?gid=1308482975#gid=1308482975">spreadsheet</a>.</li>
<ul><li>Or, optionally, use CSVs by navigating to <code>mesheryctl/templates/templates-csvs</code> directory.</li></ul>
<li>Edit the <code>models sheet</code> file to include your model definition.
   <ol>
      <li>[Annotation-only Components] For models that contain annotation components, edit the <code>components sheet</code> file to include your component definitions.</li>
   </ol>
</li>
<li>Generate your new models using a Meshery Server (optionally, use the <a href="https://playground.meshery.io">Meshery Playground</a>). Click the "Generate" button after navigating to Meshery UI --> Settings --> Registry.</li>
<li> 

<details>
   <summary>Or use Meshery CLI</summary>
   Or, opttionally, use Meshery CLI by executing the following command:<code>mesheryctl registry generate --directory templates-csvs</code>
   Review the generated components inside of the <code>server/meshmodel</code> directory under your model's name.

</details></li>
<li> Verify that your model is displayed in the Meshery UI under Settings->Registry->Models.</li>

</ol>

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

