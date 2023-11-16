import { api } from './index';

const userApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getMeshModel: builder.query({
      query: () => `meshmodels/models`,
    }),
  }),
});

export const { useGetMeshModelQuery } = userApi;
