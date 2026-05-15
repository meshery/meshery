const validLoadGenerators = new Set(['fortio', 'wrk2', 'nighthawk']);

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

const toDuration = (value: unknown) =>
  typeof value === 'string' && /^\d+[hms]$/i.test(value.trim())
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
