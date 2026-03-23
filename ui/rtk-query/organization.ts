import { api } from './index';

const organizationApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getOrgs: builder.query({
      query: (queryArgs = {}) => ({
        url: 'identity/orgs',
        params: {
          page: queryArgs.page,
          pagesize: queryArgs.pagesize,
          search: queryArgs.search,
          order: queryArgs.order,
          all: queryArgs.all,
        },
      }),
    }),
  }),
});

export const { useGetOrgsQuery } = organizationApi;
