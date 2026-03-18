import { api } from './index';

export const meshApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getMesh: builder.query({
      query: () => 'mesh',
    }),
  }),
});

export const { useGetMeshQuery } = meshApi;
