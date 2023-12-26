import { api } from './index';

const TAGS = {
  WORKSPACES: 'workspaces',
};
const workspacesApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.WORKSPACES],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getWorkspaces: builder.query({
        query: (queryArg) => ({
          url: `workspaces`,
          params: {
            search: queryArg.search,
            order: queryArg.order,
            page: queryArg.page || 0,
            pagesize: queryArg.pagesize || 'all',
            orgID: queryArg.orgId,
          },
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.WORKSPACES }],
      }),

      createWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: `workspaces`,
          method: 'POST',
          body: queryArg.workspacePayload,
        }),

        invalidatesTags: () => [{ type: TAGS.WORKSPACES }],
      }),

      updateWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: `workspaces/${queryArg.workspaceId}`,
          method: 'PUT',
          body: queryArg.workspacePayload,
        }),

        invalidatesTags: () => [{ type: TAGS.WORKSPACES }],
      }),

      deleteWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: `workspaces/${queryArg.workspaceId}`,
          method: 'DELETE',
        }),

        invalidatesTags: () => [{ type: TAGS.WORKSPACES }],
      }),
    }),
  });

export const {
  useGetWorkspacesQuery,
  useCreateWorkspaceMutation,
  useUpdateWorkspaceMutation,
  useDeleteWorkspaceMutation,
} = workspacesApi;
