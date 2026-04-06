import {
  useAddConnectionToEnvironmentMutation as useSchemasAddConnectionToEnvironmentMutation,
  useCreateEnvironmentMutation as useSchemasCreateEnvironmentMutation,
  useGetEnvironmentsQuery as useSchemasGetEnvironmentsQuery,
  useRemoveConnectionFromEnvironmentMutation as useSchemasRemoveConnectionFromEnvironmentMutation,
} from '@meshery/schemas/mesheryApi';
import { api, mesheryApiPath } from './index';

const TAGS = {
  ENVIRONMENT_CONNECTIONS: 'enivroment_connections',
};
const connectionsApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.ENVIRONMENT_CONNECTIONS],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      updateEnvironment: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`environments/${queryArg.environmentId}`),
          method: 'PUT',
          body: queryArg.environmentPayload,
        }),

        invalidatesTags: () => [{ type: TAGS.ENVIRONMENT_CONNECTIONS }],
      }),

      deleteEnvironment: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`environments/${queryArg.environmentId}`),
          method: 'DELETE',
        }),

        invalidatesTags: () => [{ type: TAGS.ENVIRONMENT_CONNECTIONS }],
      }),

      getEnvironmentConnections: builder.query({
        query: (queryArg) => ({
          url: mesheryApiPath(`environments/${queryArg.environmentId}/connections`),
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
          url: mesheryApiPath(`environments`),
          method: 'POST',
          body: queryArg.body,
        }),
        invalidatesTags: [{ type: TAGS.ENVIRONMENT_CONNECTIONS }],
      }),
    }),
  });

export const {
  useUpdateEnvironmentMutation,
  useDeleteEnvironmentMutation,
  useGetEnvironmentConnectionsQuery,
  useSaveEnvironmentMutation,
} = connectionsApi;

export const useGetEnvironmentsQuery = (queryArg, options) =>
  useSchemasGetEnvironmentsQuery(
    {
      search: queryArg?.search,
      order: queryArg?.order,
      page: queryArg?.page?.toString(),
      pagesize: queryArg?.pagesize?.toString(),
      orgId: queryArg?.orgId,
    },
    options,
  );

export const useCreateEnvironmentMutation = () => {
  const [trigger, result] = useSchemasCreateEnvironmentMutation();

  const wrappedTrigger = (queryArg) =>
    trigger({
      body: {
        name: queryArg.environmentPayload?.name,
        description: queryArg.environmentPayload?.description,
        OrganizationID:
          queryArg.environmentPayload?.OrganizationID ||
          queryArg.environmentPayload?.organization_id,
      },
    });

  return [wrappedTrigger, result] as const;
};

export const useAddConnectionToEnvironmentMutation = () => {
  const [trigger, result] = useSchemasAddConnectionToEnvironmentMutation();

  const wrappedTrigger = (queryArg) =>
    trigger({
      environmentId: queryArg.environmentId,
      connectionId: queryArg.connectionId,
    });

  return [wrappedTrigger, result] as const;
};

export const useRemoveConnectionFromEnvironmentMutation = () => {
  const [trigger, result] = useSchemasRemoveConnectionFromEnvironmentMutation();

  const wrappedTrigger = (queryArg) =>
    trigger({
      environmentId: queryArg.environmentId,
      connectionId: queryArg.connectionId,
    });

  return [wrappedTrigger, result] as const;
};
