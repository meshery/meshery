import { useSelector } from 'react-redux';
import { useGetEventFiltersQuery } from '../../../rtk-query/notificationCenter';
import TypingFilter from '@/components/shared/FormFields/typing-filter';
import { SEVERITY, STATUS } from './constants';

const useFilterSchema = () => {
  const { data } = useGetEventFiltersQuery();

  return {
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
  };
};

const Filter = ({ handleFilter }: { handleFilter: (filters: unknown) => void }) => {
  const filterSchema = useFilterSchema();
  const currentFilters = useSelector((state: any) => state.events.current_view.filters);

  const selectedFilters = Object.entries(currentFilters || {}).flatMap(([key, value]) => {
    const upperKey = key.toUpperCase();
    const schema = filterSchema[upperKey as keyof typeof filterSchema];
    if (!schema) return [];
    const values = Array.isArray(value) ? value : [value];
    return values.map((v) => ({ type: upperKey, value: v, label: `${schema.value}: ${v}` }));
  });

  return (
    <TypingFilter
      handleFilter={handleFilter}
      filterSchema={filterSchema}
      defaultFilters={selectedFilters}
      placeholder="Filter Notifications"
    />
  );
};

export default Filter;
