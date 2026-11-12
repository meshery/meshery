// @ts-check
const { defineConfig, devices } = require('@playwright/test');
const { BASE_TIMEOUT } = require('./tests/e2e/delays');
// Requiring this also loads ui/.env off CI - see tests/e2e/env.js.
const { ENV } = require('./tests/e2e/env');

/**
 * @see https://playwright.dev/docs/test-configuration
 */
module.exports = defineConfig({
  testDir: './tests/e2e',
  /* Run tests in files in parallel */
  fullyParallel: false,
  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: !!process.env.CI,
  timeout: BASE_TIMEOUT,
  expect: {
    /* The number of milliseconds the test runner will wait for the expect matchers to pass. */
    timeout: BASE_TIMEOUT,
  },
  retries: 0,
  /* Opt out of parallel tests on CI.
   *
   * This said "opt out" while running four workers - `process.env.CI ? 4 : 4`,
   * a ternary whose branches are the same number - and the CI job is the one
   * place that cannot afford them. It puts the whole stack on a single
   * 4-vCPU/16-GB hosted runner: a Kind cluster, the Meshery server container,
   * and Playwright. Four workers add four Chromium instances on top of that,
   * and the suite then fails in the shapes a starved runner produces rather
   * than the shape a broken test produces. Both are on record for 2026-08-17:
   * a test burning a whole 180s describe timeout with no assertion error
   * ("Metrics (Prometheus) page loads", run 32024269277), and the runner
   * itself dying mid-step - "The hosted runner lost communication with the
   * server ... starves it for CPU/Memory" - after 49 minutes with no logs and
   * no artifacts (run 32084526223). A different, unrelated test failed in each
   * red run that day; that spread is the signature of load and of specs racing
   * each other, not of one broken assertion.
   *
   * `fullyParallel: false` already serializes the tests within a file, so a
   * second worker buys only concurrent spec *files* - run against one Meshery
   * server, one Local-provider session and one Kind cluster, which is the
   * state those specs mutate. Serial on CI costs wall clock and buys a verdict
   * that means something. Local runs keep 4.
   */
  workers: process.env.CI ? 1 : 4,
  /* Reporter to use. See https://playwright.dev/docs/test-reporters */
  // reporter: process.env.CI ? "./tests/e2e/custom-playwright-reporter.js" : "list",
  reporter: [
    ['allure-playwright'],
    ['./tests/e2e/custom-playwright-reporter.js'],
    ['html', { outputFolder: 'playwright-report', open: 'never' }],
  ],
  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Base URL to use in actions like `await page.goto('/')`. */
    baseURL: process.env.MESHERY_SERVER_URL || 'http://localhost:9081',
    video: {
      mode: 'retain-on-failure',
    },
    /* Keep a screenshot of the final failed state. Together with trace + video
     * (all retain-on-failure), the allure-playwright reporter attaches these to
     * each failed result, so a failure is debuggable in the Allure report
     * (e.g. the Connection Lifecycle report) instead of showing only the error. */
    screenshot: 'only-on-failure',
    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'retain-on-failure',
    provider: process.env.MESHERY_PROVIDER || 'Local',
  },

  /* Configure projects for major browsers */
  projects: [
    // setup
    {
      name: 'local-setup',
      testMatch: 'tests/e2e/local.setup.js',
    },
    {
      name: 'remote-setup',
      testMatch: 'tests/e2e/remote.setup.js',
    },
    {
      name: 'chromium-meshery-provider',
      use: {
        ...devices['Desktop Chrome'],
        provider: 'Meshery',
        // Use prepared auth state.
        storageState: ENV.AUTHFILEMESHERYPROVIDER,
      },
      dependencies: ['remote-setup'],
    },
    {
      name: 'chromium-local-provider',
      use: {
        ...devices['Desktop Chrome'],
        provider: 'Local',
        // Use prepared auth state.
        storageState: ENV.AUTHFILELOCALPROVIDER,
      },
      dependencies: ['local-setup'],
    },
  ],
});
