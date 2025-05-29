import { api } from './index';

const TAGS = {
  GRAFANA: 'grafana',
};

const telemetryApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getGrafanaBoards: builder.query({
      query: ({ connectionID, grafanaBoardSearch }) => ({
        url: `telemetry/metrics/grafana/boards/${connectionID}`,
        params: { dashboardSearch: grafanaBoardSearch },
        method: 'GET',
        credentials: 'include',
      }),
      providesTags: () => [{ type: TAGS.GRAFANA }],
    }),

    getGrafanaConfig: builder.query({
      query: () => ({
        url: `telemetry/metrics/grafana/config`,
        method: 'GET',
        credentials: 'include',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
      }),
    }),

    getStaticPrometheusBoardConfig: builder.query({
      query: () => ({
        url: `telemetry/metrics/static-board`,
        method: 'GET',
        credentials: 'include',
      }),
    }),

    updateGrafanaBoards: builder.mutation({
      query: ({ connectionID, selectedBoardsConfigs }) => ({
        url: `telemetry/metrics/grafana/boards/${connectionID}`,
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json;charset=UTF-8' },
        body: selectedBoardsConfigs,
      }),
      invalidatesTags: () => [{ type: TAGS.GRAFANA }],
    }),

    configureGrafana: builder.mutation({
      query: ({ params }) => ({
        url: `telemetry/metrics/grafana/config`,
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
        body: params,
      }),
      invalidatesTags: () => [{ type: TAGS.GRAFANA }],
    }),
    getPrometheusConfig: builder.query({
      query: () => ({
        url: 'telemetry/metrics/config',
        method: 'GET',
        credentials: 'include',
      }),
    }),
    postBoardImport: builder.mutation({
      query: ({ connectionID, body }) => ({
        url: `/telemetry/metrics/board_import/${connectionID}`,
        method: 'POST',
        body,
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
      }),
    }),
    queryTemplateVars: builder.query({
      query: ({ connectionID, query }) => ({
        url: `/telemetry/metrics/query/${connectionID}?${query}`,
        method: 'GET',
        credentials: 'include',
      }),
    }),
    pingPrometheus: builder.query({
      query: ({ connectionId }) => ({
        url: `telemetry/metrics/ping/${connectionId}`,
        method: 'GET',
        credentials: 'include',
      }),
    }),
    pingGrafana: builder.query({
      query: ({ connectionId }) => ({
        url: `telemetry/metrics/grafana/ping/${connectionId}`,
        method: 'GET',
        credentials: 'include',
      }),
    }),
  }),
});

export const {
  useGetGrafanaBoardsQuery,
  useLazyGetGrafanaBoardsQuery,
  useGetGrafanaConfigQuery,
  useGetStaticPrometheusBoardConfigQuery,
  useUpdateGrafanaBoardsMutation,
  useConfigureGrafanaMutation,
  useGetPrometheusConfigQuery,
  usePostBoardImportMutation,
  useLazyQueryTemplateVarsQuery,
  useLazyPingPrometheusQuery,
  useLazyPingGrafanaQuery,
} = telemetryApi;
