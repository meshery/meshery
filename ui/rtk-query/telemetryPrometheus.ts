import { api, mesheryApiPath } from './index';

/**
 * RTK Query slice for the Prometheus-native telemetry API (backend routes under
 * /api/telemetry/prometheus/{connectionID}/*).
 *
 * Unlike Grafana telemetry (which browses dashboards on a server), this module
 * is built around metric exploration: list metric names, inspect metadata,
 * resolve label values, run PromQL, and save panels per connection.
 */

const TAGS = {
  PANELS: 'PrometheusPanels',
};

const base = (connectionID: string) => `telemetry/prometheus/${connectionID}`;

const telemetryPrometheusApi = api.injectEndpoints({
  endpoints: (builder) => ({
    // Health-check a Prometheus connection (reachability + version).
    pingPrometheusConnection: builder.query({
      query: ({ connectionID }) => ({
        url: mesheryApiPath(`${base(connectionID)}/ping`),
        method: 'GET',
        credentials: 'include',
      }),
    }),

    // List metric names, optionally filtered by a substring search.
    listPrometheusMetrics: builder.query({
      query: ({ connectionID, search, limit }) => ({
        url: mesheryApiPath(`${base(connectionID)}/metrics`),
        params: { search: search || undefined, limit: limit || undefined },
        method: 'GET',
        credentials: 'include',
      }),
    }),

    // Resolve a label's values (for building/refining PromQL matchers).
    prometheusLabelValues: builder.query({
      query: ({ connectionID, label, match }) => ({
        url: mesheryApiPath(`${base(connectionID)}/label_values`),
        params: { label, match: match || undefined },
        method: 'GET',
        credentials: 'include',
      }),
    }),

    // Metric metadata (type / help / unit) for a metric name.
    prometheusMetricMetadata: builder.query({
      query: ({ connectionID, metric }) => ({
        url: mesheryApiPath(`${base(connectionID)}/metadata`),
        params: { metric: metric || undefined },
        method: 'GET',
        credentials: 'include',
      }),
    }),

    // Single range query — used by the explorer's live preview.
    queryPrometheusRange: builder.query({
      query: ({ connectionID, query, start, end, step }) => ({
        url: mesheryApiPath(`${base(connectionID)}/query_range`),
        params: { query, start, end, step },
        method: 'GET',
        credentials: 'include',
      }),
    }),

    // Batch-proxy every saved panel's query in one request; the backend fans out
    // to Prometheus concurrently and returns one result per query id.
    queryPrometheusRangeBatch: builder.query({
      query: ({ connectionID, start, end, step, queries }) => ({
        url: mesheryApiPath(`${base(connectionID)}/query_range_batch`),
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: { start, end, step, queries },
      }),
    }),

    // Read the user's saved panels for a connection.
    getPrometheusPanels: builder.query({
      query: ({ connectionID }) => ({
        url: mesheryApiPath(`${base(connectionID)}/panels`),
        method: 'GET',
        credentials: 'include',
      }),
      providesTags: (_result, _err, arg) => [{ type: TAGS.PANELS, id: arg.connectionID }],
    }),

    // Replace the user's saved panels for a connection.
    updatePrometheusPanels: builder.mutation({
      query: ({ connectionID, panels }) => ({
        url: mesheryApiPath(`${base(connectionID)}/panels`),
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: panels,
      }),
      invalidatesTags: (_result, _err, arg) => [{ type: TAGS.PANELS, id: arg.connectionID }],
    }),
  }),
});

export const {
  usePingPrometheusConnectionQuery,
  useLazyPingPrometheusConnectionQuery,
  useListPrometheusMetricsQuery,
  useLazyListPrometheusMetricsQuery,
  usePrometheusLabelValuesQuery,
  useLazyPrometheusLabelValuesQuery,
  usePrometheusMetricMetadataQuery,
  useLazyPrometheusMetricMetadataQuery,
  useQueryPrometheusRangeQuery,
  useLazyQueryPrometheusRangeQuery,
  useQueryPrometheusRangeBatchQuery,
  useGetPrometheusPanelsQuery,
  useUpdatePrometheusPanelsMutation,
} = telemetryPrometheusApi;

export default telemetryPrometheusApi;
