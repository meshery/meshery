import {
  useDeletePerformanceProfileMutation as useSchemasDeletePerformanceProfileMutation,
  useGetPerformanceProfileQuery as useSchemasGetPerformanceProfileQuery,
  useGetPerformanceProfileResultsQuery as useSchemasGetPerformanceProfileResultsQuery,
  useGetPerformanceProfilesQuery as useSchemasGetPerformanceProfilesQuery,
  useGetPerformanceResultsQuery as useSchemasGetPerformanceResultsQuery,
  useUpsertPerformanceProfileMutation as useSchemasUpsertPerformanceProfileMutation,
} from '@meshery/schemas/mesheryApi';

const normalizePaginationParams = (queryArg) => ({
  page: queryArg?.page?.toString(),
  pagesize: queryArg?.pagesize?.toString(),
  search: queryArg?.search,
  order: queryArg?.order,
});

export const useGetPerformanceProfilesQuery = (queryArg, options) =>
  useSchemasGetPerformanceProfilesQuery(normalizePaginationParams(queryArg), options);

export const useSavePerformanceProfileMutation = () => {
  const [trigger, result] = useSchemasUpsertPerformanceProfileMutation();
  const wrappedTrigger = (queryArg) => trigger({ body: queryArg?.body });

  return [wrappedTrigger, result] as const;
};

export const useGetProfileResultsQuery = (queryArg, options) =>
  useSchemasGetPerformanceResultsQuery(
    {
      ...normalizePaginationParams(queryArg),
      from: queryArg?.from,
      to: queryArg?.to,
    },
    options,
  );

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

export const useGetProfileResultsByIdQuery = (queryArg, options) =>
  useSchemasGetPerformanceProfileResultsQuery(
    {
      performanceProfileId: queryArg?.id,
      ...normalizePaginationParams(queryArg),
    },
    options,
  );
