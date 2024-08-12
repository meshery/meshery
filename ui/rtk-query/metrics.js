import { METRICS } from '@/constants/navigator';
import { api } from './index';

const TAGS = {
  METRICS: 'metrics',
};

export const metricsApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.METRICS],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      pingGrafana: builder.query({
        query: (queryArg) => ({
          url: `telemetry/metrics/grafana/ping/${queryArg.connectionId}`,
          providesTags: () => [{ type: TAGS.METRICS, id: 'GRAFANA' }],
        }),
      }),
      pingPrometheus: builder.query({
        query: (queryArg) => ({
          url: `telemetry/metrics/ping/${queryArg.connectionId}`,
          providesTags: () => [{ type: TAGS.METRICS, id: 'PROMETHEUS' }],
        }),
      }),
      updatePrometheusBoard: builder.mutation({
        query: (queryArg) => ({
          url: `/api/telemetry/metrics/board_import/${self.props.connectionID}`,
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: queryArg.body,
        }),
        invalidatesTags: [{ type: TAGS.METRICS, id: 'PROMETHEUS' }],
      }),
      getPrometheusConfiguration: builder.query({
        query: () => '/telemetry/metrics/config',
        method: 'GET',
        providesTags: [{ type: TAGS.METRICS, id: 'PROMETHEUS' }],
      }),
      configureGrafana: builder.mutation({
        query: (queryArg) => ({
          url: 'telemetry/metrics/grafana/config',
          method: 'POST',
          credentials: 'include',
          headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
          body: queryArg.params,
        }),
        invalidatesTags: [{ type: TAGS.METRICS, id: 'GRAFANA' }],
      }),
      configurePrometheus: builder.mutation({
        query: (queryArg) => ({
          url: '/telemetry/metrics/config',
          method: 'POST',
          credentials: 'include',
          headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
          body: queryArg.params,
        }),
        invalidatesTags: [{ type: TAGS.METRICS, id: 'PROMETHEUS' }],
      }),
      deletePrometheusConfig: builder.mutation({
        query: (queryArg) => ({
          url: 'telemetry/metrics/config',
          method: 'DELETE',
        }),
        invalidatesTags: [{ type: TAGS.METRICS, id: 'PROMETHEUS' }],
      }),
      getPrometheusQuery: builder.query({
        query: (queryArg) => `telemetry/metrics/query/${queryArg.connectionId}`,
        method: 'GET',
        providesTags: [{ type: TAGS.METRICS, id: 'PROMETHEUS' }],
      }),
      getPrometheusStaticBoard: builder.query({
        query: () => 'telemetry/metrics/static-board',
        method: 'GET',
        providesTags: [{ type: TAGS.METRICS, id: 'PROMETHEUS' }],
      }),
      getGrafanBoard: builder.query({
        query: (queryArg) => `/telemetry/metrics/grafana/${queryArg.connectionID}`,
        method: 'GET',
        providesTags: [{ type: METRICS.TAGS, id: 'GRAFANA' }],
      }),
      getGrafanaConfiguration: builder.query({
        query: () => 'telemetry/metrics/grafana/config',
        method: 'GET',
        providesTags: [{ type: TAGS.METRICS, id: 'GRAFANA' }],
      }),
      deleteGrafanaConfiugration: builder.mutation({
        query: () => ({
          url: 'telemetry/metrics/grafana/config',
          method: 'DELETE',
        }),
        invalidatesTags: [{ type: TAGS.METRICS, id: 'GRAFANA' }],
      }),
    }),
  });

export const {
  usePingGrafanaQuery,
  usePingPrometheusQuery,
  useUpdatePrometheusBoardMutation,
  useConfigureGrafanaMutation,
  useGetPrometheusStaticBoardQuery,
  useDeleteGrafanaConfiugrationMutation,
  useGetGrafanBoardQuery,
  useGetGrafanaConfigurationQuery,
  useGetPrometheusConfigurationQuery,
  useGetPrometheusQueryQuery,
  useConfigurePrometheusMutation,
} = metricsApi;
