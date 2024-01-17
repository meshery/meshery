import { api } from './index';

const TAGS = {
  MESH_MODELS: 'mesh-models',
};

const meshModelApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.MESH_MODELS],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getMeshModels: builder.query({
        query: (queryArgs) => ({
          url: `meshmodels/models`,
          params: {
            page: queryArgs.page || 0,
            pagesize: queryArgs.pagesize || 'all',
            search: queryArgs.search || '',
            trim: queryArgs.trim || true,
            ...queryArgs,
          },
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getComponents: builder.query({
        query: (queryArgs) => ({
          url: `meshmodels/components`,
          params: {
            page: queryArgs.page || 0,
            pagesize: queryArgs.pagesize || 'all',
            search: queryArgs.search || '',
            trim: queryArgs.trim || true,
            ...queryArgs,
          },
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getRelationships: builder.query({
        query: (queryArgs) => ({
          url: `meshmodels/relationships`,
          params: {
            page: queryArgs.page || 0,
            pagesize: queryArgs.pagesize || 'all',
            search: queryArgs.search || '',
            ...queryArgs,
          },
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getRegistrants: builder.query({
        query: (queryArgs) => ({
          url: `meshmodels/registrants`,
          params: {
            page: queryArgs.page || 0,
            pagesize: queryArgs.pagesize || 'all',
            search: queryArgs.search || '',
            ...queryArgs,
          },
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
    }),
  });

export const {
  useGetMeshModelsQuery,
  useGetComponentsQuery,
  useGetRelationshipsQuery,
  useGetRegistrantsQuery,
} = meshModelApi;
