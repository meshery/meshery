import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';

// Mock lodash and the ability singleton so we can control them in tests.
const hoisted = vi.hoisted(() => ({
  onMock: vi.fn(() => vi.fn()),
  canMock: vi.fn(() => false),
  updateMock: vi.fn(),
}));

vi.mock('../../can', () => ({
  ability: {
    can: hoisted.canMock,
    on: hoisted.onMock,
    update: hoisted.updateMock,
  },
}));

// lodash _.lowerCase is used inside useCan; we keep it real.
// It is imported transitively and does not need mocking.

import { useCan } from '../useCan';

describe('useCan', () => {
  let registeredCallback: (() => void) | null = null;
  const unsubscribeMock = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    registeredCallback = null;

    // Capture the callback registered via ability.on('update', cb)
    hoisted.onMock.mockImplementation((event: string, cb: () => void) => {
      if (event === 'update') {
        registeredCallback = cb;
      }
      return unsubscribeMock;
    });

    hoisted.canMock.mockReturnValue(false);
  });

  afterEach(() => {
    registeredCallback = null;
  });

  it('returns false when the ability denies the action', () => {
    hoisted.canMock.mockReturnValue(false);
    const { result } = renderHook(() => useCan('some-action-id', 'Some Subject'));
    expect(result.current).toBe(false);
  });

  it('returns true when the ability allows the action', () => {
    hoisted.canMock.mockReturnValue(true);
    const { result } = renderHook(() => useCan('some-action-id', 'Some Subject'));
    expect(result.current).toBe(true);
  });

  it('subscribes to ability update events on mount', () => {
    renderHook(() => useCan('action', 'subject'));
    expect(hoisted.onMock).toHaveBeenCalledWith('update', expect.any(Function));
  });

  it('unsubscribes on unmount', () => {
    const { unmount } = renderHook(() => useCan('action', 'subject'));
    expect(unsubscribeMock).not.toHaveBeenCalled();
    unmount();
    expect(unsubscribeMock).toHaveBeenCalled();
  });

  it('re-renders with updated value when ability rules change', () => {
    // Start with denied
    hoisted.canMock.mockReturnValue(false);
    const { result } = renderHook(() => useCan('action-id', 'Subject'));
    expect(result.current).toBe(false);

    // Simulate ability.update() which fires the 'update' event
    hoisted.canMock.mockReturnValue(true);
    act(() => {
      registeredCallback?.();
    });

    expect(result.current).toBe(true);
  });

  it('lowercases the subject via lodash before calling ability.can', () => {
    hoisted.canMock.mockReturnValue(true);
    renderHook(() => useCan('action-id', 'View Performance Profiles'));

    // lodash _.lowerCase('View Performance Profiles') => 'view performance profiles'
    expect(hoisted.canMock).toHaveBeenCalledWith('action-id', 'view performance profiles');
  });

  it('re-evaluates immediately in useEffect when rules changed between render and effect', () => {
    // First call during useState init returns false
    // Second call during useEffect returns true (simulating a race)
    hoisted.canMock.mockReturnValueOnce(false).mockReturnValue(true);

    const { result } = renderHook(() => useCan('action', 'subject'));

    // After the effect runs, it should have re-evaluated to true
    expect(result.current).toBe(true);
  });
});
