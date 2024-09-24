---
layout: default
title: "Extensibility: UI"
permalink: extensibility/ui
type: Extensibility
abstract: "Meshery offers support for more adapters than any other project or product in the world. Meshery UI has a number of extension points that allow users to customize their experience with third-party plugins."
language: en
list: include
---

Meshery UI has a number of extension points that allow users to customize their experience with third-party plugins.

## Designing Custom Components for Meshery Extension Point

The Meshery extension points are the way to extend meshery and derive the more custom use-cases out of it. We already have an extension point called Kanvas.
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

### Extensibility: User Accounts

Meshery leverages remote providers for identity management. These providers can implement user_account extensions to handle custom user management scenarios.

The user avatar behavior, which changes based on the user's status, can be customized by extending the [User Component](https://github.com/meshery/meshery/blob/7de49ef4928f114080f923f2ad261f4433ca91d6/ui/components/User.js#L46).

### Extensibility: Permissions (CASL)

#### Functionality

Meshery UI uses [CASL.js](https://casl.js.org/v6/en/) to implement it's permissions framework. CASL is an isomorphic authorization JavaScript library which restricts what resources a given client is allowed to access.
Meshery's permissions framework has been designed to be very customizable and robust.

- It uses keys as it's unit checks to see if a user has permission to perform action or not.
- These keys are bundled together in keychains. Keychains are collections of similar permissions, which can themselves be grouped together and assigned to roles.
- Roles map permissions to users. Roles contain any number of keychains, which contain any number of keys (permissions). Meshery uses roles to assign permissions to users.

{% include alert.html type="info" title="Example extension" content="Meshery Provider is an extension to Meshery UI. To learn more about how permissions work in Meshery Provider in terms of Roles, keychains and keys, head over to <a href='https://docs.layer5.io/cloud/security/'>Layer5 documentation</a>." %}

#### Behavior

- Once a user has logged in, the backend will send a response with the permissions that the user has, based on the role assigned to them. Those permissions will be used to check abilities on the frontend.

{% capture code_content %}Example keys response

{
"keys": [
{"id": "382da488-9a92-4a5b-958d-c4bfe1e80253",
"function": "View All Users",},
{"id": "fa7de118-2d08-4b07-b9d7-3e0baead6d04",
"function": "View Profile",},
]
}
{% endcapture %}
{% include code.html code=code_content %}

<br>
- Meshery UI uses `permissions_constants.js` file to keep record of all the keys to check against (fully customizable if self hosting meshery). The keys are accessible using monikers (properties) of the keys object and have been setup this way to make devEx better. Anytime any key
needs to be added to an action, it needs to be added in this file with an appropriate moniker.

```
meshery
  └── ui
    └── utils
      ├── ...
      └── permission_constants.js
```

{% capture code_content %}export const keys = {
VIEW_ALL_ORGANIZATIONS: {
subject: 'View All Organizations',
action: 'e996c998-a50f-4cb8-ae7b-f2f1b523c971',
},
...
RESET_DATABASE: {
subject: 'Reset database',
action: '84fc402c-f33e-4a21-a0e3-e14f9e20b125',
},
};
{% endcapture %}
{% include code.html code=code_content %}

<br>
- Along with it, we export the `ability.can(action, subject)` provided by CASL as `CAN(action, subject)` for better devEx.

```JavaScript
export default function CAN(action, subject) {
  return ability.can(action, _.lowerCase(subject));
}
```

<br>
- Now the required actions can be easily disabled using the CAN utility.

{% capture code_content %}<Button
type="submit"
variant="contained"
color="primary"
size="large"
onClick={(e) => handleEnvironmentModalOpen(e, ACTION_TYPES.CREATE)}
style={ {
padding: '8px',
borderRadius: 5,
marginRight: '2rem',
} }
disabled={!CAN(keys.CREATE_ENVIRONMENT.action, keys.CREATE_ENVIRONMENT.subject)}
data-cy="btnResetDatabase"

> <AddIconCircleBorder style={ { width: '20px', height: '20px' } } />
> <Typography

    style={ {
      paddingLeft: '4px',
      marginRight: '4px',
    } }

>

    Create

  </Typography>
</Button>{% endcapture %}
{% include code.html code=code_content %}

### Build-Time UI Extensibility

Meshery offers powerful customization options for its web application user interface at build time. This feature allows developers to tailor the UI to specific needs by modifying component behavior, managing routes, and applying custom themes.

##### Configuration File

The build-time UI customization is controlled through the `ui.config.js` file. This file supports various configuration options to modify the Meshery UI.

##### Current Features

_Component Management_

You can control the visibility and behavior of various UI components. For example:

```javascript
module.exports = {
  components: {
    navigator: true, // Set to false to disable the navigator component (default: true)
    // Add other components here as needed
  },
};
```

_Upcoming Features_

Meshery is continuously expanding its UI extensibility capabilities. The following features are planned for future releases:

1. **Route Management**: Ability to disable or modify specific routes.
1. **Redirect Configuration**: Set up custom redirects within the application.
1. **Custom Theming**: Apply custom themes to personalize the look and feel of the Meshery UI.

#### Usage

To customize the Meshery UI:

1. Locate the `ui.config.js` file in your Meshery project.
1. Modify the configuration options according to your requirements.
1. Rebuild the Meshery application to apply your changes.

#### Passing new custom prop to forms:

1.  Pass the new prop from the Meshery Extension in the RJSF Wrapper component used like this:

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
