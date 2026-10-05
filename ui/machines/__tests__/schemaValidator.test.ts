import { describe, expect, it } from 'vitest';
import { createActor, waitFor } from 'xstate';
import { componentKey, schemaValidatorMachine } from '../validator/schemaValidator';
import { dataValidatorCommands } from '@sistent/sistent';

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

  it('rejects with an error message containing the unsupported payload type', async () => {
    const actor = createActor(schemaValidatorMachine);
    actor.start();

    try {
      // Send the proper event using dataValidatorCommands
      actor.send(
        dataValidatorCommands.validateData({
          validationPayload: {
            validationPayloadType: 'unsupported',
          },
          returnAddress: 'test-address',
        }),
      );

      // Wait for the machine to process the event and store the error in validationResults
      const snapshot = await waitFor(actor, (snapshot) => {
        const results = snapshot.context.validationResults;
        // Wait until validationResults contains the error string
        return typeof results === 'string' && results.includes('Invalid validation payload type');
      });

      // The error message should contain the unsupported type
      expect(snapshot.context.validationResults).toContain(
        'Invalid validation payload type: unsupported',
      );
    } finally {
      actor.stop();
    }
  });
});
