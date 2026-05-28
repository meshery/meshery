import { api, mesheryApiPath } from './index';

export const meshApi = api.injectEndpoints({
  overrideExisting: true,
  endpoints: (builder) => ({
    getMesh: builder.query({
      query: () => mesheryApiPath('mesh'),
    }),
  }),
});

export const { useGetMeshQuery } = meshApi;
