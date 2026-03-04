import { customizeValidator } from '@rjsf/validator-ajv8';

// Create a custom validator using ajv options and custom formats for OpenAPI/Kubernetes
const customValidator = customizeValidator({
  ajvOptionsOverrides: {
    allErrors: true,
    strict: false,
  },
  customFormats: {
    int32: (value: string) => {
      const num = Number(value);
      return Number.isInteger(num) && num >= -2147483648 && num <= 2147483647;
    },
    int64: (value: string) => {
      const num = Number(value);
      return (
        Number.isInteger(num) && num >= Number.MIN_SAFE_INTEGER && num <= Number.MAX_SAFE_INTEGER
      );
    },
    'int-or-string': (value: string) => {
      // Accept either an integer-like value or any string
      const num = Number(value);
      return Number.isInteger(num) || typeof value === 'string';
    },
  },
});

export default customValidator;
