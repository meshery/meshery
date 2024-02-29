import { api } from './index';

const TAGS = {
  DESIGNS: 'designs',
};
const designs = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.DESIGNS],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getPatterns: builder.query({
        query: (queryArg) => ({
          url: `pattern`,
          params: {
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
            visibility: queryArg.visibility,
          },
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.DESIGNS }],
      }),

      publishPattern: builder.mutation({
        query: (queryArg) => ({
          url: `pattern/catalog/publish`,
          credentials: 'include',
          method: 'POST',
          body: queryArg.patternPayload,
        }),
      }),

      unpublishPattern: builder.mutation({
        query: (queryArg) => ({
          url: `pattern/catalog/unpublish`,
          credentials: 'include',
          method: 'DELETE',
          body: queryArg.patternPayload,
        }),
      }),

      deployPattern: builder.mutation({
        query: (queryArg) => ({
          url: queryArg.url,
          credentials: 'include',
          method: 'POST',
          body: queryArg.patternPayload,
        }),
      }),

      undeployPattern: builder.mutation({
        query: (queryArg) => ({
          url: queryArg.url,
          credentials: 'include',
          method: 'DELETE',
          body: queryArg.patternPayload,
        }),
      }),
    }),
  });

export const {
  useGetPatternsQuery,
  usePublishPatternMutation,
  useUnpublishPatternMutation,
  useDeployPatternMutation,
  useUndeployPatternMutation,
} = designs;
