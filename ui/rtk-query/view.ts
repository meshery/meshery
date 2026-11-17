import { urlEncodeArrayParam, urlEncodeParams } from '@/utils/utils';
import { api, mesheryApiPath } from './index';
import { initiateQuery } from './utils';
import _ from 'lodash';

const TAGS = {
  VIEWS: 'view',
};

export const viewsApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.VIEWS],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getView: builder.query({
        query: ({ viewId }) => ({
          url: mesheryApiPath(`extensions/api/content/views/${viewId}`),
          method: 'GET',
        }),
        providesTags: (result, error, { viewId }) => [{ type: TAGS.VIEWS, id: viewId }],
      }),
      // The view mutations invalidate only the changed view's id, so an open
      // getView refetches. They deliberately leave fetchViews alone: it merges
      // every page into one cache entry, so an invalidation would refetch the
      // current page and append it a second time. List callers refetch from
      // page 0 themselves.
      deleteView: builder.mutation({
        query: ({ id }) => ({
          url: mesheryApiPath(`extensions/api/content/views/${id}`),
          method: 'DELETE',
        }),
        invalidatesTags: (result, error, { id }) => [{ type: TAGS.VIEWS, id }],
      }),
      updateViewVisibility: builder.mutation({
        query: ({ id, body }) => ({
          url: mesheryApiPath(`extensions/api/content/views/${id}`),
          method: 'PUT',
          body: body,
        }),
        invalidatesTags: (result, error, { id }) => [{ type: TAGS.VIEWS, id }],
      }),
      fetchViews: builder.query({
        query: (queryArg) =>
          mesheryApiPath(
            `extensions/api/content/views?${urlEncodeArrayParam(
              'visibility',
              queryArg.visibility,
            )}&${urlEncodeParams({
              page: queryArg.page,
              pagesize: queryArg.pagesize,
              search: queryArg.search || '',
              order: queryArg.order || '',
              trim: queryArg.trim || false,
              shared: queryArg.shared || false,
              userId: queryArg.userId,
            })}`,
          ),

        // Only have one cache entry
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
            views: [...(currentCache?.views || []), ...(newItems?.views || [])],
          };
        },
        // Refetch when the page arg changes
        forceRefetch({ currentArg, previousArg }) {
          return !_.isEqual(currentArg, previousArg);
        },
        providesTags: () => [{ type: TAGS.VIEWS }],
      }),
    }),
  });

export const getView = async ({ viewId }) => {
  return await initiateQuery(viewsApi.endpoints.getView, { viewId });
};
export const {
  useGetViewQuery,
  useUpdateViewVisibilityMutation,
  useFetchViewsQuery,
  useDeleteViewMutation,
} = viewsApi;
