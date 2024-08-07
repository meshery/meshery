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
        query: () => 'telemetry/metrics/grafana/ping',
        providesTags: () => [{ type: TAGS.METRICS, id: 'GRAFANA' }],
      }),
      pingPrometheus: builder.query({
        query: () => 'telemetry/metrics/ping', // compare the endpoint with swagger
        providesTags: () => [{ type: TAGS.METRICS, id: 'PROMETHEUS' }],
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
    }),
  });

export const {
  usePingGrafanaQuery,
  usePingPrometheusQuery,
  useConfigureGrafanaMutation,
  useConfigurePrometheusMutation,
} = metricsApi;
