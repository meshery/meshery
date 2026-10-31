import { useCallback, useState } from 'react';

export const PINNED_NAV_ITEMS_STORAGE_KEY = 'meshery_pinned_nav_items';

/**
 * Loads the user's pinned navigator item IDs from localStorage.
 *
 * Validates that the stored value is an array of strings, deduplicates entries
 * using a Set, and falls back to an empty array on corruption or in SSR.
 *
 * @returns {string[]} An array of unique pinned item IDs.
 */
export function loadPinnedNavItems(): string[] {
  if (typeof window === 'undefined') return [];
  try {
    const raw = window.localStorage.getItem(PINNED_NAV_ITEMS_STORAGE_KEY);
    if (!raw) return [];
    const parsed = JSON.parse(raw);
    if (Array.isArray(parsed)) {
      const validStrings = parsed.filter(
        (id): id is string => typeof id === 'string' && id.trim().length > 0,
      );
      return Array.from(new Set(validStrings));
    }
    return [];
  } catch {
    return [];
  }
}

/**
 * Persists the user's pinned navigator item IDs to localStorage.
 *
 * Validates and deduplicates the item IDs, silently degrading if storage
 * quota is exceeded or in private browsing mode.
 *
 * @param {string[]} itemIds Array of navigator item IDs to persist.
 */
export function savePinnedNavItems(itemIds: string[]): void {
  if (typeof window === 'undefined' || !Array.isArray(itemIds)) return;
  try {
    const valid = Array.from(new Set(itemIds.filter((id) => typeof id === 'string')));
    window.localStorage.setItem(PINNED_NAV_ITEMS_STORAGE_KEY, JSON.stringify(valid));
  } catch {
    // Quota exceeded or private browsing — degrade silently.
  }
}

/**
 * Hook to manage pinned navigator items with localStorage persistence.
 */
export function usePinnedNavItems(): {
  pinnedItemIds: string[];
  togglePin: (id: string) => void;
  isPinned: (id: string) => boolean;
} {
  const [pinnedItemIds, setPinnedItemIds] = useState<string[]>(() => loadPinnedNavItems());

  const togglePin = useCallback((id: string) => {
    if (!id || typeof id !== 'string') return;
    setPinnedItemIds((prev) => {
      const nextSet = new Set(prev);
      if (nextSet.has(id)) {
        nextSet.delete(id);
      } else {
        nextSet.add(id);
      }
      const next = Array.from(nextSet);
      savePinnedNavItems(next);
      return next;
    });
  }, []);

  const isPinned = useCallback(
    (id: string) => {
      return pinnedItemIds.includes(id);
    },
    [pinnedItemIds],
  );

  return {
    pinnedItemIds,
    togglePin,
    isPinned,
  };
}

export default usePinnedNavItems;
