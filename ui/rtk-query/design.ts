import { urlEncodeParams } from '@/utils/utils';
import { api, mesheryApiPath } from './index';
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
          return mesheryApiPath(`pattern?${params}`);
        },
        providesTags: () => [{ type: TAGS.DESIGNS }],
      }),
      getDesign: builder.query({
        query: ({ design_id }) => ({
          url: mesheryApiPath(`pattern/${design_id}`),
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
            userId: queryArg.userId,
            expandUser: queryArg.expandUser,
            metrics: queryArg.metrics,
            search: queryArg.search,
            visibility: queryArg.visibility,
            orgId: queryArg.orgId,
            shared: queryArg.shared || false,
          });
          return mesheryApiPath(`extensions/api/content/patterns?${params}`);
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
          patternFile,
          patternId,
          selectedK8sContexts,
          verify = false,
          dryRun = false,
          skipCRD = false,
        }) => ({
          url: mesheryApiPath(
            `${ctxUrl('pattern/deploy', selectedK8sContexts)}${verify ? '&verify=true' : ''}${
              dryRun ? '&dryRun=true' : ''
            }${skipCRD ? '&skipCRD=true' : ''}`,
          ),
          method: 'POST',
          body: {
            patternFile,
            patternId,
          },
        }),
        invalidatesTags: [{ type: TAGS.DESIGNS }],
      }),
      undeployPattern: builder.mutation({
        query: ({
          patternFile,
          patternId,
          selectedK8sContexts,
          verify = false,
          dryRun = false,
        }) => ({
          url: mesheryApiPath(
            `${ctxUrl('pattern/deploy', selectedK8sContexts)}${verify ? '&verify=true' : ''}${
              dryRun ? '&dryRun=true' : ''
            }`,
          ),
          method: 'DELETE',
          body: {
            patternFile,
            patternId,
          },
        }),
        invalidatesTags: [{ type: TAGS.DESIGNS }],
      }),
      clonePattern: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`pattern/clone/${queryArg.patternID}`),
          method: 'POST',
          body: queryArg.body,
        }),
        invalidatesTags: [{ type: TAGS.DESIGNS }],
      }),
      publishPattern: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`pattern/catalog/publish`),
          method: 'POST',
          body: queryArg.publishBody,
        }),
        invalidatesTags: [{ type: TAGS.DESIGNS }],
      }),
      unpublishPattern: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`pattern/catalog/unpublish`),
          method: 'DELETE',
          body: queryArg.unpublishBody,
        }),
        invalidatesTags: [{ type: TAGS.DESIGNS }],
      }),
      deletePattern: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`patterns/delete`),
          method: 'POST',
          body: queryArg.deleteBody,
        }),
        invalidatesTags: [{ type: TAGS.DESIGNS }],
      }),
      importPattern: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`pattern/import`),
          method: 'POST',
          body: queryArg.importBody,
        }),
        invalidatesTags: [{ type: TAGS.DESIGNS }],
      }),
      deletePatternFile: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`pattern/${queryArg.id}`),
          method: 'DELETE',
        }),
        invalidatesTags: [{ type: TAGS.DESIGNS }],
      }),
      updatePatternFile: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`pattern`),
          method: 'POST',
          credentials: 'include',
          body: queryArg.updateBody,
        }),
        invalidatesTags: [{ type: TAGS.DESIGNS }],
      }),
      uploadPatternFile: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`pattern/`),
          method: 'POST',
          body: queryArg.uploadBody,
        }),
        invalidatesTags: [{ type: TAGS.DESIGNS }],
      }),
      downloadPatternFile: builder.query({
        query: (queryArg) => mesheryApiPath(`pattern/${queryArg.id}`),
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
  useLazyGetDesignQuery,
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
