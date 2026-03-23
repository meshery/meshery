import { api } from './index';

const TAGS = {
  ENVIRONMENT_CONNECTIONS: 'enivroment_connections',
};
const environmentsApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.ENVIRONMENT_CONNECTIONS],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getEnvironments: builder.query({
        query: (queryArg = {}) => ({
          url: 'environments',
          params: {
            search: queryArg.search,
            order: queryArg.order,
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            orgID: queryArg.orgId,
          },
          method: 'GET',
        }),
      }),

      createEnvironment: builder.mutation({
        query: (queryArg) => ({
          url: 'environments',
          method: 'POST',
          body: {
            name: queryArg.environmentPayload?.name,
            description: queryArg.environmentPayload?.description,
            OrganizationID:
              queryArg.environmentPayload?.OrganizationID ||
              queryArg.environmentPayload?.organization_id,
          },
        }),
      }),

      addConnectionToEnvironment: builder.mutation({
        query: (queryArg) => ({
          url: `environments/${queryArg.environmentId}/connections/${queryArg.connectionId}`,
          method: 'POST',
        }),
      }),

      removeConnectionFromEnvironment: builder.mutation({
        query: (queryArg) => ({
          url: `environments/${queryArg.environmentId}/connections/${queryArg.connectionId}`,
          method: 'DELETE',
        }),
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
            filter: queryArg.filter,
          },
          method: 'GET',
        }),
        providesTags: (_result, _error, arg) => [
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

export const useGetEnvironmentsQuery = (queryArgs = {}, options = undefined) =>
  environmentsApi.endpoints.getEnvironments.useQuery(queryArgs, options);

export const useCreateEnvironmentMutation = () =>
  environmentsApi.endpoints.createEnvironment.useMutation();

export const useAddConnectionToEnvironmentMutation = () =>
  environmentsApi.endpoints.addConnectionToEnvironment.useMutation();

export const useRemoveConnectionToEnvironmentMutation = () =>
  environmentsApi.endpoints.removeConnectionFromEnvironment.useMutation();

export const useRemoveConnectionFromEnvironmentMutation = () =>
  environmentsApi.endpoints.removeConnectionFromEnvironment.useMutation();

export const useUpdateEnvironmentMutation = () =>
  environmentsApi.endpoints.updateEnvironment.useMutation();

export const useDeleteEnvironmentMutation = () =>
  environmentsApi.endpoints.deleteEnvironment.useMutation();

export const useGetEnvironmentConnectionsQuery = (queryArgs, options = undefined) =>
  environmentsApi.endpoints.getEnvironmentConnections.useQuery(queryArgs, options);

export const useSaveEnvironmentMutation = () =>
  environmentsApi.endpoints.saveEnvironment.useMutation();
