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
        gen: 'nighthawk',
      }),
    ).toEqual({
      c: 12,
      qps: 25,
      t: '5m',
      gen: 'nighthawk',
    });
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
