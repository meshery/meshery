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
    deleteAdapter: builder.mutation({
      query: (queryArg) => ({
        url: `system/adapter/manage?adapter=${encodeURIComponent(queryArg.adapterLoc)}`,
        method: 'DELETE',
      }),
    }),
    pingAdapter: builder.query({
      query: (queryArg) => `system/adapters?adapter=${encodeURIComponent(queryArg.adapterLoc)}`,
    }),
    connectAdapter: builder.mutation({
      query: (queryArg) => ({
        url: `system/adapter/manage`,
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
        body: queryArg.params,
      }),
    }),
  }),
});

export const {
  useGetAdaptersUrlQuery,
  useGetAvailableAdaptersQuery,
  useDeleteAdapterMutation,
  usePingAdapterQuery,
  useConnectAdapterMutation,
} = adapterApi;
