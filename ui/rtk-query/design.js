import { api } from './index';
import { ctxUrl } from '@/utils/multi-ctx';

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
      deployPattern: builder.mutation({
        query: ({
          pattern_file,
          pattern_id,
          selectedK8sContexts,
          verify = false,
          dryRun = false,
        }) => ({
          url: `${ctxUrl('pattern/deploy', selectedK8sContexts)}${verify ? '&verify=true' : ''}${
            dryRun ? '&dryRun=true' : ''
          }`,
          method: 'POST',
          body: {
            pattern_file,
            pattern_id,
          },
        }),
      }),
      undeployPattern: builder.mutation({
        query: ({
          pattern_file,
          pattern_id,
          selectedK8sContexts,
          verify = false,
          dryRun = false,
        }) => ({
          url: `${ctxUrl('pattern/deploy', selectedK8sContexts)}${verify ? '&verify=true' : ''}${
            dryRun ? '&dryRun=true' : ''
          }`,
          method: 'DELETE',
          body: {
            pattern_file,
            pattern_id,
          },
        }),
      }),
    }),
  });

export const { useGetPatternsQuery, useDeployPatternMutation, useUndeployPatternMutation } =
  designs;
