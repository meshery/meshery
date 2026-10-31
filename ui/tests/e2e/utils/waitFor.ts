import { expect } from '@playwright/test';

/**
 * Polls `func` using Playwright's native `expect.poll` until it returns a
 * truthy value, then resolves with that value.
 *
 * Replaces a legacy hand-rolled `setInterval`/`setTimeout` loop so that
 * poll failures surface inside Playwright's step log and trace viewer instead
 * of throwing generic rejection strings.
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
