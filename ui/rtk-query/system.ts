import { api } from './index';

const TAGS = {
  SYSTEM: 'system',
  ADAPTERS: 'adapters',
  SYNC: 'sync',
};

const systemApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getDatabaseSummary: builder.query({
      query: (queryArg) => ({
        url: `system/database`,
        params: {
          page: queryArg.page,
          pagesize: queryArg.pagesize,
          search: queryArg.search,
          order: queryArg.order,
        },
        method: 'GET',
      }),
      providesTags: () => [{ type: TAGS.SYSTEM }],
    }),
    getAdapters: builder.query({
      query: () => ({
        url: 'system/adapters',
        method: 'GET',
        credentials: 'include',
      }),
      providesTags: [TAGS.ADAPTERS],
    }),

    getAvailableAdapters: builder.query({
      query: () => ({
        url: 'system/availableAdapters',
        method: 'GET',
        credentials: 'include',
      }),
      providesTags: [TAGS.ADAPTERS],
    }),

    pingAdapter: builder.query({
      query: (adapterLoc) => ({
        url: `system/adapters`,
        params: { adapter: adapterLoc },
        credentials: 'include',
      }),
      providesTags: (result, error, adapterLoc) => [{ type: TAGS.ADAPTERS, id: adapterLoc }],
    }),
    getSystemSync: builder.query({
      query: () => ({
        url: 'system/sync',
        method: 'GET',
        credentials: 'include',
      }),
      providesTags: [TAGS.SYNC],
    }),

    getKubernetesContexts: builder.query({
      query: (queryArg) => ({
        url: 'system/kubernetes/contexts',
        params: {
          pagesize: queryArg?.pagesize || 10,
          search: queryArg?.search || '',
        },
        method: 'GET',
      }),
      providesTags: [TAGS.SYSTEM],
    }),

    manageAdapter: builder.mutation({
      query: (queryArg) => {
        if (queryArg.method === 'DELETE') {
          return {
            url: `system/adapter/manage`,
            method: 'DELETE',
            credentials: 'include',
            params: { adapter: queryArg.adapter },
          };
        }

        return {
          url: 'system/adapter/manage',
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
  useGetDatabaseSummaryQuery,
  useGetAdaptersQuery,
  useGetAvailableAdaptersQuery,
  useLazyPingAdapterQuery,
  useManageAdapterMutation,
  useGetSystemSyncQuery,
  useLazyGetSystemSyncQuery,
  useGetKubernetesContextsQuery,
  useLazyGetKubernetesContextsQuery,
} = systemApi;
