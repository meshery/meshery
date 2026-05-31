import { api, mesheryApiPath } from './index';
import { shouldOverrideExisting } from './utils';

export const meshApi = api.injectEndpoints({
  overrideExisting: shouldOverrideExisting,
  endpoints: (builder) => ({
    getMesh: builder.query({
      query: () => mesheryApiPath('mesh'),
    }),
  }),
});

export const { useGetMeshQuery } = meshApi;
