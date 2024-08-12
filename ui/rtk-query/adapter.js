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
        url: `/api/system/adapter/manage?adapter=${encodeURIComponent(queryArg.adapterLoc)}`,
        method: 'DELETE',
      }),
    }),
  }),
});

export const { useGetAdaptersUrlQuery, useGetAvailableAdaptersQuery, useDeleteAdaptersMutation } =
  adapterApi;
