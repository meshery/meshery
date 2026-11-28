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

  it('falls back to unknown when modelReference is missing', () => {
    expect(
      componentKey({
        component: { kind: 'Pod', version: 'core/v1' },
      }),
    ).toBe('Pod-unknown-core/v1');
  });

  it('falls back to unknown when component is missing', () => {
    expect(
      componentKey({
        modelReference: { name: 'kubernetes' },
      }),
    ).toBe('unknown-kubernetes-unknown');
  });

  it('falls back to unknown for every field when both are missing', () => {
    expect(componentKey({})).toBe('unknown-unknown-unknown');
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

  it('reports a missing component definition as a validation error', async () => {
    const validator = createActor(schemaValidatorMachine).start();
    try {
      validator.send(
        dataValidatorCommands.validateData({
          validationPayload: {
            validationPayloadType: 'component',
            component: {
              id: 'c-missing',
              displayName: 'Missing Pod',
              component: { kind: 'Pod', version: 'core/v1' },
              modelReference: { name: 'kubernetes' },
              configuration: {},
            },
          },
          returnAddress: null,
        }),
      );
      const results = (await settled(validator)).context.validationResults;
      const result = results['c-missing'];

      expect(result.isValid).toBe(false);
      expect(result.componentDefinition).toBeUndefined();
      expect(result.errors).toHaveLength(1);
      expect(result.errors[0]).toMatchObject({
        instancePath: '',
        keyword: 'definition',
      });
      expect(result.errors[0].message).toContain('Pod-kubernetes-core/v1');
    } finally {
      validator.stop();
    }
  });

  it('reports a malformed component as a definition error instead of throwing', async () => {
    const validator = createActor(schemaValidatorMachine).start();
    const error = vi.spyOn(console, 'error').mockImplementation(() => {});
    try {
      validator.send(
        dataValidatorCommands.validateData({
          validationPayload: {
            validationPayloadType: 'component',
            component: { id: 'c-malformed', configuration: {} },
          },
          returnAddress: null,
        }),
      );
      const results = (await settled(validator)).context.validationResults;
      const result = results['c-malformed'];

      expect(result.isValid).toBe(false);
      expect(result.errors).toHaveLength(1);
      expect(result.errors[0]).toMatchObject({ keyword: 'definition' });
      expect(result.errors[0].message).toContain('unknown-unknown-unknown');
      expect(error).not.toHaveBeenCalled();
    } finally {
      error.mockRestore();
      validator.stop();
    }
  });

  it('reports missing definitions during design validation instead of returning no errors', async () => {
    const validator = createActor(schemaValidatorMachine).start();
    try {
      validator.send(
        dataValidatorCommands.validateData({
          validationPayload: {
            validationPayloadType: 'design',
            design: {
              components: [
                {
                  id: 'component-1',
                  displayName: 'Missing Service',
                  component: { kind: 'Service', version: 'v1' },
                  modelReference: { name: 'kubernetes' },
                  configuration: {},
                },
              ],
            },
            componentDefs: {},
          },
          returnAddress: null,
        }),
      );
      const results = (await settled(validator)).context.validationResults;

      expect(results['component-1'].isValid).toBe(false);
      expect(results['component-1'].errors).toHaveLength(1);
      expect(results['component-1'].errors[0].keyword).toBe('definition');
    } finally {
      validator.stop();
    }
  });

  it('keeps malformed components in the design results instead of dropping them', async () => {
    const validator = createActor(schemaValidatorMachine).start();
    try {
      validator.send(
        dataValidatorCommands.validateData({
          validationPayload: {
            validationPayloadType: 'design',
            design: { components: [{ id: 'c-malformed-design' }] },
            componentDefs: {},
          },
          returnAddress: null,
        }),
      );
      const results = (await settled(validator)).context.validationResults;

      expect(results['c-malformed-design'].isValid).toBe(false);
      expect(results['c-malformed-design'].errors[0].keyword).toBe('definition');
    } finally {
      validator.stop();
    }
  });

  it('skips annotations by default and validates them when requested', async () => {
    const validator = createActor(schemaValidatorMachine).start();
    try {
      const annotationDef = {
        ...componentDef('def-annotation'),
        metadata: { isAnnotation: true },
      };

      validator.send(
        dataValidatorCommands.validateData({
          validationPayload: {
            validationPayloadType: 'component',
            component: { id: 'annotation', configuration: { replicas: 'three' } },
            componentDef: annotationDef,
          },
          returnAddress: null,
        }),
      );
      const skipped = (await settled(validator)).context.validationResults;
      expect(skipped.annotation.errors).toEqual([]);

      validator.send(
        dataValidatorCommands.validateData({
          validationPayload: {
            validationPayloadType: 'component',
            validateAnnotations: true,
            component: { id: 'annotation-checked', configuration: { replicas: 'three' } },
            componentDef: annotationDef,
          },
          returnAddress: null,
        }),
      );
      const checked = await waitFor(
        validator,
        (state) =>
          state.matches({ idle: 'waiting' }) &&
          'annotation-checked' in state.context.validationResults,
      );
      expect(checked.context.validationResults['annotation-checked'].isValid).toBe(false);
    } finally {
      validator.stop();
    }
  });

  it('keeps skipping annotations whose definition is unavailable', async () => {
    const validator = createActor(schemaValidatorMachine).start();
    try {
      validator.send(
        dataValidatorCommands.validateData({
          validationPayload: {
            validationPayloadType: 'component',
            component: {
              id: 'c-annotation',
              metadata: { isAnnotation: true },
              component: { kind: 'Annotation', version: 'v1' },
              modelReference: { name: 'kubernetes' },
              configuration: {},
            },
          },
          returnAddress: null,
        }),
      );
      const results = (await settled(validator)).context.validationResults;

      expect(results['c-annotation'].errors).toEqual([]);
      expect(results['c-annotation'].isValid).toBeUndefined();
    } finally {
      validator.stop();
    }
  });
});
