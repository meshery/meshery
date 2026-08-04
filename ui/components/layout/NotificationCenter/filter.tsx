import { useMemo } from 'react';
import { useSelector } from 'react-redux';
import { useGetEventFiltersQuery } from '../../../rtk-query/notificationCenter';
import TypingFilter from '@/components/shared/FormFields/typing-filter';
import { SEVERITY, STATUS } from './constants';

const DEFAULT_STATUS_CHIP = {
  type: 'STATUS',
  value: STATUS.UNREAD,
  label: `status: ${STATUS.UNREAD}`,
};

const useFilterSchema = () => {
  const { data } = useGetEventFiltersQuery();

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
 * Falls back to unread (product default) for unset / `{ initial: true }` views.
 */
export const filtersToChips = (filters, filterSchema) => {
  if (!filters || filters.initial) {
    return [DEFAULT_STATUS_CHIP];
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

  return chips.length > 0 ? chips : [DEFAULT_STATUS_CHIP];
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
const Filter = ({ handleFilter }: { handleFilter: (filters: unknown) => void }) => {
  const filterSchema = useFilterSchema();
  const currentFilters = useSelector((state: any) => state.events.current_view.filters);
  const selectedFilters = useMemo(
    () => filtersToChips(currentFilters, filterSchema),
    [currentFilters, filterSchema],
  );
  const filtersKey = useMemo(() => JSON.stringify(selectedFilters), [selectedFilters]);

  // Clear restores unread — the product default — not an unfiltered fetch.
  const onFilterChange = (filters: Record<string, unknown>) => {
    const normalized = normalizeFilterPayload(filters || {}, filterSchema);
    if (Object.keys(normalized).length === 0) {
      handleFilter({ status: STATUS.UNREAD });
      return;
    }
    handleFilter(normalized);
  };

  return (
    <TypingFilter
      key={filtersKey}
      handleFilter={onFilterChange}
      filterSchema={filterSchema}
      defaultFilters={selectedFilters}
      placeholder="Filter Notifications"
    />
  );
};

export default Filter;
