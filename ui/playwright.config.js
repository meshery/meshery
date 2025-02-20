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
  /* Retry on CI only */
  retries: 0,
  /* Opt out of parallel tests on CI. */
  workers: process.env.CI ? 1 : 1,
  /* Reporter to use. See https://playwright.dev/docs/test-reporters */
  // reporter: process.env.CI ? "./tests/e2e/custom-playwright-reporter.js" : "list",
  reporter: [['./tests/e2e/custom-playwright-reporter.js']],
  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Base URL to use in actions like `await page.goto('/')`. */
    // baseURL: 'http://127.0.0.1:3000',
    video: {
      mode: 'retain-on-failure',
    },
    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'retain-on-first-failure',
  },

  /* Configure projects for major browsers */
  projects: [
    // setup
    {
      name: 'setup',
      testMatch: 'tests/e2e/*.setup.js',
    },
    {
      name: 'chromium-meshery-provider',
      use: {
        ...devices['Desktop Chrome'],
        provider: 'Meshery',
        // Use prepared auth state.
        storageState: ENV.AUTHFILEMESHERYPROVIDER,
      },
      dependencies: ['setup'],
    },
    {
      name: 'chromium-local-provider',
      use: {
        ...devices['Desktop Chrome'],
        provider: 'None',
        // Use prepared auth state.
        storageState: ENV.AUTHFILELOCALPROVIDER,
      },
      dependencies: ['setup'],
    },
    {
      name: 'firefox',
      use: {
        ...devices['Desktop Firefox'],
        // Use prepared auth state.
        storageState: ENV.AUTHFILEMESHERYPROVIDER,
      },
      dependencies: ['setup'],
    } /* Test against mobile viewports. */,
    // {
    //   name: 'Mobile Chrome',
    //   use: { ...devices['Pixel 5'] },
    // },
    // {
    //   name: 'Mobile Safari',
    //   use: { ...devices['iPhone 12'] },
    // },

    /* Test against branded browsers. */
    // {
    //   name: 'Microsoft Edge',
    //   use: { ...devices['Desktop Edge'], channel: 'msedge' },
    // },
    // {
    //   name: 'Google Chrome',
    //   use: { ...devices['Desktop Chrome'], channel: 'chrome' },
    // },
  ],

  /* Run your local dev server before starting the tests */
  // webServer: {
  //   command: 'npm run start',
  //   url: 'http://127.0.0.1:3000',
  //   reuseExistingServer: !process.env.CI,
  // },
});
