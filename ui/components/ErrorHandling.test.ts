import { renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const notify = vi.fn();

vi.mock('../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('../lib/event-types', () => ({
  EVENT_TYPES: {
    ERROR: { type: 'error' },
  },
}));

import HandleError from './ErrorHandling';

describe('HandleError (ErrorHandling)', () => {
  beforeEach(() => {
    notify.mockClear();
    vi.spyOn(console, 'error').mockImplementation(() => {});
  });

  it('returns a stable error handler function', () => {
    const { result } = renderHook(() => HandleError());
    expect(typeof result.current).toBe('function');
  });

  it('forwards Error instances by message and surfaces notify metadata', () => {
    const { result } = renderHook(() => HandleError());
    const err = new Error('boom');
    result.current(err, 'failed to load', 'error');

    expect(notify).toHaveBeenCalledTimes(1);
    const [payload] = notify.mock.calls[0];
    expect(payload.message).toBe('failed to load: boom');
    expect(payload.event_type).toEqual({ type: 'error' });
    expect(payload.details).toBe(err.toString());
  });

  it('forwards a plain string error message verbatim', () => {
    const { result } = renderHook(() => HandleError());
    result.current('something went wrong', 'oops', 'warning');

    expect(notify).toHaveBeenCalledTimes(1);
    const [payload] = notify.mock.calls[0];
    expect(payload.message).toBe('oops: something went wrong');
    expect(payload.details).toBe('something went wrong');
  });

  it('logs the severity and the error to the console', () => {
    const consoleSpy = vi.spyOn(console, 'error');
    consoleSpy.mockClear();
    const { result } = renderHook(() => HandleError());
    result.current(new Error('e1'), 'prefix', 'warning');

    // Find the explicit "an error occured" call (avoid React internal noise).
    const ours = consoleSpy.mock.calls.find(
      (args) => typeof args[0] === 'string' && args[0].includes('an error occured with severity'),
    );
    expect(ours).toBeDefined();
    expect(ours![1]).toBe('warning');
  });
});
