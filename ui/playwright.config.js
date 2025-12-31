// @ts-check
const { defineConfig, devices } = require('@playwright/test');
const { BASE_TIMEOUT } = require('./tests/e2e/delays');
const { ENV } = require('./tests/e2e/env');
const dotenv = require('dotenv');
const path = require('path');

/**
 * Read environment variables from file.
 * https://github.com/motdotla/dotenv
 */
if (!process.env.CI) {
  dotenv.config({ path: path.resolve(__dirname, '.env') });
}

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
  /* Opt out of parallel tests on CI. */
  workers: process.env.CI ? 4 : 4,
  /* Reporter to use. See https://playwright.dev/docs/test-reporters */
  // reporter: process.env.CI ? "./tests/e2e/custom-playwright-reporter.js" : "list",
  reporter: [['./tests/e2e/custom-playwright-reporter.js']],
  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Base URL to use in actions like `await page.goto('/')`. */
    baseURL: process.env.MESHERY_SERVER_URL || 'http://localhost:9081',
    video: {
      mode: 'retain-on-failure',
    },
    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'retain-on-failure',
    provider: process.env.MESHERY_PROVIDER || 'None',
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
        provider: 'None',
        // Use prepared auth state.
        storageState: ENV.AUTHFILELOCALPROVIDER,
      },
      dependencies: ['local-setup'],
    },
  ],
});
