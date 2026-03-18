import { api } from './index';

const TAGS = {
  USER_KEYS: 'user-keys',
};

const userKeysApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.USER_KEYS],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getUserKeys: builder.query({
        query: (queryArgs) => ({
          url: `identity/orgs/${queryArgs.orgId}/users/keys`,
          params: {
            page: queryArgs.page || 0,
            pagesize: queryArgs.pagesize || 10,
          },
        }),
        providesTags: () => [{ type: TAGS.USER_KEYS }],
      }),
    }),
  });

export const { useLazyGetUserKeysQuery, useGetUserKeysQuery } = userKeysApi;
