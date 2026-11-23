import { describe, expect, it, vi, beforeEach } from 'vitest';
import { bindDock, publishDockState, restoreConsoles, useConsoleDock } from '../dock-store';
import { renderHook, act } from '@testing-library/react';

describe('dock-store', () => {
  beforeEach(() => {
    publishDockState({ count: 0, minimized: false });
  });

  it('publishes dock state to subscribers', () => {
    const { result } = renderHook(() => useConsoleDock());

    act(() => publishDockState({ count: 2, minimized: true }));

    expect(result.current).toEqual({ count: 2, minimized: true });
  });

  it('keeps the snapshot identity stable when nothing changed', () => {
    const { result } = renderHook(() => useConsoleDock());

    act(() => publishDockState({ count: 1, minimized: false }));
    const first = result.current;

    // useSyncExternalStore re-renders on any snapshot whose identity differs, so
    // republishing an equal state must not hand back a new object.
    act(() => publishDockState({ count: 1, minimized: false }));

    expect(result.current).toBe(first);
  });

  it('restores through the bound provider', () => {
    const restore = vi.fn();
    bindDock(restore);

    restoreConsoles();

    expect(restore).toHaveBeenCalledOnce();
  });

  it('clears published state when the provider unmounts', () => {
    const unbind = bindDock(vi.fn());
    publishDockState({ count: 3, minimized: true });

    unbind();

    const { result } = renderHook(() => useConsoleDock());
    expect(result.current).toEqual({ count: 0, minimized: false });
    // A button that outlived its panel must not be able to restore it.
    expect(restoreConsoles()).toBeUndefined();
  });

  it('leaves a newer provider bound when an older one unmounts after it', () => {
    // StrictMode double-mounts: the second provider binds before the first one's
    // cleanup runs, and that cleanup must not unbind the live provider.
    const stale = vi.fn();
    const live = vi.fn();

    const unbindStale = bindDock(stale);
    bindDock(live);
    unbindStale();

    restoreConsoles();

    expect(stale).not.toHaveBeenCalled();
    expect(live).toHaveBeenCalledOnce();
  });
});
