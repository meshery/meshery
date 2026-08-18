import { api, mesheryApiPath } from '.';

const resourceApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getAccessActorsInfoOfResource: builder.query({
      query: (queryArg) => ({
        url: mesheryApiPath(
          `extensions/api/resource/${queryArg.resourceType}/share/${queryArg.resourceId}/${queryArg.actorType}`,
        ),
      }),
      providesTags: ['access_update'],
    }),
    createAndRevokeResourceAccessRecord: builder.mutation({
      query: (queryArg) => ({
        url: mesheryApiPath(
          `extensions/api/resource/${queryArg.resourceType}/share/${queryArg.resourceId}`,
        ),
        method: 'POST',
        // Sistent builds this payload and it goes over the wire verbatim. The Cloud
        // API reads camelCase and silently drops anything else, so a re-cased or
        // re-wrapped body grants nothing while still returning success. Guarded by
        // components/workspaces/__tests__/ShareWorkspaceModal.payload.test.ts.
        body: queryArg.resourceAccessMappingPayload,
      }),
      invalidatesTags: ['access_update'],
    }),
  }),
});

export const {
  useGetAccessActorsInfoOfResourceQuery,
  useCreateAndRevokeResourceAccessRecordMutation,
} = resourceApi;
