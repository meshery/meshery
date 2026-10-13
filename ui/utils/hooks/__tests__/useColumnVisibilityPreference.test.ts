import { describe, expect, it, beforeEach, afterEach, vi } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useColumnVisibilityPreference } from '../useColumnVisibilityPreference';

const KEY = 'test-table';
const STORAGE_KEY = `meshery_cols_v2_${KEY}`;
const WIDE = { name: true, created_at: true, owner: true };
const NARROW = { name: true, created_at: false, owner: false };

const stored = () => JSON.parse(localStorage.getItem(STORAGE_KEY) ?? 'null');

describe('useColumnVisibilityPreference', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('starts from the responsive defaults when nothing is stored', () => {
    const { result } = renderHook(() => useColumnVisibilityPreference(KEY, WIDE));

    expect(result.current.columnVisibility).toEqual(WIDE);
    expect(stored()).toBeNull();
  });

  it('persists only the column the user toggled', () => {
    const { result } = renderHook(() => useColumnVisibilityPreference(KEY, WIDE));

    act(() => {
      result.current.setColumnVisibilityByUser((prev) => ({ ...prev, owner: false }));
    });

    expect(result.current.columnVisibility).toEqual({ ...WIDE, owner: false });
    expect(stored()).toEqual({ owner: false });
  });

  it('keeps earlier toggles when the user toggles another column', () => {
    const { result } = renderHook(() => useColumnVisibilityPreference(KEY, WIDE));

    act(() => {
      result.current.setColumnVisibilityByUser((prev) => ({ ...prev, owner: false }));
    });
    act(() => {
      result.current.setColumnVisibilityByUser((prev) => ({ ...prev, name: false }));
    });

    expect(stored()).toEqual({ owner: false, name: false });
  });

  it('lets untouched columns follow responsive defaults after a toggle', () => {
    const { result } = renderHook(() => useColumnVisibilityPreference(KEY, WIDE));

    act(() => {
      result.current.setColumnVisibilityByUser((prev) => ({ ...prev, owner: false }));
    });
    act(() => {
      result.current.setColumnVisibilityByResponsive(NARROW);
    });

    expect(result.current.columnVisibility.created_at).toBe(false);
  });

  it('keeps a user override when the responsive defaults change', () => {
    const { result } = renderHook(() => useColumnVisibilityPreference(KEY, NARROW));

    act(() => {
      result.current.setColumnVisibilityByUser((prev) => ({ ...prev, created_at: true }));
    });
    act(() => {
      result.current.setColumnVisibilityByResponsive(WIDE);
    });
    act(() => {
      result.current.setColumnVisibilityByResponsive(NARROW);
    });

    expect(result.current.columnVisibility).toEqual({ ...NARROW, created_at: true });
  });

  it('rehydrates stored preferences on top of responsive defaults', () => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify({ owner: false }));

    const { result } = renderHook(() => useColumnVisibilityPreference(KEY, WIDE));

    expect(result.current.columnVisibility).toEqual({ ...WIDE, owner: false });
  });

  it('ignores maps saved under the v1 key', () => {
    localStorage.setItem(`meshery_cols_${KEY}`, JSON.stringify(WIDE));

    const { result } = renderHook(() => useColumnVisibilityPreference(KEY, NARROW));

    expect(result.current.columnVisibility).toEqual(NARROW);
  });

  it('removes the v1 map from storage', () => {
    localStorage.setItem(`meshery_cols_${KEY}`, JSON.stringify(WIDE));

    renderHook(() => useColumnVisibilityPreference(KEY, NARROW));

    expect(localStorage.getItem(`meshery_cols_${KEY}`)).toBeNull();
  });

  it('keeps a column toggled back to its default as an explicit choice', () => {
    const { result } = renderHook(() => useColumnVisibilityPreference(KEY, WIDE));

    act(() => {
      result.current.setColumnVisibilityByUser((prev) => ({ ...prev, owner: false }));
    });
    act(() => {
      result.current.setColumnVisibilityByUser((prev) => ({ ...prev, owner: true }));
    });
    act(() => {
      result.current.setColumnVisibilityByResponsive(NARROW);
    });

    expect(stored()).toEqual({ owner: true });
    expect(result.current.columnVisibility.owner).toBe(true);
  });

  it('still loads stored preferences when localStorage access throws', () => {
    vi.spyOn(Storage.prototype, 'removeItem').mockImplementation(() => {
      throw new Error('SecurityError');
    });

    const { result } = renderHook(() => useColumnVisibilityPreference(KEY, WIDE));

    expect(result.current.columnVisibility).toEqual(WIDE);
  });

  it('falls back to responsive defaults when storage holds malformed JSON', () => {
    localStorage.setItem(STORAGE_KEY, '{not json');

    const { result } = renderHook(() => useColumnVisibilityPreference(KEY, WIDE));

    expect(result.current.columnVisibility).toEqual(WIDE);
  });
});
