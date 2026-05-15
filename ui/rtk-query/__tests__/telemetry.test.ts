import { describe, expect, it, vi, beforeAll, beforeEach, afterEach } from 'vitest';

beforeAll(() => {
  process.env.RTK_MESHERY_ENDPOINT_PREFIX = 'http://localhost';
});

const okResponse = (body: unknown = {}) => ({
  ok: true,
  status: 200,
  redirected: false,
  headers: new Headers({ 'content-type': 'application/json' }),
  url: '',
  text: () => Promise.resolve(JSON.stringify(body)),
  json: () => Promise.resolve(body),
  clone() {
    return this;
  },
});

const setupStore = async () => {
  vi.resetModules();
  const apiMod = await import('../index');
  await import('../telemetry');
  const { configureStore } = await import('@reduxjs/toolkit');
  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (g) => g().concat(apiMod.api.middleware),
  });
  return { api: apiMod.api, store };
};

describe('telemetry endpoints', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn().mockResolvedValue(okResponse({}));
    global.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('exports all expected hooks', async () => {
    const mod = await import('../telemetry');
    expect(mod.useGetGrafanaBoardsQuery).toBeTypeOf('function');
    expect(mod.useLazyGetGrafanaBoardsQuery).toBeTypeOf('function');
    expect(mod.useGetGrafanaConfigQuery).toBeTypeOf('function');
    expect(mod.useGetStaticPrometheusBoardConfigQuery).toBeTypeOf('function');
    expect(mod.useUpdateGrafanaBoardsMutation).toBeTypeOf('function');
    expect(mod.useConfigureGrafanaMutation).toBeTypeOf('function');
    expect(mod.useGetPrometheusConfigQuery).toBeTypeOf('function');
    expect(mod.usePostBoardImportMutation).toBeTypeOf('function');
    expect(mod.useLazyQueryTemplateVarsQuery).toBeTypeOf('function');
    expect(mod.useLazyPingPrometheusQuery).toBeTypeOf('function');
    expect(mod.useLazyPingGrafanaQuery).toBeTypeOf('function');
    expect(mod.useLazyQueryRangeQuery).toBeTypeOf('function');
  });

  it('getGrafanaBoards GETs the grafana boards URL with dashboardSearch param', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.getGrafanaBoards.initiate({
        connectionID: 'conn-123',
        grafanaBoardSearch: 'node-exporter',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/telemetry/metrics/grafana/boards/conn-123');
    expect(req.url).toContain('dashboardSearch=node-exporter');
    expect(req.credentials).toBe('include');
  });

  it('getGrafanaConfig GETs the grafana config URL with form encoded content-type', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getGrafanaConfig.initiate({}));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/telemetry/metrics/grafana/config');
    expect(req.headers.get('content-type')).toContain('application/x-www-form-urlencoded');
  });

  it('getStaticPrometheusBoardConfig GETs the static-board URL', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getStaticPrometheusBoardConfig.initiate('conn-7'));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/telemetry/metrics/static-board/conn-7');
    expect(req.credentials).toBe('include');
  });

  it('updateGrafanaBoards POSTs JSON body to the boards URL', async () => {
    const { api, store } = await setupStore();
    const selectedBoardsConfigs = [{ board_id: 1, dashboard_uid: 'a' }];
    await store.dispatch(
      api.endpoints.updateGrafanaBoards.initiate({
        connectionID: 'conn-9',
        selectedBoardsConfigs,
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/telemetry/metrics/grafana/boards/conn-9');
    expect(req.headers.get('content-type')).toContain('application/json');
    expect(JSON.parse(await req.text())).toEqual(selectedBoardsConfigs);
  });

  it('configureGrafana POSTs form-encoded params to the config URL', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.configureGrafana.initiate({ params: 'grafanaURL=http%3A%2F%2Fg' }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/telemetry/metrics/grafana/config');
    expect(req.headers.get('content-type')).toContain('application/x-www-form-urlencoded');
    expect(await req.text()).toBe('grafanaURL=http%3A%2F%2Fg');
  });

  it('getPrometheusConfig GETs /api/telemetry/metrics/config', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.getPrometheusConfig.initiate({}));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/telemetry/metrics/config');
    expect(req.credentials).toBe('include');
  });

  it('postBoardImport POSTs JSON body to the board_import URL', async () => {
    const { api, store } = await setupStore();
    const body = { dashboards: ['n'] };
    await store.dispatch(api.endpoints.postBoardImport.initiate({ connectionID: 'conn-3', body }));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('POST');
    expect(req.url).toContain('/api/telemetry/metrics/board_import/conn-3');
    expect(req.headers.get('content-type')).toContain('application/json');
    expect(JSON.parse(await req.text())).toEqual(body);
  });

  it('queryTemplateVars GETs the templated query URL', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.queryTemplateVars.initiate({ connectionID: 'c-1', query: 'q=$var' }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/telemetry/metrics/query/c-1');
    expect(req.url).toContain('q=$var');
  });

  it('pingPrometheus GETs the prometheus ping URL', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.pingPrometheus.initiate({ connectionId: 'p-1' }));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/telemetry/metrics/ping/p-1');
  });

  it('queryRange GETs the {type}/query_range URL', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(
      api.endpoints.queryRange.initiate({
        type: 'telemetry/metrics',
        connectionID: 'conn-5',
        queryParams: 'q=foo',
      }),
    );
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/telemetry/metrics/query_range/conn-5');
    expect(req.url).toContain('q=foo');
  });

  it('pingGrafana GETs the grafana ping URL', async () => {
    const { api, store } = await setupStore();
    await store.dispatch(api.endpoints.pingGrafana.initiate({ connectionId: 'g-1' }));
    const req = fetchMock.mock.calls[0][0] as Request;
    expect(req.method).toBe('GET');
    expect(req.url).toContain('/api/telemetry/metrics/grafana/ping/g-1');
  });
});
