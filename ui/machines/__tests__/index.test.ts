import { describe, expect, it, vi } from 'vitest';

// Heavy machinery deps need stubbing exactly like the individual machine tests
// or the import graph will pull in the relay env, the redux store, the RTK
// query API, and the events SSE subscription bootstrap.
vi.mock('@/rtk-query/meshModel', () => ({ getComponentDefinition: vi.fn() }));
vi.mock('@/rtk-query/design', () => ({
  designsApi: {
    endpoints: { deployPattern: { initiate: vi.fn() }, undeployPattern: { initiate: vi.fn() } },
  },
}));
vi.mock('@/rtk-query/utils', () => ({ initiateQuery: vi.fn() }));
vi.mock('@/utils/utils', () => ({
  encodeDesignFile: vi.fn(),
  processDesign: vi.fn(() => ({ components: [] })),
}));
vi.mock('lib/eventsSubscription', () => ({
  subscribeToEvents: () => ({ dispose: vi.fn() }),
}));
vi.mock('../../store', () => ({ store: { dispatch: vi.fn() } }));
vi.mock('@/store/slices/events', () => ({ pushEvent: vi.fn() }));
vi.mock('../../rtk-query', () => ({ api: { util: { invalidateTags: vi.fn() } } }));
vi.mock('@/rtk-query/notificationCenter', () => ({ PROVIDER_TAGS: { EVENT: 'Event' } }));
vi.mock('@/components/layout/NotificationCenter/constants', () => ({
  SEVERITY_TO_NOTIFICATION_TYPE_MAPPING: {},
}));
vi.mock('../../lib/relayEnvironment', () => ({
  createRelayEnvironment: () => ({}),
}));

import * as machines from '../index';

describe('machines/index re-exports', () => {
  it('re-exports the operationsCenterActor and its event constants', () => {
    expect(machines.operationsCenterActor).toBeDefined();
    expect(machines.OPERATION_CENTER_EVENTS).toEqual({
      EVENT_RECEIVED_FROM_SERVER: 'EVENT_RECEIVED_FROM_SERVER',
      ERROR_OCCURRED_IN_SUBSCRIPTION: 'ERROR_OCCURRED_IN_SUBSCRIPTION',
    });
  });

  it('re-exports designValidator pieces', () => {
    expect(machines.designValidationMachine).toBeDefined();
    expect(machines.DESIGN_VALIDATOR_EVENTS).toBeDefined();
    expect(typeof machines.designValidatorCommands.validateDesignSchema).toBe('function');
    expect(typeof machines.designValidatorEvents.tapOnError).toBe('function');
    expect(typeof machines.formatDryRunResponse).toBe('function');
    expect(typeof machines.selectValidator).toBe('function');
    expect(typeof machines.useDesignSchemaValidationResults).toBe('function');
    expect(typeof machines.useDryRunValidationResults).toBe('function');
    expect(typeof machines.useIsValidatingDesign).toBe('function');
    expect(typeof machines.useIsValidatingDesignSchema).toBe('function');
    expect(typeof machines.useIsValidatingDryRun).toBe('function');
  });

  it('re-exports schemaValidatorMachine', () => {
    expect(machines.schemaValidatorMachine).toBeDefined();
  });
});
