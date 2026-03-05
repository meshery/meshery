---
title: "Extensibility: UI"
description: Meshery offers support for more adapters than any other project or product in the world. Meshery UI has a number of extension points that allow users to customize their experience with third-party plugins.
---

Meshery UI has a number of extension points that allow you to greatly customize both its functional behavior and visual appearance. These extension points come in different types, and this document describes each type, its example use, best practices to consider, and caveats of which to be aware.

### Extensibility: Customizing Text-based Forms using RJSF Custom Component

RJSFWrapperComponent provides customizations for RJSF forms, overriding the default behavior of meshery-ui rjsf forms. The [Rjsf forms are wrapped](https://github.com/meshery/meshery/blob/0bc68d1cd0ba80a565afa68bce80899c22db9a2e/ui/components/MesheryMeshInterface/PatternService/RJSF.js#L66) in this component to receive custom props from a Meshery extension.
{{< code code=`<RJSFWrapperComponent {...props}>
  <RJSFForm
    isLoading={isLoading}
    schema={schema}
    data={data}
    onChange={(e) => {
      setData(e.formData)
    }}
    jsonSchema={jsonSchema}
  />
</RJSFWrapperComponent>` >}}
See this [RJSF Component](https://github.com/meshery/meshery/blob/0bc68d1cd0ba80a565afa68bce80899c22db9a2e/ui/components/MesheryMeshInterface/PatternService/RJSF.js#L91) as an example of how these properties are received.

### Extensibility: User Accounts

Meshery Server uses [Providers](/extensibility/providers) for identity management. Providers can implement the `user_account` extension to handle custom user management scenarios. The user avatar behavior, which changes based on the user's status, can be customized by extending the [User Component](https://github.com/meshery/meshery/blob/7de49ef4928f114080f923f2ad261f4433ca91d6/ui/components/User.js#L46).


###  Build-Time UI Extensibility

Meshery offers powerful customization options for its web application user interface at build time. This feature allows developers to tailor the UI to specific needs by modifying component behavior, managing routes, and applying custom themes.

##### Configuration File

The build-time UI customization is controlled through the `ui.config.js` file. This file supports various configuration options to modify the Meshery UI.

##### Current Features

*Component Management*

You can control the visibility and behavior of various UI components. For example:

```javascript
module.exports = {
  components: {
    navigator: true, // Set to false to disable the navigator component (default: true)
    // Add other components here as needed
  },
};
```

*Upcoming Features*

Meshery is continuously expanding its UI extensibility capabilities. The following features are planned for future releases:

1. **Route Management**: Ability to disable or modify specific routes.
1. **Redirect Configuration**: Set up custom redirects within the application.
1. **Custom Theming**: Apply custom themes to personalize the look and feel of the Meshery UI.

#### Usage

To customize the Meshery UI:

1. Locate the `ui.config.js` file in your Meshery project.
1. Modify the configuration options according to your requirements.
1. Rebuild the Meshery application to apply your changes.

### Loading Screen Message Persistence

Meshery UI displays a randomly selected loading message while the application and extensions load. To ensure a consistent user experience, the same loading message is displayed across all loading screens (main UI and extensions) during a single session.

#### How It Works

The loading message is selected once and stored in the browser's `window` object (`window.__mesheryLoadingMessage`). All subsequent loading screens retrieve this persisted message, preventing jarring re-renders with different messages.

#### Using the Persisted Loading Message in UI Plugins

**No changes needed** for plugins that:
- Use `DynamicFullScreenLoader` (already imports `randomLoadingMessage`)
- Are loaded through the main Meshery UI framework

**For standalone loaders**, import the persisted message:

```javascript
import { randomLoadingMessage } from '@/components/LoadingComponents/loadingMessages';
// or
import { getPersistedLoadingMessage } from '@/components/LoadingComponents/loadingMessages';

// Then use in your loading component:
<LoadingScreen message={randomLoadingMessage} />
```

The `randomLoadingMessage` export automatically retrieves the persisted value from `window.__mesheryLoadingMessage`, ensuring consistency across all loaders in your plugin.

#### Passing new custom prop to forms:

1.  Pass the new prop from the Meshery Extension in the RJSF Wrapper component used like this:

{{< code code=`function RJSFWrapperComponent(props) {
// Clone the child to pass in additional props
const children = React.cloneElement(props.children, {
...(props.children?.props || {}),
customComponent: YOUR_NEW_CUSTOM_COMPONENT_OR_PROP
});

return children
}` >}}
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

{{< code code=`<PatternServiceFormCore
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

</PatternServiceFormCore>` >}}

Meshery UI's RJSF form accepts two props:

- `RJSFWrapperComponent`: The wrapper component gets all of the props that are passed to the underlying form, allowing to inspect the props and changing the behaviour based on them.
- `RJSFFormChildComponent`: This component will customize the internals of the RJSF form.

With both of these props, Remote Providers can customize the wrapper and can also customize the body of the form. This allows full customization of the form.

from _ui/components/MesheryMeshInterface/PatternService/index.js_
{{< code code="function PatternService({ formData, jsonSchema, onChange, type, onSubmit, onDelete, RJSFWrapperComponent, RJSFFormChildComponent })" >}}
