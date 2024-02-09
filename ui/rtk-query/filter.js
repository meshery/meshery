import { api } from './index';

const TAGS = {
  FILTERS: 'filters',
};
const filters = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.FILTERS],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getFilters: builder.query({
        query: (queryArg) => ({
          url: `filter`,
          params: {
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
            visibility: queryArg.visibility,
          },
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.FILTERS }],
      }),
    }),
  });

export const { useGetFiltersQuery } = filters;
