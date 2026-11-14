import { fromPromise } from 'xstate';
import Ajv from 'ajv';
import addFormats from 'ajv-formats';

import { dataValidatorMachine } from '@sistent/sistent';

const ajv = new Ajv({
  allErrors: true,
  strict: false, // allow additional properties like x-kubernetes-attributes ( this is safe the schema is sourced from the component definition and is not ours)
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
  let validate = ajv.getSchema(id);
  if (!validate) {
    ajv.addSchema(schema, id);
    validate = ajv.getSchema(id);
  }
  const valid = validate(data);

  return {
    isValid: valid,
    errors: validate.errors,
  };
};

const validateComponent = (component, validateAnnotations = false, componentDef) => {
  if (!componentDef || (componentDef?.metadata?.isAnnotation && !validateAnnotations)) {
    // skip validation for annotations
    return {
      errors: [],
      componentDefinition: componentDef,
      component,
    };
  }
  const schema = JSON.parse(componentDef.component.schema);
  const results = validateSchema(schema, component.configuration || {}, componentDef.id);

  const validationResults = {
    ...results,
    component,
  };
  return validationResults;
};

export const componentKey = (component) =>
  `${component.component.kind}-${component.modelReference.name}-${component.component.version}`;

const validateDesign = (design, componentDefsStore) => {
  const configurableComponents = design.components;
  const validationResults = {};

  for (const configurableComponent of configurableComponents) {
    try {
      const componentDef = componentDefsStore?.[componentKey(configurableComponent)];
      const componentValidationResults = validateComponent(
        configurableComponent,
        false,
        componentDef,
      );
      validationResults[configurableComponent.id] = componentValidationResults;
    } catch (error) {
      console.log('Error validating component', error, design, componentDefsStore);
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

    const previousResults = prevValidationResults || {};
    // A string is the message of a failed validation. It stands until a design
    // validation replaces it rather than being overwritten by one component's result.
    if (typeof previousResults !== 'object') {
      return { validationResults: previousResults };
    }

    // Merge into a copy. The previous results are this validator's current
    // context; mutating them hands every consumer the same object back, so a
    // selector compared by reference never sees the component's new result.
    return {
      validationResults: { ...previousResults, [component.id]: validationResults },
    };
  }

  throw new Error(`Invalid validation payload type: ${validationPayloadType}`);
});

export const schemaValidatorMachine = dataValidatorMachine.provide({
  actors: {
    ValidateActor: SchemaValidateDesignActor,
  },
});
