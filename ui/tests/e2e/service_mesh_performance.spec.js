import { expect, test } from '@playwright/test';
import fs from 'fs';
const { ENV } = require('./env');

let adapterStatus = {};

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

async function readAdapterStatus() {
  try {
    return JSON.parse(fs.readFileSync('./adapterStatus.json', 'utf8'));
  } catch (error) {
    console.error('Error reading adapter status file:', error);
    return {};
  }
}

test.describe('Service Mesh Performance Management Tests', () => {
  test.beforeEach(async ({ page }) => {
    adapterStatus = await readAdapterStatus();
    await navigateToPage(page, '', 'text=Performance');
  });

  test('Run a performance test through profile', async ({ page }) => {
    await clickButton(page, 'text=Performance');
    await clickButton(page, 'text=Profiles');

    await expect(page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
    await clickButton(page, 'button[aria-label="Add Performance Profile"]');
    await fillInput(page, '#profileName', 'Sample-test');

    await clickButton(page, 'label[for="meshName"] + div');

    if (adapterStatus.Istio) {
      await clickButton(page, 'li[data-value="istio"]');
    } else {
      console.error('Istio adapter is not accessible.');
      test.fail();
    }

    await fillInput(page, '#url', 'https://layer5.io/');
    await fillInput(page, '#c', '5');
    await fillInput(page, '#qps', '5');

    await page.waitForSelector('button:has-text("Run Test")', { state: 'visible' });
    await clickButton(page, 'button:has-text("Run Test")', { force: true });
  });

  test('View detailed result of a performance profile (Graph Visualiser)', async ({ page }) => {
    await clickButton(page, 'text=Performance');
    await expect(page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/performance`);
    await clickButton(page, 'button:has-text("Manage Profiles")');
    await expect(page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);

    await page.waitForSelector('text=Sample-test', { state: 'visible' });
    await clickButton(page, 'button:has-text("View Results")');
    await clickButton(page, 'button[aria-label="more"]');

    const graphVisible = await page.isVisible('div[class*="bb"]');
    expect(graphVisible).toBeTruthy();
  });

  test('Run a performance test', async ({ page }) => {
    await clickButton(page, 'text=Performance');
    await expect(page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/performance`);

    await clickButton(page, 'button:has-text("Run Test")');
    await fillInput(page, '#profileName', 'Sample-test');

    await clickButton(page, 'label[for="meshName"] + div');
    await clickButton(page, 'li[data-value="None"]');

    await fillInput(page, '#url', 'https://layer5.io/');
    await fillInput(page, '#c', '5');
    await fillInput(page, '#qps', '5');

    await page.waitForSelector('button:has-text("Run Test")', { state: 'visible' });
    await clickButton(page, 'button:has-text("Run Test")', { force: true });
  });

  test('View Results from a performance profile', async ({ page }) => {
    await clickButton(page, 'text=Performance');
    await clickButton(page, 'text=Profiles');

    await expect(page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
    await clickButton(page, 'button:has-text("View Results")');

    const resultsVisible = await page
      .locator('tr[id^="MUIDataTableBodyRow-"][id$="-0"]')
      .isVisible();
    expect(resultsVisible);
  });

  test('View/Edit the configuration of a performance profile', async ({ page }) => {
    await clickButton(page, 'text=Performance');
    await clickButton(page, 'text=Profiles');
    await expect(page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);

    await clickButton(page, 'div:has-text("Sample-test")');

    await clickButton(
      page,
      'button.MuiIconButton-roothas(svg path[d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34a.9959.9959 0 00-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"])',
    );
    await clickButton(page, 'label[for="meshName"] + div');
    await clickButton(page, 'li[data-value="None"]');

    await fillInput(page, '#c', '6');
    await page.waitForSelector('button:has-text("Run Test")', { state: 'visible' });
    await clickButton(page, 'button:has-text("Run Test")', { force: true });
  });
});
