import { api, mesheryApiPath } from './index';
import { normalizeKubernetesContextsResponse } from './transforms';
import {
  mesheryApi,
  useGetSystemDatabaseQuery,
  useGetSystemSyncQuery,
  useResetSystemDatabaseMutation,
} from '@meshery/schemas/mesheryApi';

const TAGS = {
  SYSTEM: 'system',
  ADAPTERS: 'adapters',
};

const systemApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getAdapters: builder.query({
      query: () => ({
        url: mesheryApiPath('system/adapters'),
        method: 'GET',
        credentials: 'include',
      }),
      providesTags: [TAGS.ADAPTERS],
    }),

    getAvailableAdapters: builder.query({
      query: () => ({
        url: mesheryApiPath('system/availableAdapters'),
        method: 'GET',
        credentials: 'include',
      }),
      providesTags: [TAGS.ADAPTERS],
    }),

    pingAdapter: builder.query({
      query: (adapterLoc) => ({
        url: mesheryApiPath(`system/adapters`),
        params: { adapter: adapterLoc },
        credentials: 'include',
      }),
      providesTags: (result, error, adapterLoc) => [{ type: TAGS.ADAPTERS, id: adapterLoc }],
    }),

    getKubernetesContexts: builder.query({
      query: (queryArg) => ({
        url: mesheryApiPath('system/kubernetes/contexts'),
        params: {
          pagesize: queryArg?.pagesize || 10,
          search: queryArg?.search || '',
        },
        method: 'GET',
      }),
      transformResponse: normalizeKubernetesContextsResponse,
      providesTags: [TAGS.SYSTEM],
    }),

    adapterOperation: builder.mutation({
      query: (queryArg) => ({
        url: mesheryApiPath(queryArg.url || 'system/adapter/operation'),
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
        body: queryArg.body,
      }),
    }),

    getSmiResults: builder.query({
      query: (queryArg) => ({
        url: mesheryApiPath('smi/results'),
        params: {
          page: queryArg?.page,
          pagesize: queryArg?.pagesize,
          search: queryArg?.search,
          order: queryArg?.order,
        },
        method: 'GET',
        credentials: 'include',
      }),
    }),

    manageAdapter: builder.mutation({
      query: (queryArg) => {
        if (queryArg.method === 'DELETE') {
          return {
            url: mesheryApiPath(`system/adapter/manage`),
            method: 'DELETE',
            credentials: 'include',
            params: { adapter: queryArg.adapter },
          };
        }

        return {
          url: mesheryApiPath('system/adapter/manage'),
          method: 'POST',
          credentials: 'include',
          headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
          body: `meshLocationURL=${encodeURIComponent(queryArg.meshLocationURL)}`,
        };
      },
      invalidatesTags: [TAGS.ADAPTERS],
    }),
  }),
});

export const {
  useGetAdaptersQuery,
  useGetAvailableAdaptersQuery,
  useLazyPingAdapterQuery,
  useManageAdapterMutation,
  useGetKubernetesContextsQuery,
  useLazyGetKubernetesContextsQuery,
  useAdapterOperationMutation,
  useLazyGetSmiResultsQuery,
} = systemApi;

export { useGetSystemDatabaseQuery, useGetSystemSyncQuery, useResetSystemDatabaseMutation };

export const useLazyGetSystemSyncQuery = mesheryApi.endpoints.getSystemSync.useLazyQuery;
