import { expect } from '@playwright/test';
import { ENV } from '../env';

export class PerformancePage {
  /**
   * @param {import('@playwright/test').Page} page
   */
  constructor(page) {
    this.page = page;
  }

  async navigate() {
    await this.page.goto(`${ENV.MESHERY_SERVER_URL}/performance`);
    await expect(this.page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/performance`);
  }

  async goToProfiles() {
    await this.page.goto(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
    await expect(this.page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/performance/profiles`);
  }

  async createPerformanceProfile(profileName) {
    if (!(await this.doesProfileExist(profileName))) {
      this.page.setDefaultTimeout(90 * 1000); // 90 seconds
      await this.page.getByRole('button', { name: 'Add Performance Profile', exact: true }).click();
      await this.fillInput(this.page.getByLabel('Profile Name'), profileName);
      await this.page.locator('[aria-labelledby="meshName-label meshName"]').click();
      await this.page.locator('[data-value="istio"]').click();
      await this.fillInput(this.page.getByRole('textbox', { name: 'url' }), 'https://layer5.io/');
      await this.fillInput(this.page.getByRole('spinbutton', { name: 'Concurrent requests' }), '5');
      await this.fillInput(this.page.getByRole('spinbutton', { name: 'Queries per second' }), '5');
      await this.fillInput(this.page.getByRole('textbox', { name: 'Duration' }), '15s');
      await expect(this.page.getByRole('button', { name: 'Run Test', exact: true })).toBeVisible();
      await this.page.getByRole('button', { name: 'Run Test', exact: true }).click();
      // Check for notification visibility
      const notification = await this.page.locator('text=Initiating load test . . .').first();
      await expect(notification).toBeVisible();
    } else {
      console.log(`Profile “${profileName}” already exists.`);
    }
  }

  async doesProfileExist(profileName) {
    await this.goToProfiles();
    return this.page
      .locator(`text=${profileName}`)
      .count()
      .then((count) => count > 0);
  }

  async viewPerformanceProfileResult(profileName) {
    await this.page.waitForSelector(`text=${profileName}`, { state: 'visible' });
    await this.page.click(`button:has-text("View Results")`);
    await this.page.waitForSelector('button[aria-label="more"]', { state: 'visible' });
    await this.page.click('button[aria-label="more"]');

    await this.page.evaluate(() => {
      const sentinelStartDiv = document.querySelector('div[data-testid="sentinelStart"]');
      if (sentinelStartDiv) {
        sentinelStartDiv.setAttribute('data-testid', 'sentinel-graph');
      }
    });

    const graphVisible = await this.page.getByTestId('sentinel-graph').isVisible();
    expect(graphVisible);
  }

  async runPerformanceTest(profileName) {
    await this.page.goto(`${ENV.MESHERY_SERVER_URL}/performance`);
    await this.page.getByRole('button', { name: 'Run Test', exact: true }).click();
    await this.fillInput(this.page.getByLabel('Profile Name'), profileName);
    await this.page.locator('[aria-labelledby="meshName-label meshName"]').click();
    await this.page.locator('[data-value="istio"]').click();
    await this.fillInput(this.page.getByRole('textbox', { name: 'url' }), 'https://layer5.io/');
    await this.fillInput(this.page.getByRole('spinbutton', { name: 'Concurrent requests' }), '5');
    await this.fillInput(this.page.getByRole('spinbutton', { name: 'Queries per second' }), '5');
    await this.fillInput(this.page.getByRole('textbox', { name: 'Duration' }), '15s');
    await expect(this.page.getByRole('button', { name: 'Run Test', exact: true })).toBeVisible();
    await this.page.getByRole('button', { name: 'Run Test', exact: true }).click();
    // Check for notification visibility
    const notification = await this.page.locator('text=Initiating load test . . .').first();
    await expect(notification).toBeVisible();
  }

  async viewPerformanceProfileConfiguration(profileName) {
    await this.page.click(`div:has-text("${profileName}")`);
    await this.page.evaluate(() => {
      const editIcon = document.querySelector(
        'svg path[d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34a.9959.9959 0 00-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"]',
      );
      if (editIcon) {
        editIcon.closest('svg').setAttribute('data-testid', 'edit-icon');
      }
    });
    await this.page.getByTestId('edit-icon').click();
  }

  async fillInput(locator, value) {
    await locator.fill(value);
  }
}
