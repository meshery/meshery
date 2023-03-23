---
layout: default
title: "Extensibility: UI"
permalink: extensibility/ui
type: Extensibility
abstract: "Meshery offers support for more adapters than any other project or product in the world. Meshery UI has a number of extension points that allow for users to customize their experience with third-party plugins."
language: en
list: include
---
Meshery UI has a number of extension points that allow for users to customize their experience with third-party plugins.

## Designing Custom Components for Meshery Extension Point

The Meshery extension points are the way to extend meshery and derive the more custom use-cases out of it. We already have an extension point called MeshMap. 
Meshery can provide the extension point in various ways by providing the feature of custom-components. One of the example where these custom components are used is RJSF forms in meshery-extension

### Extensibility: RJSF Custom Component
RJSFWrapperComponent are the customizations done on RJSF forms overriding the default behaviour of meshery-ui rjsf forms.
The [Rjsf forms are wrapped](https://github.com/meshery/meshery/blob/0bc68d1cd0ba80a565afa68bce80899c22db9a2e/ui/components/MesheryMeshInterface/PatternService/RJSF.js#L66) under these component to receive the custom-props from Meshery-extension.
{% capture code_content %} <RJSFWrapperComponent {...props}>
      <RJSFForm
        isLoading={isLoading}
        schema={schema}
        data={data}
        onChange={(e) => {
          setData(e.formData)
        }}
        jsonSchema={jsonSchema}
      />
 </RJSFWrapperComponent>{% endcapture %}
{% include code.html code=code_content %}
These props are received in the RJSF forms like this: [RJSF Component](https://github.com/meshery/meshery/blob/0bc68d1cd0ba80a565afa68bce80899c22db9a2e/ui/components/MesheryMeshInterface/PatternService/RJSF.js#L91) 

#### Passing new custom prop to forms:

 1. Pass the new prop from the Meshery Extension in the RJSF Wrapper component used like this:

{% capture code_content %} function RJSFWrapperComponent(props) {
  // Clone the child to pass in additional props
  const children = React.cloneElement(props.children, {
    ...(props.children?.props || {}),
    customComponent: YOUR_NEW_CUSTOM_COMPONENT_OR_PROP
  });

  return children
}{% endcapture %}
{% include code.html code=code_content %}
Extract the props in the [RJSFForm Component](https://github.com/meshery/meshery/blob/0bc68d1cd0ba80a565afa68bce80899c22db9a2e/ui/components/MesheryMeshInterface/PatternService/RJSF.js#L91)

_Extensibility documentation missing?_
Submit an issue to request more documentation.

### Extension Points by File

1. **/ui/components/NavigatorExtension.js** - add custom menu items in Meshery's main navigation menu.
1. **/ui/pages/extension/[component].js** - optionally, define the state of a parent menu item: expanded or collapsed.
1. **ui/remote-component.config.js** - list of Material UI components made available to Remote Providers and their plugins.
1. **ui/components/MesheryMeshInterface/PatternServiceFormCore.js** - PatternServiceFormCore component which decouples the SettingsForm and TraitsForm from their UI representation while keeping the logic coupled. This design lets the Remote Provider amend the design of the components without interfering with Meshery UI's core logic.

### Using React JSON Schema Form
 Meshery exposes a custom RJSF Form component which is capable of generating "Pattern" YAMLs without being extremely opinionated about the UI. This custom component is available at `ui/components/MesheryMeshInterface/PatternServiceFormCore.js`. An example usage of the component which will render the logically coupled `SettingsForm` and `TraitsForm` in a Material UI `TabPanel`:

{% capture code_content %} <PatternServiceFormCore
    formData={formData}
    schemaSet={schemaSet}
    onSubmit={onSubmit}
    onDelete={onDelete}
    reference={reference}
    namespace={namespace}
  >
	{(SettingsForm, TraitsForm) => {
		return (
			<div>
				<TabPanel value={tab} index={0} style={ { marginTop: "1.1rem" } }>
					<SettingsForm />
				</TabPanel>
				<TabPanel value={tab} index={1} style={ { marginTop: "1.1rem" } }>
					<TraitsForm />
				</TabPanel>
			</div>
		);
	}}
  </PatternServiceFormCore>{% endcapture %}
{% include code.html code=code_content %}

 Meshery UI's RJSF form accepts two props:

- `RJSFWrapperComponent`: The wrapper component gets all of the props that are passed to the underlying form, allowing to inspect the props and changing the behaviour based on them.
- `RJSFFormChildComponent`: This component will customize the internals of the RJSF form.

With both of these props, Remote Providers can customize the wrapper and can also customize the body of the form. This allows full customization of the form.

from _ui/components/MesheryMeshInterface/PatternService/index.js_
{% include code.html code="function PatternService({ formData, jsonSchema, onChange, type, onSubmit, onDelete, RJSFWrapperComponent, RJSFFormChildComponent })" %}
