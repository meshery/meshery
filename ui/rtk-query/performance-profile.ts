import { api } from './index';

const TAGS = {
  PERFORMANCE_PROFILE: 'performance-profile',
};

const performanceProfile = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.PERFORMANCE_PROFILE],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getPerformanceProfiles: builder.query({
        query: (queryArg) => ({
          url: `user/performance/profiles`,
          params: {
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
          },
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.PERFORMANCE_PROFILE }],
      }),
      savePerformanceProfile: builder.mutation({
        query: (queryArg) => ({
          url: `user/performance/profiles`,
          method: 'POST',
          body: queryArg.body,
        }),
        invalidatesTags: [{ type: TAGS.PERFORMANCE_PROFILE }],
      }),
      getProfileResults: builder.query({
        query: (queryArg) => ({
          url: `user/performance/profiles/results`,
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
        providesTags: () => [{ type: TAGS.PERFORMANCE_PROFILE }],
      }),
      getPerformanceProfileById: builder.query({
        query: (queryArg) => ({
          url: `user/performance/profiles/${queryArg.id}`,
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.PERFORMANCE_PROFILE }],
      }),
      deletePerformanceProfile: builder.mutation({
        query: (queryArg) => ({
          url: `user/performance/profiles/${queryArg.id}`,
          method: 'DELETE',
        }),
        invalidatesTags: [{ type: TAGS.PERFORMANCE_PROFILE }],
      }),
      getProfileResultsById: builder.query({
        query: (queryArg) => ({
          url: `user/performance/profiles/${queryArg.id}/results`,
          params: {
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
          },
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.PERFORMANCE_PROFILE }],
      }),
    }),
  });

export const {
  useGetPerformanceProfilesQuery,
  useSavePerformanceProfileMutation,
  useGetProfileResultsQuery,
  useGetPerformanceProfileByIdQuery,
  useDeletePerformanceProfileMutation,
  useGetProfileResultsByIdQuery,
} = performanceProfile;
