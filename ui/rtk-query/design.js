import { urlEncodeParams } from '@/utils/utils';
import { api } from './index';
import { ctxUrl } from '@/utils/multi-ctx';
import { initiateQuery } from './utils';
import _ from 'lodash';

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
        query: (queryArg) => {
          const params = urlEncodeParams({
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
            visibility: queryArg.visibility,
            populate: queryArg.populate,
          });
          return `pattern?${params}`;
        },
        providesTags: () => [{ type: TAGS.DESIGNS }],
      }),
      getDesign: builder.query({
        query: ({ design_id }) => ({
          url: `pattern/${design_id}`,
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.DESIGNS }],
      }),
      getUserDesigns: builder.query({
        query: (queryArg) => {
          const params = urlEncodeParams({
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            order: queryArg.order,
            user_id: queryArg.user_id,
            expandUser: queryArg.expandUser,
            metrics: queryArg.metrics,
            search: queryArg.search,
            visibility: queryArg.visibility,
            orgID: queryArg.orgId,
            shared: queryArg.shared || false,
          });
          return `extensions/api/content/patterns?${params}`;
        },
        serializeQueryArgs: ({ endpointName }) => {
          return endpointName;
        },

        // Always merge incoming data to the cache entry
        merge: (currentCache, newItems, { arg }) => {
          if (arg.page === 0) {
            return newItems;
          }
          return {
            ...(currentCache || {}),
            ...(newItems || {}),
            patterns: [...(currentCache?.patterns || []), ...(newItems?.patterns || [])],
          };
        },

        // Refetch when any arg changes
        forceRefetch({ currentArg, previousArg }) {
          return !_.isEqual(currentArg, previousArg);
        },
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
          url: `pattern/import`,
          method: 'POST',
          body: queryArg.importBody,
        }),
        invalidatesTags: () => [{ type: TAGS.DESIGNS }],
      }),
      deletePatternFile: builder.mutation({
        query: (queryArg) => ({
          url: `pattern/${queryArg.id}`,
          method: 'DELETE',
        }),
        providesTags: () => [{ type: TAGS.DESIGNS }],
        invalidatesTags: () => [{ type: TAGS.DESIGNS }],
      }),
      updatePatternFile: builder.mutation({
        query: (queryArg) => ({
          url: `pattern`,
          method: 'POST',
          credentials: 'include',
          body: queryArg.updateBody,
        }),
        providesTags: () => [{ type: TAGS.DESIGNS }],
      }),
      uploadPatternFile: builder.mutation({
        query: (queryArg) => ({
          url: `pattern/`,
          method: 'POST',
          body: queryArg.uploadBody,
        }),
      }),
      downloadPatternFile: builder.query({
        query: (queryArg) => `pattern/${queryArg.id}`,
      }),
    }),
  });

export const getDesign = async ({ design_id }) => {
  return await initiateQuery(
    designsApi.endpoints.getDesign,
    { design_id },
    {
      forceRefetch: true,
    },
  );
};

export const {
  useGetPatternsQuery,
  useGetDesignQuery,
  useGetUserDesignsQuery,
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
  useDownloadPatternFileQuery,
} = designsApi;
