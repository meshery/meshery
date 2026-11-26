import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';

// Mock the next/router singleton so we can observe the in-place `change`
// override the hook installs on `SingletonRouter.router`.

const hoisted = vi.hoisted(() => {
  const routerInstance: { change?: (...args: unknown[]) => unknown } = {};
  const singleton = { router: routerInstance };
  const RouterPrototype: { change: (...args: unknown[]) => Promise<boolean> } = {
    change: vi.fn().mockResolvedValue(true),
  };
  return { routerInstance, singleton, RouterPrototype };
});
const routerInstance = hoisted.routerInstance;
const RouterPrototype = hoisted.RouterPrototype;

vi.mock('next/router', () => ({
  default: hoisted.singleton,
  Router: { prototype: hoisted.RouterPrototype },
}));

import usePreventUserFromLeavingPage from '../usePreventUserFromLeavingPage';

describe('usePreventUserFromLeavingPage', () => {
  let originalOnBeforeUnload: typeof window.onbeforeunload;
  let originalConfirm: typeof window.confirm;

  beforeEach(() => {
    originalOnBeforeUnload = window.onbeforeunload;
    originalConfirm = window.confirm;
    window.onbeforeunload = null;
    window.confirm = vi.fn(() => true);
    routerInstance.change = undefined;
    (RouterPrototype.change as unknown as ReturnType<typeof vi.fn>).mockClear();
  });

  afterEach(() => {
    window.onbeforeunload = originalOnBeforeUnload;
    window.confirm = originalConfirm;
  });

  it('installs a beforeunload prompt and a router.change interceptor when enabled', () => {
    renderHook(() => usePreventUserFromLeavingPage(true));

    expect(typeof window.onbeforeunload).toBe('function');
    const msg = (window.onbeforeunload as () => string)();
    expect(msg).toMatch(/unsaved changes/i);

    expect(typeof routerInstance.change).toBe('function');
  });

  it('skips installation when starting disabled and leaves handler null on unmount', () => {
    const { unmount } = renderHook(() => usePreventUserFromLeavingPage(false));

    expect(window.onbeforeunload).toBeNull();
    expect(routerInstance.change).toBeUndefined();

    unmount();
    expect(window.onbeforeunload).toBeNull();
  });

  it('routes navigation through prototype.change when the user confirms', async () => {
    renderHook(() => usePreventUserFromLeavingPage(true));
    (window.confirm as ReturnType<typeof vi.fn>).mockReturnValue(true);

    const ret = await routerInstance.change!('/new', '/new', {});
    expect(window.confirm).toHaveBeenCalled();
    expect(RouterPrototype.change).toHaveBeenCalledTimes(1);
    expect(ret).toBe(true);
  });

  it('blocks navigation when the user cancels', async () => {
    renderHook(() => usePreventUserFromLeavingPage(true));
    (window.confirm as ReturnType<typeof vi.fn>).mockReturnValue(false);

    const ret = await routerInstance.change!('/new');
    expect(RouterPrototype.change).not.toHaveBeenCalled();
    expect(ret).toBe(false);
  });

  it('blocks navigation fail-safe when window.confirm is unavailable', async () => {
    renderHook(() => usePreventUserFromLeavingPage(true));
    // @ts-expect-error simulating missing confirm function
    window.confirm = undefined;

    const ret = await routerInstance.change!('/new');
    expect(RouterPrototype.change).not.toHaveBeenCalled();
    expect(ret).toBe(false);
  });

  it('does not clear external handlers when unmounting a disabled or superseded instance', () => {
    const first = renderHook(() => usePreventUserFromLeavingPage(true));
    const externalUnload = vi.fn();
    const externalChange = vi.fn();
    window.onbeforeunload = externalUnload;
    routerInstance.change = externalChange;

    first.unmount();
    expect(window.onbeforeunload).toBe(externalUnload);
    expect(routerInstance.change).toBe(externalChange);

    const disabled = renderHook(() => usePreventUserFromLeavingPage(false));
    routerInstance.change = externalChange;
    window.onbeforeunload = externalUnload;
    disabled.unmount();
    expect(window.onbeforeunload).toBe(externalUnload);
    expect(routerInstance.change).toBe(externalChange);
  });

  it('clears the router.change override on unmount', () => {
    const { unmount } = renderHook(() => usePreventUserFromLeavingPage(true));
    expect(routerInstance.change).toBeDefined();
    unmount();
    expect(routerInstance.change).toBeUndefined();
  });

  it('clears window.onbeforeunload on unmount', () => {
    const { unmount } = renderHook(() => usePreventUserFromLeavingPage(true));
    expect(typeof window.onbeforeunload).toBe('function');
    unmount();
    expect(window.onbeforeunload).toBeNull();
  });

  it('safely handles missing or nullish SingletonRouter.router without throwing', () => {
    const originalRouter = hoisted.singleton.router;
    try {
      // @ts-expect-error testing runtime missing router
      hoisted.singleton.router = null;
      expect(() => {
        const { unmount } = renderHook(() => usePreventUserFromLeavingPage(true));
        unmount();
      }).not.toThrow();
    } finally {
      hoisted.singleton.router = originalRouter;
    }
  });

  it('returns a setter that toggles the prevention state and leaves handler null after disabling and unmount', () => {
    const { result, unmount } = renderHook(() => usePreventUserFromLeavingPage(false));
    expect(routerInstance.change).toBeUndefined();
    expect(window.onbeforeunload).toBeNull();

    act(() => {
      result.current(true);
    });

    expect(typeof routerInstance.change).toBe('function');
    expect(typeof window.onbeforeunload).toBe('function');

    act(() => {
      result.current(false);
    });
    expect(routerInstance.change).toBeUndefined();
    expect(window.onbeforeunload).toBeNull();

    unmount();
    expect(window.onbeforeunload).toBeNull();
  });
});
