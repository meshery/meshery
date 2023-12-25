import { api } from './index';

const organizationsApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getOrgs: builder.query({
      query: ({ page = 0, pagesize = 25 }) => ({
        url: `identity/orgs`,
        params: {
          page: page,
          pagesize: pagesize,
        },
      }),
      providesTags: ['organization'],
    }),
  }),
});

export const { useGetOrgsQuery } = organizationsApi;
