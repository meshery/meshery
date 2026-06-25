import { useCallback, useRef, useState, type SetStateAction } from 'react';

const STORAGE_KEY_PREFIX = 'meshery_cols_';

function loadFromStorage(tableKey: string): Record<string, boolean> {
  if (typeof window === 'undefined') return {};
  try {
    const raw = window.localStorage.getItem(`${STORAGE_KEY_PREFIX}${tableKey}`);
    return raw ? (JSON.parse(raw) as Record<string, boolean>) : {};
  } catch {
    return {};
  }
}

function saveToStorage(tableKey: string, prefs: Record<string, boolean>): void {
  if (typeof window === 'undefined') return;
  try {
    window.localStorage.setItem(`${STORAGE_KEY_PREFIX}${tableKey}`, JSON.stringify(prefs));
  } catch {
    // Quota exceeded or private browsing — degrade silently.
  }
}

/**
 * Persistent column visibility for data tables.
 *
 * User-initiated column toggles (via `setColumnVisibilityByUser`) are saved to
 * localStorage so they survive navigation and refresh. Responsive-layout
 * updates (via `setColumnVisibilityByResponsive`) are NOT saved — the
 * user's explicit overrides are layered on top of each responsive update.
 *
 * @param tableKey - Stable, unique key for this table (e.g. "connections").
 * @param responsiveDefaults - Breakpoint-computed defaults from
 *   `getResponsiveColumnVisibility` / `updateVisibleColumns`.
 */
export function useColumnVisibilityPreference(
  tableKey: string,
  responsiveDefaults: Record<string, boolean | undefined>,
): {
  columnVisibility: Record<string, boolean | undefined>;
  /** Call from the column-visibility toolbar. Persists to localStorage. */
  setColumnVisibilityByUser: (updater: SetStateAction<Record<string, boolean | undefined>>) => void;
  /** Call from width-change effects. Merges responsive defaults with stored user prefs. */
  setColumnVisibilityByResponsive: (next: Record<string, boolean | undefined>) => void;
} {
  // Stored user preferences — only the columns the user explicitly toggled.
  const storedPrefs = useRef<Record<string, boolean>>(loadFromStorage(tableKey));

  const [columnVisibility, setColumnVisibility] = useState<Record<string, boolean | undefined>>(
    () => ({ ...responsiveDefaults, ...storedPrefs.current }),
  );

  const setColumnVisibilityByUser = useCallback(
    (updater: SetStateAction<Record<string, boolean | undefined>>) => {
      setColumnVisibility((prev) => {
        const next = typeof updater === 'function' ? updater(prev) : updater;
        // Persist only defined (non-undefined) values.
        const toStore: Record<string, boolean> = {};
        for (const [col, val] of Object.entries(next)) {
          if (val !== undefined) toStore[col] = val;
        }
        storedPrefs.current = toStore;
        saveToStorage(tableKey, toStore);
        return next;
      });
    },
    [tableKey],
  );

  const setColumnVisibilityByResponsive = useCallback(
    (next: Record<string, boolean | undefined>) => {
      setColumnVisibility((prev) => {
        const merged = { ...next, ...storedPrefs.current };
        const keys = new Set([...Object.keys(prev), ...Object.keys(merged)]);
        for (const key of keys) {
          if (prev[key] !== merged[key]) return merged;
        }
        return prev;
      });
    },
    [],
  );

  return { columnVisibility, setColumnVisibilityByUser, setColumnVisibilityByResponsive };
}
