import { useCallback, useEffect, useState } from 'react';
import _ from 'lodash';
import {
  useGetRegistryModelsQuery,
  useLazyGetRegistryModelsQuery,
  useGetRegistryComponentsQuery,
  useLazyGetRegistryComponentsQuery,
  useGetRegistryRelationshipsQuery,
  useLazyGetRegistryRelationshipsQuery,
  useGetRegistryRegistrantsQuery,
  useLazyGetRegistryRegistrantsQuery,
  useListConnectionDefinitionsQuery,
  useLazyListConnectionDefinitionsQuery,
  useGetRegistryComponentsByModelQuery,
  useLazyGetRegistryComponentsByModelQuery,
  useGetRegistryRelationshipsByModelQuery,
  useLazyGetRegistryRelationshipsByModelQuery,
  useGetRegistryCategoriesQuery,
  useGetRegistryModelsByCategoryQuery,
  useLazyGetRegistryModelsByCategoryQuery,
  useGetRegistryModelsByNameQuery,
  useLazyGetRegistryModelsByNameQuery,
  useGetRegistryComponentsByNameQuery,
  useGetRegistryComponentsByModelAndNameQuery,
  useLazyExportRegistryModelQuery,
  useUpdateRegistryEntityStatusMutation,
  useRegisterRegistryEntitiesMutation,
} from '@meshery/schemas/mesheryApi';
import { api } from './index';
import { initiateQuery } from './utils';

// Capabilities-registry endpoints are generated from the schemas registry
// construct (schemas/constructs/v1beta1/registry/api.yml) and call the
// canonical /api/registry routes. This module only adapts the legacy hook
// names and argument shapes ({ params: { pagesize, ... } }) still used
// across the UI to the generated hooks' flat, camelCase arguments.

const REGISTRY_TAG = 'Registry_Registry';
const CONNECTION_DEFINITIONS_TAG = 'Connection_API_ConnectionDefinitions';

// Registry imports and status changes can add or retire connection
// definitions, which the generated endpoints track under a separate tag;
// invalidate both so every registry view refetches (the legacy hand-rolled
// endpoints shared one blanket tag).
api.enhanceEndpoints({
  endpoints: {
    registerRegistryEntities: {
      invalidatesTags: [REGISTRY_TAG, CONNECTION_DEFINITIONS_TAG],
    },
    updateRegistryEntityStatus: {
      invalidatesTags: [REGISTRY_TAG, CONNECTION_DEFINITIONS_TAG],
    },
  },
});

const defaultOptions = {
  trim: false,
  // annotations: false,
  search: '',
  page: 0,
  pagesize: 'all',
};

// Merges the legacy defaults into a legacy `params` object and renames the
// legacy `pagesize` key to the canonical `pageSize` wire param.
const listParams = (params?: Record<string, any>): any => {
  const { pagesize, pageSize, ...rest } = _.merge({}, defaultOptions, params);
  return { ...rest, pageSize: pageSize ?? pagesize };
};

const listArg = (queryArgs?: { params?: Record<string, any> }): any =>
  listParams(queryArgs?.params);

type LazyHook = () => [any, ...any[]];

const useWrappedLazy = (useLazyHook: LazyHook, mapArg: (queryArgs?: any) => any) => {
  const [trigger, ...rest] = useLazyHook();
  const wrapped = useCallback(
    (queryArgs?: any, preferCacheValue?: boolean) => trigger(mapArg(queryArgs), preferCacheValue),
    [trigger, mapArg],
  );
  return [wrapped, ...rest] as const;
};

const modelScopedArg = (queryArgs: any): any => ({
  model: queryArgs.model,
  ...listArg(queryArgs),
});

const modelByNameArg = (queryArgs: any): any => ({
  model: queryArgs.name,
  ...listArg(queryArgs),
});

const componentByNameArg = (queryArgs: any): any => ({
  name: queryArgs.name,
  ...listArg(queryArgs),
});

const categoryScopedArg = (queryArgs: any): any => ({
  category: queryArgs.category,
  ...listArg(queryArgs),
});

const componentsByModelAndKindArg = (queryArg: any): any => ({
  model: queryArg.model,
  name: queryArg.component,
  ...listArg(queryArg),
});

const exportArg = (queryArg?: { params?: Record<string, any> }): any => ({
  ...(queryArg?.params ?? {}),
});

export const useGetMeshModelsQuery = (queryArgs?: any, options?: any) =>
  useGetRegistryModelsQuery(listArg(queryArgs), options);
export const useLazyGetMeshModelsQuery = () =>
  useWrappedLazy(useLazyGetRegistryModelsQuery, listArg);

export const useGetComponentsQuery = (queryArgs?: any, options?: any) =>
  useGetRegistryComponentsQuery(listArg(queryArgs), options);
export const useLazyGetComponentsQuery = () =>
  useWrappedLazy(useLazyGetRegistryComponentsQuery, listArg);

export const useGetRelationshipsQuery = (queryArgs?: any, options?: any) =>
  useGetRegistryRelationshipsQuery(listArg(queryArgs), options);
export const useLazyGetRelationshipsQuery = () =>
  useWrappedLazy(useLazyGetRegistryRelationshipsQuery, listArg);

export const useGetRegistrantsQuery = (queryArgs?: any, options?: any) =>
  useGetRegistryRegistrantsQuery(listArg(queryArgs), options);
export const useLazyGetRegistrantsQuery = () =>
  useWrappedLazy(useLazyGetRegistryRegistrantsQuery, listArg);

export const useGetConnectionDefinitionsQuery = (queryArgs?: any, options?: any) =>
  useListConnectionDefinitionsQuery(listArg(queryArgs), options);
export const useLazyGetConnectionDefinitionsQuery = () =>
  useWrappedLazy(useLazyListConnectionDefinitionsQuery, listArg);

export const useGetComponentsFromModalQuery = (queryArgs: any, options?: any) =>
  useGetRegistryComponentsByModelQuery(modelScopedArg(queryArgs), options);
export const useLazyGetComponentsFromModalQuery = () =>
  useWrappedLazy(useLazyGetRegistryComponentsByModelQuery, modelScopedArg);

export const useGetRelationshipsFromModalQuery = (queryArgs: any, options?: any) =>
  useGetRegistryRelationshipsByModelQuery(modelScopedArg(queryArgs), options);
export const useLazyGetRelationshipsFromModalQuery = () =>
  useWrappedLazy(useLazyGetRegistryRelationshipsByModelQuery, modelScopedArg);

export const useGetModelCategoriesQuery = (queryArgs?: any, options?: any) =>
  useGetRegistryCategoriesQuery(queryArgs, options);

export const useGetModelFromCategoryQuery = (queryArgs: any, options?: any) =>
  useGetRegistryModelsByCategoryQuery(categoryScopedArg(queryArgs), options);
export const useLazyGetModelFromCategoryQuery = () =>
  useWrappedLazy(useLazyGetRegistryModelsByCategoryQuery, categoryScopedArg);

export const useGetModelByNameQuery = (queryArgs: any, options?: any) =>
  useGetRegistryModelsByNameQuery(modelByNameArg(queryArgs), options);
export const useLazyGetModelByNameQuery = () =>
  useWrappedLazy(useLazyGetRegistryModelsByNameQuery, modelByNameArg);

export const useGetComponentByNameQuery = (queryArgs: any, options?: any) =>
  useGetRegistryComponentsByNameQuery(componentByNameArg(queryArgs), options);

export const useGetComponentsByModelAndKindQuery = (queryArg: any, options?: any) =>
  useGetRegistryComponentsByModelAndNameQuery(componentsByModelAndKindArg(queryArg), options);

export const useLazyExportModelQuery = () =>
  useWrappedLazy(useLazyExportRegistryModelQuery, exportArg);

export const useUpdateEntityStatusMutation = () => {
  const [updateStatus, ...rest] = useUpdateRegistryEntityStatusMutation();
  const wrapped = useCallback(
    ({ entityType, body }: { entityType: string; body: Record<string, any> }) => {
      // displayName is the canonical wire key; displayname is the legacy
      // spelling some call sites still pass.
      const { displayname, displayName, ...restBody } = body ?? {};
      return updateStatus({
        entityType,
        body: { ...restBody, displayName: displayName ?? displayname } as any,
      });
    },
    [updateStatus],
  );
  return [wrapped, ...rest] as const;
};

export const useImportMeshModelMutation = () => {
  const [registerEntities, ...rest] = useRegisterRegistryEntitiesMutation();
  const wrapped = useCallback(
    (queryArgs: { importBody: any }) => registerEntities({ body: queryArgs.importBody }),
    [registerEntities],
  );
  return [wrapped, ...rest] as const;
};

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
      categoryMap[category.name] = data?.totalCount ?? data?.total_count ?? 0;
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
  const res = await initiateQuery(api.endpoints.getRegistryComponentsByModelAndName, {
    model,
    name: component,
    ...listParams(_.omit(params, ['apiVersion'])),
    annotations: 'include',
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
