import { Delimiter } from './state';

// returns the filter object from the filterSchema
const getFilterByValue = (value, filterSchema) => {
  return Object.values(filterSchema).find((filter) => filter.value == value);
};

/**
 * Parses a filter string and returns a filter object.
 *
 * @param {string} filterString - The input filter string of the form "type:value type2:value2 type:value2".
 * @returns {Object} - The filter object with types as keys and arrays of values as values.
 */
export const getFilters = (filterString, filterSchema) => {
  const filters = {};
  const filterValuePairs = filterString.split(Delimiter.FILTER);
  filterValuePairs.forEach((filterValuePair) => {
    const [filter, value] = filterValuePair.split(Delimiter.FILTER_VALUE);

    if (getFilterByValue(filter, filterSchema)?.multiple == false) {
      filters[filter] = value;
      return;
    }

    if (filter && value) {
      filters[filter] = filters[filter] || [];
      if (!filters[filter].includes(value)) {
        filters[filter].push(value);
      }
    }
  });

  return filters;
};

// return a filter string of form "type:value type2:value2 type:value2"
// from a filter object of form { type : {values} , type2 : {values}  }
export const getFilterString = (filters) => {
  return Object.entries(filters).reduce((filterString, [filter, values]) => {
    return (
      filterString +
      [...values].map((value) => `${filter}${Delimiter.FILTER_VALUE}${value}`).join(' ')
    );
  }, '');
};

export const getCurrentFilterAndValue = (filteringState) => {
  const { context } = filteringState;
  const currentFilterValue = context.value.split(Delimiter.FILTER).at(-1);
  const currentFilter = currentFilterValue.split(Delimiter.FILTER_VALUE)?.[0] || '';
  const currentValue = currentFilterValue.split(Delimiter.FILTER_VALUE)?.[1] || '';
  return {
    filter: currentFilter,
    value: currentValue,
  };
};
