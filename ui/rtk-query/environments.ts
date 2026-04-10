import {
  useAddConnectionToEnvironmentMutation as useSchemasAddConnectionToEnvironmentMutation,
  useCreateEnvironmentMutation as useSchemasCreateEnvironmentMutation,
  useDeleteEnvironmentMutation as useSchemasDeleteEnvironmentMutation,
  useGetEnvironmentConnectionsQuery as useSchemasGetEnvironmentConnectionsQuery,
  useGetEnvironmentsQuery as useSchemasGetEnvironmentsQuery,
  useRemoveConnectionFromEnvironmentMutation as useSchemasRemoveConnectionFromEnvironmentMutation,
  useUpdateEnvironmentMutation as useSchemasUpdateEnvironmentMutation,
} from '@meshery/schemas/mesheryApi';

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

export const useGetEnvironmentConnectionsQuery = (queryArg, options) =>
  useSchemasGetEnvironmentConnectionsQuery(
    {
      environmentId: queryArg?.environmentId,
      search: queryArg?.search,
      order: queryArg?.order,
      page: queryArg?.page?.toString(),
      pagesize: queryArg?.pagesize?.toString(),
      filter: queryArg?.filter,
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
        organization_id:
          queryArg.environmentPayload?.organization_id ||
          queryArg.environmentPayload?.OrganizationID,
      },
    });

  return [wrappedTrigger, result] as const;
};

export const useSaveEnvironmentMutation = () => {
  const [trigger, result] = useSchemasCreateEnvironmentMutation();

  const wrappedTrigger = (queryArg) =>
    trigger({
      body: queryArg.body,
    });

  return [wrappedTrigger, result] as const;
};

export const useUpdateEnvironmentMutation = () => {
  const [trigger, result] = useSchemasUpdateEnvironmentMutation();

  const wrappedTrigger = (queryArg) =>
    trigger({
      environmentId: queryArg.environmentId,
      body: {
        name: queryArg.environmentPayload?.name,
        description: queryArg.environmentPayload?.description,
        organization_id:
          queryArg.environmentPayload?.organization_id ||
          queryArg.environmentPayload?.OrganizationID,
      },
    });

  return [wrappedTrigger, result] as const;
};

export const useDeleteEnvironmentMutation = () => {
  const [trigger, result] = useSchemasDeleteEnvironmentMutation();

  const wrappedTrigger = (queryArg) =>
    trigger({
      environmentId: queryArg.environmentId,
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
