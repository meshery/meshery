import { customizeValidator } from '@rjsf/validator-ajv8';

// Create a custom validator using the RJSF v6 customizeValidator API.
// Custom formats and Ajv options are passed via the options object.
const customValidator = customizeValidator({
  ajvOptionsOverrides: {
    allErrors: true,
    strict: false,
    validateFormats: true,
  },
  // Enable standard formats (date-time, email, uri, etc.) from ajv-formats
  ajvFormatOptions: {},
  // Add custom format validators for OpenAPI/Kubernetes specific formats
  customFormats: {
    int32: {
      type: 'number',
      validate: (value: number) => {
        return Number.isInteger(value) && value >= -2147483648 && value <= 2147483647;
      },
    },
    int64: {
      type: 'number',
      validate: (value: number) => {
        return (
          Number.isInteger(value) &&
          value >= Number.MIN_SAFE_INTEGER &&
          value <= Number.MAX_SAFE_INTEGER
        );
      },
    },
    'int-or-string': (value: string) => {
      return typeof value === 'string' || Number.isInteger(Number(value));
    },
  },
});

export default customValidator;
