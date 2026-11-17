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
 *
 * A `timeout` of 0 rejects immediately without polling, matching the legacy
 * hand-rolled loop. `description` optionally names the awaited condition in
 * the timeout message; when omitted the condition's source is used.
 */
export const waitFor = async <T>(
  func: () => T | Promise<T>,
  timeout = 5000,
  description?: string,
): Promise<T> => {
  const condition = description ?? func.toString();

  if (timeout === 0) {
    throw new Error(`Condition did not resolve within 0 ms: ${condition}`);
  }

  let lastResult: T;

  await expect
    .poll(
      async () => {
        lastResult = await func();
        return Boolean(lastResult);
      },
      {
        message: `Condition did not resolve within ${timeout} ms: ${condition}`,
        timeout,
        intervals: [100, 250, 500],
      },
    )
    .toBe(true);

  return lastResult!;
};
