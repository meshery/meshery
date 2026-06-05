import {
  useDeletePerformanceProfileMutation as useSchemasDeletePerformanceProfileMutation,
  useGetPerformanceProfileQuery as useSchemasGetPerformanceProfileQuery,
  useGetPerformanceProfileResultsQuery as useSchemasGetPerformanceProfileResultsQuery,
  useGetPerformanceProfilesQuery as useSchemasGetPerformanceProfilesQuery,
  useGetPerformanceResultsQuery as useSchemasGetPerformanceResultsQuery,
  useUpsertPerformanceProfileMutation as useSchemasUpsertPerformanceProfileMutation,
} from '@meshery/schemas/mesheryApi';

const stripNullishParams = (params) =>
  Object.fromEntries(
    Object.entries(params).filter(([, value]) => value !== undefined && value !== null),
  );

const normalizePaginationParams = (queryArg) =>
  stripNullishParams({
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
    stripNullishParams({
      ...normalizePaginationParams(queryArg),
      from: queryArg?.from,
      to: queryArg?.to,
    }),
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
    stripNullishParams({
      performanceProfileId: queryArg?.id,
      ...normalizePaginationParams(queryArg),
    }),
    options,
  );
