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
  }),
});

export const {
  useGetGrafanaBoardsQuery,
  useLazyGetGrafanaBoardsQuery,
  useGetGrafanaConfigQuery,
  useUpdateGrafanaBoardsMutation,
  useConfigureGrafanaMutation,
} = telemetryApi;
