import { api, mesheryApiPath } from './index';

export const meshApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getMesh: builder.query({
      query: () => mesheryApiPath('mesh'),
    }),
  }),
});

export const { useGetMeshQuery } = meshApi;
