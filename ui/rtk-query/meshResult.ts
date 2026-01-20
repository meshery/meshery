import { api } from './index';

const TAGS = {
  MESHERY_RESULTS: 'meshery-results',
};

const resultsApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getResults: builder.query({
      query: ({ endpoint, ...queryArg }) => ({
        url: endpoint || '/perf/profile/result',
        params: {
          page: queryArg.page,
          pagesize: queryArg.pagesize,
          search: queryArg.search,
          order: queryArg.sortOrder,
        },
        method: 'GET',
      }),
      providesTags: () => [{ type: TAGS.MESHERY_RESULTS }],
    }),
  }),
});

export const { useGetResultsQuery, useLazyGetResultsQuery } = resultsApi;
