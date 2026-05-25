import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import telemetryReducer, {
  updateGrafanaConfig,
  updatePrometheusConfig,
  updateStaticPrometheusBoardConfig,
} from '../telemetry';

describe('telemetry slice', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2024-06-01T00:00:00Z'));
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('returns the initial state', () => {
    const state = telemetryReducer(undefined, { type: 'init' } as any);
    expect(state.grafana).toEqual({
      grafanaURL: '',
      grafanaAPIKey: '',
      grafanaBoardSearch: '',
      grafanaBoards: [],
      selectedBoardsConfigs: [],
      ts: -8640000000000000,
    });
    expect(state.prometheus).toEqual({
      prometheusURL: '',
      selectedPrometheusBoardsConfigs: [],
      ts: -8640000000000000,
    });
    expect(state.staticPrometheusBoardConfig).toEqual({});
  });

  describe('updateGrafanaConfig', () => {
    it('merges payload into grafana and refreshes ts', () => {
      const initial = telemetryReducer(undefined, { type: 'init' } as any);
      vi.setSystemTime(new Date('2024-06-01T00:01:00Z'));
      const state = telemetryReducer(
        initial,
        updateGrafanaConfig({
          grafanaURL: 'http://grafana',
          grafanaAPIKey: 'key-1',
        }),
      );
      expect(state.grafana.grafanaURL).toBe('http://grafana');
      expect(state.grafana.grafanaAPIKey).toBe('key-1');
      expect(state.grafana.grafanaBoardSearch).toBe(''); // preserved
      expect(state.grafana.grafanaBoards).toEqual([]); // preserved
      expect(state.grafana.ts).toBeGreaterThan(initial.grafana.ts);
    });

    it('replaces lists when payload contains them', () => {
      const initial = telemetryReducer(undefined, { type: 'init' } as any);
      const boards = [{ id: 'b1' }];
      const state = telemetryReducer(initial, updateGrafanaConfig({ grafanaBoards: boards }));
      expect(state.grafana.grafanaBoards).toBe(boards);
    });

    it('overwrites prior grafana fields with new payload values', () => {
      let state = telemetryReducer(undefined, { type: 'init' } as any);
      state = telemetryReducer(state, updateGrafanaConfig({ grafanaURL: 'http://1' }));
      state = telemetryReducer(state, updateGrafanaConfig({ grafanaURL: 'http://2' }));
      expect(state.grafana.grafanaURL).toBe('http://2');
    });
  });

  describe('updatePrometheusConfig', () => {
    it('merges payload into prometheus and refreshes ts', () => {
      const initial = telemetryReducer(undefined, { type: 'init' } as any);
      vi.setSystemTime(new Date('2024-06-01T00:02:00Z'));
      const state = telemetryReducer(
        initial,
        updatePrometheusConfig({ prometheusURL: 'http://prom' }),
      );
      expect(state.prometheus.prometheusURL).toBe('http://prom');
      expect(state.prometheus.selectedPrometheusBoardsConfigs).toEqual([]); // preserved
      expect(state.prometheus.ts).toBeGreaterThan(initial.prometheus.ts);
    });

    it('replaces selectedPrometheusBoardsConfigs when provided', () => {
      const initial = telemetryReducer(undefined, { type: 'init' } as any);
      const configs = [{ id: 'p1' }];
      const state = telemetryReducer(
        initial,
        updatePrometheusConfig({ selectedPrometheusBoardsConfigs: configs }),
      );
      expect(state.prometheus.selectedPrometheusBoardsConfigs).toBe(configs);
    });
  });

  describe('updateStaticPrometheusBoardConfig', () => {
    it('replaces staticPrometheusBoardConfig with the payload', () => {
      const initial = telemetryReducer(undefined, { type: 'init' } as any);
      const payload = { board: 'cpu' };
      const state = telemetryReducer(initial, updateStaticPrometheusBoardConfig(payload));
      expect(state.staticPrometheusBoardConfig).toEqual(payload);
    });
  });

  it('action types have the expected slice prefix', () => {
    expect(updateGrafanaConfig({}).type).toBe('telemetry/updateGrafanaConfig');
    expect(updatePrometheusConfig({}).type).toBe('telemetry/updatePrometheusConfig');
    expect(updateStaticPrometheusBoardConfig({}).type).toBe(
      'telemetry/updateStaticPrometheusBoardConfig',
    );
  });

  it('reducers produce a new state object (immutable)', () => {
    const initial = telemetryReducer(undefined, { type: 'init' } as any);
    const next = telemetryReducer(initial, updateGrafanaConfig({ grafanaURL: 'a' }));
    expect(next).not.toBe(initial);
    expect(initial.grafana.grafanaURL).toBe('');
  });
});
