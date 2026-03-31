import {
  mesheryApi,
  useGetConnectionsQuery as useSchemasGetConnectionsQuery,
} from '@meshery/schemas/dist/mesheryApi';
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

export const {
  useGetCredentialsQuery,
  useVerifyAndRegisterConnectionMutation,
  useConnectToConnectionMutation,
  useLazyGetConnectionDetailsQuery,
  useVerifyConnectionURLMutation,
  useConnectionMetaDataMutation,
  useConfigureConnectionMutation,
  useUpdateConnectionByIdMutation,
  useCancelConnectionRegisterMutation,
  useAddKubernetesConfigMutation,
} = connectionsApi;

export const useGetConnectionsQuery = (queryArg, options) =>
  useSchemasGetConnectionsQuery(
    {
      page: queryArg?.page?.toString(),
      pagesize: queryArg?.pagesize?.toString(),
      search: queryArg?.search,
      order: queryArg?.order,
      status: queryArg?.status,
      kind: queryArg?.kind,
      type: queryArg?.type,
      name: queryArg?.name,
    },
    options,
  );

export const useLazyGetConnectionsQuery = () => {
  const [trigger, result, lastPromiseInfo] = mesheryApi.endpoints.getConnections.useLazyQuery();

  const wrappedTrigger = (queryArg, preferCacheValue) =>
    trigger(
      {
        page: queryArg?.page?.toString(),
        pagesize: queryArg?.pagesize?.toString(),
        search: queryArg?.search,
        order: queryArg?.order,
        status: queryArg?.status,
        kind: queryArg?.kind,
        type: queryArg?.type,
        name: queryArg?.name,
      },
      preferCacheValue,
    );

  return [wrappedTrigger, result, lastPromiseInfo] as const;
};
