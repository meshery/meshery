import { api } from './index';

export const adapterApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getAdapters: builder.query({
      query: () => 'system/adapters',
      method: 'GET',
    }),
  }),
});

export const { useGetAdaptersQuery } = adapterApi;
