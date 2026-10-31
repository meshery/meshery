import { renderHook, act } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import usePinnedNavItems, {
  loadPinnedNavItems,
  savePinnedNavItems,
  PINNED_NAV_ITEMS_STORAGE_KEY,
} from '../usePinnedNavItems';

describe('pinned nav items storage helpers', () => {
  beforeEach(() => {
    window.localStorage.clear();
  });

  describe('loadPinnedNavItems', () => {
    it('returns empty array when nothing is stored', () => {
      expect(loadPinnedNavItems()).toEqual([]);
    });

    it('loads valid array of string IDs', () => {
      window.localStorage.setItem(
        PINNED_NAV_ITEMS_STORAGE_KEY,
        JSON.stringify(['dashboard', 'designs']),
      );
      expect(loadPinnedNavItems()).toEqual(['dashboard', 'designs']);
    });

    it('deduplicates stored array entries', () => {
      window.localStorage.setItem(
        PINNED_NAV_ITEMS_STORAGE_KEY,
        JSON.stringify(['dashboard', 'designs', 'dashboard']),
      );
      expect(loadPinnedNavItems()).toEqual(['dashboard', 'designs']);
    });

    it('falls back to empty array on corrupt/non-JSON string', () => {
      window.localStorage.setItem(PINNED_NAV_ITEMS_STORAGE_KEY, 'not-valid-json{[');
      expect(loadPinnedNavItems()).toEqual([]);
    });

    it('falls back to empty array if stored data is not an array', () => {
      window.localStorage.setItem(
        PINNED_NAV_ITEMS_STORAGE_KEY,
        JSON.stringify({ id: 'dashboard' }),
      );
      expect(loadPinnedNavItems()).toEqual([]);

      window.localStorage.setItem(PINNED_NAV_ITEMS_STORAGE_KEY, JSON.stringify(12345));
      expect(loadPinnedNavItems()).toEqual([]);
    });

    it('filters out non-string entries from array instead of discarding the whole list', () => {
      window.localStorage.setItem(
        PINNED_NAV_ITEMS_STORAGE_KEY,
        JSON.stringify(['dashboard', 123, null, 'designs', undefined, {}]),
      );
      expect(loadPinnedNavItems()).toEqual(['dashboard', 'designs']);

      window.localStorage.setItem(PINNED_NAV_ITEMS_STORAGE_KEY, JSON.stringify([1, 2, 3]));
      expect(loadPinnedNavItems()).toEqual([]);
    });

    it('handles localStorage read errors gracefully', () => {
      const getItemSpy = vi.spyOn(Storage.prototype, 'getItem').mockImplementationOnce(() => {
        throw new Error('Access denied');
      });
      expect(loadPinnedNavItems()).toEqual([]);
      getItemSpy.mockRestore();
    });
  });

  describe('savePinnedNavItems', () => {
    it('persists unique array of strings to localStorage', () => {
      savePinnedNavItems(['dashboard', 'catalog']);
      const raw = window.localStorage.getItem(PINNED_NAV_ITEMS_STORAGE_KEY);
      expect(JSON.parse(raw!)).toEqual(['dashboard', 'catalog']);
    });

    it('deduplicates duplicate IDs before saving', () => {
      savePinnedNavItems(['catalog', 'designs', 'catalog']);
      const raw = window.localStorage.getItem(PINNED_NAV_ITEMS_STORAGE_KEY);
      expect(JSON.parse(raw!)).toEqual(['catalog', 'designs']);
    });

    it('filters out any non-string elements', () => {
      savePinnedNavItems(['dashboard', 123 as any, null as any, 'designs']);
      const raw = window.localStorage.getItem(PINNED_NAV_ITEMS_STORAGE_KEY);
      expect(JSON.parse(raw!)).toEqual(['dashboard', 'designs']);
    });

    it('ignores non-array inputs', () => {
      savePinnedNavItems('not-an-array' as any);
      expect(window.localStorage.getItem(PINNED_NAV_ITEMS_STORAGE_KEY)).toBeNull();
    });

    it('gracefully handles quota or storage write errors', () => {
      const setItemSpy = vi.spyOn(Storage.prototype, 'setItem').mockImplementationOnce(() => {
        throw new Error('QuotaExceededError');
      });
      expect(() => savePinnedNavItems(['dashboard'])).not.toThrow();
      setItemSpy.mockRestore();
    });
  });
});

describe('usePinnedNavItems hook', () => {
  beforeEach(() => {
    window.localStorage.clear();
  });

  it('initializes with items from localStorage', () => {
    window.localStorage.setItem(
      PINNED_NAV_ITEMS_STORAGE_KEY,
      JSON.stringify(['dashboard', 'environments']),
    );

    const { result } = renderHook(() => usePinnedNavItems());
    expect(result.current.pinnedItemIds).toEqual(['dashboard', 'environments']);
    expect(result.current.isPinned('dashboard')).toBe(true);
    expect(result.current.isPinned('environments')).toBe(true);
    expect(result.current.isPinned('catalog')).toBe(false);
  });

  it('pins an unpinned item and updates localStorage', () => {
    const { result } = renderHook(() => usePinnedNavItems());
    expect(result.current.pinnedItemIds).toEqual([]);

    act(() => {
      result.current.togglePin('designs');
    });

    expect(result.current.pinnedItemIds).toEqual(['designs']);
    expect(result.current.isPinned('designs')).toBe(true);
    expect(JSON.parse(window.localStorage.getItem(PINNED_NAV_ITEMS_STORAGE_KEY)!)).toEqual([
      'designs',
    ]);
  });

  it('unpins an already pinned item and updates localStorage', () => {
    window.localStorage.setItem(
      PINNED_NAV_ITEMS_STORAGE_KEY,
      JSON.stringify(['designs', 'catalog']),
    );

    const { result } = renderHook(() => usePinnedNavItems());
    expect(result.current.isPinned('designs')).toBe(true);

    act(() => {
      result.current.togglePin('designs');
    });

    expect(result.current.pinnedItemIds).toEqual(['catalog']);
    expect(result.current.isPinned('designs')).toBe(false);
    expect(result.current.isPinned('catalog')).toBe(true);
    expect(JSON.parse(window.localStorage.getItem(PINNED_NAV_ITEMS_STORAGE_KEY)!)).toEqual([
      'catalog',
    ]);
  });

  it('toggles item off when called a second time', () => {
    const { result } = renderHook(() => usePinnedNavItems());

    act(() => {
      result.current.togglePin('connections');
    });
    expect(result.current.pinnedItemIds).toEqual(['connections']);

    act(() => {
      result.current.togglePin('connections');
    });
    expect(result.current.pinnedItemIds).toEqual([]);
  });

  it('handles two togglePin calls in one act() correctly', () => {
    const { result } = renderHook(() => usePinnedNavItems());

    act(() => {
      result.current.togglePin('designs');
      result.current.togglePin('catalog');
    });

    expect(result.current.pinnedItemIds).toEqual(['designs', 'catalog']);
    expect(result.current.isPinned('designs')).toBe(true);
    expect(result.current.isPinned('catalog')).toBe(true);
  });

  it('handles two togglePin calls on the same item in one act()', () => {
    const { result } = renderHook(() => usePinnedNavItems());

    act(() => {
      result.current.togglePin('designs');
      result.current.togglePin('designs');
    });

    expect(result.current.pinnedItemIds).toEqual([]);
    expect(result.current.isPinned('designs')).toBe(false);
  });

  it('ignores invalid non-string or empty input to togglePin', () => {
    const { result } = renderHook(() => usePinnedNavItems());

    act(() => {
      result.current.togglePin('');
      result.current.togglePin(null as any);
      result.current.togglePin(undefined as any);
    });

    expect(result.current.pinnedItemIds).toEqual([]);
  });
});
