// fortio is the only supported load generator. wrk2 was removed; any stored
// "wrk2" preference is coerced to the fortio default by toLoadGenerator.
const validLoadGenerators = new Set(['fortio']);

export const DEFAULT_LOAD_TEST_PREFS = {
  c: 0,
  qps: 0,
  t: '30s',
  gen: 'fortio',
};

type LoadTestPrefsInput =
  | {
      c?: unknown;
      qps?: unknown;
      t?: unknown;
      gen?: unknown;
    }
  | null
  | undefined;

const toNonNegativeNumber = (value: unknown, fallback: number) => {
  const parsedValue =
    typeof value === 'number'
      ? value
      : typeof value === 'string' && value.trim() !== ''
        ? Number(value)
        : NaN;

  return Number.isFinite(parsedValue) && parsedValue >= 0 ? parsedValue : fallback;
};

// Share the duration syntax between read-back and submission. Stored zero
// durations are preserved on read-back, but the form requires a positive value.
const loadTestDurationPattern = /^(?:(?:\d+(?:\.\d*)?|\.\d+)(?:ms|h|m|s))+$/;

export const isValidLoadTestDuration = (value: string) =>
  loadTestDurationPattern.test(value) && /[1-9]/.test(value);

export const isZeroLoadTestDuration = (value: string) =>
  loadTestDurationPattern.test(value) && !/[1-9]/.test(value);

const toDuration = (value: unknown) =>
  typeof value === 'string' && loadTestDurationPattern.test(value.trim())
    ? value.trim()
    : DEFAULT_LOAD_TEST_PREFS.t;

const toLoadGenerator = (value: unknown) =>
  typeof value === 'string' && validLoadGenerators.has(value.trim())
    ? value.trim()
    : DEFAULT_LOAD_TEST_PREFS.gen;

export const normalizeLoadTestPrefs = (loadTestPrefs: LoadTestPrefsInput = {}) => ({
  c: toNonNegativeNumber(loadTestPrefs?.c, DEFAULT_LOAD_TEST_PREFS.c),
  qps: toNonNegativeNumber(loadTestPrefs?.qps, DEFAULT_LOAD_TEST_PREFS.qps),
  t: toDuration(loadTestPrefs?.t),
  gen: toLoadGenerator(loadTestPrefs?.gen),
});
