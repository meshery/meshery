import { expect, test } from '@playwright/test';
import { DYNAMIC_TIMEOUTS } from './delays';
const { ENV } = require('./env');

async function navigateToPage(page, path, waitForSelector = '') {
  await page.goto(`${ENV.MESHERY_SERVER_URL}/${path}`);
  if (waitForSelector) {
    await page.waitForSelector(waitForSelector, { state: 'visible' });
  }
}

async function clickButton(page, selector) {
  await page.click(selector, { force: true });
}

async function fillInput(page, selector, value) {
  await page.fill(selector, value);
}

async function doesProfileExist(page, profileName) {
  await navigateToPage(page, 'performance/profiles');
  return page
    .locator(`text=${profileName}`)
    .count()
    .then((count) => count > 0);
}

// This function first checks if a profile with the given name already exists using the `doesProfileExist` function.
// If not, it sets a longer timeout for the profile creation process and proceeds to add a new performance profile.
async function createPerformanceProfile(page, profileName) {
  if (!(await doesProfileExist(page, profileName))) {
    page.setDefaultTimeout(90 * 1000); // 90 seconds
    await clickButton(page, 'button[aria-label="Add Performance Profile"]');
    await fillInput(page, '#profileName', profileName);

    await clickButton(page, 'label[for="meshName"] + div');
    await clickButton(page, 'li[data-value="istio"]');

    await fillInput(page, '#url', 'https://layer5.io/');
    await fillInput(page, '#c', '5');
    await fillInput(page, '#qps', '5');
    await page.fill('#t', '15s');
    await page.waitForSelector('button:has-text("Run Test")', { state: 'visible' });

    await page.getByRole('button', { name: 'Run Test', exact: true }).click();
    await page.waitForTimeout(DYNAMIC_TIMEOUTS.PERFORMANCE_RUN());
  } else {
    console.log(`Profile "${profileName}" already exists.`);
  }
}

test.describe('Service Mesh Performance Management Tests', () => {
  const profileName = 'Sample-test';

  test.beforeEach(async ({ page }) => {
    await navigateToPage(page, '', 'text=Performance');
  });

  // Test to create and run a performance test through a new profile
  test('Run a performance test through profile', async ({ page }) => {
    await clickButton(page, 'text=Performance');
    await page.waitForSelector('text=Profiles');
    await clickButton(page, 'text=Profiles');
    await createPerformanceProfile(page, profileName); // Creates a performance profile if it doesn't exist and runs a test
    await expect(page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
  });

  // Test to view the detailed results of a performance profile and check graph visibility
  test('View detailed result of a performance profile (Graph Visualiser)', async ({ page }) => {
    await clickButton(page, 'text=Performance');
    await expect(page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/performance`);
    await clickButton(page, 'button:has-text("Manage Profiles")');
    await expect(page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);

    await page.waitForSelector(`text=${profileName}`, { state: 'visible' }); // Waits until the profile name is visible
    await clickButton(page, 'button:has-text("View Results")'); // Opens the results for the performance profile
    await clickButton(page, 'button[aria-label="more"]');

    const graphVisible = await page.isVisible('div[class*="bb"]'); // Check if the graph is visible
    expect(graphVisible);
  });

  // Test to run a performance test directly from the performance page
  test('Run a performance test', async ({ page }) => {
    await clickButton(page, 'text=Performance');
    await expect(page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/performance`);

    await clickButton(page, 'button:has-text("Run Test")'); // Initiates a performance test run
    await fillInput(page, '#profileName', profileName); // Fills out the performance profile form

    await clickButton(page, 'label[for="meshName"] + div');
    await clickButton(page, 'li[data-value="None"]');

    await fillInput(page, '#url', 'https://layer5.io/');
    await fillInput(page, '#c', '5');
    await fillInput(page, '#qps', '5');

    await page.waitForSelector('button:has-text("Run Test")', { state: 'visible' });
    await page.getByRole('button', { name: 'Run Test', exact: true }).click();
  });

  // Test to view results from an existing performance profile
  test('View Results from a performance profile', async ({ page }) => {
    await clickButton(page, 'text=Performance');
    await page.waitForSelector('text=Profiles');
    await clickButton(page, 'text=Profiles');

    await expect(page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
    await clickButton(page, 'button:has-text("View Results")'); // Views the results for the selected profile

    const resultsVisible = await page
      .locator('tr[id^="MUIDataTableBodyRow-"][id$="-0"]')
      .isVisible(); // Verifies if the results are visible
    expect(resultsVisible);
  });

  // Test to view and edit the configuration of an existing performance profile
  test('View/Edit the configuration of a performance profile', async ({ page }) => {
    page.setDefaultTimeout(90 * 1000); // Set a longer timeout for configuration updates
    await clickButton(page, 'text=Performance');
    await page.waitForSelector('text=Profiles');
    await clickButton(page, 'text=Profiles');
    await expect(page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);

    await clickButton(page, `div:has-text("${profileName}")`); // Opens the profile configuration for editing

    await clickButton(
      page,
      'button.MuiIconButton-root:has(svg path[d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34a.9959.9959 0 00-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"])', // Clicks on the edit button
    );
    await clickButton(page, 'label[for="meshName"] + div');
    await clickButton(page, 'li[data-value="None"]');

    await fillInput(page, '#c', '6'); // Updates the configuration values
    await page.waitForSelector('button:has-text("Run Test")', { state: 'visible' }); // Runs the test with the updated configuration
    await page.getByRole('button', { name: 'Run Test', exact: true }).click();
    await page.waitForTimeout(DYNAMIC_TIMEOUTS.PERFORMANCE_RUN()); // Waits for the custom timeout period to ensure the test completes
  });
});
