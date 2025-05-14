import { urlEncodeArrayParam, urlEncodeParams } from '@/utils/utils';
import { api } from './index';
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
          url: `extensions/api/content/views/${viewId}`,
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.VIEWS }],
      }),
      deleteView: builder.mutation({
        query: ({ id }) => ({
          url: `extensions/api/content/views/${id}`,
          method: 'DELETE',
        }),
        providesTags: () => [{ type: TAGS.VIEWS }],
        invalidatesTags: () => [{ type: TAGS.VIEWS }],
      }),
      updateViewVisibility: builder.mutation({
        query: ({ id, body }) => ({
          url: `extensions/api/content/views/${id}`,
          method: 'PUT',
          body: body,
        }),
        providesTags: () => [{ type: TAGS.VIEWS }],
      }),
      fetchViews: builder.query({
        query: (queryArg) =>
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
            user_id: queryArg.user_id,
          })}`,

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
