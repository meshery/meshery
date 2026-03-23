import { api } from './index';

const userKeysApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getUserKeys: builder.query({
      query: (queryArgs = {}) => ({
        url: `identity/orgs/${queryArgs.orgId}/users/keys`,
      }),
    }),
  }),
});

export const useGetUserKeysQuery = (queryArgs = {}, options = undefined) =>
  userKeysApi.endpoints.getUserKeys.useQuery(queryArgs, options);

export const useLazyGetUserKeysQuery = () => userKeysApi.endpoints.getUserKeys.useLazyQuery();
