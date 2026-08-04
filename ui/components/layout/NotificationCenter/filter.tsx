import { memo, useCallback, useMemo, useState } from 'react';
import { useGetEventFiltersQuery } from '../../../rtk-query/notificationCenter';
import TypingFilter from '@/components/shared/FormFields/typing-filter';
import { SEVERITY, STATUS } from './constants';

const DEFAULT_STATUS_CHIP = {
  type: 'STATUS',
  value: STATUS.UNREAD,
  label: `status: ${STATUS.UNREAD}`,
};

const useFilterSchema = () => {
  // Only subscribe to `data` so isFetching/etc. don't re-render Filter.
  const { data } = useGetEventFiltersQuery(undefined, {
    selectFromResult: ({ data }) => ({ data }),
  });

  return useMemo(
    () => ({
      SEVERITY: {
        value: 'severity',
        description: 'Filter by severity',
        values: Object.values(SEVERITY),
      },
      STATUS: {
        value: 'status',
        description: 'Filter by status',
        values: Object.values(STATUS),
        multiple: false,
      },
      ACTION: {
        value: 'action',
        values: data?.action || [],
        description: 'Filter by type',
      },
      AUTHOR: {
        value: 'author',
        description: 'Filter by any user or system',
      },
      CATEGORY: {
        value: 'category',
        description: 'Filter by category',
        values: data?.category || [],
      },
    }),
    [data],
  );
};

/**
 * Map Redux filter object → TypingFilter chips.
 * Unread chip only for the pre-fetch `{ initial: true }` sentinel — not after user clear.
 */
export const filtersToChips = (filters, filterSchema) => {
  if (filters?.initial) {
    return [DEFAULT_STATUS_CHIP];
  }

  if (!filters || Object.keys(filters).length === 0) {
    return [];
  }
  const valueToSchemaKey = Object.entries(filterSchema).reduce((acc, [schemaKey, schema]) => {
    acc[schema.value] = schemaKey;
    return acc;
  }, {});

  const chips = [];
  Object.entries(filters).forEach(([key, value]) => {
    const schemaKey = valueToSchemaKey[key];
    if (!schemaKey || value == null || value === '') {
      return;
    }
    const values = Array.isArray(value) ? value : [value];
    const schema = filterSchema[schemaKey];
    const limited = schema.multiple === false ? values.slice(0, 1) : values;
    limited.forEach((v) => {
      chips.push({
        type: schemaKey,
        value: v,
        label: `${schema.value}: ${v}`,
      });
    });
  });

  return chips;
};

/** Coerce array-valued `multiple: false` filters to the last selected scalar. */
export const normalizeFilterPayload = (filters, filterSchema) => {
  const normalized = { ...filters };
  Object.values(filterSchema).forEach((schema) => {
    if (schema.multiple !== false) {
      return;
    }
    const value = normalized[schema.value];
    if (Array.isArray(value)) {
      if (value.length > 0) {
        normalized[schema.value] = value[value.length - 1];
      } else {
        delete normalized[schema.value];
      }
    }
  });
  return normalized;
};

/**
 * Notification Center filter bar: NC schema, Redux→chip mapping, unread default.
 * TypingFilter is left unchanged (shared). Remount via `key` when chips change so
 * its local state initializes from the new `defaultFilters` without a Redux↔
 * useEffect sync loop.
 */
const Filter = memo(
  ({
    handleFilter,
    currentFilters,
  }: {
    handleFilter: (filters: unknown) => void;
    currentFilters: Record<string, unknown>;
  }) => {
    const filterSchema = useFilterSchema();
    const [resetVersion, setResetVersion] = useState(0);
    const selectedFilters = useMemo(
      () => filtersToChips(currentFilters, filterSchema),
      [currentFilters, filterSchema],
    );
    const filtersKey = useMemo(
      () => JSON.stringify({ selectedFilters, resetVersion }),
      [selectedFilters, resetVersion],
    );

    // User clear removes all filters (including unread); remount so TypingFilter resyncs.
    const onFilterChange = useCallback(
      (filters: Record<string, unknown>) => {
        const normalized = normalizeFilterPayload(filters || {}, filterSchema);
        if (Object.keys(normalized).length === 0) {
          setResetVersion((version) => version + 1);
          handleFilter({});
          return;
        }
        handleFilter(normalized);
      },
      [filterSchema, handleFilter],
    );

    return (
      <TypingFilter
        key={filtersKey}
        handleFilter={onFilterChange}
        filterSchema={filterSchema}
        defaultFilters={selectedFilters}
        placeholder="Filter Notifications"
      />
    );
  },
);

Filter.displayName = 'Filter';

export default Filter;
