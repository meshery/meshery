---
layout: page
title: Contributing to Model Components
permalink: project/contributing/contributing-components
redirect_from: project/contributing/contributing-components/
abstract: How to contribute to Meshery Model Components
language: en
type: project
category: contributing
list: include
---

In Meshery, a [Components](/concepts/logical/components) is a fundamental building block used to represent and define the infrastructure under management. Each component provides granular and specific support for your infrastructure and applications. Once registered with Meshery Server (in the [Registry](/concepts/logical/registry)), components are available for inclusion in [Designs](/concepts/logical/designs) that you create. Components can be created and published by anyone, allowing you to share you custom extensions with the community.

## Overview of Steps to Create Components

<!-- **Prework:**

1. [Component Identification](#component-identification)
2. [Component Classification](#component-visualizations) -->

**Development:**
1. [Component Definition](#component-definitions)
2. [Component Customization](#component-customization)

**Postwork:**
3. [Component Testing](#component-testing)
4. [Component Contribution](#component-contribution)

<!-- ## Prework

<a id="component-identification"></a>

### 1. Characterize the component

#### Understand Component Packaging -->

<a id="component-customization"></a>

### 1. Customize Component Form-based Representation

Each component has any number of customizable properties in the form of metadata. The metadata provides additional details about the component, enhancing its capabilities. Metadata can be attached to components to customize their behavior. The metadata can be used to define the component's behavior, appearance, and interactions with other components.

Components optionally have a UI schema, which is a JSON object that describes how a form should be rendered. It can be passed to JSON forms and is categorized into controls or layouts. The UI schema object follows the form field hierarchy's tree structure and defines how each property should be rendered. For example, it can specify the order of controls, their visibility, and the layout.

Some UI schema elements have an options property that allows for further configuration of the rendering result. These configuration options are often renderer specific and need to be looked up. Some commonly used options include:

`ui:order`: An array of field names in the order in which they should appear
`ui:widget`: The name of an alternative widget to use for the field
`ui:field`: The name of a custom field
`classNames`: The class names to put on the component

Meshery UI uses the `react-jsonschema-form` library to render forms. See the [RJSF documentation](https://rjsf-team.github.io/react-jsonschema-form/docs/) for more information.

### 2. Customize Component Visual Representation

PrimaryColor
SecondaryColor
LogoColor
LogoWhite

## Development

<a id="component-definitions"></a>

### 3. Create a Component Definition as a JSON file

Create a relationship definition as a JSON file, placing this new definition file into its respective model folder (see [Contributing to Models](./models)). Relationship definition files are commonly named  `relationships.yaml` as a convention, however, this name is not required. A model may include any number of relationship definitions. Include the following attributes in your relationship definition:

- `kind`: The genre of component (e.g., `Pod`).
- `model`: The model to which the component belongs (e.g., `kubernetes`).
- `version`: The version of the component definition (e.g., `v1.0.0`).
- `description`: A characterization of the component, its purpose and behavior.

{% include alert.html title="Use Existing Components as Examples" type="info" content="Browse the <a href='https://github.com/meshery/meshery/tree/master/server/meshmodel'>existing components in the Meshery repository</a> to find examples of how to existing component, using them as a template." %}

## Postwork

<a id="component-testing"></a>

#### 4. Component Authoring Best Practices and Considerations

##### General

1. Use camelCasing as the formatting convention.
2. Increment the component definition using the `version` property with each change published.

<a class="anchorjs-link" id="component-contribution"></a>

#### 5. Contribute your component to the Meshery project

Submit a pull request to the Meshery repository with your new component definition, so that all users can benefit from the component(s) you have defined.

Keeping your component definition in a separate file allows for easier management and review of the component(s) you have defined.

{% include alert.html title="Keeping your custom Components private" type="info" content="Alternatively, if you would like to keep the component definition private, you can bundle your component(s) in a custom model, import the custom model into your Meshery deployment. Your private component definition(s) will be registered in your Meshery Server's <a href='/concepts/logical/registry'>registry</a> and available for use within your Meshery deployment." %}
