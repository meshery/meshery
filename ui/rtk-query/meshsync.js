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
            ...(queryArg.model ? { model: queryArg.model } : {}),
            ...(queryArg.namespace ? { namespace: queryArg.namespace } : {}),
            clusterIds: queryArg.clusterIds,
            label: queryArg.label,
            status: queryArg.status ?? false,
            annotations: queryArg.annotations ?? false,
            spec: queryArg.spec ?? false,
            labels: queryArg.labels ?? false,
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
      getMeshSyncResourceDetails: builder.query({
        query: ({ resourceId, status = true, spec = true, annotations = true, labels = true }) => ({
          url: `system/meshsync/resources/${resourceId}`,
          params: {
            status,
            spec,
            annotations,
            labels,
          },
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.MESH_SYNC }],
      }),

      getMeshSyncResourcesWithDetails: builder.query({
        query: (queryArg) => ({
          url: `system/meshsync/resources`,
          params: {
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
            ...(queryArg.kind ? { kind: queryArg.kind } : {}),
            ...(queryArg.model ? { model: queryArg.model } : {}),
            ...(queryArg.namespace ? { namespace: queryArg.namespace } : {}),
            clusterIds: queryArg.clusterIds,
            label: queryArg.label,
            status: queryArg.status !== undefined ? queryArg.status : true,
            annotations: queryArg.annotations !== undefined ? queryArg.annotations : true,
            spec: queryArg.spec !== undefined ? queryArg.spec : true,
            labels: queryArg.labels !== undefined ? queryArg.labels : true,
            apiVersion: queryArg.apiVersion,
          },
          method: 'GET',
        }),
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
  useGetMeshSyncResourceDetailsQuery,
  useGetMeshSyncResourcesWithDetailsQuery,
  useDeleteMeshsyncResourceMutation,
} = meshSyncApi;
