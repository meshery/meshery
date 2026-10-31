import { expect } from '@playwright/test';

export const waitFor = async <T>(func: () => T | Promise<T>, timeout = 5000): Promise<T> => {
  if (timeout <= 0) {
    throw new Error(`Condition did not resolve within ${timeout}ms`);
  }

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
        intervals: [200],
      },
    )
    .toBe(true);

  return lastResult!;
};
