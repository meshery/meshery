import { api } from './index';
import { ctxUrl } from '@/utils/multi-ctx';

const TAGS = {
  DESIGNS: 'designs',
};

export const designsApi = api
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
          skipCRD = false,
        }) => ({
          url: `${ctxUrl('pattern/deploy', selectedK8sContexts)}${verify ? '&verify=true' : ''}${
            dryRun ? '&dryRun=true' : ''
          }${skipCRD ? '&skipCRD=true' : ''}`,
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
      clonePattern: builder.mutation({
        query: (queryArg) => ({
          url: `pattern/clone/${queryArg.patternID}`,
          method: 'POST',
          body: queryArg.body,
        }),
      }),
      publishPattern: builder.mutation({
        query: (queryArg) => ({
          url: `pattern/catalog/publish`,
          method: 'POST',
          body: queryArg.publishBody,
        }),
      }),
      unpublishPattern: builder.mutation({
        query: (queryArg) => ({
          url: `pattern/catalog/unpublish`,
          method: 'DELETE',
          body: queryArg.unpublishBody,
        }),
      }),
      deletePattern: builder.mutation({
        query: (queryArg) => ({
          url: `patterns/delete`,
          method: 'POST',
          body: queryArg.deleteBody,
        }),
      }),
      importPattern: builder.mutation({
        query: (queryArg) => ({
          url: `pattern/${queryArg.type}`,
          method: 'POST',
          body: queryArg.importBody,
        }),
      }),
      deletePatternFile: builder.mutation({
        query: (queryArg) => ({
          url: `pattern/${queryArg.id}`,
          method: 'DELETE',
        }),
      }),
      updatePatternFile: builder.mutation({
        query: (queryArg) => ({
          url: `pattern`,
          method: 'POST',
          body: queryArg.updateBody,
        }),
      }),
      uploadPatternFile: builder.mutation({
        query: (queryArg) => ({
          url: `pattern/`,
          method: 'POST',
          body: queryArg.uploadBody,
        }),
      }),
    }),
  });

export const {
  useGetPatternsQuery,
  useDeployPatternMutation,
  useUndeployPatternMutation,
  useClonePatternMutation,
  usePublishPatternMutation,
  useUnpublishPatternMutation,
  useDeletePatternMutation,
  useImportPatternMutation,
  useUpdatePatternFileMutation,
  useUploadPatternFileMutation,
  useDeletePatternFileMutation,
} = designsApi;
