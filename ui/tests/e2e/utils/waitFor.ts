import { expect } from '@playwright/test';

export const waitFor = async <T>(
  func: () => T | Promise<T>,
  timeout = 5000,
): Promise<T> => {
  let lastResult: T;
  await expect
    .poll(
      async () => {
        lastResult = await func();
        return Boolean(lastResult);
      },
      {
        message: `Condition did not resolve within ${timeout}ms`,
        timeout,
        intervals: [100, 250, 500],
      },
    )
    .toBe(true);

  return lastResult!;
};
