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

// Durations the performance preferences form accepts: one or more
// `<number><h|m|s>` segments with lowercase units and at least one nonzero
// digit, so "0h0m" is rejected. This is the single definition shared by the
// submit validation and by toDuration below, so a duration the form saves
// is never silently replaced with the default when it is read back.
export const isValidLoadTestDuration = (value: string) =>
  /^(?:(?:\d+(?:\.\d*)?|\.\d+)[hms])+$/.test(value) && /[1-9]/.test(value);

const toDuration = (value: unknown) =>
  typeof value === 'string' && isValidLoadTestDuration(value.trim())
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
