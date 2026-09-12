import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import prefTestReducer, {
  updateLoadTestPref,
  updateLoadTest,
  updateResultsSelection,
  clearResultsSelection,
} from '../prefTest';

describe('prefTest slice', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2024-06-01T00:00:00Z'));
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('returns the initial state for an unknown action', () => {
    const state = prefTestReducer(undefined, { type: 'init' } as any);
    expect(state.loadTest).toEqual({
      testName: '',
      meshName: '',
      url: '',
      qps: 0,
      c: 0,
      t: '30s',
      result: {},
    });
    expect(state.loadTestPref.gen).toBe('fortio');
    expect(state.loadTestPref.t).toBe('30s');
    expect(state.loadTestPref.qps).toBe(0);
    expect(state.loadTestPref.c).toBe(0);
    expect(typeof state.loadTestPref.ts).toBe('number');
    expect(state.results_selection).toEqual({});
    expect(state.results).toEqual({ startKey: '', results: [] });
  });

  describe('updateLoadTestPref', () => {
    it('replaces loadTestPref with the payload', () => {
      const initial = prefTestReducer(undefined, { type: 'init' } as any);
      const ts = Date.now();
      const state = prefTestReducer(
        initial,
        updateLoadTestPref({
          loadTestPref: { qps: 50, t: '1m', c: 10, gen: 'fortio', ts },
        }),
      );
      expect(state.loadTestPref).toEqual({
        qps: 50,
        t: '1m',
        c: 10,
        gen: 'fortio',
        ts,
      });
    });

    it('converts a Date instance for ts to a timestamp', () => {
      const initial = prefTestReducer(undefined, { type: 'init' } as any);
      const dateObj = new Date('2024-06-01T00:01:00Z');
      const state = prefTestReducer(
        initial,
        updateLoadTestPref({
          loadTestPref: { qps: 1, t: '15s', c: 1, gen: 'fortio', ts: dateObj },
        }),
      );
      expect(state.loadTestPref.ts).toBe(dateObj.getTime());
      expect(typeof state.loadTestPref.ts).toBe('number');
    });

    it('keeps ts as a numeric value when given a number', () => {
      const initial = prefTestReducer(undefined, { type: 'init' } as any);
      const ts = 1700000000000;
      const state = prefTestReducer(
        initial,
        updateLoadTestPref({
          loadTestPref: { qps: 5, t: '30s', c: 1, gen: 'fortio', ts },
        }),
      );
      expect(state.loadTestPref.ts).toBe(ts);
    });
  });

  describe('updateLoadTest', () => {
    it('replaces loadTest with the payload', () => {
      const initial = prefTestReducer(undefined, { type: 'init' } as any);
      const loadTest = {
        testName: 'My Test',
        meshName: 'istio',
        url: 'https://example.com',
        qps: 100,
        c: 5,
        t: '1m',
        result: { foo: 'bar' },
      };
      const state = prefTestReducer(initial, updateLoadTest({ loadTest }));
      expect(state.loadTest).toEqual(loadTest);
    });
  });

  describe('updateResultsSelection', () => {
    it('sets selection for a page when results are non-empty', () => {
      const initial = prefTestReducer(undefined, { type: 'init' } as any);
      const state = prefTestReducer(
        initial,
        updateResultsSelection({ page: 0, results: { a: 'b' } }),
      );
      expect(state.results_selection).toEqual({ 0: { a: 'b' } });
    });

    it('removes selection for a page when results object is empty', () => {
      let state = prefTestReducer(undefined, { type: 'init' } as any);
      state = prefTestReducer(state, updateResultsSelection({ page: 0, results: { a: 'b' } }));
      state = prefTestReducer(state, updateResultsSelection({ page: 0, results: {} }));
      expect(state.results_selection).toEqual({});
    });

    it('overwrites existing selection on the same page', () => {
      let state = prefTestReducer(undefined, { type: 'init' } as any);
      state = prefTestReducer(state, updateResultsSelection({ page: 1, results: { x: 1 } }));
      state = prefTestReducer(state, updateResultsSelection({ page: 1, results: { y: 2 } }));
      expect(state.results_selection).toEqual({ 1: { y: 2 } });
    });

    it('supports selections across multiple pages', () => {
      let state = prefTestReducer(undefined, { type: 'init' } as any);
      state = prefTestReducer(state, updateResultsSelection({ page: 1, results: { a: 1 } }));
      state = prefTestReducer(state, updateResultsSelection({ page: 2, results: { b: 2 } }));
      expect(state.results_selection).toEqual({ 1: { a: 1 }, 2: { b: 2 } });
    });
  });

  describe('clearResultsSelection', () => {
    it('clears all entries in results_selection', () => {
      let state = prefTestReducer(undefined, { type: 'init' } as any);
      state = prefTestReducer(state, updateResultsSelection({ page: 0, results: { a: 1 } }));
      state = prefTestReducer(state, updateResultsSelection({ page: 1, results: { b: 2 } }));
      state = prefTestReducer(state, clearResultsSelection());
      expect(state.results_selection).toEqual({});
    });
  });

  it('actions are immutable', () => {
    const initial = prefTestReducer(undefined, { type: 'init' } as any);
    const next = prefTestReducer(
      initial,
      updateLoadTest({
        loadTest: { testName: 'x', meshName: '', url: '', qps: 0, c: 0, t: '30s', result: {} },
      }),
    );
    expect(next).not.toBe(initial);
    expect(initial.loadTest.testName).toBe('');
  });

  it('action types have the expected slice prefix', () => {
    expect(updateLoadTest({ loadTest: {} as any }).type).toBe('prefTest/updateLoadTest');
    expect(updateLoadTestPref({ loadTestPref: {} as any }).type).toBe(
      'prefTest/updateLoadTestPref',
    );
    expect(updateResultsSelection({ page: 0, results: {} }).type).toBe(
      'prefTest/updateResultsSelection',
    );
    expect(clearResultsSelection().type).toBe('prefTest/clearResultsSelection');
  });
});
