import { api } from './index';

const organizationsApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getOrgs: builder.query({
      query: (queryArg) => ({
        url: `identity/orgs`,
        params: {
          page: queryArg.page,
          pagesize: queryArg.pagesize,
          search: queryArg.search,
          order: queryArg.order,
        },
      }),
    }),
  }),
});

export const { useGetOrgsQuery } = organizationsApi;
