import { expect } from '@playwright/test';

/**
 * Polls `func` until it returns a truthy value or `timeout` expires, then
 * resolves with that value.
 *
 * Uses Playwright's native `expect.poll` so failures surface inside the step
 * log and trace viewer rather than as generic rejection strings.
 *
 * **Retry-on-throw:** if `func` throws, the error is caught and polling
 * continues until the timeout expires, rather than rejecting immediately.
 * Callers that need fail-fast error propagation should handle errors inside
 * `func` or not use this helper.
 */
export const waitFor = async <T>(func: () => T | Promise<T>, timeout = 5000): Promise<T> => {
  let lastResult: T;

  await expect
    .poll(
      async () => {
        lastResult = await func();
        return Boolean(lastResult);
      },
      {
        message: `Condition did not resolve within ${timeout} ms`,
        timeout,
        intervals: [100, 200, 500],
      },
    )
    .toBe(true);

  return lastResult!;
};
