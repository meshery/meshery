import { api } from './index';

export const adapterApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getAdaptersUrl: builder.query({
      query: () => 'system/adapters',
      method: 'GET',
    }),
    getAvailableAdapters: builder.query({
      query: () => 'system/availableAdapters',
      method: 'GET',
    }),
    deleteAdapters: builder.mutation({
      query: (queryArg) => ({
        url: `system/adapter/manage?adapter=${encodeURIComponent(queryArg.adapterLoc)}`,
        method: 'DELETE',
      }),
    }),
    pingAdapter: builder.query({
      query: (queryArg) => `system/adapters?adapter=${encodeURIComponent(queryArg.adapterLoc)}`,
    }),
  }),
});

export const {
  useGetAdaptersUrlQuery,
  useGetAvailableAdaptersQuery,
  useDeleteAdaptersMutation,
  usePingAdapterQuery,
} = adapterApi;
