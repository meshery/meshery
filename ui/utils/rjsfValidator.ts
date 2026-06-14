import Ajv from 'ajv';
import addFormats from 'ajv-formats';
import { customizeValidator } from '@rjsf/validator-ajv8';

// Create a custom ajv instance with OpenAPI/Kubernetes format support
const ajv = new Ajv({
  allErrors: true,
  strict: false,
  validateFormats: true,
});

// Add standard formats (date-time, email, etc.)
addFormats(ajv);

// Add custom format validators for OpenAPI/Kubernetes specific formats
ajv.addFormat('int32', {
  type: 'number',
  validate: (value) => {
    return Number.isInteger(value) && value >= -2147483648 && value <= 2147483647;
  },
});

ajv.addFormat('int64', {
  type: 'number',
  validate: (value) => {
    return (
      Number.isInteger(value) &&
      value >= Number.MIN_SAFE_INTEGER &&
      value <= Number.MAX_SAFE_INTEGER
    );
  },
});

ajv.addFormat('int-or-string', {
  validate: (value) => {
    return typeof value === 'string' || Number.isInteger(value);
  },
});

// Create a custom validator using the configured ajv instance
const customValidator = customizeValidator({}, ajv);

export default customValidator;
