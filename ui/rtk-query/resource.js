import { api } from '.';

const resourceApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getAccessActorsInfoOfResource: builder.query({
      query: (queryArg) => ({
        url: `extensions/api/resource/${queryArg.resourceType}/share/${queryArg.resourceId}/${queryArg.actorType}`,
      }),
      providesTags: ['access_update'],
    }),
    createAndRevokeResourceAccessRecord: builder.mutation({
      query: (queryArg) => ({
        url: `extensions/api/resource/${queryArg.resourceType}/share/${queryArg.resourceId}`,
        method: 'POST',
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
