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
    }),
  });

export const {
  useGetEnvironmentConnectionsQuery,
  useAddConnectionToEnvironmentMutation,
  useRemoveConnectionFromEnvironmentMutation,
} = connectionsApi;
