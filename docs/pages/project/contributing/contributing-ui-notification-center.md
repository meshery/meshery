---
layout: page
title: Contributing to Meshery UI - Notification Center
permalink: project/contributing/contributing-ui-notification-center
abstract: How to contribute to the Notification Center in Meshery's web-based UI.
language: en
display-title: false
type: project
category: contributing
list: include
---

<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisite Reading</strong></p>
  <ol><li><a ahref="contributing-ui">Contributing to Meshery UI</a></li></ol>
</div>

## <a name="contributing-ui-notification-center">Contributing to Meshery UI - Notification Center</a>

<video style="width:min(100%,750px)" height="auto" autoplay muted loop>
  <source src="https://github.com/meshery/meshery/assets/65964225/345672de-3f61-4be0-b3c8-0e7480cc496c" type="video/mp4">
 Your browser does not support the video tag
</video>

The `NotificationCenter` component of Meshery UI Switching to Graphql subscriptions and implementing robust filtering. Events are persisted in Meshery Server and state management on client is done using Redux Toolkit and RTK.

#### User-facing Features

- Robust filtering support inspired by GitHub's notification filtering style.
  - Search is also included.
- Proper hierarchial presentation of error details, including probable cause and suggested remeditation.
- Suport for notification status (notifications can be marked as read and unread)
  - *Future: Notifications can be acknowledged or resolved.*
- Event-based notification via Graphql subscription (provided by Meshery Server and any upstream components or externally managed systems, like Kubernetes)
- Infinite scroll for pagination.

#### State Management and Internal Details

- The State on client is managed using `Redux Tooltik` and `Rtk-query`
- Update and Delete operations are optimistically handled.
- Network Request are cached and are invalidated when new events come or events are deleted/updated.
- Due to need for infinite scroll and optimistic update the events are stored globally in Redux.

### Notification Severities and Colors

Notification severities and colors are defined in the constants file, `ui/components/NotificationCenter/constants.js`.

### Notification Filtering and Searching

**Table of Contents**

- [Usage](#usage)
- [Props](#props)
- [Examples](#examples)

The Notfication Center includes a resusable component, `TypingFilter`, for sophisticated filtering and searching of notifications based on their attributes. It adheres to the GitHub-style syntax for filtering, offering a straight-forward and adaptable way to filter and search notification details. The `TypingFilter` component is a customizable React component that enables real-time filtering and selection based on user input.

The state for filtering is managed by a state machine created using a reducer. `TypingFilter` supports multiple filters, suggestions and completions.

### Usage

The `TypingFilter` component is designed to provide an interactive filtering experience in your application. Here's how you can use it:

```javascript
import React from 'react';
import TypingFilter from './path-to-TypingFilter';

function MyComponent() {
  // Define a filter schema that describes the available filter options.
  const filterSchema = {
    // Define your filter categories here
    // Example:
    SEVERITY: {
      value: "severity",
      description: "Filter by severity",
      values: ["Low", "Medium", "High"],
      multiple : true // default
    },
    // Add more filter categories as needed
  };

  // Define a callback function to handle filter changes.
  const handleFilterChange = (filteredData) => {
    // Implement your logic to react to the filtered data.
    // This function will be called when the user applies a filter. ( on presing enter in input)
    console.log("Filtered data:", filteredData);
  };

  return (
    <div>
      <TypingFilter
        filterSchema={filterSchema}
        handleFilter={handleFilterChange}
      />
      {/* Your other components */}
    </div>
  );
}

export default MyComponent;
```

### Props

The `TypingFilter` component accepts the following props:

- `filterSchema` (object, required): An object that defines available filter options. Each property of this object represents a filter category with the following properties:
  - `value` (string, required): The filter name used for filtering within the category.
  - `description` (string, required): Description of the filter category.
  - `type` (string, optional): The data type of the filter (e.g., "string", "number").
  - `values` (array, optional): Possible values for the filter.

- `handleFilter` (function, required): A callback function that is called when the user applies a filter. This function receives the filtered data as an argument.


## Finite State Machine (FSM) for `TypingFilter` Component

This section provides an overview of the Finite State Machine (FSM) implementation used to manage the state of the `TypingFilter` component. The FSM is responsible for handling user interactions, such as selecting filters, entering values, and clearing the filter, within the component. The FSM implementation within the `TypingFilter` component ensures that user interactions are correctly processed and managed, resulting in a smooth and intuitive filtering experience.

**Table of Contents**

- [Overview](#overview)
- [State Definitions](#state-definitions)
- [Reducers](#reducers)
- [State Transitions](#state-transitions)
- [Initial State Handling](#initial-state-handling)

### State Definitions

The FSM code defines three sets of constants to represent important elements within the state management:

#### 1. `FILTERING_STATE`

Defines the possible states that the `TypingFilter` component can be in. These states include:
- `IDLE`: Represents the initial state when the component is not actively filtering.
- `SELECTING_FILTER`: Indicates that the user is selecting a filter.
- `SELECTING_VALUE`: Indicates that the user is entering a filter value.

#### 2. `FILTER_EVENTS`

Represents the events that trigger state transitions within the FSM. Some of the events include:
- `START`: Initiates the filtering process.
- `SELECT`: Indicates the selection of a filter.
- `INPUT_CHANGE`: Represents a change in the filter input.
- `CLEAR`: Clears the filter.
- `EXIT`: Exits the filtering process.

#### 3. `Delimiter`

Defines delimiters used to separate filter and value entries within the component. Delimiters include:
- `FILTER`: Separates multiple filters.
- `FILTER_VALUE`: Separates filters from their corresponding values.

### Reducers

The FSM implementation includes two key reducer functions:

#### 1. `commonReducer`

This common reducer function handles events that are common across all states. It includes logic to handle "CLEAR" and "EXIT" events, which reset the component's state and clear any entered values.

#### 2. `filterSelectionReducer`

The `filterSelectionReducer` is a specific reducer used to manage transitions between "SELECTING_FILTER" and "SELECTING_VALUE" states. It handles events related to selecting filters and entering values. The logic ensures that delimiters are appropriately added or removed when the user interacts with the filter.

### State Transitions

State transitions are managed based on user actions and the current state of the component. For example, when the user selects a filter, the state transitions from "SELECTING_FILTER" to "SELECTING_VALUE." When the user inputs values or clears the filter, the state transitions are managed accordingly.

### Initial State Handling

The FSM implementation includes handling for the initial state, where it listens for the "START" event to transition from "IDLE" to "SELECTING_FILTER." This ensures that the filtering process is initiated when the user interacts with the component.

## Bulk Operations

Bulk operations in the Notification Center allow users to perform actions like deleting multiple notifications or changing the status of multiple notifications in a batch. This documentation outlines the key features and functionality of bulk operations, including the restriction of performing only one bulk operation at a time, the disabling of buttons during ongoing operations, and the display of a loading icon to indicate ongoing activity.

**How It Works**

### Initiating a Bulk Operation

- Users select the notifications they want to include in the bulk operation. This is typically done by checking checkboxes next to each notification.
- After selecting notifications, users trigger the desired bulk operation (e.g., delete or change status) by clicking the corresponding action button.
- Once initiated, the bulk operation begins processing the selected notifications.

### Further Enchacements

To ensure a seamless user experience with bulk operations in the Notification Center, consider the following best practices:

- Provide clear and concise messaging to inform users about the progress and outcome of the bulk operation.
- Implement error handling to gracefully handle any issues that may arise during the operation and communicate errors to the user.
- Consider offering a confirmation dialog before initiating a bulk operation to prevent accidental actions.


## Metadata Formatter

When an event is received from the server, it adheres to a fixed schema containing information that is valuable for presentation to the user. This information typically includes details such as the description, date, user_id, system_id, action, and acted-upon resources. Additionally, sometimes there may be a detailed traceback, a summary, or a comprehensive error log, all of which are dynamically generated data encapsulated within the metadata of the event. Presenting this structured data in a user-friendly manner is a crucial task because it contains valuable insights into ongoing operations.

To accomplish this task, we employ metadata formatters that transform structured data into visually appealing formats. There are currently two types of formatters in use:

1. **Metadata Specific Formatters:** These formatters are specifically designed for particular types of metadata, such as Error and DryRunResponse. Metadata Specific Formatters are implemented as React components that take the metadata as input and render it within the component.
2. **Dynamic Formatter:** Since metadata can vary significantly in structure, it is not practical to create a specific formatter for each kind. Dynamic formatters analyze the schema's structure and apply custom-defined rules for formatting:
     - Text strings are rendered using the BodySectionRenderer (more on this later).
     - Arrays are rendered using the ArrayRenderer.
     - Key-value pairs are rendered using the KeyValueRenderer.
     - Nested objects are recursively rendered.

### BodySectionRenderer
The BodySectionRenderer is responsible for formatting and rendering raw text strings into React components. During this process, it parses the string to replace external links with `<Link>` components and checks if the link matches predefined sites to render the link accordingly.

### ArrayRenderer
The ArrayRenderer is responsible for rendering an array of items in a recursive manner, presenting them as a bulletized list using the MetdataFormatter.

### KeyValueRenderer
Object properties with string values are considered key-value pairs and are rendered as such.

### The Metadata Specific Formatter
Certain metadata, such as Kubernetes responses and Errors, hold high importance and have dedicated renderers. These dedicated renderers can still utilize the dynamic formatter to format specific parts of the response, such as DryRunResponse.

### Reusability
While this system was initially developed for our events and notification center, the components it comprises are highly reusable and can be employed in other contexts where dynamic formatting of structured data is required.

{% include suggested-reading.html %}
