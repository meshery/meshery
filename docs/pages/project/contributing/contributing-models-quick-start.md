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

# Instructions for creating a new Model
Meshery Models are a way to represent the architecture of a system or application. Models are defined in JSON and can be used to visualize the components and relationships between them. This guide will walk you through the process of creating a new model.

Meshery Components are the building blocks of a model. Each component represents a different part of the system or application. Components can be anything from a database to a microservice to a server. Relationships define how components interact with each other. For example, a database component might have a relationship with a microservice component that represents the microservice's dependency on the database.

## Quick Start

### Pre-requisites

1. Install the Meshery CLI by following the [installation instructions](https://docs.meshery.io/installation/).
1. Fork the [meshery/meshery](https://github.com/meshery/meshery) repository.

### Create a Model Definition

1. In your meshery/meshery fork, navigate to `mesheryctl/templates/templates-csvs` directory.
1. Edit the `models.csv` file to include your model definition.
1. [Annotations] For models that contain annotation components, edit the `components.csv` file to include your component definitions.
1. Create a new model by executing the following command:
{% capture code_content %}$ mesheryctl registry generate --directory templates-csvs{% endcapture %}
 {% include code.html code=code_content %}
1. Review the generated components inside of the `server/meshmodel` directory under your model's name.
1. Import your model definition into a Meshery Server.
{% capture code_content %}$ mesheryctl model import -f `server/meshmodel/<model-name>`{% endcapture %}
 {% include code.html code=code_content %}
   1. You can deploy Meshery Server using the following command:
   {% capture code_content %}$ mesheryctl system start{% endcapture %}
   {% include code.html code=code_content %}
   or use the [Meshery Playground](https://playground.meshery.io).
1. Verify that your model is displayed in the Meshery UI under Settings->Registry->Models.



-----



{% include alert.html type="light" title="Model Source Code" content="See examples of <a href='https://github.com/meshery/meshery/tree/master/server/meshmodel'>Models defined in JSON in meshery/meshery</a>." %}

To add or update a model, follow these steps:

1. **Create a Model Definition.** Open the <a href='https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw/edit#'>Meshery Integrations spreadsheet</a>. Create a new row (or comment to suggest a new row) to capture the specific details of your model. As you fill-in model details, referernce each column's notes and comments as instructions and an explanation of their purpose.
2. **Generate Components.** Once you have entered values into the required columns, either execute step 2.a. or 2.b.
   1. Execute the following command to generate components for your model.
{% capture code_content %}$ mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred “${{SPREADSHEET_CRED}}"{% endcapture %}
 {% include code.html code=code_content %}
   1. Ask a maintainer to invoke the [Model Generator workflow](https://github.com/meshery/meshery/actions/workflows/model-generator.yml).
1. **Enhance Component details.** While the default shape for new components is a circle, each component should be considered for its best-fit shape.
2. Review and familiarize with the available set of predefined relationship types. Refer to the Cytoscape [node types](https://js.cytoscape.org/demos/node-types/) for a list of possible shapes.
3. Propose a specific shape, best-suited to visually represent the Component. _Example - Deployment as a pentagon._
4. Proposee a specific icon, best-suited to visually represent the Component. _Example - DaemonSet as a skull icon._

{% include alert.html type="info" title="Using Meshery CLI with the Registry (models)" content="Create new and list existing models by using <code>mesheryctl registry</code> to interact with the Meshery Registry and the <a href='https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw/edit#'>Meshery Integrations spreadsheet</a>." %}

### Instructions for creating a new Component

See the [Contributing to Components]({{site.baseurl}}/project/contributing/contributing-components) for detailed instructions.

### Instructions for creating a new Relationship

See the [Contributing to Relationships]({{site.baseurl}}/project/contributing/contributing-relationships) for detailed instructions.

{% include alert.html type="info" title="Generating Models does not require Meshery Server" content="Meshery Server is not required to generate models. The Meshery CLI can be used to generate models. Model and Component generation logic is MeshKit. `mesheryctl` and Meshery Server both utilize MeshKit’s libraries for ongoing programmatic generation of models and components." %}


<!-- ### Instructions for Creating a New Connection

### Managed and Unmanaged Connections

Each Meshery Model can contain one more ConnectionDefinitions (files), each Definition representing one Connection, and also, (as a matter of convenience multiple Connections can be described in the same ConnectionDefinition file).

Connections can be:

1. a ConnectionDefinition based Meshery's [Connection Schema](https://github.com/meshery/schemas/) with hand-curated Connection attributes.
2. a custom ConnectionDefinition based Meshery's Connection Schema that references an existing Component within the same Model. -->

## Next Steps

The Meshery team is currently working on the following:

* Extending the model to support additional entities
* Improving the tooling for working with models
* Defining relationships between components and embedding those policies within models

We encourage you to get involved in the development of Meshery Models and to share your feedback.
  
  {% include alert.html type="info" title="Meshery Models are extensible" content="Meshery Models are designed to be extensible, allowing you to define new components as needed. If you have an idea for a new component, please create one and share it with the Meshery community." %}
