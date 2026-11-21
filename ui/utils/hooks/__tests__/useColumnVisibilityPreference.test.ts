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
    const responsiveDefaults: Record<string, boolean | undefined> = Object.create(null);
    next.__proto__ = false;
    responsiveDefaults.__proto__ = true;

    const { result, unmount } = renderHook(() =>
      useColumnVisibilityPreference('prototype-columns', responsiveDefaults),
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
      useColumnVisibilityPreference('prototype-columns', responsiveDefaults),
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

  it('does not restore persisted columns excluded by current responsive defaults', () => {
    window.localStorage.setItem(
      'meshery_cols_hidden-layout',
      JSON.stringify({ available: false, visibility: false }),
    );

    const remounted = renderHook(() =>
      useColumnVisibilityPreference('hidden-layout', { available: true }),
    );

    expect(remounted.result.current.columnVisibility.available).toBe(false);
    expect(
      Object.prototype.hasOwnProperty.call(remounted.result.current.columnVisibility, 'visibility'),
    ).toBe(false);

    act(() => {
      remounted.result.current.setColumnVisibilityByUser({
        available: true,
      });
    });

    const stored = JSON.parse(
      window.localStorage.getItem('meshery_cols_hidden-layout') ?? '{}',
    ) as Record<string, boolean>;
    expect(Object.prototype.hasOwnProperty.call(stored, 'visibility')).toBe(false);
    expect(
      Object.prototype.hasOwnProperty.call(remounted.result.current.columnVisibility, 'visibility'),
    ).toBe(false);

    act(() => {
      remounted.result.current.setColumnVisibilityByResponsive({
        available: true,
      });
    });

    expect(remounted.result.current.columnVisibility.available).toBe(true);
    expect(
      Object.prototype.hasOwnProperty.call(remounted.result.current.columnVisibility, 'visibility'),
    ).toBe(false);
  });
});
