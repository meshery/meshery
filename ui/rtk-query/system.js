import { api } from './index';

const TAGS = {
  SYSTEM: 'system',
};

const systemApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getDatabaseSummary: builder.query({
      query: (queryArg) => ({
        url: `system/database`,
        params: {
          page: queryArg.page,
          pagesize: queryArg.pagesize,
          search: queryArg.search,
        },
        method: 'GET',
      }),
      providesTags: () => [{ type: TAGS.SYSTEM }],
    }),
    getVersionDetails: builder.query({
      query: () => 'system/version',
      method: 'GET',
      providesTags: () => [{ type: TAGS.SYSTEM }],
    }),
  }),
});

export const { useGetDatabaseSummaryQuery, useGetVersionDetailsQuery } = systemApi;
