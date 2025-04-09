# TypingFilter Component

A flexible, searchable filtering component that allows users to filter data based on multiple criteria with an intuitive typing interface.

## Features

- Real-time filtering as users type
- Support for multiple filter categories
- Predefined filter values or free-form input
- Chip-based selected filter display
- Keyboard navigation support

## Installation

The TypingFilter component is part of the Meshery UI library. No additional installation is required if you're working within the Meshery project.

## Props

| Prop             | Type     | Description                                                                     |
| ---------------- | -------- | ------------------------------------------------------------------------------- |
| `filterSchema`   | Object   | Configuration object defining available filter options (see schema below)       |
| `placeholder`    | String   | Placeholder text for the input field                                            |
| `handleFilter`   | Function | Callback function that receives the transformed filter data when filters change |
| `defaultFilters` | Array    | Optional array of default filters to initialize the component                   |

## FilterSchema Structure

```javascript
{
  FILTER_KEY: {
    value: "filterName",          // The filter name used in the UI
    description: "Description",   // Description shown in the dropdown
    type: "string",               // Optional data type
    values: ["Value1", "Value2"], // Optional predefined values
    multiple: false               // Optional - set to false to allow only one value per type
  },
  // More filters...
}
```

## Usage Example

```jsx
import TypingFilter from 'components/TypingFilter';
import React, { useState } from 'react';

const MyComponent = () => {
  const [filters, setFilters] = useState({});

  // Define your filter schema
  const filterSchema = {
    SEVERITY: {
      value: 'severity',
      description: 'Filter by severity level',
      values: ['Low', 'Medium', 'High', 'Critical'],
    },
    STATUS: {
      value: 'status',
      description: 'Filter by status',
      values: ['Active', 'Inactive', 'Pending'],
      multiple: false, // Only one status can be selected
    },
    AUTHOR: {
      value: 'author',
      description: 'Filter by content author',
      // No predefined values, accepts free text input
    },
  };

  // Handle filter changes
  const handleFilterChange = (newFilters) => {
    setFilters(newFilters);
    // Use the filters to update your data/UI
    console.log('Applied filters:', newFilters);
  };

  return (
    <div>
      <TypingFilter
        filterSchema={filterSchema}
        placeholder="Search or filter..."
        handleFilter={handleFilterChange}
        defaultFilters={[]} // Optional: pre-selected filters
      />

      {/* Rest of your component */}
    </div>
  );
};
```

## How to Use

1. **Basic Filtering**:

   - Start typing to see available filter options
   - Select a filter category from the dropdown

2. **Using Filter Categories**:

   - After selecting a category, type a colon (`:`) to see available values
   - For categories with predefined values, select from the dropdown
   - For free-text categories, type your value and press Enter

3. **Managing Filters**:

   - Selected filters appear as chips
   - Click the 'x' on a chip to remove that filter group
   - Click the clear icon to remove all filters

4. **Advanced Usage**:
   - Combine multiple filter types for complex queries
   - Set `multiple: false` in your schema for exclusive filter options

## Filter Transformation

The component transforms selected filters into a structured object for easy consumption in your application:

```javascript
// Example output from handleFilter callback:
{
  severity: ["High", "Critical"],  // Array for multi-value filters
  status: "Active",                // Single value when multiple: false
  author: ["John Doe"]             // Single item array for text inputs
}
```
