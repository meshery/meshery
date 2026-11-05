import { renderHook, act } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const routerState = {
  query: {} as Record<string, any>,
  pathname: '/registry',
  route: '/registry',
  push: vi.fn(),
};

vi.mock('next/router', () => ({
  useRouter: () => routerState,
}));

vi.mock('../../constants/navigator', () => ({
  OVERVIEW: 'Overview',
  MODELS: 'Models',
  GRAFANA: 'Grafana',
  PROMETHEUS: 'Prometheus',
  REGISTRY: 'Registry',
}));

import { useMeshModelComponentRouter, useInfiniteScrollRef, useRegistryRouter } from './hooks';

describe('useMeshModelComponentRouter', () => {
  beforeEach(() => {
    routerState.query = {};
    routerState.push = vi.fn();
  });

  it('redirects to tab=Models when settingsCategory is Registry and no tab is set', () => {
    routerState.query = { settingsCategory: 'Registry' };

    renderHook(() => useMeshModelComponentRouter());

    expect(routerState.push).toHaveBeenCalledWith(
      expect.objectContaining({
        query: expect.objectContaining({ tab: 'Models' }),
      }),
    );
  });

  it('does not redirect when a tab is already set', () => {
    routerState.query = { settingsCategory: 'Registry', tab: 'Components' };

    renderHook(() => useMeshModelComponentRouter());

    expect(routerState.push).not.toHaveBeenCalled();
  });

  it('returns OVERVIEW when tab is GRAFANA or PROMETHEUS', () => {
    routerState.query = { tab: 'Grafana' };
    const { result } = renderHook(() => useMeshModelComponentRouter());
    expect(result.current.selectedTab).toBe('Overview');

    routerState.query = { tab: 'Prometheus' };
    const { result: result2 } = renderHook(() => useMeshModelComponentRouter());
    expect(result2.current.selectedTab).toBe('Overview');
  });

  it('returns the tab as-is otherwise', () => {
    routerState.query = { tab: 'Models' };
    const { result } = renderHook(() => useMeshModelComponentRouter());
    expect(result.current.selectedTab).toBe('Models');
  });

  it('defaults page size to 25 when not in query', () => {
    routerState.query = {};
    const { result } = renderHook(() => useMeshModelComponentRouter());
    expect(result.current.selectedPageSize).toBe(25);
  });

  it('uses query.pagesize when set', () => {
    routerState.query = { pagesize: 100 };
    const { result } = renderHook(() => useMeshModelComponentRouter());
    expect(result.current.selectedPageSize).toBe(100);
  });

  it('returns null searchQuery by default', () => {
    routerState.query = {};
    const { result } = renderHook(() => useMeshModelComponentRouter());
    expect(result.current.searchQuery).toBeNull();
  });

  it('uses query.searchText for searchQuery when present', () => {
    routerState.query = { searchText: 'kube' };
    const { result } = renderHook(() => useMeshModelComponentRouter());
    expect(result.current.searchQuery).toBe('kube');
  });
});

describe('useRegistryRouter', () => {
  beforeEach(() => {
    routerState.query = {};
    routerState.push = vi.fn();
    routerState.route = '/registry';
  });

  it('exposes query values directly', () => {
    routerState.query = {
      settingsCategory: 'Cat',
      tab: 'TabA',
      selectedItemUUID: 'u-1',
      searchText: 'find me',
    };

    const { result } = renderHook(() => useRegistryRouter());

    expect(result.current.settingsCategory).toBe('Cat');
    expect(result.current.tab).toBe('TabA');
    expect(result.current.selectedItemUUID).toBe('u-1');
    expect(result.current.filters).toEqual({ searchText: 'find me' });
  });

  it('defaults selectedItemUUID and searchText when not set', () => {
    const { result } = renderHook(() => useRegistryRouter());

    expect(result.current.selectedItemUUID).toBe('');
    expect(result.current.filters).toEqual({ searchText: null });
  });

  it('handleUpdateSelectedRoute pushes a constructed query string', () => {
    routerState.query = { settingsCategory: 'Reg', tab: 'Models' };

    const { result } = renderHook(() => useRegistryRouter());

    act(() => {
      result.current.handleUpdateSelectedRoute(['node-1'], { searchText: 'q' });
    });

    expect(routerState.push).toHaveBeenCalledTimes(1);
    const [path, , opts] = routerState.push.mock.calls[0];
    expect(typeof path).toBe('string');
    expect(path).toContain('settingsCategory=Reg');
    expect(path).toContain('tab=Models');
    expect(path).toContain('selectedItemUUID=node-1');
    expect(path).toContain('searchText=q');
    expect(opts).toEqual({ shallow: true });
  });

  it('falls back to Registry/Models when no settingsCategory/tab set', () => {
    routerState.query = {};

    const { result } = renderHook(() => useRegistryRouter());

    act(() => {
      result.current.handleUpdateSelectedRoute(['node-1'], {});
    });

    const [path] = routerState.push.mock.calls[0];
    expect(path).toContain('settingsCategory=Registry');
    expect(path).toContain('tab=Models');
  });
});

describe('useInfiniteScrollRef', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach?.(() => {
    vi.useRealTimers();
  });

  it('returns a triggerRef without crashing when no DOM target is attached', () => {
    const cb = vi.fn();
    const { result } = renderHook(() => useInfiniteScrollRef(cb));

    expect(result.current).toEqual({ current: null });
    // Run pending setTimeout
    act(() => {
      vi.advanceTimersByTime(0);
    });
    // callback should not yet be invoked because triggerRef is unattached
    expect(cb).not.toHaveBeenCalled();
    vi.useRealTimers();
  });

  it('observes the trigger element when it is set', () => {
    const observe = vi.fn();
    const disconnect = vi.fn();
    const recorded: Array<(entries: any[]) => void> = [];
    class ObserverMock {
      observe = observe;
      disconnect = disconnect;
      _cb: (entries: any[]) => void;
      constructor(cbFn: (entries: any[]) => void) {
        this._cb = cbFn;
        recorded.push(cbFn);
      }
    }
    // @ts-expect-error – overriding for the test
    globalThis.IntersectionObserver = ObserverMock;

    const cb = vi.fn();
    const { result } = renderHook(() => useInfiniteScrollRef(cb));

    const el = document.createElement('div');
    document.body.appendChild(el);
    (result.current as any).current = el;

    // Effect uses a 0ms setTimeout to ensure the DOM has been attached.
    act(() => {
      vi.advanceTimersByTime(1);
    });

    expect(observe).toHaveBeenCalledWith(el);

    // Simulate intersection -> callback fires.
    recorded[0]([{ isIntersecting: true }]);
    expect(cb).toHaveBeenCalled();

    document.body.removeChild(el);
    vi.useRealTimers();
  });
});
