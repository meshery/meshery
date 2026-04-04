import { useEffect, useState } from 'react';
import { api, mesheryApiPath } from './index';
import _ from 'lodash';
import { initiateQuery } from './utils';

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
          url: mesheryApiPath(`meshmodels/models`),
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getComponents: builder.query({
        query: (queryArgs) => ({
          url: mesheryApiPath(`meshmodels/components`),
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getRelationships: builder.query({
        query: (queryArgs) => ({
          url: mesheryApiPath(`meshmodels/relationships`),
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getRegistrants: builder.query({
        query: (queryArgs) => ({
          url: mesheryApiPath(`meshmodels/registrants`),
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getComponentsFromModal: builder.query({
        query: (queryArgs) => ({
          url: mesheryApiPath(`meshmodels/models/${queryArgs.model}/components`),
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getRelationshipsFromModal: builder.query({
        query: (queryArgs) => ({
          url: mesheryApiPath(`meshmodels/models/${queryArgs.model}/relationships`),
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      updateEntityStatus: builder.mutation({
        query: (queryArgs) => ({
          url: mesheryApiPath(`meshmodels/${queryArgs.entityType}/status`),
          method: 'POST',
          body: queryArgs.body,
        }),
        invalidatesTags: [TAGS.MESH_MODELS],
      }),
      getModelCategories: builder.query({
        query: () => ({
          url: mesheryApiPath(`meshmodels/categories`),
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getModelFromCategory: builder.query({
        query: (queryArgs) => ({
          url: mesheryApiPath(`meshmodels/categories/${queryArgs.category}/models`),
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getModelByName: builder.query({
        query: (queryArgs) => ({
          url: mesheryApiPath(`meshmodels/models/${queryArgs.name}`),
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getComponentByName: builder.query({
        query: (queryArgs) => ({
          url: mesheryApiPath(`meshmodels/components/${queryArgs.name}`),
          params: _.merge({}, defaultOptions, queryArgs.params),
        }),
        providesTags: () => [{ type: TAGS.MESH_MODELS }],
      }),
      getComponentsByModelAndKind: builder.query({
        query: (queryArg) => ({
          url: mesheryApiPath(
            `meshmodels/models/${queryArg.model}/components/${queryArg.component}`,
          ),
          params: _.merge({}, defaultOptions, queryArg.params),
        }),
      }),
      exportModel: builder.query({
        query: (queryArg) => ({
          url: mesheryApiPath(`meshmodels/export`),
          params: _.merge({}, defaultOptions, queryArg.params),
          providesTags: () => [{ type: TAGS.MESH_MODELS }],
        }),
      }),
      importMeshModel: builder.mutation({
        query: (queryArgs) => {
          return {
            url: mesheryApiPath(`meshmodels/register`),
            method: 'POST',
            body: queryArgs.importBody,
          };
        },
        invalidatesTags: [TAGS.MESH_MODELS],
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
  useGetComponentsFromModalQuery,
  useLazyGetComponentsFromModalQuery,
  useGetRelationshipsFromModalQuery,
  useLazyExportModelQuery,
  useLazyGetRelationshipsFromModalQuery,
  useUpdateEntityStatusMutation,
  useGetModelCategoriesQuery,
  useLazyGetModelFromCategoryQuery,
  useGetModelByNameQuery,
  useLazyGetModelByNameQuery,
  useGetMeshModelsQuery,
  useGetComponentByNameQuery,
  useGetModelFromCategoryQuery,
  useGetComponentsByModelAndKindQuery,
  useImportMeshModelMutation,
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

  useEffect(() => {
    const fetchData = async () => {
      const categoryMap = await fetchModelsForCategories();
      setCategoryMap(categoryMap);
    };
    fetchData();
  }, [categories]);
  return categoryMap;
};

export const getComponentDefinition = async (component, model, params = {}) => {
  const res = await initiateQuery(meshModelApi.endpoints.getComponentsByModelAndKind, {
    component,
    model,
    params: _.omit({ params, annotations: 'include' }, ['apiVersion']),
  });

  if (params.apiVersion) {
    return res?.data?.components?.find((c) => c.component.version === params.apiVersion);
  }
  return res?.data?.components?.[0];
};

export const modelUniqueKey = (model) => `${model.name}-${model.version}`;
export const componentUniqueKey = (component) =>
  `${component.component.kind}-${component.component.version}-${component.version}-${modelUniqueKey(
    component.model,
  )}`;
