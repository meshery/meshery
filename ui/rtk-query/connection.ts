import { api } from './index';

const TAGS = {
  CONNECTIONS: 'connections',
  CREDENTIALS: 'credentials',
};

const connectionsApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.CONNECTIONS, TAGS.CREDENTIALS],
  })
  .injectEndpoints({
    overrideExisting: true,
    endpoints: (builder) => ({
      getConnections: builder.query({
        query: (queryArg = {}) => ({
          url: 'integrations/connections',
          params: {
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
            filter: queryArg.filter,
            kind: queryArg.kind,
            status: queryArg.status,
            type: queryArg.type,
            name: queryArg.name,
          },
        }),
        providesTags: [TAGS.CONNECTIONS],
      }),

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
      updateConnectionById: builder.mutation({
        query: (queryArg) => ({
          url: `integrations/connections/${queryArg.connectionId}`,
          method: 'PUT',
          body: {
            status: queryArg.body?.status,
            metadata: queryArg.body?.metadata,
          },
        }),
        invalidatesTags: () => [{ type: TAGS.CONNECTIONS }],
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

export const useGetConnectionsQuery = (queryArgs = {}, options = undefined) =>
  connectionsApi.endpoints.getConnections.useQuery(queryArgs, options);

export const useLazyGetConnectionsQuery = () =>
  connectionsApi.endpoints.getConnections.useLazyQuery();

export const useGetCredentialsQuery = (queryArgs = undefined, options = undefined) =>
  connectionsApi.endpoints.getCredentials.useQuery(queryArgs, options);

export const useVerifyAndRegisterConnectionMutation = () =>
  connectionsApi.endpoints.verifyAndRegisterConnection.useMutation();

export const useConnectToConnectionMutation = () =>
  connectionsApi.endpoints.connectToConnection.useMutation();

export const useLazyGetConnectionDetailsQuery = () =>
  connectionsApi.endpoints.getConnectionDetails.useLazyQuery();

export const useVerifyConnectionURLMutation = () =>
  connectionsApi.endpoints.verifyConnectionURL.useMutation();

export const useConnectionMetaDataMutation = () =>
  connectionsApi.endpoints.connectionMetaData.useMutation();

export const useConfigureConnectionMutation = () =>
  connectionsApi.endpoints.configureConnection.useMutation();

export const useUpdateConnectionByIdMutation = () =>
  connectionsApi.endpoints.updateConnectionById.useMutation();

export const useCancelConnectionRegisterMutation = () =>
  connectionsApi.endpoints.cancelConnectionRegister.useMutation();

export const useAddKubernetesConfigMutation = () =>
  connectionsApi.endpoints.addKubernetesConfig.useMutation();
