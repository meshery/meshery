import { useGetEventFiltersQuery } from '../../rtk-query/notificationCenter';
import TypingFilter from '../TypingFilter';
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

const Filter = ({ handleFilter, currentFilters }) => {
  const filterSchema = useFilterSchema();

  const transformToFilterArray = (filters) => {
    if (!filters || typeof filters !== 'object') return [];

    const filterArray = [];
    Object.entries(filters).forEach(([key, value]) => {
      const schemaKey = Object.keys(filterSchema).find(
        (schemaKey) => filterSchema[schemaKey].value === key,
      );

      if (schemaKey && filterSchema[schemaKey]) {
        if (Array.isArray(value)) {
          value.forEach((v) => {
            filterArray.push({
              type: schemaKey,
              value: v,
              label: `${filterSchema[schemaKey].value}: ${v}`,
            });
          });
        } else {
          filterArray.push({
            type: schemaKey,
            value: value,
            label: `${filterSchema[schemaKey].value}: ${value}`,
          });
        }
      }
    });

    return filterArray;
  };

  const defaultFilters = currentFilters
    ? transformToFilterArray(currentFilters)
    : [
        {
          type: 'STATUS',
          value: 'unread',
          label: 'status: unread',
        },
      ];

  const filterKey = currentFilters ? JSON.stringify(currentFilters) : 'default';

  return (
    <TypingFilter
      key={filterKey}
      handleFilter={handleFilter}
      filterSchema={filterSchema}
      defaultFilters={defaultFilters}
      placeholder="Filter Notifications"
    />
  );
};

export default Filter;
