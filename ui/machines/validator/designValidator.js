import { createMachine, fromPromise, sendTo, assign, emit } from 'xstate';
import { getComponentDefinition } from '@/rtk-query/meshModel';
import {
  dataValidatorCommands,
  dataValidatorMachine,
  selectIsValidating,
  selectValidationResults,
} from '@layer5/sistent';
import { useSelector } from '@xstate/react';
import { encodeDesignFile, processDesign } from '@/utils/utils';
import { designsApi } from '@/rtk-query/design';
import { initiateQuery } from '@/rtk-query/utils';

import { componentKey } from './schemaValidator';
import { fromWorkerfiedActor } from '@layer5/sistent';

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
  dryRunDesignDeployment: ({ design, k8sContexts, includeDependencies, returnAddress }) => ({
    type: DESIGN_VALIDATOR_COMMANDS.DRY_RUN_DESIGN,
    returnAddress,
    data: { design, k8sContexts, dryRunType: DRY_RUN_TYPE.DEPLOY, includeDependencies },
  }),
  dryRunDesignUnDeployment: ({ design, k8sContexts, includeDependencies, returnAddress }) => ({
    type: DESIGN_VALIDATOR_COMMANDS.DRY_RUN_DESIGN,
    returnAddress,
    data: { design, k8sContexts, dryRunType: DRY_RUN_TYPE.UNDEPLOY, includeDependencies },
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
            component: errorAndMeta.component,
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
  const { design, k8sContexts, dryRunType, includeDependencies } = validationPayload;
  // const { pattern_file, pattern_id } = design;
  const dryRunEndpoint =
    dryRunType === DRY_RUN_TYPE.DEPLOY
      ? designsApi.endpoints.deployPattern
      : designsApi.endpoints.undeployPattern;

  const dryRunResults = await initiateQuery(dryRunEndpoint, {
    pattern_file: encodeDesignFile(design),
    pattern_id: design.id,
    skipCRD: !includeDependencies,
    selectedK8sContexts: k8sContexts,
    dryRun: true,
    verify: false,
  });

  const validationResults = formatDryRunResponse(dryRunResults.data?.dryRunResponse);
  return {
    validationResults: validationResults,
  };
});

const dryRunValidatorMachine = dataValidatorMachine.provide({
  actors: {
    ValidateActor: DryRunDesignActor,
  },
});

const getAllComponentsDefsInDesign = async (design) => {
  const { components } = processDesign(design);
  const componentDefs = (
    await Promise.allSettled(
      components.map(async (component) =>
        getComponentDefinition(component.component.kind, component.model.name, {
          apiVersion: component.component.version,
          annotations: 'include',
        }),
      ),
    )
  )
    .filter((result) => result.status === 'fulfilled' && result.value)
    .map((result) => result.value);

  const componentStore = componentDefs.reduce((acc, componentDef) => {
    const key = componentKey(componentDef);
    acc[key] = componentDef;
    return acc;
  }, {});

  return componentStore;
};

export const designValidationMachine = createMachine({
  id: 'designValidationMachine',
  initial: 'init',
  context: {},

  states: {
    init: {
      entry: [
        () => console.log('spawning design validation actor wooo'),
        assign({
          schemaValidator: ({ spawn }) =>
            spawn(
              fromWorkerfiedActor(
                new Worker(new URL('./schemaValidatorWorker', import.meta.url), { type: 'module' }),
              ),
              {
                name: 'schemaValidator',
                id: 'schemaValidator',
                syncSnapshot: true,
              },
            ),

          dryRunValidator: ({ spawn }) =>
            spawn(dryRunValidatorMachine, {
              name: 'dryRunValidator',
              id: 'dryRunValidator',
              syncSnapshot: true,
            }),
        }),
      ],
      always: 'idle',
    },

    idle: {
      on: {
        [DESIGN_VALIDATOR_COMMANDS.VALIDATE_DESIGN_SCHEMA]: {
          target: 'validateDesignSchema',
        },
        [DESIGN_VALIDATOR_COMMANDS.VALIDATE_DESING_COMPONENT]: {
          target: 'validateComponentSchema',
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

    validateComponentSchema: {
      invoke: {
        input: ({ context, event }) => ({ context, event }),
        src: fromPromise(async ({ input }) => {
          const { component } = input.event.data;

          const def = await getComponentDefinition(component.component.kind, component.model.name, {
            apiVersion: component.component.version,
            annotations: 'include',
          });

          return {
            validationPayload: {
              ...input.event.data,
              componentDef: def,
              component: component,
            },
            returnAddress: input.event.returnAddress,
          };
        }),
        onDone: {
          actions: [
            sendTo('schemaValidator', ({ event }) =>
              dataValidatorCommands.validateData(event.output),
            ),
          ],
          target: 'idle',
        },
        onError: {
          target: 'idle',
          actions: ({ event }) =>
            console.log('error while relaying validateComponentSchema', event),
        },
      },
    },

    validateDesignSchema: {
      invoke: {
        input: ({ context, event }) => ({ context, event }),
        src: fromPromise(async ({ input }) => {
          const { event } = input;
          const def = await getAllComponentsDefsInDesign(event.data.design);
          return {
            validationPayload: {
              ...event.data,
              componentDefs: def,
            },
            returnAddress: event.returnAddress,
          };
        }),

        onDone: {
          actions: [
            sendTo('schemaValidator', ({ event }) =>
              dataValidatorCommands.validateData(event.output),
            ),
          ],
          target: 'idle',
        },
        onError: {
          target: 'idle',
          actions: ({ event }) => console.log('error while relaying validateDesignSchema', event),
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
