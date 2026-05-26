import {
  useDeletePerformanceProfileMutation as useSchemasDeletePerformanceProfileMutation,
  useGetPerformanceProfileQuery as useSchemasGetPerformanceProfileQuery,
  useGetPerformanceProfilesQuery as useSchemasGetPerformanceProfilesQuery,
  useUpsertPerformanceProfileMutation as useSchemasUpsertPerformanceProfileMutation,
} from '@meshery/schemas/mesheryApi';

import { api } from './index';

const PERFORMANCE_PROFILE_RESULTS_TAG = 'performance-profile-results';

const performanceProfileResults = api
  .enhanceEndpoints({
    addTagTypes: [PERFORMANCE_PROFILE_RESULTS_TAG],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getProfileResults: builder.query({
        query: (queryArg) => ({
          url: `/api/user/performance/profiles/results`,
          params: {
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
            from: queryArg.from,
            to: queryArg.to,
          },
          method: 'GET',
        }),
        providesTags: () => [{ type: PERFORMANCE_PROFILE_RESULTS_TAG }],
      }),
      getProfileResultsById: builder.query({
        query: (queryArg) => ({
          url: `/api/user/performance/profiles/${queryArg.id}/results`,
          params: {
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
          },
          method: 'GET',
        }),
        providesTags: () => [{ type: PERFORMANCE_PROFILE_RESULTS_TAG }],
      }),
    }),
  });

export const { useGetProfileResultsQuery, useGetProfileResultsByIdQuery } =
  performanceProfileResults;

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
