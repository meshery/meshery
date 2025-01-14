import { urlEncodeParams } from '@/utils/utils';
import { api } from './index';

const TAGS = {
  MESH_SYNC: 'meshsync',
};

const meshSyncApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.MESH_SYNC],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getMeshSyncResources: builder.query({
        query: (queryArg) => ({
          url: `system/meshsync/resources`,
          params: {
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
            ...(queryArg.kind ? { kind: queryArg.kind } : {}), // empty object to avoid undefined as no kind signifies all resources
            clusterIds: queryArg.clusterIds,
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
        query: ({ clusterIds = ['all'], namespaces = [], pagesize, order }) => {
          const params = urlEncodeParams({
            clusterId: clusterIds,
            namespace: namespaces,
            pagesize,
            order,
          });
          return `system/meshsync/resources/summary?${params}`;
        },

        providesTags: () => [{ type: TAGS.MESH_SYNC }],
      }),
      deleteMeshsyncResource: builder.mutation({
        query: (resourceId) => ({
          url: `system/meshsync/resources/${resourceId}`,
          method: 'DELETE',
          credentials: 'include',
        }),
        invalidatesTags: [{ type: TAGS.MESH_SYNC }],
      }),
    }),
  });

export const {
  useGetMeshSyncResourcesQuery,
  useGetMeshSyncResourceKindsQuery,
  useDeleteMeshsyncResourceMutation,
} = meshSyncApi;
