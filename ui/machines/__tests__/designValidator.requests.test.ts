import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { createActor, createMachine } from 'xstate';

vi.mock('@/rtk-query/meshModel', () => ({
  getComponentDefinition: vi.fn(),
}));
vi.mock('@/rtk-query/design', () => ({
  designsApi: {
    endpoints: { deployPattern: { initiate: vi.fn() }, undeployPattern: { initiate: vi.fn() } },
  },
}));
vi.mock('@/rtk-query/utils', () => ({
  initiateQuery: vi.fn(),
}));
vi.mock('@/utils/utils', () => ({
  encodeDesignFile: vi.fn(() => 'encoded'),
  processDesign: vi.fn((design) => ({
    components: design?.components || [],
    configurableComponents: design?.components || [],
    annotationComponents: [],
  })),
}));

import { getComponentDefinition } from '@/rtk-query/meshModel';
import { initiateQuery } from '@/rtk-query/utils';
import { designValidationMachine, designValidatorCommands } from '../validator/designValidator';

const TEST_TIMEOUT = 15000;

type Deferred = {
  promise: Promise<unknown>;
  resolve: (value: unknown) => void;
  reject: (reason: unknown) => void;
};

const deferred = (): Deferred => {
  let resolve;
  let reject;
  const promise = new Promise((res, rej) => {
    resolve = res;
    reject = rej;
  });
  return { promise, resolve, reject };
};

const until = async (condition: () => boolean, timeout = 4000) => {
  const start = Date.now();
  while (!condition()) {
    if (Date.now() - start > timeout) {
      throw new Error('condition not met before timeout');
    }
    await new Promise((resolve) => setTimeout(resolve, 10));
  }
};

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

let definitionCalls: { kind: string; deferred: Deferred }[] = [];

const deploymentDefinition = {
  id: 'def-deployment',
  component: {
    kind: 'Deployment',
    version: 'apps/v1',
    schema: JSON.stringify({
      type: 'object',
      properties: {
        spec: { type: 'object', properties: { replicas: { type: 'integer' } } },
      },
    }),
  },
  modelReference: { name: 'kubernetes' },
  metadata: {},
};

const deployment = (id: string, replicas: unknown) => ({
  id,
  component: { kind: 'Deployment', version: 'apps/v1' },
  modelReference: { name: 'kubernetes' },
  configuration: { spec: { replicas } },
});

const design = (id: string, components: ReturnType<typeof deployment>[]) => ({
  id,
  name: id,
  components,
});

const createSender = () => {
  const replies: any[] = [];
  const actor = createActor(
    createMachine({
      on: {
        '*': {
          actions: ({ event }) => {
            if (event.type.startsWith('DESIGN_VALIDATION')) {
              replies.push(event);
            }
          },
        },
      },
    }),
  ).start();
  return { actor, replies };
};

const payloadOf = (reply) => reply.data.validationPayload;

const schemaResults = (validator) =>
  validator.getSnapshot().context.schemaValidator.getSnapshot().context.validationResults;

const resolveDefinitions = (indexes: number[]) => {
  for (const index of indexes) {
    definitionCalls[index].deferred.resolve(deploymentDefinition);
  }
};

describe('designValidationMachine overlapping requests', () => {
  let validator;
  let senders: ReturnType<typeof createSender>[] = [];
  const sender = () => {
    const created = createSender();
    senders.push(created);
    return created;
  };

  beforeEach(() => {
    definitionCalls = [];
    vi.mocked(getComponentDefinition).mockImplementation((kind) => {
      const call = { kind, deferred: deferred() };
      definitionCalls.push(call);
      return call.deferred.promise as any;
    });
    validator = createActor(designValidationMachine).start();
  });

  afterEach(() => {
    validator.stop();
    senders.forEach(({ actor }) => actor.stop());
    senders = [];
    vi.mocked(getComponentDefinition).mockReset();
    vi.mocked(initiateQuery).mockReset();
    vi.restoreAllMocks();
  });

  it(
    'validates a component and then a design sent before the component definition resolved',
    async () => {
      const { actor: kanvas, replies } = sender();
      const component = deployment('c1', 'two');
      const newest = design('d1', [deployment('c1', 3), deployment('c2', 1)]);
      validator.send(
        designValidatorCommands.validateDesignComponent({ component, returnAddress: kanvas }),
      );
      expect(validator.getSnapshot().matches('validateComponentSchema')).toBe(true);
      validator.send(
        designValidatorCommands.validateDesignSchema({ design: newest, returnAddress: kanvas }),
      );
      expect(validator.getSnapshot().matches('validateDesignSchema')).toBe(true);
      await until(() => definitionCalls.length === 3);
      resolveDefinitions([1, 2]);
      resolveDefinitions([0]);
      await until(() => replies.length === 2);
      expect(replies.map((reply) => reply.type)).toEqual([
        'DESIGN_VALIDATION_DONE',
        'DESIGN_VALIDATION_DONE',
      ]);
      expect(payloadOf(replies[0]).validationPayloadType).toBe('component');
      expect(payloadOf(replies[0]).component).toBe(component);
      expect(payloadOf(replies[1]).validationPayloadType).toBe('design');
      expect(payloadOf(replies[1]).design).toBe(newest);
      const results = schemaResults(validator);
      expect(Object.keys(results).sort()).toEqual(['c1', 'c2']);
      expect(results.c1.isValid).toBe(true);
      expect(results.c1.component).toBe(newest.components[0]);
      expect(validator.getSnapshot().matches('idle')).toBe(true);
    },
    TEST_TIMEOUT,
  );

  it(
    'merges a component sent after a design into that design, even when it resolves first',
    async () => {
      const { actor: kanvas, replies } = sender();
      const first = design('d1', [deployment('c1', 'two'), deployment('c2', 1)]);
      const edited = deployment('c1', 5);
      validator.send(
        designValidatorCommands.validateDesignSchema({ design: first, returnAddress: kanvas }),
      );
      validator.send(
        designValidatorCommands.validateDesignComponent({
          component: edited,
          returnAddress: kanvas,
        }),
      );
      await until(() => definitionCalls.length === 3);
      resolveDefinitions([2]);
      await sleep(50);
      expect(replies).toHaveLength(0);
      resolveDefinitions([0, 1]);
      await until(() => replies.length === 2);
      expect(payloadOf(replies[0]).validationPayloadType).toBe('design');
      expect(payloadOf(replies[1]).validationPayloadType).toBe('component');
      const results = schemaResults(validator);
      expect(Object.keys(results).sort()).toEqual(['c1', 'c2']);
      expect(results.c1.isValid).toBe(true);
      expect(results.c1.component).toBe(edited);
    },
    TEST_TIMEOUT,
  );

  it(
    'keeps the newest of two designs when the older one resolves last',
    async () => {
      const autoValidation = sender();
      const modal = sender();
      const older = design('d-old', [deployment('c1', 'two')]);
      const newer = design('d-new', [deployment('c1', 3)]);
      validator.send(
        designValidatorCommands.validateDesignSchema({
          design: older,
          returnAddress: autoValidation.actor,
        }),
      );
      validator.send(
        designValidatorCommands.validateDesignSchema({ design: newer, returnAddress: modal.actor }),
      );
      await until(() => definitionCalls.length === 2);
      resolveDefinitions([1]);
      await until(() => modal.replies.length === 1 && autoValidation.replies.length === 1);
      expect(payloadOf(modal.replies[0]).design).toBe(newer);
      expect(payloadOf(autoValidation.replies[0]).design).toBe(newer);
      expect(schemaResults(validator).c1.isValid).toBe(true);
      resolveDefinitions([0]);
      await sleep(500);
      expect(modal.replies).toHaveLength(1);
      expect(autoValidation.replies).toHaveLength(1);
      expect(schemaResults(validator).c1.isValid).toBe(true);
      expect(schemaResults(validator).c1.component).toBe(newer.components[0]);
    },
    TEST_TIMEOUT,
  );

  it(
    'answers a sender once when it supersedes its own earlier design request',
    async () => {
      const { actor: kanvas, replies } = sender();
      const older = design('d-old', [deployment('c1', 'two')]);
      const newer = design('d-new', [deployment('c1', 3)]);
      validator.send(
        designValidatorCommands.validateDesignSchema({ design: older, returnAddress: kanvas }),
      );
      validator.send(
        designValidatorCommands.validateDesignSchema({ design: newer, returnAddress: kanvas }),
      );
      await until(() => definitionCalls.length === 2);
      resolveDefinitions([1, 0]);
      await until(() => replies.length === 1);
      await sleep(500);
      expect(replies).toHaveLength(1);
      expect(payloadOf(replies[0]).design).toBe(newer);
    },
    TEST_TIMEOUT,
  );

  it(
    'sends each sender exactly its own reply',
    async () => {
      const autoValidation = sender();
      const modal = sender();
      const component = deployment('c1', 4);
      const current = design('d1', [deployment('c1', 4), deployment('c2', 'x')]);
      validator.send(
        designValidatorCommands.validateDesignComponent({
          component,
          returnAddress: autoValidation.actor,
        }),
      );
      validator.send(
        designValidatorCommands.validateDesignSchema({
          design: current,
          returnAddress: modal.actor,
        }),
      );
      await until(() => definitionCalls.length === 3);
      resolveDefinitions([1, 2, 0]);
      await until(() => autoValidation.replies.length === 1 && modal.replies.length === 1);
      await sleep(500);
      expect(autoValidation.replies).toHaveLength(1);
      expect(modal.replies).toHaveLength(1);
      expect(payloadOf(autoValidation.replies[0]).validationPayloadType).toBe('component');
      expect(payloadOf(autoValidation.replies[0]).component).toBe(component);
      expect(payloadOf(modal.replies[0]).validationPayloadType).toBe('design');
      expect(payloadOf(modal.replies[0]).design).toBe(current);
      const results = schemaResults(validator);
      expect(results.c1.isValid).toBe(true);
      expect(results.c2.isValid).toBe(false);
    },
    TEST_TIMEOUT,
  );

  it(
    'still answers an earlier sender when a later request names no return address',
    async () => {
      const autoValidation = sender();
      const first = design('d1', [deployment('c1', 1)]);
      const fromModal = design('d2', [deployment('c1', 2)]);
      validator.send(
        designValidatorCommands.validateDesignSchema({
          design: first,
          returnAddress: autoValidation.actor,
        }),
      );
      await until(() => definitionCalls.length === 1);
      resolveDefinitions([0]);
      await until(() => Boolean(validator.getSnapshot().context.schemaRequestInFlight));
      validator.send(designValidatorCommands.validateDesignSchema({ design: fromModal }));
      await until(() => definitionCalls.length === 2);
      resolveDefinitions([1]);
      await until(() => autoValidation.replies.length === 1);
      await sleep(500);
      expect(autoValidation.replies).toHaveLength(1);
      expect(payloadOf(autoValidation.replies[0]).design).toBe(first);
      expect(schemaResults(validator).c1.component).toBe(fromModal.components[0]);
    },
    TEST_TIMEOUT,
  );

  it(
    'accepts a dry run while definitions load and replies to its own sender',
    async () => {
      const autoValidation = sender();
      const dryRunModal = sender();
      vi.mocked(initiateQuery).mockResolvedValue({ data: { dryRunResponse: {} } } as any);
      const current = design('d1', [deployment('c1', 1)]);
      validator.send(
        designValidatorCommands.validateDesignSchema({
          design: current,
          returnAddress: autoValidation.actor,
        }),
      );
      expect(validator.getSnapshot().matches('validateDesignSchema')).toBe(true);
      validator.send(
        designValidatorCommands.dryRunDesignDeployment({
          design: current,
          k8sContexts: ['ctx-1'],
          includeDependencies: false,
          returnAddress: dryRunModal.actor,
        }),
      );
      await until(() => dryRunModal.replies.length === 1);
      expect(payloadOf(dryRunModal.replies[0]).dryRunType).toBe('Deploy');
      expect(autoValidation.replies).toHaveLength(0);
      resolveDefinitions([0]);
      await until(() => autoValidation.replies.length === 1);
      expect(dryRunModal.replies).toHaveLength(1);
    },
    TEST_TIMEOUT,
  );

  it(
    'warns about an unregistered string address and still replies to the other recipients',
    async () => {
      const warn = vi.spyOn(console, 'warn').mockImplementation(() => {});
      const { actor: kanvas, replies } = sender();
      const older = design('d-old', [deployment('c1', 1)]);
      const newer = design('d-new', [deployment('c1', 2)]);
      validator.send(
        designValidatorCommands.validateDesignSchema({
          design: older,
          returnAddress: 'unregisteredListener',
        }),
      );
      validator.send(
        designValidatorCommands.validateDesignSchema({ design: newer, returnAddress: kanvas }),
      );
      await until(() => definitionCalls.length === 2);
      resolveDefinitions([0, 1]);
      await until(() => replies.length === 1);
      expect(payloadOf(replies[0]).design).toBe(newer);
      expect(warn).toHaveBeenCalledWith(expect.stringContaining('"unregisteredListener"'));
    },
    TEST_TIMEOUT,
  );

  it(
    'replies with a failure when component definitions cannot be resolved',
    async () => {
      const errorLog = vi.spyOn(console, 'error').mockImplementation(() => {});
      const { actor: kanvas, replies } = sender();
      const component = deployment('c1', 1);
      validator.send(
        designValidatorCommands.validateDesignComponent({ component, returnAddress: kanvas }),
      );
      await until(() => definitionCalls.length === 1);
      definitionCalls[0].deferred.reject(new Error('registry unavailable'));
      await until(() => replies.length === 1);
      expect(replies[0].type).toBe('DESIGN_VALIDATION_FAILED');
      expect(payloadOf(replies[0]).validationPayloadType).toBe('component');
      expect(replies[0].data.systemErrors).toBeInstanceOf(Error);
      expect(errorLog).toHaveBeenCalled();
      expect(validator.getSnapshot().matches('idle')).toBe(true);
      validator.send(
        designValidatorCommands.validateDesignComponent({ component, returnAddress: kanvas }),
      );
      await until(() => definitionCalls.length === 2);
      resolveDefinitions([1]);
      await until(() => replies.length === 2);
      expect(replies[1].type).toBe('DESIGN_VALIDATION_DONE');
    },
    TEST_TIMEOUT,
  );

  it(
    'runs dry-run requests for different k8s context sets concurrently without superseding each other',
    async () => {
      // Two deploy dry-runs aimed at *different* clusters must both complete.
      // Before the fix, both shared the key `DRY_RUN:Deploy` so the second
      // request superseded the first and ctx-A's caller never got a reply.
      const callerA = sender();
      const callerB = sender();
      const d = design('d1', [deployment('c1', 1)]);
      let callCount = 0;
      vi.mocked(initiateQuery).mockImplementation(async () => {
        callCount++;
        return { data: { dryRunResponse: {} } } as any;
      });

      validator.send(
        designValidatorCommands.dryRunDesignDeployment({
          design: d,
          k8sContexts: ['ctx-A'],
          includeDependencies: false,
          returnAddress: callerA.actor,
        }),
      );
      validator.send(
        designValidatorCommands.dryRunDesignDeployment({
          design: d,
          k8sContexts: ['ctx-B'],
          includeDependencies: false,
          returnAddress: callerB.actor,
        }),
      );

      await until(() => callerA.replies.length === 1 && callerB.replies.length === 1);
      expect(callerA.replies).toHaveLength(1);
      expect(callerB.replies).toHaveLength(1);
      expect(callCount).toBe(2);
      expect(payloadOf(callerA.replies[0]).k8sContexts).toEqual(['ctx-A']);
      expect(payloadOf(callerB.replies[0]).k8sContexts).toEqual(['ctx-B']);
    },
    TEST_TIMEOUT,
  );

  it(
    'supersedes a queued dry-run when a third request has the same type and context set',
    async () => {
      // With d1 in-flight (blocked), d2 queues behind it. Then d3 arrives with
      // the same supersedeKey and replaces d2 in the queue, inheriting d2's replyTo.
      // After d1 finishes the caller gets a reply for d1; then d3 runs and the caller
      // gets a reply for d3. d2 is never dispatched.
      const caller = sender();
      const d1 = design('d1', [deployment('c1', 1)]);
      const d2 = design('d2', [deployment('c1', 2)]);
      const d3 = design('d3', [deployment('c1', 3)]);
      let resolveFirst: (v: unknown) => void;
      let callCount = 0;
      vi.mocked(initiateQuery).mockImplementation(async (_ep, args) => {
        callCount++;
        if (args.patternId === d1.id) {
          // Block d1 so d2 and d3 can queue up before it finishes
          await new Promise((r) => {
            resolveFirst = r as (v: unknown) => void;
          });
        }
        return { data: { dryRunResponse: {} } } as any;
      });

      // d1 → goes in-flight immediately (dryRun is READY, no definition resolution)
      validator.send(
        designValidatorCommands.dryRunDesignDeployment({
          design: d1,
          k8sContexts: ['ctx-A'],
          includeDependencies: false,
          returnAddress: caller.actor,
        }),
      );
      await until(() => callCount === 1); // d1 is now in-flight

      // d2 → queues behind d1 (same supersedeKey, but d1 is in-flight not in queue)
      validator.send(
        designValidatorCommands.dryRunDesignDeployment({
          design: d2,
          k8sContexts: ['ctx-A'],
          includeDependencies: false,
          returnAddress: caller.actor,
        }),
      );

      // d3 → supersedes d2 in the queue (d2 has same supersedeKey and is queued)
      validator.send(
        designValidatorCommands.dryRunDesignDeployment({
          design: d3,
          k8sContexts: ['ctx-A'],
          includeDependencies: false,
          returnAddress: caller.actor,
        }),
      );

      // Unblock d1
      resolveFirst!(undefined);

      // Caller gets reply for d1, then d3 starts and finishes (d2 was dropped)
      await until(() => caller.replies.length === 2);
      await sleep(200);
      expect(caller.replies).toHaveLength(2);
      expect(callCount).toBe(2); // d1 and d3 ran; d2 was superseded before dispatch
      expect(payloadOf(caller.replies[0]).design).toBe(d1);
      expect(payloadOf(caller.replies[1]).design).toBe(d3);
    },
    TEST_TIMEOUT,
  );

  it(
    'removes the definition resolver from its children once the request settles',
    async () => {
      vi.spyOn(console, 'error').mockImplementation(() => {});
      const resolverChildren = () =>
        Object.keys(validator.getSnapshot().children).filter((id) =>
          id.startsWith('resolveDefinitions-'),
        );
      const { actor: kanvas, replies } = sender();
      const component = deployment('c1', 1);

      validator.send(
        designValidatorCommands.validateDesignComponent({ component, returnAddress: kanvas }),
      );
      await until(() => definitionCalls.length === 1);
      expect(resolverChildren()).toEqual(['resolveDefinitions-1']);
      resolveDefinitions([0]);
      await until(() => replies.length === 1);
      expect(replies[0].type).toBe('DESIGN_VALIDATION_DONE');
      expect(resolverChildren()).toEqual([]);

      validator.send(
        designValidatorCommands.validateDesignComponent({ component, returnAddress: kanvas }),
      );
      await until(() => definitionCalls.length === 2);
      expect(resolverChildren()).toEqual(['resolveDefinitions-2']);
      definitionCalls[1]?.deferred.reject(new Error('registry unavailable'));
      await until(() => replies.length === 2);
      expect(replies[1].type).toBe('DESIGN_VALIDATION_FAILED');
      expect(resolverChildren()).toEqual([]);
    },
    TEST_TIMEOUT,
  );
});
