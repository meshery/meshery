import { api } from './index';

const TAGS = {
  CONNECTIONS: 'connections',
  CREDENTIALS: 'credentials',
};

const connectionsApi = api.injectEndpoints({
  overrideExisting: true,
  endpoints: (builder) => ({
    getCredentials: builder.query({
      query: () => ({
        url: 'integrations/credentials',
        method: 'GET',
      }),
      providesTags: [TAGS.CREDENTIALS],
    }),

    verifyAndRegisterConnection: builder.mutation({
      query: (queryArg) => ({
        url: 'integrations/connections/register',
        method: 'POST',
        body: queryArg.body,
      }),
      invalidatesTags: [TAGS.CONNECTIONS],
    }),

    connectToConnection: builder.mutation({
      query: (queryArg) => ({
        url: 'integrations/connections/register',
        method: 'POST',
        body: queryArg.body,
      }),
      invalidatesTags: [TAGS.CONNECTIONS],
    }),
    getConnections: builder.query({
      query: (queryArg) => ({
        url: `integrations/connections`,
        params: {
          page: queryArg.page,
          pagesize: queryArg.pagesize,
          search: queryArg.search,
          order: queryArg.order,
          status: queryArg.status,
          kind: queryArg.kind,
        },
        method: 'GET',
      }),
      providesTags: () => [{ type: TAGS.CONNECTIONS }],
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
    getAllConnectionStatus: builder.query({
      query: () => ({
        url: `integrations/connections/status`,
        method: 'GET',
      }),
    }),
    getConnectionByKind: builder.query({
      query: (queryArg) => ({
        url: `integrations/connections/${queryArg.connectionKind}`,
        params: {
          page: queryArg.page,
          pagesize: queryArg.pagesize,
          search: queryArg.search,
          order: queryArg.order,
        },
        method: 'GET',
      }),
    }),
    cancelConnectionRegister: builder.mutation({
      query: (queryArg) => ({
        url: `integrations/connections/register`,
        method: 'DELETE',
        body: queryArg.body,
      }),
    }),
    addKubernetesConfig: builder.mutation({
      query: (queryArg) => ({
        url: `system/kubernetes`,
        method: 'POST',
        body: queryArg.body,
      }),
      invalidatesTags: () => [{ type: TAGS.CONNECTIONS }],
    }),
  }),
});

export const {
  useGetCredentialsQuery,
  useVerifyAndRegisterConnectionMutation,
  useConnectToConnectionMutation,
  useGetConnectionsQuery,
  useGetConnectionStatusQuery,
  useLazyGetConnectionDetailsQuery,
  useVerifyConnectionURLMutation,
  useConnectionMetaDataMutation,
  useConfigureConnectionMutation,
  useUpdateConnectionMutation,
  useGetAllConnectionStatusQuery,
  useGetConnectionByKindQuery,
  useCancelConnectionRegisterMutation,
  useAddKubernetesConfigMutation,
} = connectionsApi;
