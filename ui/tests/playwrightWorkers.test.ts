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

describe('Playwright worker allocation', () => {
  beforeEach(() => {
    vi.stubEnv('CI', undefined);
  });

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  it('runs the suite serially on CI', async () => {
    vi.stubEnv('CI', 'true');

    await expect(loadConfig()).resolves.toMatchObject({ workers: 1 });
  });

  it('keeps a local run parallel', async () => {
    await expect(loadConfig()).resolves.toMatchObject({ workers: 4 });
  });

  it('keeps tests within a spec file serial, so workers only add concurrent files', async () => {
    vi.stubEnv('CI', 'true');

    await expect(loadConfig()).resolves.toMatchObject({ fullyParallel: false });
  });
});
