import { api, mesheryApiPath } from './index';

/**
 * RTK Query slice for the clean, connection-driven Grafana telemetry API
 * (backend routes under /api/telemetry/grafana/{connectionID}/*).
 *
 * This is intentionally separate from the legacy ./telemetry.ts slice, which
 * targets the deprecated /api/telemetry/metrics/grafana/* endpoints.
 */

const TAGS = {
  PINNED_BOARDS: 'GrafanaPinnedBoards',
};

const base = (connectionID: string) => `telemetry/grafana/${connectionID}`;

const telemetryGrafanaApi = api.injectEndpoints({
  endpoints: (builder) => ({
    // Health-check a Grafana connection (reachability + credential).
    pingGrafanaConnection: builder.query({
      query: ({ connectionID }) => ({
        url: mesheryApiPath(`${base(connectionID)}/ping`),
        method: 'GET',
        credentials: 'include',
      }),
    }),

    // Search dashboards on the Grafana instance.
    searchGrafanaBoards: builder.query({
      query: ({ connectionID, search }) => ({
        url: mesheryApiPath(`${base(connectionID)}/boards`),
        params: search ? { search } : undefined,
        method: 'GET',
        credentials: 'include',
      }),
    }),

    // Fetch a single dashboard's panels.
    getGrafanaBoard: builder.query({
      query: ({ connectionID, uid }) => ({
        url: mesheryApiPath(`${base(connectionID)}/boards/${uid}`),
        method: 'GET',
        credentials: 'include',
      }),
    }),

    // List the instance's datasources.
    getGrafanaDatasources: builder.query({
      query: ({ connectionID }) => ({
        url: mesheryApiPath(`${base(connectionID)}/datasources`),
        method: 'GET',
        credentials: 'include',
      }),
    }),

    // Proxy a Prometheus-style range query for rendering a panel.
    queryGrafanaRange: builder.query({
      query: ({ connectionID, ds, query, start, end, step }) => ({
        url: mesheryApiPath(`${base(connectionID)}/query_range`),
        params: { ds, query, start, end, step },
        method: 'GET',
        credentials: 'include',
      }),
    }),

    // Batch-proxy every query for a board in a single request; the backend fans
    // out to Grafana concurrently and returns one result per query id. This keeps
    // a board's render to one round trip instead of one request per panel target.
    queryGrafanaRangeBatch: builder.query({
      query: ({ connectionID, start, end, step, queries }) => ({
        url: mesheryApiPath(`${base(connectionID)}/query_range_batch`),
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: { start, end, step, queries },
      }),
    }),

    // Read the user's pinned (added) boards for a connection.
    getPinnedBoards: builder.query({
      query: ({ connectionID }) => ({
        url: mesheryApiPath(`${base(connectionID)}/pinned`),
        method: 'GET',
        credentials: 'include',
      }),
      providesTags: (_result, _err, arg) => [{ type: TAGS.PINNED_BOARDS, id: arg.connectionID }],
    }),

    // Replace the user's pinned boards for a connection.
    updatePinnedBoards: builder.mutation({
      query: ({ connectionID, boards }) => ({
        url: mesheryApiPath(`${base(connectionID)}/pinned`),
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: boards,
      }),
      invalidatesTags: (_result, _err, arg) => [{ type: TAGS.PINNED_BOARDS, id: arg.connectionID }],
    }),
  }),
});

export const {
  usePingGrafanaConnectionQuery,
  useLazyPingGrafanaConnectionQuery,
  useSearchGrafanaBoardsQuery,
  useGetGrafanaBoardQuery,
  useGetGrafanaDatasourcesQuery,
  useLazyQueryGrafanaRangeQuery,
  useQueryGrafanaRangeBatchQuery,
  useGetPinnedBoardsQuery,
  useUpdatePinnedBoardsMutation,
} = telemetryGrafanaApi;

export default telemetryGrafanaApi;
