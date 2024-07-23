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

**Prework:**

<!-- 1. [Component Identification](#component-identification) -->
1. [Component Classification](#component-classification) 

**Development:**
2. [Component Definition](#component-definitions)
3. [Component Customization](#component-customization)

**Postwork:**
4. [Component Testing](#component-testing)
5. [Component Contribution](#component-contribution)

## Prework

<a id="component-classification"></a>

### 1. Understand Model Generation and Packaging

{% include alert.html title="Prerequisite Reading" type="warning" content="Components are defined and packaged in context of a Model. Be sure to first read and understand <a href='/project/contributing/contributing-models'>how models are created and packaged</a> before attempting to create a new component. Without a model defined first, they component will be homeless." %}

<a id="component-customization"></a>

### 2. Customize Component Form-based Representation

Each component has any number of customizable properties in the form of metadata. The metadata provides additional details about the component, enhancing its capabilities. Metadata can be attached to components to customize their behavior. The metadata can be used to define the component's behavior, appearance, and interactions with other components.

#### Customize Component Form-based Representation

Components optionally have a UI schema, which is a JSON object that describes how a form should be rendered. It can be passed to JSON forms and is categorized into controls or layouts. The UI schema object follows the form field hierarchy's tree structure and defines how each property should be rendered. For example, it can specify the order of controls, their visibility, and the layout.

Some UI schema elements have an options property that allows for further configuration of the rendering result. These configuration options are often renderer specific and need to be looked up. Some commonly used options include:

`ui:order`: An array of field names in the order in which they should appear
`ui:widget`: The name of an alternative widget to use for the field
`ui:field`: The name of a custom field
`classNames`: The class names to put on the component

Meshery UI uses the `react-jsonschema-form` library to render forms. See the [RJSF documentation](https://rjsf-team.github.io/react-jsonschema-form/docs/) for more information.

#### Customize Component Visual Representation

Meshery contributors who want to customize the visual representation of a Meshery component can do so by modifying the component's metadata. This metadata includes fields for specifying the component's:

* **SVG**: The SVG file used to represent the component visually.
* **Primary Color**: The primary color used for the component in hex format.
* **Secondary Color**: The secondary color used for the component in hex format.
* **Shape**: The basic shape of the component, such as a circle or square.
* **isAnnotation**: A boolean value that indicates whether the component is an annotation or is [semantically meaningful]({{site.baseurl}}/concepts/logical/components).

<details>
<summary>Example Component Metadata</summary>
<pre><code class="highlight-json">
"metadata": {
   "capabilities": "",
   "defaultData": "",
   "genealogy": "",
   "isAnnotation": false,
   "isCustomResource": false,
   "isModelAnnotation": "FALSE",
   "isNamespaced": false,
   "logoURL": "Created by Meshery Authors",
   "model": "kubernetes",
   "modelDisplayName": "Kubernetes",
   "primaryColor": "#326CE5",
   "secondaryColor": "#7aa1f0",
   "shape": "round-rectangle",
   "shapePolygonPoints": "",
   "status": "enabled",
   "styleOverrides": "{\"height\":\"22px\", \"width\":\"22px\", \"x\":\"8.5\", \"y\":\"7.5\"}",
   "styles": "{\"height\":\"22px\", \"width\":\"22px\", \"x\":\"8.5\", \"y\":\"7.5\"}",
   "subCategory": "",
   "svgColor": "\u003csvg width=\"90\" height=\"90\" viewBox=\"0 0 90 90\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\"\u003e\n\u003cpath fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M41.0114 45.015C51.699 45.015 60.363 36.3512 60.363 25.6637C60.363 14.9763 51.699 6.3125 41.0114 6.3125C30.3239 6.3125 21.6599 14.9763 21.6599 25.6637C21.6599 36.3512 30.3239 45.015 41.0114 45.015ZM26.9841 46.0338C27.3733 45.9771 27.7668 46.0783 28.0996 46.2881C31.5844 48.484 36.0555 49.8036 40.931 49.8036C45.8346 49.8036 50.3289 48.4688 53.8223 46.2502C54.148 46.0434 54.5325 45.9405 54.9151 45.9897C64.2686 47.1911 71.4966 55.1832 71.4966 64.8635V72.3594C71.4966 78.7567 66.3106 83.9428 59.9132 83.9428H22.308C15.9107 83.9428 10.7246 78.7567 10.7246 72.3594V64.8635C10.7246 55.2946 17.7872 47.3752 26.9841 46.0338Z\" fill=\"#326CE5\"/\u003e\n\u003cpath fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M72.0356 82.4937C73.7201 80.3512 74.7249 77.649 74.7249 74.7122V64.0123C74.7249 55.9577 70.1245 48.9776 63.408 45.5561C72.7393 47.0963 79.8167 54.6562 79.7238 63.5783V72.0106C79.7238 76.766 76.5504 80.8211 72.0356 82.4937ZM63.6189 25.6633C63.6189 30.5769 62.1331 35.1434 59.5862 38.9381C66.2233 36.7133 71.0055 30.4442 71.0055 23.0581C71.0055 13.8103 63.5088 6.31348 54.261 6.31348C54.0822 6.31348 53.904 6.31628 53.7266 6.32184C59.7194 10.6542 63.6189 17.7034 63.6189 25.6633Z\" fill=\"#326CE5\"/\u003e\n\u003c/svg\u003e\n",
   "svgComplete": "",
   "svgWhite": "\u003csvg width=\"90\" height=\"90\" viewBox=\"0 0 90 90\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\"\u003e\n\u003cpath fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M41.0114 45.015C51.699 45.015 60.363 36.3512 60.363 25.6637C60.363 14.9763 51.699 6.3125 41.0114 6.3125C30.3239 6.3125 21.6599 14.9763 21.6599 25.6637C21.6599 36.3512 30.3239 45.015 41.0114 45.015ZM26.9841 46.0338C27.3733 45.9771 27.7668 46.0783 28.0996 46.2881C31.5844 48.484 36.0555 49.8036 40.931 49.8036C45.8346 49.8036 50.3289 48.4688 53.8223 46.2502C54.148 46.0434 54.5325 45.9405 54.9151 45.9897C64.2686 47.1911 71.4966 55.1832 71.4966 64.8635V72.3594C71.4966 78.7567 66.3106 83.9428 59.9132 83.9428H22.308C15.9107 83.9428 10.7246 78.7567 10.7246 72.3594V64.8635C10.7246 55.2946 17.7872 47.3752 26.9841 46.0338Z\" fill=\"white\"/\u003e\n\u003cpath fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M72.0356 82.4937C73.7201 80.3512 74.7249 77.649 74.7249 74.7122V64.0123C74.7249 55.9577 70.1245 48.9776 63.408 45.5561C72.7393 47.0963 79.8167 54.6562 79.7238 63.5783V72.0106C79.7238 76.766 76.5504 80.8211 72.0356 82.4937ZM63.6189 25.6633C63.6189 30.5769 62.1331 35.1434 59.5862 38.9381C66.2233 36.7133 71.0055 30.4442 71.0055 23.0581C71.0055 13.8103 63.5088 6.31348 54.261 6.31348C54.0822 6.31348 53.904 6.31628 53.7266 6.32184C59.7194 10.6542 63.6189 17.7034 63.6189 25.6633Z\" fill=\"white\"/\u003e\n\u003c/svg\u003e\n"
  },
</code></pre>
</details>

Contributors can also override these settings for individual components if needed. For example, if a component has a different SVG color than its model, you can specify a custom SVG color for that component.

Component icons will be written to the `/meshmodel/components/<model-name>/icon/...` directory upon registration.

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
