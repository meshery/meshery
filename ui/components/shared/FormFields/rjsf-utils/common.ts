/**
 * Utility functions for RJSF (React JSON Schema Form) schema generation
 */

/**
 * JSON Schema property structure for select components
 */
interface SelectSchemaProperty {
  description: string;
  enum?: string[];
  items?: {
    enum: string[];
    type: string;
  };
  minItems?: number;
  type?: string;
  title: string;
  uniqueItems: boolean;
  'x-rjsf-grid-area': number;
}

/**
 * JSON Schema structure for select components
 */
interface SelectCompSchema {
  properties: {
    [key: string]: SelectSchemaProperty;
  };
  required: string[];
  type: 'object';
}

/**
 * Single select schema structure
 */
interface SingleSelectSchema {
  enum: string[];
}

/**
 * Multi-select schema structure
 */
interface MultiSelectSchema {
  items: {
    enum: string[];
    type: 'string';
  };
  minItems: number;
  type: 'array';
}

/**
 * Creates a JSON Schema for a select component (single or multi-select)
 *
 * @param enums - Array of enum values for the select options
 * @param description - Description text for the select component
 * @param title - Title text for the select component
 * @param name - Property name for the schema
 * @param multiSelect - Whether to create a multi-select (default: false)
 * @returns JSON Schema object for the select component
 *
 * @example
 * // Single select
 * const schema = selectCompSchema(
 *   ['Option 1', 'Option 2'],
 *   'Choose an option',
 *   'Select Option',
 *   'mySelect'
 * );
 *
 * @example
 * // Multi select
 * const schema = selectCompSchema(
 *   ['Option 1', 'Option 2', 'Option 3'],
 *   'Choose multiple options',
 *   'Select Options',
 *   'myMultiSelect',
 *   true
 * );
 */
export const selectCompSchema = (
  enums: string[],
  description: string,
  title: string,
  name: string,
  multiSelect: boolean = false,
): SelectCompSchema => {
  const intermediary: SingleSelectSchema | MultiSelectSchema =
    multiSelect === true ? getMultiselectSchema(enums) : getSingleSelectSchema(enums);

  return {
    properties: {
      [name]: {
        description: description,
        ...intermediary,
        title: title,
        uniqueItems: true,
        'x-rjsf-grid-area': 12,
      },
    },
    required: [name],
    type: 'object',
  };
};

/**
 * Creates a single select schema structure
 *
 * @param enums - Array of enum values
 * @returns Single select schema
 */
const getSingleSelectSchema = (enums: string[]): SingleSelectSchema => {
  return {
    enum: enums,
  };
};

/**
 * Creates a multi-select schema structure
 *
 * @param enums - Array of enum values
 * @returns Multi-select schema
 */
const getMultiselectSchema = (enums: string[]): MultiSelectSchema => {
  return {
    items: {
      enum: enums,
      type: 'string',
    },
    minItems: 1,
    type: 'array',
  };
};
