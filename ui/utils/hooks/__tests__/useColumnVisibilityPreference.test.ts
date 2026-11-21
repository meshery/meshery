import { act, renderHook } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { useColumnVisibilityPreference } from '../useColumnVisibilityPreference';

describe('useColumnVisibilityPreference', () => {
  beforeEach(() => {
    window.localStorage.clear();
  });

  afterEach(() => {
    window.localStorage.clear();
  });

  it('persists and restores an own __proto__ visibility key', () => {
    const next: Record<string, boolean | undefined> = Object.create(null);
    next.__proto__ = false;

    const { result, unmount } = renderHook(() =>
      useColumnVisibilityPreference('prototype-columns', { __proto__: true }),
    );

    act(() => {
      result.current.setColumnVisibilityByUser(next);
    });

    const stored = JSON.parse(
      window.localStorage.getItem('meshery_cols_prototype-columns') ?? '{}',
    ) as Record<string, boolean>;
    expect(Object.prototype.hasOwnProperty.call(stored, '__proto__')).toBe(true);
    expect(stored.__proto__).toBe(false);

    unmount();
    const remounted = renderHook(() =>
      useColumnVisibilityPreference('prototype-columns', { __proto__: true }),
    );

    expect(
      Object.prototype.hasOwnProperty.call(remounted.result.current.columnVisibility, '__proto__'),
    ).toBe(true);
    expect(remounted.result.current.columnVisibility.__proto__).toBe(false);
  });

  it('ignores malformed persisted values and falls back to responsive defaults', () => {
    window.localStorage.setItem(
      'meshery_cols_existing',
      JSON.stringify({ name: 'stored-as-is', enabled: null, active: false }),
    );

    const remounted = renderHook(() =>
      useColumnVisibilityPreference('existing', { name: true, enabled: true, active: true }),
    );

    expect(remounted.result.current.columnVisibility.name).toBe(true);
    expect(remounted.result.current.columnVisibility.enabled).toBe(true);
    expect(remounted.result.current.columnVisibility.active).toBe(false);
  });
});
