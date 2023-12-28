import { api } from './index';

const TAGS = {
  ENVIRONMENT_CONNECTIONS: 'enivroment_connections',
};
const connectionsApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.ENVIRONMENT_CONNECTIONS],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getEnvironments: builder.query({
        query: (queryArg) => ({
          url: `environments`,
          params: {
            search: queryArg.search,
            order: queryArg.order,
            page: queryArg.page || 0,
            pagesize: queryArg.pagesize || 'all',
            orgID: queryArg.orgId,
          },
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.ENVIRONMENT_CONNECTIONS }],
      }),

      createEnvironment: builder.mutation({
        query: (queryArg) => ({
          url: `environments`,
          method: 'POST',
          body: queryArg.environmentPayload,
        }),

        invalidatesTags: () => [{ type: TAGS.ENVIRONMENT_CONNECTIONS }],
      }),

      updateEnvironment: builder.mutation({
        query: (queryArg) => ({
          url: `environments/${queryArg.environmentId}`,
          method: 'PUT',
          body: queryArg.environmentPayload,
        }),

        invalidatesTags: () => [{ type: TAGS.ENVIRONMENT_CONNECTIONS }],
      }),

      deleteEnvironment: builder.mutation({
        query: (queryArg) => ({
          url: `environments/${queryArg.environmentId}`,
          method: 'DELETE',
        }),

        invalidatesTags: () => [{ type: TAGS.ENVIRONMENT_CONNECTIONS }],
      }),

      getEnvironmentConnections: builder.query({
        query: (queryArg) => ({
          url: `environments/${queryArg.environmentId}/connections`,
          params: {
            page: queryArg.page || 0,
            per_page: queryArg.per_page,
            pagesize: queryArg.pagesize || 'all',
          },
          method: 'GET',
        }),
        providesTags: (_result, _error, arg) => [
          { type: TAGS.ENVIRONMENT_CONNECTIONS, id: arg.environmentId },
        ],
      }),

      addConnectionToEnvironment: builder.mutation({
        query: (queryArg) => ({
          url: `environments/${queryArg.environmentId}/connections/${queryArg.connectionId}`,
          method: 'POST',
          body: {},
        }),

        invalidatesTags: (_result, _error, arg) => [
          { type: TAGS.ENVIRONMENT_CONNECTIONS, id: arg.environmentId },
        ],
      }),

      removeConnectionFromEnvironment: builder.mutation({
        query: (queryArg) => ({
          url: `environments/${queryArg.environmentId}/connections/${queryArg.connectionId}`,
          method: 'DELETE',
          body: {},
        }),

        invalidatesTags: (_result, _error, arg) => [
          { type: TAGS.ENVIRONMENT_CONNECTIONS, id: arg.environmentId },
        ],
      }),

      saveEnvironment: builder.mutation({
        query: (queryArg) => ({
          url: `environments`,
          method: 'POST',
          body: queryArg.body,
        }),
        invalidatesTags: [{ type: TAGS.ENVIRONMENT_CONNECTIONS }],
      }),
    }),
  });

export const {
  useGetEnvironmentsQuery,
  useCreateEnvironmentMutation,
  useUpdateEnvironmentMutation,
  useDeleteEnvironmentMutation,
  useGetEnvironmentConnectionsQuery,
  useAddConnectionToEnvironmentMutation,
  useRemoveConnectionFromEnvironmentMutation,
  useSaveEnvironmentMutation,
} = connectionsApi;
