import { createMachine, fromPromise, assign, emit, enqueueActions } from 'xstate';
import { getComponentDefinition } from '@/rtk-query/meshModel';
import {
  DATA_VALIDATOR_EVENTS,
  dataValidatorCommands,
  dataValidatorEvents,
  dataValidatorMachine,
  selectIsValidating,
  selectValidationResults,
} from '@sistent/sistent';
import { useSelector } from '@xstate/react';
import { encodeDesignFile, processDesign } from '@/utils/utils';
import { designsApi } from '@/rtk-query/design';
import { initiateQuery } from '@/rtk-query/utils';

import { componentKey, schemaValidatorMachine } from './schemaValidator';

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
            type: 'ComponentError',
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
    patternFile: encodeDesignFile(design),
    patternId: design.id,
    skipCRD: !includeDependencies,
    selectedK8sContexts: k8sContexts,
    dryRun: true,
    verify: false,
  });

  const requestErrors = dryRunResults?.error
    ? [
        {
          type: 'RequestError',
          errors: [dryRunResults.error],
        },
      ]
    : [];
  const validationResults = formatDryRunResponse(dryRunResults.data?.dryRunResponse);
  const matrixEntry = dryRunResults.data?.dryRunResponse
    ? [
        {
          type: 'CompatibilityMatrix',
          rawDryRunResponse: dryRunResults.data?.dryRunResponse,
        },
      ]
    : [];
  return {
    validationResults: [...validationResults, ...matrixEntry, ...requestErrors],
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
        getComponentDefinition(component.component.kind, component.modelReference.name, {
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

const REQUEST_STATUS = {
  RESOLVING: 'resolving',
  READY: 'ready',
  FAILED: 'failed',
};

const REQUEST_KIND = {
  DESIGN: 'design',
  COMPONENT: 'component',
  DRY_RUN: 'dryRun',
};

const VALIDATOR_LANES = {
  schema: {
    validatorId: 'schemaValidator',
    queueKey: 'schemaRequests',
    inFlightKey: 'schemaRequestInFlight',
  },
  dryRun: {
    validatorId: 'dryRunValidator',
    queueKey: 'dryRunRequests',
    inFlightKey: 'dryRunRequestInFlight',
  },
};

const RESOLVE_DEFINITIONS_ACTOR_PREFIX = 'resolveDefinitions-';

const resolveDefinitionsActorId = (seq) => `${RESOLVE_DEFINITIONS_ACTOR_PREFIX}${seq}`;

const seqFromResolveDefinitionsActorId = (actorId) => {
  if (typeof actorId !== 'string' || !actorId.startsWith(RESOLVE_DEFINITIONS_ACTOR_PREFIX)) {
    return null;
  }
  const seq = Number(actorId.slice(RESOLVE_DEFINITIONS_ACTOR_PREFIX.length));
  return Number.isInteger(seq) ? seq : null;
};

const isDefinitionResolverEvent = ({ event }) =>
  seqFromResolveDefinitionsActorId(event.actorId) !== null;

const supersedeKeyFor = (kind, data) => {
  if (kind === REQUEST_KIND.COMPONENT) {
    return `${REQUEST_KIND.COMPONENT}:${data?.component?.id}`;
  }
  if (kind === REQUEST_KIND.DRY_RUN) {
    return `${REQUEST_KIND.DRY_RUN}:${data?.dryRunType}`;
  }
  return REQUEST_KIND.DESIGN;
};

const ResolveDefinitionsActor = fromPromise(async ({ input: { kind, data } }) => {
  if (kind === REQUEST_KIND.COMPONENT) {
    const { component } = data;
    const def = await getComponentDefinition(
      component.component.kind,
      component.modelReference.name,
      {
        apiVersion: component.component.version,
        annotations: 'include',
      },
    );
    return {
      validationPayload: {
        ...data,
        componentDef: def,
        validationType: 'schema',
        component: component,
      },
    };
  }
  const defs = await getAllComponentsDefsInDesign(data.design);
  return {
    validationPayload: {
      ...data,
      validationType: 'schema',
      componentDefs: defs,
    },
  };
});

const resolveRecipient = (system, address) =>
  typeof address === 'string' ? system.get(address) : address;

const recipientsOf = (system, request) => {
  const recipients = [];
  for (const address of request?.replyTo || []) {
    const recipient = resolveRecipient(system, address);
    if (recipient && !recipients.includes(recipient)) {
      recipients.push(recipient);
    }
  }
  return recipients;
};

const planLane = ({ lane, queue, inFlight, self, system }) => {
  const { validatorId } = VALIDATOR_LANES[lane];
  const sends = [];
  let pending = queue;
  let current = inFlight;
  while (!current && pending.length > 0) {
    const [head, ...rest] = pending;
    if (head.status === REQUEST_STATUS.RESOLVING) {
      break;
    }
    pending = rest;
    if (head.status === REQUEST_STATUS.FAILED) {
      const failure = dataValidatorEvents.dataValidationFailed({
        validationPayload: { ...head.data, validationType: 'schema' },
        systemErrors: head.error,
      });
      for (const recipient of recipientsOf(system, head)) {
        sends.push({ to: recipient, event: failure });
      }
      continue;
    }
    current = head;
    sends.push({
      to: validatorId,
      event: dataValidatorCommands.validateData({
        validationPayload: head.validationPayload,
        returnAddress: self,
      }),
    });
  }
  return { queue: pending, inFlight: current, sends };
};

const commitLane = (enqueue, lane, plan) => {
  const { queueKey, inFlightKey } = VALIDATOR_LANES[lane];
  enqueue.assign({ [queueKey]: plan.queue, [inFlightKey]: plan.inFlight });
  for (const { to, event } of plan.sends) {
    enqueue.sendTo(to, event);
  }
};

const acceptRequest = (lane, kind) =>
  enqueueActions(({ context, event, enqueue, self, system }) => {
    const { queueKey, inFlightKey } = VALIDATOR_LANES[lane];
    const seq = (context.requestSeq || 0) + 1;
    const supersedeKey = supersedeKeyFor(kind, event.data);
    const needsDefinitions = kind !== REQUEST_KIND.DRY_RUN;
    const queue = context[queueKey] || [];
    const superseded = queue.filter((request) => request.supersedeKey === supersedeKey);
    const request = {
      seq,
      kind,
      supersedeKey,
      status: needsDefinitions ? REQUEST_STATUS.RESOLVING : REQUEST_STATUS.READY,
      data: event.data,
      validationPayload: needsDefinitions ? null : event.data,
      error: null,
      replyTo: [...superseded.flatMap((older) => older.replyTo), event.returnAddress].filter(
        Boolean,
      ),
    };
    enqueue.assign({ requestSeq: seq });
    for (const older of superseded) {
      if (older.status === REQUEST_STATUS.RESOLVING) {
        enqueue.stopChild(resolveDefinitionsActorId(older.seq));
      }
    }
    if (needsDefinitions) {
      enqueue.spawnChild(ResolveDefinitionsActor, {
        id: resolveDefinitionsActorId(seq),
        input: { kind, data: event.data },
      });
    }
    commitLane(
      enqueue,
      lane,
      planLane({
        lane,
        queue: [...queue.filter((older) => older.supersedeKey !== supersedeKey), request],
        inFlight: context[inFlightKey],
        self,
        system,
      }),
    );
  });

const settleRequest = (settle) =>
  enqueueActions(({ context, event, enqueue, self, system }) => {
    const seq = seqFromResolveDefinitionsActorId(event.actorId);
    const { queueKey, inFlightKey } = VALIDATOR_LANES.schema;
    const queue = context[queueKey] || [];
    if (seq === null || !queue.some((request) => request.seq === seq)) {
      return;
    }
    commitLane(
      enqueue,
      'schema',
      planLane({
        lane: 'schema',
        queue: queue.map((request) => (request.seq === seq ? settle(request, event) : request)),
        inFlight: context[inFlightKey],
        self,
        system,
      }),
    );
  });

const laneOfReply = (context, event) => {
  const payload = event.data?.validationPayload;
  const lanes = Object.keys(VALIDATOR_LANES);
  const byPayload = lanes.find(
    (lane) => context[VALIDATOR_LANES[lane].inFlightKey]?.validationPayload === payload,
  );
  if (byPayload) {
    return byPayload;
  }
  const lane = payload?.validationType === 'schema' ? 'schema' : 'dryRun';
  return context[VALIDATOR_LANES[lane].inFlightKey] ? lane : null;
};

const relayValidatorReply = enqueueActions(({ context, event, enqueue, self, system }) => {
  const lane = laneOfReply(context, event);
  if (!lane) {
    return;
  }
  const { queueKey, inFlightKey } = VALIDATOR_LANES[lane];
  for (const recipient of recipientsOf(system, context[inFlightKey])) {
    enqueue.sendTo(recipient, event);
  }
  commitLane(
    enqueue,
    lane,
    planLane({ lane, queue: context[queueKey] || [], inFlight: null, self, system }),
  );
});

const hasResolvingRequest = (context, kind) =>
  (context[VALIDATOR_LANES.schema.queueKey] || []).some(
    (request) => request.kind === kind && request.status === REQUEST_STATUS.RESOLVING,
  );

export const designValidationMachine = createMachine({
  id: 'designValidationMachine',
  initial: 'init',
  context: {
    requestSeq: 0,
    schemaRequests: [],
    schemaRequestInFlight: null,
    dryRunRequests: [],
    dryRunRequestInFlight: null,
  },
  on: {
    [DESIGN_VALIDATOR_COMMANDS.VALIDATE_DESIGN_SCHEMA]: {
      actions: acceptRequest('schema', REQUEST_KIND.DESIGN),
    },
    [DESIGN_VALIDATOR_COMMANDS.VALIDATE_DESING_COMPONENT]: {
      actions: acceptRequest('schema', REQUEST_KIND.COMPONENT),
    },
    [DESIGN_VALIDATOR_COMMANDS.DRY_RUN_DESIGN]: {
      actions: acceptRequest('dryRun', REQUEST_KIND.DRY_RUN),
    },
    'xstate.done.actor.*': {
      guard: isDefinitionResolverEvent,
      actions: settleRequest((request, event) => ({
        ...request,
        status: REQUEST_STATUS.READY,
        validationPayload: event.output.validationPayload,
      })),
    },
    'xstate.error.actor.*': {
      guard: isDefinitionResolverEvent,
      actions: [
        ({ event }) =>
          console.error('Failed to resolve component definitions for validation', event.error),
        settleRequest((request, event) => ({
          ...request,
          status: REQUEST_STATUS.FAILED,
          error: event.error,
        })),
      ],
    },
    [DATA_VALIDATOR_EVENTS.DESIGN_VALIDATION_DONE]: {
      actions: relayValidatorReply,
    },
    [DATA_VALIDATOR_EVENTS.DESIGN_VALIDATION_FAILED]: {
      actions: relayValidatorReply,
    },
  },

  states: {
    init: {
      entry: [
        assign({
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
      ],
      always: 'idle',
    },

    idle: {
      always: [
        {
          guard: ({ context }) => hasResolvingRequest(context, REQUEST_KIND.DESIGN),
          target: 'validateDesignSchema',
        },
        {
          guard: ({ context }) => hasResolvingRequest(context, REQUEST_KIND.COMPONENT),
          target: 'validateComponentSchema',
        },
      ],
      on: {
        [DESIGN_VALIDATOR_EVENTS.TAP_ON_ERROR]: {
          actions: [emit(({ event }) => event)],
        },
      },
    },

    validateDesignSchema: {
      always: {
        guard: ({ context }) => !hasResolvingRequest(context, REQUEST_KIND.DESIGN),
        target: 'idle',
      },
    },

    validateComponentSchema: {
      always: [
        {
          guard: ({ context }) => hasResolvingRequest(context, REQUEST_KIND.DESIGN),
          target: 'validateDesignSchema',
        },
        {
          guard: ({ context }) => !hasResolvingRequest(context, REQUEST_KIND.COMPONENT),
          target: 'idle',
        },
      ],
    },
  },
});

export const selectValidator = (state, validator) => {
  try {
    return state.context?.[validator] || null;
  } catch (error) {
    console.warn('Error accessing validator from state context:', error);
    return null;
  }
};

const useSelectValidator = (validationMachine, validatorName, selector) => {
  const validator = useSelector(validationMachine, (s) => selectValidator(s, validatorName));
  const data = useSelector(validator, (state) => {
    // If validator is null/undefined, return null instead of crashing
    if (!validator || !state) return null;
    try {
      return selector(state);
    } catch (error) {
      console.warn('Error calling selector on validator state:', error);
      return null;
    }
  });

  // Return data only if validator exists
  return validator ? data : null;
};

export const useDesignSchemaValidationResults = (validationMachine) =>
  useSelectValidator(validationMachine, 'schemaValidator', selectValidationResults);

export const useDryRunValidationResults = (validationMachine) =>
  useSelectValidator(validationMachine, 'dryRunValidator', selectValidationResults);

export const selectComponentValidationResults = (state, componentId) => {
  const schemaValidatorActor = selectValidator(state, 'schemaValidator');
  if (!schemaValidatorActor) {
    return null;
  }

  try {
    const designValidationResults = selectValidationResults(schemaValidatorActor.getSnapshot());
    if (!designValidationResults) return null;
    const componentResults = Object.values(designValidationResults).find(
      (result) => result?.component?.traits?.meshmap?.id === componentId,
    );
    return componentResults;
  } catch (error) {
    console.warn('Error getting component validation results:', error);
    return null;
  }
};

export const selectComponentDryRunResults = (state, componentName) => {
  const dryRunValidatorActor = selectValidator(state, 'dryRunValidator');
  if (!dryRunValidatorActor) {
    return null;
  }

  try {
    const designValidationResults = selectValidationResults(dryRunValidatorActor.getSnapshot());
    return designValidationResults?.find((result) => result.compName === componentName);
  } catch (error) {
    console.warn('Error getting component dry run results:', error);
    return null;
  }
};

export const useIsValidatingDesign = (validationMachine, validatorName) => {
  const validator = useSelector(validationMachine, (s) => selectValidator(s, validatorName));
  const isValidating = useSelector(validator, (state) => {
    // If validator is null/undefined, return false instead of crashing
    if (!validator || !state) return false;
    try {
      return selectIsValidating(state);
    } catch (error) {
      console.warn('Error checking validation state:', error);
      return false;
    }
  });

  // Return isValidating only if validator exists, otherwise false
  return validator ? isValidating : false;
};

export const useIsValidatingDesignSchema = (validationMachine) =>
  useIsValidatingDesign(validationMachine, 'schemaValidator');
export const useIsValidatingDryRun = (validationMachine) =>
  useIsValidatingDesign(validationMachine, 'dryRunValidator');

export const useDryRunRawResponse = (validationMachine) => {
  const results = useDryRunValidationResults(validationMachine);
  if (!Array.isArray(results)) return null;
  const matrixEntry = results.find((r) => r?.type === 'CompatibilityMatrix');
  return matrixEntry?.rawDryRunResponse || null;
};
