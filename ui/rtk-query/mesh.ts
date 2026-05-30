import { api, mesheryApiPath } from './index';

export const meshApi = api.injectEndpoints({
  overrideExisting: module.hot?.status() === 'apply',
  endpoints: (builder) => ({
    getMesh: builder.query({
      query: () => mesheryApiPath('mesh'),
    }),
  }),
});

export const { useGetMeshQuery } = meshApi;
