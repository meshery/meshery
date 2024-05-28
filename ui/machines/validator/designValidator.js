import { createMachine, fromPromise, sendTo, assign, emit } from 'xstate';
import { getComponentDefinition } from '@/rtk-query/meshModel';
import Ajv from 'ajv';
import _ from 'lodash';
import {
  dataValidatorCommands,
  dataValidatorMachine,
  selectIsValidating,
  selectValidationResults,
} from './dataValidator';
import { useSelector } from '@xstate/react';
import { processDesign } from '@/utils/utils';
import { designsApi } from '@/rtk-query/design';
import { initiateQuery } from '@/rtk-query/utils';

const ajv = new Ajv({
  allErrors: true,
  strict: false, // allow additional properties like x-kubernetes-attributes ( this is safe the schema is sourced from the component definition and is not ours)
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

const DESIGN_VALIDATOR_COMMANDS = {
  VALIDATE_DESIGN_SCHEMA: 'VALIDATE_DESIGN_SCHEMA',
  VALIDATE_DESING_COMPONENT: 'VALIDATE_DESING_COMPONENT',
  DRY_RUN_DESIGN: 'DRY_RUN_DESIGN',
};

export const DESIGN_VALIDATOR_EVENTS = {
  DESIGN_SCHEMA_VALIDATION_DONE: 'DESIGN_SCHEMA_VALIDATION_DONE',
  DESIGN_SCHEMA_VALIDATION_FAILED: 'DESIGN_SCHEMA_VALIDATION_FAILED',
  DESIGN_DRY_RUN_DONE: 'DESIGN_DRY_RUN_DONE',
  DESIGN_DRY_RUN_FAILED: 'DESIGN_DRY_RUN_FAILED',
  TAP_ON_ERROR: 'TAP_ON_ERROR',
};

export const designValidatorCommands = {
  validateDesignSchema: ({ design, returnAddress }) => ({
    type: DESIGN_VALIDATOR_COMMANDS.VALIDATE_DESIGN_SCHEMA,
    returnAddress,
    data: {
      design,
      validationPayloadType: 'design',
    },
  }),
  validateDesignComponent: ({ component, returnAddress }) => ({
    type: DESIGN_VALIDATOR_COMMANDS.VALIDATE_DESING_COMPONENT,
    returnAddress,
    data: {
      component,
      validationPayloadType: 'component',
    },
  }),
  dryRunDesignDeployment: ({ design, k8sContexts, returnAddress }) => ({
    type: DESIGN_VALIDATOR_COMMANDS.DRY_RUN_DESIGN,
    returnAddress,
    data: { design, k8sContexts, dryRunType: DRY_RUN_TYPE.DEPLOY },
  }),
  dryRunDesignUnDeployment: ({ design, k8sContexts, returnAddress }) => ({
    type: DESIGN_VALIDATOR_COMMANDS.DRY_RUN_DESIGN,
    returnAddress,
    data: { design, k8sContexts, dryRunType: DRY_RUN_TYPE.UNDEPLOY },
  }),
};

export const designValidatorEvents = {
  designSchemaValidated: ({ design, validationResults }) => ({
    type: DESIGN_VALIDATOR_EVENTS.DESIGN_SCHEMA_VALIDATION_DONE,
    data: { design, validationResults },
  }),
  tapOnError: ({ error, type, component }) => ({
    type: DESIGN_VALIDATOR_EVENTS.TAP_ON_ERROR,
    data: { error, type, component },
  }),
};

const validateComponent = async (component, validateAnnotations = false) => {
  const componentDef = await getComponentDefinition(component.type, component.model, {
    apiVersion: component.apiVersion,
    annotations: 'include',
  });

  if (!componentDef || (componentDef?.metadata?.isAnnotation && !validateAnnotations)) {
    // skip validation for annotations
    return {
      errors: [],
      componentDefinition: componentDef,
      component,
    };
  }
  const schema = JSON.parse(componentDef.component.schema);
  const results = validateSchema(schema, component.settings || {}, componentDef.id);

  const validationResults = {
    ...results,
    componentDefinition: componentDef,
    component,
  };

  return validationResults;
};

const validateDesign = async (design) => {
  const { configurableComponents } = processDesign(design);

  const validationResults = {};

  for (const configurableComponent of configurableComponents) {
    try {
      const componentValidationResults = await validateComponent(configurableComponent);
      validationResults[configurableComponent.name] = componentValidationResults;
    } catch (error) {
      console.error('Error validating component', error);
    }
  }
  return validationResults;
};

const SchemaValidateDesignActor = fromPromise(async ({ input }) => {
  const { validationPayload, prevValidationResults } = input;
  const { validationPayloadType } = validationPayload;

  if (validationPayloadType === 'design') {
    const { design } = validationPayload;
    const validationResults = await validateDesign(design);
    return {
      validationResults,
    };
  }

  if (validationPayloadType === 'component') {
    const { component } = validationPayload;
    const validationResults = await validateComponent(
      component,
      validationPayload.validateAnnotations,
    );

    return {
      validationResults: _.set(prevValidationResults || {}, component.name, validationResults),
    };
  }

  throw new Error('Invalid validation payload type', validationPayloadType);
});

const DRY_RUN_TYPE = {
  DEPLOY: 'Deploy',
  UNDEPLOY: 'Undeploy',
};

export const formatDryRunResponse = (dryRunResponse) => {
  function getErrors(error) {
    if (error?.Causes) {
      return error.Causes;
    }
    // if causes aren't present use the status
    // status are strings, and causes are objects
    // so they're handled seprately in UI
    if (error?.Status) {
      return [error.Status];
    }
    return [];
  }

  let errorList = [];
  if (dryRunResponse) {
    Object.keys(dryRunResponse).forEach((compName) => {
      const contextsErrors = dryRunResponse?.[compName];

      if (!contextsErrors) {
        return;
      }

      Object.keys(contextsErrors).forEach((contextKey) => {
        const errorAndMeta = contextsErrors[contextKey];

        if (!errorAndMeta.success) {
          errorList.push({
            compName,
            contextId: contextKey,
            errors: getErrors(errorAndMeta.error),
          });
        }
      });
    });
  }

  return errorList;
};

const DryRunDesignActor = fromPromise(async ({ input: { validationPayload } }) => {
  const { design, k8sContexts, dryRunType } = validationPayload;
  const { pattern_file, pattern_id } = design;
  const dryRunEndpoint =
    dryRunType === DRY_RUN_TYPE.DEPLOY
      ? designsApi.endpoints.deployPattern
      : designsApi.endpoints.undeployPattern;

  const dryRunResults = await initiateQuery(dryRunEndpoint, {
    pattern_file,
    pattern_id,
    selectedK8sContexts: k8sContexts,
    dryRun: true,
    verify: false,
  });

  const validationResults = formatDryRunResponse(dryRunResults.data?.dryRunResponse);
  return {
    validationResults: validationResults,
  };
});

const schemaValidatorMachine = dataValidatorMachine.provide({
  actors: {
    ValidateActor: SchemaValidateDesignActor,
  },
});

const dryRunValidatorMachine = dataValidatorMachine.provide({
  actors: {
    ValidateActor: DryRunDesignActor,
  },
});

export const designValidationMachine = createMachine({
  id: 'designValidationMachine',
  initial: 'init',
  context: {},

  states: {
    init: {
      entry: assign({
        schemaValidator: ({ spawn }) =>
          spawn(schemaValidatorMachine, {
            name: 'schemaValidator',
            id: 'schemaValidator',
            syncSnapshot: true,
          }),
        dryRunValidator: ({ spawn }) =>
          spawn(dryRunValidatorMachine, {
            name: 'dryRunValidator',
            id: 'dryRunValidator',
            syncSnapshot: true,
          }),
      }),
      always: 'idle',
    },

    idle: {
      on: {
        [DESIGN_VALIDATOR_COMMANDS.VALIDATE_DESIGN_SCHEMA]: {
          actions: sendTo('schemaValidator', ({ event }) =>
            dataValidatorCommands.validateData({
              validationPayload: event.data,
              returnAddress: event.returnAddress, // directly return the response event from schemaValidator
            }),
          ),
        },
        [DESIGN_VALIDATOR_COMMANDS.VALIDATE_DESING_COMPONENT]: {
          actions: sendTo('schemaValidator', ({ event }) =>
            dataValidatorCommands.validateData({
              validationPayload: event.data,
              returnAddress: event.returnAddress,
            }),
          ),
        },
        [DESIGN_VALIDATOR_COMMANDS.DRY_RUN_DESIGN]: {
          actions: sendTo('dryRunValidator', ({ event }) =>
            dataValidatorCommands.validateData({
              validationPayload: event.data,
              returnAddress: event.returnAddress,
            }),
          ),
        },

        [DESIGN_VALIDATOR_EVENTS.TAP_ON_ERROR]: {
          actions: [emit(({ event }) => event)],
        },
      },
    },
  },
});

export const selectValidator = (state, validator) => state.context[validator];

const useSelectValidator = (validationMachine, validatorName, selector) => {
  const validator = useSelector(validationMachine, (s) => selectValidator(s, validatorName));
  const data = useSelector(validator, selector);
  return data;
};

export const useDesignSchemaValidationResults = (validationMachine) =>
  useSelectValidator(validationMachine, 'schemaValidator', selectValidationResults);

export const useDryRunValidationResults = (validationMachine) =>
  useSelectValidator(validationMachine, 'dryRunValidator', selectValidationResults);

export const selectComponentValidationResults = (state, componentId) => {
  const designValidationResults = selectValidationResults(
    selectValidator(state, 'schemaValidator').getSnapshot(),
  );
  if (!designValidationResults) return null;
  const componentResults = Object.values(designValidationResults).find(
    (result) => result?.component?.traits?.meshmap?.id === componentId,
  );
  return componentResults;
};

export const selectComponentDryRunResults = (state, componentName) => {
  const designValidationResults = selectValidationResults(
    selectValidator(state, 'dryRunValidator').getSnapshot(),
  );
  return designValidationResults?.find((result) => result.compName === componentName);
};

export const useIsValidatingDesign = (validationMachine, validatorName) => {
  const validator = useSelector(validationMachine, (s) => selectValidator(s, validatorName));
  const isValidating = useSelector(validator, selectIsValidating);
  return isValidating;
};

export const useIsValidatingDesignSchema = (validationMachine) =>
  useIsValidatingDesign(validationMachine, 'schemaValidator');
export const useIsValidatingDryRun = (validationMachine) =>
  useIsValidatingDesign(validationMachine, 'dryRunValidator');
