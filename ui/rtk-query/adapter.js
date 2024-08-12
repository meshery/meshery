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
  }),
});

export const { useGetAdaptersUrlQuery, useGetAvailableAdaptersQuery } = adapterApi;
