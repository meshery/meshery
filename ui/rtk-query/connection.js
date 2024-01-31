import { api } from './index';

const TAGS = {
  CONNECTIONS: 'connections',
};

const connectionsApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getConnections: builder.query({
      query: (queryArg) => ({
        url: `integrations/connections`,
        params: {
          page: queryArg.page,
          pagesize: queryArg.pagesize,
          search: queryArg.search,
          order: queryArg.order,
        },
        method: 'GET',
      }),
      providesTags: () => [{ type: TAGS.CONNECTIONS }],
    }),
    getDatabaseSummary: builder.query({
      query: (queryArg) => ({
        url: `system/database`,
        params: {
          page: queryArg.page,
          pagesize: queryArg.pagesize,
          search: queryArg.search,
        },
        method: 'GET',
      }),
    }),
    getConnectionStatus: builder.query({
      query: (queryArg) => ({
        url: `integrations/connections/${queryArg.connectionKind}/status`,
        params: { id: queryArg.repoURL },
      }),
    }),
    getConnectionDetails: builder.query({
      query: (queryArg) => ({
        url: `integrations/connections/${queryArg.connectionKind}/details`,
        params: { id: queryArg.repoURL },
      }),
    }),
    verifyConnectionURL: builder.mutation({
      query: (queryArg) => ({
        url: `integrations/connections/${queryArg.connectionKind}/verify`,
        method: 'POST',
        params: { id: queryArg.repoURL },
      }),
    }),
    connectionMetaData: builder.mutation({
      query: (queryArg) => ({
        url: `integrations/connections/${queryArg.connectionKind}/metadata`,
        method: 'POST',
        body: queryArg.body,
      }),
    }),
    configureConnection: builder.mutation({
      query: (queryArg) => ({
        url: `integrations/connections/${queryArg.connectionKind}/configure`,
        method: 'POST',
        body: queryArg.body,
      }),
    }),
    updateConnection: builder.mutation({
      query: (queryArg) => ({
        url: `integrations/connections/${queryArg.connectionKind}/status`,
        method: 'PUT',
        body: queryArg.connectionPayload,
      }),
      invalidatesTags: () => [{ type: TAGS.CONNECTIONS }],
    }),
  }),
});

export const {
  useGetConnectionsQuery,
  useGetDatabaseSummaryQuery,
  useGetConnectionStatusQuery,
  useLazyGetConnectionDetailsQuery,
  useVerifyConnectionURLMutation,
  useConnectionMetaDataMutation,
  useConfigureConnectionMutation,
  useUpdateConnectionMutation,
} = connectionsApi;
