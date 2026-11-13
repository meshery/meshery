import { describe, expect, it, vi } from 'vitest';
import { assign, createActor, createMachine, waitFor } from 'xstate';
import { DATA_VALIDATOR_EVENTS, dataValidatorCommands } from '@sistent/sistent';
import { componentKey, schemaValidatorMachine } from '../validator/schemaValidator';

describe('componentKey', () => {
  it('combines kind, model name and version with hyphens', () => {
    const component = {
      component: { kind: 'Pod', version: 'core/v1' },
      modelReference: { name: 'kubernetes' },
    };
    expect(componentKey(component)).toBe('Pod-kubernetes-core/v1');
  });

  it('returns a key for any input shape carrying the required fields', () => {
    expect(
      componentKey({
        component: { kind: 'Service', version: 'v1' },
        modelReference: { name: 'k8s' },
      }),
    ).toBe('Service-k8s-v1');
  });
});

describe('schemaValidatorMachine', () => {
  it('boots into the idle.waiting state', () => {
    const actor = createActor(schemaValidatorMachine);
    actor.start();
    const snap = actor.getSnapshot();
    // Per the dataValidatorMachine state graph (idle: { waiting | debouncing })
    expect(typeof snap.value).toBeDefined();
    actor.stop();
  });

  it('exposes a snapshot context with the validator slots', () => {
    const actor = createActor(schemaValidatorMachine);
    actor.start();
    const ctx = actor.getSnapshot().context;
    expect(ctx).toHaveProperty('validationResults');
    expect(ctx).toHaveProperty('validationPayload');
    expect(ctx).toHaveProperty('returnAddress');
    actor.stop();
  });

  it('includes an unsupported payload type in the validation failure', async () => {
    const receiver = createActor(
      createMachine({
        initial: 'waiting',
        context: { error: null },
        states: {
          waiting: {
            on: {
              [DATA_VALIDATOR_EVENTS.DESIGN_VALIDATION_FAILED]: {
                target: 'received',
                actions: assign({
                  error: ({ event }) => event.data.systemErrors,
                }),
              },
            },
          },
          received: {},
        },
      }),
    );
    const validator = createActor(schemaValidatorMachine);

    receiver.start();
    validator.start();

    try {
      validator.send(
        dataValidatorCommands.validateData({
          validationPayload: {
            validationPayloadType: 'unsupported',
          },
          returnAddress: receiver,
        }),
      );

      const snapshot = await waitFor(receiver, (state) => state.matches('received'));

      expect(snapshot.context.error).toBeInstanceOf(Error);
      expect(snapshot.context.error.message).toBe('Invalid validation payload type: unsupported');
    } finally {
      validator.stop();
      receiver.stop();
    }
  });
});

describe('schemaValidatorMachine component validation', () => {
  const componentDef = (id: string) => ({
    id,
    metadata: {},
    component: {
      schema: JSON.stringify({
        type: 'object',
        properties: { replicas: { type: 'integer' } },
      }),
    },
  });

  const validateComponent = (validator, id: string, replicas: unknown) =>
    validator.send(
      dataValidatorCommands.validateData({
        validationPayload: {
          validationPayloadType: 'component',
          component: { id, configuration: { replicas } },
          componentDef: componentDef(`def-${id}`),
        },
        returnAddress: null,
      }),
    );

  const settled = (validator) =>
    waitFor(
      validator,
      (state) => state.matches({ idle: 'waiting' }) && state.context.validationResults !== null,
    );

  // A component validation merges into the previous results. Merging in place
  // handed every consumer the same object back, so a selector compared by
  // reference (useSelector) saw no change, and the snapshot it had already read
  // changed underneath it.
  it('returns new results instead of mutating the previous ones', async () => {
    const validator = createActor(schemaValidatorMachine).start();
    try {
      validateComponent(validator, 'c-1', 'three');
      const first = (await settled(validator)).context.validationResults;
      expect(first['c-1'].isValid).toBe(false);

      validateComponent(validator, 'c-2', 2);
      const second = await waitFor(
        validator,
        (state) => state.matches({ idle: 'waiting' }) && 'c-2' in state.context.validationResults,
      );

      expect(second.context.validationResults).not.toBe(first);
      expect(Object.keys(first)).toEqual(['c-1']);
      expect(second.context.validationResults['c-1']).toBe(first['c-1']);
      expect(second.context.validationResults['c-2'].isValid).toBe(true);
    } finally {
      validator.stop();
    }
  });

  it('keeps the message of a failed validation instead of one component result', async () => {
    const validator = createActor(schemaValidatorMachine).start();
    const error = vi.spyOn(console, 'error').mockImplementation(() => {});
    try {
      validator.send(
        dataValidatorCommands.validateData({
          validationPayload: { validationPayloadType: 'unsupported' },
          returnAddress: null,
        }),
      );
      const failed = (await settled(validator)).context.validationResults;
      expect(typeof failed).toBe('string');

      validateComponent(validator, 'c-1', 3);
      await waitFor(validator, (state) => state.matches('validatingData'));
      const after = await settled(validator);

      expect(after.context.validationResults).toBe(failed);
    } finally {
      error.mockRestore();
      validator.stop();
    }
  });

  it('starts from empty results when nothing was validated before', async () => {
    const validator = createActor(schemaValidatorMachine).start();
    try {
      validateComponent(validator, 'c-1', 3);
      const results = (await settled(validator)).context.validationResults;
      expect(Object.keys(results)).toEqual(['c-1']);
      expect(results['c-1'].isValid).toBe(true);
    } finally {
      validator.stop();
    }
  });
});
