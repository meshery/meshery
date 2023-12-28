import { api } from './index';

const TAGS = {
  ORGANIZATION: 'organization',
};

const organizationsApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.ORGANIZATION],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getOrgs: builder.query({
        query: (queryArgs) => ({
          url: `identity/orgs`,
          params: {
            page: queryArgs.page || 0,
            pagesize: queryArgs.pagesize || 10,
          },
        }),
        providesTags: () => [{ type: TAGS.ORGANIZATION }],
      }),
    }),
  });

export const { useGetOrgsQuery } = organizationsApi;
