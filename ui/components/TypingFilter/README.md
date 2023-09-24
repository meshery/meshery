# TypingFilter Component

The `TypingFilter` component is a customizable React component that enables real-time filtering and selection based on user input.
It provides a user-friendly interface for filtering data in your application.

## Table of Contents

- [Usage](#usage)
- [Props](#props)
- [Examples](#examples)

## Usage

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
      value: 'severity',
      description: 'Filter by severity',
      values: ['Low', 'Medium', 'High'],
      multiple: true, // default
    },
    // Add more filter categories as needed
  };

  // Define a callback function to handle filter changes.
  const handleFilterChange = (filteredData) => {
    // Implement your logic to react to the filtered data.
    // This function will be called when the user applies a filter. ( on presing enter in input)
    console.log('Filtered data:', filteredData);
  };

  return (
    <div>
      <TypingFilter filterSchema={filterSchema} handleFilter={handleFilterChange} />
      {/* Your other components */}
    </div>
  );
}

export default MyComponent;
```

## Props

The `TypingFilter` component accepts the following props:

- `filterSchema` (object, required): An object that defines available filter options. Each property of this object represents a filter category with the following properties:

  - `value` (string, required): The filter name used for filtering within the category.
  - `description` (string, required): Description of the filter category.
  - `type` (string, optional): The data type of the filter (e.g., "string", "number").
  - `values` (array, optional): Possible values for the filter.

- `handleFilter` (function, required): A callback function that is called when the user applies a filter. This function receives the filtered data as an argument.

- `autoFilter` (boolean,optional) :  A boolean to indicate if the filter should be applied automatically (on user input) .
 

# Finite State Machine (FSM) for `TypingFilter` Component

This README provides an overview of the Finite State Machine (FSM) implementation used to manage the state of the `TypingFilter` component.
The FSM is responsible for handling user interactions, such as selecting filters, entering values, and clearing the filter, within the component.

## Table of Contents

- [Overview](#overview)
- [State Definitions](#state-definitions)
- [Reducers](#reducers)
- [State Transitions](#state-transitions)
- [Initial State Handling](#initial-state-handling)

## Overview

The FSM implementation within the `TypingFilter` component ensures that user interactions are correctly processed and managed, resulting in a smooth and intuitive filtering experience.

## State Definitions

The FSM code defines three sets of constants to represent important elements within the state management:

### 1. `FILTERING_STATE`

Defines the possible states that the `TypingFilter` component can be in. These states include:

- `IDLE`: Represents the initial state when the component is not actively filtering.
- `SELECTING_FILTER`: Indicates that the user is selecting a filter.
- `SELECTING_VALUE`: Indicates that the user is entering a filter value.

### 2. `FILTER_EVENTS`

Represents the events that trigger state transitions within the FSM. Some of the events include:

- `START`: Initiates the filtering process.
- `SELECT`: Indicates the selection of a filter.
- `INPUT_CHANGE`: Represents a change in the filter input.
- `CLEAR`: Clears the filter.
- `EXIT`: Exits the filtering process.

### 3. `Delimiter`

Defines delimiters used to separate filter and value entries within the component. Delimiters include:

- `FILTER`: Separates multiple filters.
- `FILTER_VALUE`: Separates filters from their corresponding values.

## Reducers

The FSM implementation includes two key reducer functions:

### 1. `commonReducer`

This common reducer function handles events that are common across all states. It includes logic to handle "CLEAR" and "EXIT" events, which reset the component's state and clear any entered values.

### 2. `filterSelectionReducer`

The `filterSelectionReducer` is a specific reducer used to manage transitions between "SELECTING_FILTER" and "SELECTING_VALUE" states. It handles events related to selecting filters and entering values. The logic ensures that delimiters are appropriately added or removed when the user interacts with the filter.

## State Transitions

State transitions are managed based on user actions and the current state of the component. For example, when the user selects a filter, the state transitions from "SELECTING_FILTER" to "SELECTING_VALUE." When the user inputs values or clears the filter, the state transitions are managed accordingly.

## Initial State Handling

The FSM implementation includes handling for the initial state, where it listens for the "START" event to transition from "IDLE" to "SELECTING_FILTER." This ensures that the filtering process is initiated when the user interacts with the component.
