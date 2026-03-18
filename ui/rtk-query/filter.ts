import { api } from './index';

const TAGS = {
  FILTERS: 'filters',
};
const filters = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.FILTERS],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getFilters: builder.query({
        query: (queryArg) => ({
          url: `filter`,
          params: {
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
            visibility: queryArg.visibility,
          },
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.FILTERS }],
      }),
      cloneFilter: builder.mutation({
        query: (queryArg) => ({
          url: `filter/clone/${queryArg.filterID}`,
          method: 'POST',
          body: queryArg.body,
        }),
      }),
      publishFilter: builder.mutation({
        query: (queryArg) => ({
          url: `filter/catalog/publish`,
          method: 'POST',
          body: queryArg.publishBody,
        }),
      }),
      unpublishFilter: builder.mutation({
        query: (queryArg) => ({
          url: `filter/catalog/unpublish`,
          method: 'DELETE',
          body: queryArg.unpublishBody,
        }),
      }),
      deleteFilter: builder.mutation({
        query: (queryArg) => ({
          url: `filter/${queryArg.id}`,
          method: 'DELETE',
        }),
      }),
      updateFilterFile: builder.mutation({
        query: (queryArg) => ({
          url: `filter`,
          method: 'POST',
          body: queryArg.updateBody,
        }),
      }),
      uploadFilterFile: builder.mutation({
        query: (queryArg) => ({
          url: `filter`,
          headers: {
            'Content-Type': 'application/octet-stream', // Set appropriate content type for binary data
          },
          method: 'POST',
          body: queryArg.uploadBody,
        }),
      }),
    }),
  });

export const {
  useGetFiltersQuery,
  useCloneFilterMutation,
  usePublishFilterMutation,
  useUnpublishFilterMutation,
  useDeleteFilterMutation,
  useUpdateFilterFileMutation,
  useUploadFilterFileMutation,
} = filters;
