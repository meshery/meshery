import { describe, expect, it, vi, beforeEach } from 'vitest';

// The `can.ts` module wires `@casl/ability` together with the redux store and
// sistent's `createCanShow`. We mock the store + sistent's create helper so we
// can unit-test the `CAN` permission function and verify that `CanShow` is
// built with the right callbacks.

const hoisted = vi.hoisted(() => ({
  createCanShowMock: vi.fn(() => 'mocked-can-show'),
  getStateMock: vi.fn(() => ({ providerCapabilities: { capabilities: [] } })),
  eventBus: { publish: vi.fn() },
}));
const createCanShowMock = hoisted.createCanShowMock;
const getStateMock = hoisted.getStateMock;
const eventBus = hoisted.eventBus;

vi.mock('@sistent/sistent', async (importOriginal) => {
  const actual = await importOriginal<Record<string, unknown>>();
  return {
    ...actual,
    createCanShow: (...args: unknown[]) => hoisted.createCanShowMock(...args),
  };
});

vi.mock('../../store', () => ({
  store: { getState: hoisted.getStateMock },
}));

vi.mock('../eventBus', () => ({
  mesheryEventBus: hoisted.eventBus,
}));

import CAN, { CanShow, ability } from '../can';

describe('ability + CAN', () => {
  beforeEach(() => {
    // reset rules between tests
    ability.update([]);
  });

  it('denies all actions when no rules are defined', () => {
    expect(CAN('read', 'Pattern')).toBe(false);
    expect(CAN('manage', 'Workspace')).toBe(false);
  });

  it('lowercases the subject before delegating to ability.can', () => {
    ability.update([{ action: 'read', subject: 'pattern' }]);
    // exact match: should be allowed because we lowercase subject
    expect(CAN('read', 'Pattern')).toBe(true);
    expect(CAN('read', 'PATTERN')).toBe(true);
  });

  it('uses lodash _.lowerCase semantics which splits camelCase / spaces', () => {
    // _.lowerCase('PatternName') => 'pattern name', _.lowerCase('Pattern Name') => 'pattern name'
    ability.update([{ action: 'read', subject: 'pattern name' }]);
    expect(CAN('read', 'Pattern Name')).toBe(true);
    expect(CAN('read', 'PatternName')).toBe(true);
  });

  it('returns false when the rule does not match the action', () => {
    ability.update([{ action: 'read', subject: 'pattern' }]);
    expect(CAN('write', 'pattern')).toBe(false);
  });

  it('supports manage as an aggregate', () => {
    ability.update([{ action: 'manage', subject: 'design' }]);
    expect(CAN('read', 'design')).toBe(true);
    expect(CAN('update', 'design')).toBe(true);
    expect(CAN('delete', 'design')).toBe(true);
  });
});

describe('CanShow', () => {
  it('is a value (function or component) produced by sistent.createCanShow', () => {
    // The real `createCanShow` returns a React component; under our mock it
    // returns the string sentinel. We just guarantee something was wired up.
    expect(CanShow).toBeDefined();
  });
});
