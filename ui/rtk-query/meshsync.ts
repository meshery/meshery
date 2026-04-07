import { urlEncodeParams } from '@/utils/utils';
import { api, mesheryApiPath } from './index';

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
          url: mesheryApiPath(`system/meshsync/resources`),
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
            labels: queryArg.labels,
            status: queryArg.status,
            annotation: queryArg.annotation,
            annotations: queryArg.annotations,
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
          return mesheryApiPath(`system/meshsync/resources/summary?${params}`);
        },

        providesTags: () => [{ type: TAGS.MESH_SYNC }],
      }),
      deleteMeshsyncResource: builder.mutation({
        query: (resourceId) => ({
          url: mesheryApiPath(`system/meshsync/resources/${resourceId}`),
          method: 'DELETE',
          credentials: 'include',
        }),
        invalidatesTags: [{ type: TAGS.MESH_SYNC }],
      }),
    }),
  });

export const {
  useGetMeshSyncResourcesQuery,
  useLazyGetMeshSyncResourcesQuery,
  useGetMeshSyncResourceKindsQuery,
  useDeleteMeshsyncResourceMutation,
} = meshSyncApi;
