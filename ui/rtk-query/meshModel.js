import { useEffect, useState } from 'react';
import { api } from './index';
import _ from 'lodash';

const TAGS = {
  MESH_MODELS: 'mesh-models',
};

const defaultOptions = {
  trim: false,
  // annotations: false,
  search: '',
  page: 0,
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
      updateEntityStatus: builder.mutation({
        query: (queryArgs) => ({
          url: `meshmodels/${queryArgs.entityType}/status`,
          method: 'POST',
          body: queryArgs.body,
        }),
        invalidatesTags: [TAGS.MESH_MODELS],
      }),
      getModelCategories: builder.query({
        query: () => ({
          url: `meshmodels/categories`,
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getModelFromCategory: builder.query({
        query: (queryArgs) => ({
          url: `meshmodels/categories/${queryArgs.category}/models`,
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getModelByName: builder.query({
        query: (queryArgs) => ({
          url: `meshmodels/models/${queryArgs.name}`,
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getComponentByName: builder.query({
        query: (queryArgs) => ({
          url: `meshmodels/components/${queryArgs.name}`,
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
    }),
  });

export const {
  useLazyGetMeshModelsQuery,
  useLazyGetComponentsQuery,
  useGetComponentsQuery,
  useLazyGetRelationshipsQuery,
  useGetRegistrantsQuery,
  useGetRelationshipsQuery,
  useLazyGetRegistrantsQuery,
  useLazyGetComponentsFromModalQuery,
  useLazyGetRelationshipsFromModalQuery,
  useUpdateEntityStatusMutation,
  useGetModelCategoriesQuery,
  useLazyGetModelFromCategoryQuery,
  useGetModelByNameQuery,
  useGetMeshModelsQuery,
  useGetComponentByNameQuery,
  useGetModelFromCategoryQuery,
} = meshModelApi;

export const useGetCategoriesSummary = () => {
  const [getModelFromCategory] = useLazyGetModelFromCategoryQuery();
  const { data: categories } = useGetModelCategoriesQuery();
  const [categoryMap, setCategoryMap] = useState({});

  const fetchModelsForCategories = async () => {
    const categoryMap = {};
    if (!categories) return categoryMap;

    const requests = categories.categories.map(async (category) => {
      const { data } = await getModelFromCategory(
        { category: category.name, params: { page: 1, pagesize: 1 } },
        true,
      );
      categoryMap[category.name] = data?.total_count || 0;
    });
    await Promise.allSettled(requests);
    return categoryMap;
  };

  useEffect(async () => {
    const categoryMap = await fetchModelsForCategories();
    setCategoryMap(categoryMap);
  }, [categories]);
  return categoryMap;
};
