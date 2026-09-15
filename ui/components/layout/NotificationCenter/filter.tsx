import { memo, useCallback, useMemo } from 'react';
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
    const limited = schema.multiple === false ? values.slice(-1) : values;
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

/**
 * Collapse an array-valued `multiple: false` filter to its last scalar, and drop
 * any filter whose value is an empty array (so an emptied multi-select clears).
 */
export const normalizeFilterPayload = (filters, filterSchema) => {
  const normalized = { ...filters };
  Object.values(filterSchema).forEach((schema) => {
    const value = normalized[schema.value];
    if (!Array.isArray(value)) {
      return;
    }
    if (value.length === 0) {
      delete normalized[schema.value];
      return;
    }
    if (schema.multiple === false) {
      normalized[schema.value] = value[value.length - 1];
    }
  });
  return normalized;
};

/**
 * Notification Center filter bar: NC schema, Redux→chip mapping, unread default.
 * Chips flow one way: Redux → `filtersToChips` → TypingFilter's `defaultFilters`,
 * which re-seeds its local state whenever that set changes.
 */
const Filter = memo(function Filter({
  handleFilter,
  currentFilters,
}: {
  handleFilter: (filters: unknown) => void;
  currentFilters: Record<string, unknown>;
}) {
  const filterSchema = useFilterSchema();
  const selectedFilters = useMemo(
    () => filtersToChips(currentFilters, filterSchema),
    [currentFilters, filterSchema],
  );

  // An empty normalized payload clears every filter, including the unread default.
  const onFilterChange = useCallback(
    (filters: Record<string, unknown>) => {
      const normalized = normalizeFilterPayload(filters || {}, filterSchema);
      handleFilter(Object.keys(normalized).length === 0 ? {} : normalized);
    },
    [filterSchema, handleFilter],
  );

  return (
    <TypingFilter
      handleFilter={onFilterChange}
      filterSchema={filterSchema}
      defaultFilters={selectedFilters}
      placeholder="Filter Notifications"
    />
  );
});

export default Filter;
