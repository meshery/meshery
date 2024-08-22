import { expect } from '@playwright/test';

export class MeshMapTutorial {
  /**
   * @param {import("@playwright/test").Page} page
   */
  constructor(page) {
    this.page = page;
  }

  async goto() {
    await this.page.goto('/extension/meshmap');
  }

  async skipTutorial() {
    await expect(this.page.getByText('Tutorial')).toBeVisible();

    await this.page.getByLabel('Do not show again').check();

    await this.page.locator('#tutorial-dialog-title svg').nth(2).click();
  }

  async turnOnTutorialSettings() {
    await this.page.getByLabel('CANVAS_CONTROLS').nth(1).hover();

    await this.page.getByRole('menuitem', { name: 'Options' }).click();

    await expect(this.page.getByText('Show Tutorial')).toBeVisible();

    await this.page
      .locator('li')
      .filter({ hasText: 'Show Tutorial' })
      .getByRole('checkbox')
      .click();

    await this.page.getByRole('button').first().click();
  }
}
