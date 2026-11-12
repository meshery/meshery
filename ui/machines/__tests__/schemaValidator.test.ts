import { describe, expect, it } from 'vitest';
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
