import { api } from './index';

const TAGS = {
  MESH_SYNC: 'meshsync',
};

const meshSyncApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getMeshSyncResources: builder.query({
      query: (queryArg) => ({
        url: `system/meshsync/resources`,
        params: {
          page: queryArg.page,
          pagesize: queryArg.pagesize,
          search: queryArg.search,
          order: queryArg.order,
          kind: queryArg.kind,
          clusterId: queryArg.clusterId,
          label: queryArg.label,
          status: queryArg.status,
          annotation: queryArg.annotation,
          spec: queryArg.spec,
          apiVersion: queryArg.apiVersion,
        },
        method: 'GET',
      }),
      providesTags: () => [{ type: TAGS.MESH_SYNC }],
    }),
    getMeshSyncResourceKinds: builder.query({
      query: (queryArg) => ({
        url: `system/meshsync/resources/kinds`,
        params: {
          page: queryArg.page,
          pagesize: queryArg.pagesize,
          search: queryArg.search,
          order: queryArg.order,
          clusterId: queryArg.clusterId,
        },
        method: 'GET',
      }),
      providesTags: () => [{ type: TAGS.MESH_SYNC }],
    }),
  }),
});

export const { useGetMeshSyncResourcesQuery, useGetMeshSyncResourceKindsQuery } = meshSyncApi;
