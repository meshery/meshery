import { api } from './index';
import _ from 'lodash';

const TAGS = {
  MESH_MODELS: 'mesh-models',
};

const defaultOptions = {
  trim: false,
  annotations: false,
  search: '',
  page: 1,
  pagesize: 'all',
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
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getComponents: builder.query({
        query: (queryArgs) => ({
          url: `meshmodels/components`,
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getRelationships: builder.query({
        query: (queryArgs) => ({
          url: `meshmodels/relationships`,
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getRegistrants: builder.query({
        query: (queryArgs) => ({
          url: `meshmodels/registrants`,
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getComponentsFromModal: builder.query({
        query: (queryArgs) => ({
          url: `meshmodels/models/${queryArgs.model}/components`,
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getRelationshipsFromModal: builder.query({
        query: (queryArgs) => ({
          url: `meshmodels/models/${queryArgs.model}/relationships`,
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
    }),
  });

export const {
  useLazyGetMeshModelsQuery,
  useLazyGetComponentsQuery,
  useLazyGetRelationshipsQuery,
  useLazyGetRegistrantsQuery,
  useLazyGetComponentsFromModalQuery,
  useLazyGetRelationshipsFromModalQuery,
} = meshModelApi;
