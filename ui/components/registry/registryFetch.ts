import { toLower } from 'lodash';
import {
  MODELS,
  COMPONENTS,
  RELATIONSHIPS,
  REGISTRANTS,
  CONNECTIONS,
} from '../../constants/navigator';
import { removeDuplicateVersions } from './helper';

export type TrackedQueryHandle = Promise<any> & {
  unsubscribe?: () => void;
  abort?: () => void;
};

type FetchRegistrantsParams = {
  searchText: string | null;
  registrantPage: number;
  modelsPage?: number;
  getRegistrantsData: (args: any, preferCacheValue?: boolean) => TrackedQueryHandle;
  getMeshModelsData: (args: any, preferCacheValue?: boolean) => TrackedQueryHandle;
  runTrackedQuery: <T>(
    queryPromise: Promise<T> & { unsubscribe?: () => void; abort?: () => void },
  ) => Promise<T>;
  isStale: () => boolean;
};

type FetchRegistryViewParams = {
  view: string;
  searchText: string | null;
  checked: boolean;
  rowsPerPage: number;
  page: { [key: string]: number };
  modelFilters: { page: number };
  componentsFilters: { page: number };
  relationshipsFilters: { page: number };
  registrantFilters: { page: number };
  connectionsFilters: { page: number };
  getMeshModelsData: (args: any, preferCacheValue?: boolean) => TrackedQueryHandle;
  getComponentsData: (args: any, preferCacheValue?: boolean) => TrackedQueryHandle;
  getRelationshipsData: (args: any, preferCacheValue?: boolean) => TrackedQueryHandle;
  getRegistrantsData: (args: any, preferCacheValue?: boolean) => TrackedQueryHandle;
  getConnectionDefinitionsData: (args: any, preferCacheValue?: boolean) => TrackedQueryHandle;
  runTrackedQuery: <T>(
    queryPromise: Promise<T> & { unsubscribe?: () => void; abort?: () => void },
  ) => Promise<T>;
  isStale: () => boolean;
};

/**
 * Fetches registrants along with their deduplicated models while checking
 * cancellation between sequential per-registrant queries.
 */
export const fetchRegistrantsWithModels = async ({
  searchText,
  registrantPage,
  modelsPage,
  getRegistrantsData,
  getMeshModelsData,
  runTrackedQuery,
  isStale,
}: FetchRegistrantsParams) => {
  let response;
  const registrantResponse = await runTrackedQuery(
    getRegistrantsData(
      {
        params: {
          page: searchText ? 0 : registrantPage,
          pagesize: searchText ? 'all' : 25,
          search: searchText || '',
        },
      },
      true,
    ),
  );
  if (isStale()) {
    return null;
  }
  if (registrantResponse.data && registrantResponse.data.registrants) {
    const registrants = registrantResponse.data.registrants;
    const tempResourcesDetail = [];

    for (const registrant of registrants) {
      if (isStale()) {
        return null;
      }
      const hostname = toLower(registrant?.hostname);
      const { data: modelRes } = await runTrackedQuery(
        getMeshModelsData(
          {
            params: {
              page: modelsPage,
              pagesize: 'all',
              registrant: hostname,
              components: false,
              relationships: false,
            },
          },
          true,
        ),
      );
      if (isStale()) {
        return null;
      }
      if (modelRes?.models && modelRes.models.length > 0) {
        tempResourcesDetail.push({
          ...registrant,
          models: removeDuplicateVersions(modelRes.models) || [],
        });
      }
    }
    response = {
      data: {
        registrants: tempResourcesDetail,
      },
    };
  }
  return response;
};

/**
 * Dispatches the tracked RTK lazy query for the active registry resource tab.
 */
export const fetchRegistryViewData = async ({
  view,
  searchText,
  checked,
  rowsPerPage,
  page,
  modelFilters,
  componentsFilters,
  relationshipsFilters,
  registrantFilters,
  connectionsFilters,
  getMeshModelsData,
  getComponentsData,
  getRelationshipsData,
  getRegistrantsData,
  getConnectionDefinitionsData,
  runTrackedQuery,
  isStale,
}: FetchRegistryViewParams) => {
  switch (view) {
    case MODELS:
      return runTrackedQuery(
        getMeshModelsData(
          {
            params: {
              page: searchText || checked ? 0 : modelFilters.page,
              pagesize: searchText || checked ? 'all' : 25,
              components: false,
              relationships: false,
              search: searchText || '',
            },
          },
          true,
        ),
      );
    case COMPONENTS:
      return runTrackedQuery(
        getComponentsData(
          {
            params: {
              page: searchText ? 0 : componentsFilters.page,
              pagesize: searchText ? 'all' : rowsPerPage,
              search: searchText || '',
              trim: true,
            },
          },
          true,
        ),
      );
    case RELATIONSHIPS:
      return runTrackedQuery(
        getRelationshipsData(
          {
            params: {
              page: searchText ? 0 : relationshipsFilters.page,
              pagesize: 'all',
              search: searchText || '',
            },
          },
          true,
        ),
      );
    case REGISTRANTS:
      return fetchRegistrantsWithModels({
        searchText,
        registrantPage: registrantFilters.page,
        modelsPage: page?.Models,
        getRegistrantsData,
        getMeshModelsData,
        runTrackedQuery,
        isStale,
      });
    case CONNECTIONS: {
      const res = await runTrackedQuery(
        getConnectionDefinitionsData(
          {
            params: {
              page: searchText ? 0 : connectionsFilters.page,
              pagesize: searchText ? 'all' : rowsPerPage,
              search: searchText || '',
            },
          },
          true,
        ),
      );
      return res?.data
        ? { ...res, data: { ...res.data, connections: res.data.connectionDefinitions || [] } }
        : res;
    }
    default:
      return undefined;
  }
};
