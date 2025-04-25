import { api } from './index';
import { initiateQuery } from './utils';

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
        providesTags: [{ type: TAGS.VIEWS }],
      }),
      updateViewVisibility: builder.mutation({
        query: ({ id, body }) => ({
          url: `extensions/api/content/views/${id}`,
          method: 'PUT',
          body: body,
        }),
        invalidatesTags: [{ type: TAGS.VIEWS }],
      }),
    }),
  });

export const getView = async ({ viewId }) => {
  return await initiateQuery(viewsApi.endpoints.getView, { viewId });
};
export const { useUpdateViewVisibilityMutation } = viewsApi;
