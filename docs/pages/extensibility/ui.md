---
layout: default
title: "Extensibility: UI"
permalink: extensibility/ui
redirect_from: extensibility/ui/
type: Extensibility
abstract: "Meshery offers support for more adapters than any other project or product in the world. Meshery UI has a number of extension points that allow for users to customize their experience with third-party plugins."
language: en
list: include
---
Meshery UI has a number of extension points that allow for users to customize their experience with third-party plugins.

## Designing custom plugins for Meshery UI

_documentation coming soon..._
<!-- Description of the process for designing plugins -->

### Extension Points by File

1. **/ui/components/NavigatorExtension.js** - add custom menu items in Meshery's main navigation menu.
1. **/ui/pages/extension/[component].js** - optionally, define the state of a parent menu item: expanded or collapsed.
1. **ui/remote-component.config.js** - list of Material UI components made available to Remote Providers and their plugins.
1. **ui/components/MesheryMeshInterface/PatternServiceFormCore.js** - PatternServiceFormCore component which decouples the SettingsForm and TraitsForm from their UI representation while keeping the logic coupled. This design lets the Remote Provider amend the design of the components without interfering with Meshery UI's core logic.

### Using React JSON Schema Form
 Meshery exposes a custom RJSF Form component which is capable of generating "Pattern" YAMLs without being extremely opiniated about the UI. This custom component is available at `ui/components/MesheryMeshInterface/PatternServiceFormCore.js`. An example usage of the component which will render the logically coupled `SettingsForm` and `TraitsForm` in a Material UI `TabPanel`:

 ```js
  <PatternServiceFormCore
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
  </PatternServiceFormCore>
 ```

 Meshery UI's RJSF form accepts two props:

- `RJSFWrapperComponent`: The wrapper component gets all of the props that are passed to the underlying form, allowing to inspect the props and changing the behaviour based on them.
- `RJSFFormChildComponent`: This component will customize the internals of the RJSF form.

With both of these props, Remote Providers can customize the wrapper and can also customize the body of the form. This allows full customization of the form.

from _ui/components/MesheryMeshInterface/PatternService/index.js_
```
function PatternService({ formData, jsonSchema, onChange, type, onSubmit, onDelete, RJSFWrapperComponent, RJSFFormChildComponent }) 
```
