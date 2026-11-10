import { describe, expect, it } from 'vitest';
import { DEFAULT_LOAD_TEST_PREFS, normalizeLoadTestPrefs } from '../load-test-prefs';

describe('normalizeLoadTestPrefs', () => {
  it('returns defaults when preferences are missing', () => {
    expect(normalizeLoadTestPrefs()).toEqual(DEFAULT_LOAD_TEST_PREFS);
    expect(normalizeLoadTestPrefs(null)).toEqual(DEFAULT_LOAD_TEST_PREFS);
  });

  it('keeps valid values and coerces numeric strings', () => {
    expect(
      normalizeLoadTestPrefs({
        c: '12',
        qps: 25,
        t: '5m',
        gen: 'fortio',
      }),
    ).toEqual({
      c: 12,
      qps: 25,
      t: '5m',
      gen: 'fortio',
    });
  });

  it('coerces the removed wrk2 generator to fortio', () => {
    expect(
      normalizeLoadTestPrefs({
        c: 1,
        qps: 1,
        t: '5m',
        gen: 'wrk2',
      }),
    ).toEqual({
      c: 1,
      qps: 1,
      t: '5m',
      gen: 'fortio',
    });
    expect(normalizeLoadTestPrefs({ gen: ' wrk2 ' }).gen).toBe('fortio');
  });

  it('falls back when duration or generator values are invalid', () => {
    expect(
      normalizeLoadTestPrefs({
        c: -1,
        qps: 'bad-value',
        t: 0,
        gen: 0,
      }),
    ).toEqual(DEFAULT_LOAD_TEST_PREFS);
  });

  it('keeps supported compound durations when reading stored preferences', () => {
    for (const duration of ['1h30m', '0h45m', '0h30m', '1h.5m', '1.5h']) {
      expect(normalizeLoadTestPrefs({ t: duration }).t).toBe(duration);
    }
  });

  it('falls back for durations the form no longer accepts', () => {
    for (const duration of ['24H', '1H30m', '0h0m', '1d', '1d1h']) {
      expect(normalizeLoadTestPrefs({ t: duration }).t).toBe(DEFAULT_LOAD_TEST_PREFS.t);
    }
  });

  it('falls back for durations missing a unit', () => {
    expect(normalizeLoadTestPrefs({ t: '30' })).toEqual({
      ...DEFAULT_LOAD_TEST_PREFS,
      t: DEFAULT_LOAD_TEST_PREFS.t,
    });
  });

  it('trims whitespace on generator values', () => {
    expect(normalizeLoadTestPrefs({ gen: ' fortio ' })).toEqual({
      ...DEFAULT_LOAD_TEST_PREFS,
      gen: 'fortio',
    });
  });
});
