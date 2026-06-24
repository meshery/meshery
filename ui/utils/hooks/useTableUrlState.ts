import { useCallback, useMemo, useRef } from 'react';
import { useRouter } from 'next/router';

/**
 * Persistent table state backed by URL query parameters.
 *
 * Each table gets a short `tableKey` prefix (e.g. "con" for connections,
 * "des" for designs) so multiple tables can co-exist on the same page without
 * colliding. Pagination, sort, and search are serialised as:
 *
 *   <key>_page  – current page (omitted when 0)
 *   <key>_ps    – page size   (omitted when equal to the default)
 *   <key>_sort  – sort string, e.g. "name asc"
 *   <key>_q     – search query
 *
 * Additional per-table filter fields are passed through `defaults.filters`
 * and appear as `<key>_<fieldName>` params.
 *
 * Deep-links to a specific row use a caller-provided `rowParam` (defaults to
 * `<key>_row`). The parent component writes the param when a row is expanded;
 * `copyRowDeepLink` writes the URL to the clipboard.
 */

export interface TableUrlState<F extends Record<string, string> = Record<string, string>> {
  page: number;
  pageSize: number;
  sortOrder: string;
  search: string;
  filters: F;
}

interface UseTableUrlStateOptions<F extends Record<string, string>> {
  /** Short, stable identifier, e.g. "con", "des", "fil". */
  tableKey: string;
  /**
   * Optional name for the row-id URL parameter used by deeplinks.
   * Defaults to `<tableKey>_row`.
   */
  rowParam?: string;
  defaults: {
    page?: number;
    pageSize?: number;
    sortOrder?: string;
    search?: string;
    /** Default values for each custom filter field. */
    filters?: F;
  };
}

function parseIntSafe(value: string | string[] | undefined, fallback: number): number {
  const n = parseInt(value as string, 10);
  return Number.isFinite(n) ? n : fallback;
}

function decodeParam(value: string | string[] | undefined): string {
  return typeof value === 'string' ? decodeURIComponent(value) : '';
}

export function useTableUrlState<F extends Record<string, string> = Record<string, string>>(
  options: UseTableUrlStateOptions<F>,
) {
  const router = useRouter();
  const { tableKey, defaults, rowParam } = options;
  const prefix = `${tableKey}_`;
  const resolvedRowParam = rowParam ?? `${prefix}row`;

  // Keep a stable ref of the router so callbacks don't re-create on every render.
  const routerRef = useRef(router);
  routerRef.current = router;

  const tableState = useMemo<TableUrlState<F>>(() => {
    const { query } = router;

    const page = parseIntSafe(query[`${prefix}page`], defaults.page ?? 0);
    const pageSize = parseIntSafe(query[`${prefix}ps`], defaults.pageSize ?? 10);
    const sortOrder = decodeParam(query[`${prefix}sort`]) || defaults.sortOrder || '';
    const search = decodeParam(query[`${prefix}q`]) || '';

    const filters: Record<string, string> = {};
    if (defaults.filters) {
      for (const key of Object.keys(defaults.filters)) {
        const raw = decodeParam(query[`${prefix}${key}`]);
        filters[key] = raw || defaults.filters[key] || '';
      }
    }

    return { page, pageSize, sortOrder, search, filters: filters as F };
  }, [router.query, tableKey]);

  const updateTableState = useCallback(
    (updates: Partial<TableUrlState<F>>) => {
      const { query: currentQuery } = routerRef.current;
      const next: Record<string, string | undefined> = { ...currentQuery } as Record<
        string,
        string | undefined
      >;

      if (updates.page !== undefined) {
        if (updates.page === 0) delete next[`${prefix}page`];
        else next[`${prefix}page`] = String(updates.page);
      }

      if (updates.pageSize !== undefined) {
        if (updates.pageSize === (defaults.pageSize ?? 10)) delete next[`${prefix}ps`];
        else next[`${prefix}ps`] = String(updates.pageSize);
      }

      if (updates.sortOrder !== undefined) {
        if (!updates.sortOrder || updates.sortOrder === defaults.sortOrder)
          delete next[`${prefix}sort`];
        else next[`${prefix}sort`] = updates.sortOrder;
      }

      if (updates.search !== undefined) {
        if (!updates.search) delete next[`${prefix}q`];
        else next[`${prefix}q`] = updates.search;
      }

      if (updates.filters) {
        for (const [key, value] of Object.entries(updates.filters)) {
          const paramKey = `${prefix}${key}`;
          if (!value || value === (defaults.filters as Record<string, string>)?.[key])
            delete next[paramKey];
          else next[paramKey] = value;
        }
      }

      // Remove undefined entries so router.replace gets a clean query object.
      const cleanQuery = Object.fromEntries(
        Object.entries(next).filter(([, v]) => v !== undefined),
      );

      routerRef.current.replace({ query: cleanQuery }, undefined, { shallow: true });
    },
    [prefix, defaults],
  );

  /**
   * Builds a shareable URL that points directly to `rowId` in this table.
   * All current table-state params are preserved so the recipient lands on
   * the same view, with just the row expanded.
   */
  const buildRowDeepLink = useCallback(
    (rowId: string): string => {
      const { query } = routerRef.current;
      const next: Record<string, string> = {};
      for (const [k, v] of Object.entries(query)) {
        if (typeof v === 'string') next[k] = v;
      }
      if (rowId) next[resolvedRowParam] = rowId;
      else delete next[resolvedRowParam];

      const url = new URL(window.location.href);
      url.search = new URLSearchParams(next).toString();
      return url.toString();
    },
    [resolvedRowParam],
  );

  /**
   * Copies a deeplink for the given row to the clipboard.
   */
  const copyRowDeepLink = useCallback(
    (rowId: string) => {
      const url = buildRowDeepLink(rowId);
      navigator.clipboard.writeText(url).catch(() => {
        // Fallback for non-secure contexts (http in dev)
        const ta = document.createElement('textarea');
        ta.value = url;
        document.body.appendChild(ta);
        ta.select();
        document.execCommand('copy');
        document.body.removeChild(ta);
      });
    },
    [buildRowDeepLink],
  );

  return { tableState, updateTableState, copyRowDeepLink, buildRowDeepLink };
}
