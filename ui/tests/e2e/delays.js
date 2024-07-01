export const DYNAMIC_TIMEOUTS = {
  DESING_IMPORT: (numberOfComponents) => 30 * 1000 + 1000 * numberOfComponents,
  DESING_LOAD: (numberOfComponents) => 30 * 1000 + 500 * numberOfComponents,
};

export const SLOW_TEST_DURATION = 60 * 1000;
export const BASE_TIMEOUT = 30 * 1000;
