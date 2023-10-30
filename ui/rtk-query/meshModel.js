import { api } from './index';

const userApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getMeshModel: builder.query({
      query: (queryArg) => ({
        url: `meshmodels/models`,
        params: {
          // queryArg.page,
          // queryArg.pageSize,
        },
      }),
    }),
  }),
});

export const { useGetMeshModelQuery } = userApi;
