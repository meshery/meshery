import {
  useDeletePerformanceProfileMutation as useSchemasDeletePerformanceProfileMutation,
  useGetPerformanceProfileQuery as useSchemasGetPerformanceProfileQuery,
  useGetPerformanceProfilesQuery as useSchemasGetPerformanceProfilesQuery,
  useUpsertPerformanceProfileMutation as useSchemasUpsertPerformanceProfileMutation,
} from '@meshery/schemas/mesheryApi';

export const useGetPerformanceProfilesQuery = (queryArg, options) =>
  useSchemasGetPerformanceProfilesQuery(
    {
      page: queryArg?.page?.toString(),
      pagesize: queryArg?.pagesize?.toString(),
      search: queryArg?.search,
      order: queryArg?.order,
    },
    options,
  );

export const useSavePerformanceProfileMutation = () => {
  const [trigger, result] = useSchemasUpsertPerformanceProfileMutation();
  const wrappedTrigger = (queryArg) => trigger({ body: queryArg?.body });

  return [wrappedTrigger, result] as const;
};

export const useGetPerformanceProfileByIdQuery = (queryArg, options) =>
  useSchemasGetPerformanceProfileQuery(
    {
      performanceProfileId: queryArg?.id,
    },
    options,
  );

export const useDeletePerformanceProfileMutation = () => {
  const [trigger, result] = useSchemasDeletePerformanceProfileMutation();
  const wrappedTrigger = (queryArg) =>
    trigger({
      performanceProfileId: queryArg?.id,
    });

  return [wrappedTrigger, result] as const;
};
