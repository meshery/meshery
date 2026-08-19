import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

/**
 * The E2E job runs Kind, the Meshery server container and Playwright on one
 * 4-vCPU hosted runner. Every extra worker is another Chromium on that runner,
 * and the failures it produces do not look like test failures: a spec blows a
 * whole 180s describe timeout with no assertion error, or the runner itself
 * dies mid-step ("The hosted runner lost communication with the server").
 *
 * `workers` is what Playwright actually reads, so this loads the real config
 * module the CI command loads and asserts the value it resolves to, rather
 * than the expression that produced it - `process.env.CI ? 4 : 4` read as
 * "opt out of parallel tests on CI" for as long as nobody evaluated it.
 */
const loadConfig = async () => {
  vi.resetModules();
  return (await import('../playwright.config.js')).default;
};

/**
 * Loading the real config pulls `@playwright/test` in with it, and every case
 * here re-imports it after `vi.resetModules()`. Under a full `vitest run` that
 * import competes with several hundred other files for the same cores, and the
 * default 5s budget turns these into load tests: they time out at the `it`
 * boundary with no assertion evaluated, which reads as a worker-allocation
 * defect that is not there. What is asserted is deterministic; only how long
 * the import takes is not, so the budget is set by the import.
 */
const IMPORTS_PLAYWRIGHT = { timeout: 60_000 };

describe('Playwright worker allocation', () => {
  beforeEach(() => {
    vi.stubEnv('CI', undefined);
  });

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  it('runs the suite serially on CI', IMPORTS_PLAYWRIGHT, async () => {
    vi.stubEnv('CI', 'true');

    await expect(loadConfig()).resolves.toMatchObject({ workers: 1 });
  });

  it('keeps a local run parallel', IMPORTS_PLAYWRIGHT, async () => {
    await expect(loadConfig()).resolves.toMatchObject({ workers: 4 });
  });

  it(
    'keeps tests within a spec file serial, so workers only add concurrent files',
    IMPORTS_PLAYWRIGHT,
    async () => {
      vi.stubEnv('CI', 'true');

      await expect(loadConfig()).resolves.toMatchObject({ fullyParallel: false });
    },
  );
});
