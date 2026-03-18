import { fromPromise } from 'xstate';
import Ajv from 'ajv';
import addFormats from 'ajv-formats';
import _ from 'lodash';

import { dataValidatorMachine } from '@sistent/sistent';

const ajv = new Ajv({
  allErrors: true,
  strict: false,
  validateSchema: false, // Kubernetes/OpenAPI component schemas use extensions (x-kubernetes-*, non-standard format values) that fail AJV's meta-schema validation
  validateFormats: false, // OpenAPI schemas use format values (e.g. "int-or-string") as hints, not strict validators
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

// dynamically add schemas to ajv to avoid recompiling the same schema and cache it
const validateSchema = (schema, data, id) => {
  try {
    let validate = ajv.getSchema(id);
    if (!validate) {
      ajv.addSchema(schema, id);
      validate = ajv.getSchema(id);
    }
    if (!validate) {
      return { isValid: true, errors: [] };
    }
    const valid = validate(data);
    return {
      isValid: valid,
      errors: validate.errors || [],
    };
  } catch (e) {
    console.warn(`Schema validation setup failed for "${id}":`, e.message);
    return { isValid: true, errors: [] };
  }
};

const validateComponent = (component, validateAnnotations = false, componentDef) => {
  if (!componentDef || (componentDef?.metadata?.isAnnotation && !validateAnnotations)) {
    return {
      errors: [],
      componentDefinition: componentDef,
      component,
    };
  }

  const rawSchema = componentDef.component?.schema;
  if (!rawSchema) {
    return { errors: [], componentDefinition: componentDef, component };
  }

  let schema;
  try {
    schema = typeof rawSchema === 'string' ? JSON.parse(rawSchema) : rawSchema;
  } catch (e) {
    console.warn('Failed to parse component schema', componentDef.displayName, e);
    return { errors: [], componentDefinition: componentDef, component };
  }

  const results = validateSchema(schema, component.configuration || {}, componentDef.id);

  const validationResults = {
    ...results,
    component,
  };
  return validationResults;
};

export const componentKey = (component) => {
  const modelName = component.modelReference?.name || component.model?.name || '';
  return `${component.component.kind}-${modelName}-${component.component.version}`;
};

const validateDesign = (design, componentDefsStore) => {
  const configurableComponents = design?.components || [];
  const validationResults = {};

  for (const configurableComponent of configurableComponents) {
    try {
      const key = componentKey(configurableComponent);
      const componentDef = componentDefsStore?.[key];
      const componentValidationResults = validateComponent(
        configurableComponent,
        false,
        componentDef,
      );
      validationResults[configurableComponent.id] = componentValidationResults;
    } catch (error) {
      console.error('Error validating component', error, configurableComponent);
      validationResults[configurableComponent.id] = {
        errors: [],
        component: configurableComponent,
      };
    }
  }

  return validationResults;
};

const SchemaValidateDesignActor = fromPromise(async ({ input }) => {
  const { validationPayload, prevValidationResults } = input;
  const { validationPayloadType } = validationPayload;

  if (validationPayloadType === 'design') {
    const { design, componentDefs } = validationPayload;
    const validationResults = validateDesign(design, componentDefs);
    return {
      validationResults,
    };
  }

  if (validationPayloadType === 'component') {
    const { component, componentDef } = validationPayload;
    const validationResults = validateComponent(
      component,
      validationPayload.validateAnnotations || false,
      componentDef,
    );

    return {
      validationResults: _.set(prevValidationResults || {}, component.id, validationResults),
    };
  }

  throw new Error('Invalid validation payload type', validationPayloadType);
});

export const schemaValidatorMachine = dataValidatorMachine.provide({
  actors: {
    ValidateActor: SchemaValidateDesignActor,
  },
});
