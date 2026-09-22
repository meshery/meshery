import { act, renderHook } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { EVENT_TYPES } from '../../../lib/event-types';
import useDelayedYamlErrorNotify, { DESIGN_YAML_ERROR_IDLE_MS } from './useDelayedYamlErrorNotify';

const notify = vi.fn();

vi.mock('../../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

describe('useDelayedYamlErrorNotify', () => {
  beforeEach(() => {
    notify.mockReset();
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('notifies once after typing pauses on still-invalid YAML', () => {
    const { result } = renderHook(() => useDelayedYamlErrorNotify());

    act(() => {
      result.current.scheduleInvalidYamlCheck('name: [');
      result.current.scheduleInvalidYamlCheck('name: [a');
    });

    expect(notify).not.toHaveBeenCalled();

    act(() => {
      vi.advanceTimersByTime(DESIGN_YAML_ERROR_IDLE_MS);
    });

    expect(notify).toHaveBeenCalledTimes(1);
    expect(notify).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Invalid Yaml Data',
        event_type: EVENT_TYPES.ERROR,
      }),
    );
  });

  it('does not notify when YAML becomes valid before the idle window ends', () => {
    const { result } = renderHook(() => useDelayedYamlErrorNotify());

    act(() => {
      result.current.scheduleInvalidYamlCheck('name: [');
      result.current.cancelPendingYamlError();
    });

    act(() => {
      vi.advanceTimersByTime(DESIGN_YAML_ERROR_IDLE_MS);
    });

    expect(notify).not.toHaveBeenCalled();
  });

  it('cancels a pending notification when cancel is called after scheduling', () => {
    const { result } = renderHook(() => useDelayedYamlErrorNotify());

    act(() => {
      result.current.scheduleInvalidYamlCheck('name: [');
    });

    act(() => {
      vi.advanceTimersByTime(DESIGN_YAML_ERROR_IDLE_MS - 1);
      result.current.cancelPendingYamlError();
      vi.advanceTimersByTime(1);
    });

    expect(notify).not.toHaveBeenCalled();
  });
});
