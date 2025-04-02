import { DYNAMIC_TIMEOUTS } from '../delays';
import { ENV } from '../env';

const base = require('@playwright/test');

// import types
/** @typedef {import("@playwright/test")} PlaywrightTest
 * @typedef {import("@playwright/test").Page} PlaywrightPage
 * @typedef {import("@playwright/test").ElementHandle} PlaywrightElementHandle
 * @typedef {import("@playwright/test").JSHandle} JSHandle
 */

/**
 * @typedef {import("@playwright/test").JSHandle<CytoscapeCore>} CytoscapeHandle
 */

// Extend base test by providing "todoPage" and "settingsPage".
// This new "test" can be used in multiple test files, and each of them will get the fixtures.

// Follows Page Object pattern
export class DesignsPage {
  /**
   * @param {import("@playwright/test").Page} page
   */
  constructor(page) {
    this.page = page;
  }

  async loadDefault() {
    await this.page.goto(ENV.MESHMAP_DESIGN_URL());
    await expect(this.page.getByTestId('DesignerContainer')).toBeVisible();
  }

  async loadDesign(id) {
    await this.page.goto(ENV.MESHMAP_DESIGN_URL(id));
    await expect(this.page.getByTestId('DesignerContainer')).toBeVisible();
  }

  // NOTE: Refactor the import form to use accessible labels
  // and proper roles selectors to make it easier to test
  async importDesignFromURI(type, uri, name, expectedNumberOfComponents = 1) {
    const timeout = DYNAMIC_TIMEOUTS.DESING_IMPORT(expectedNumberOfComponents);
    const page = this.page;
    await page.getByRole('tab', { name: 'Designs' }).click();
    await page.getByRole('button', { name: 'Import Design' }).click();
    await page.getByLabel('Meshery Design').click();
    await page.locator('li').filter({ hasText: type }).click();
    await page.getByLabel('URL *').fill(uri);
    await page.getByLabel('Design file name').fill(name);
    await page.getByRole('button', { name: 'Import' }).click();
    const designsTable = page.getByRole('table', { name: 'sidepabel-tab-content' });
    await expect(designsTable).toBeVisible();

    // check if imported design is shown in getSidePanelTab
    await expect(designsTable.getByTitle(`designs-table-row-${name}`)).toBeVisible({
      timeout,
    });
  }

  async loadDesignFromSidePanel(name) {
    const page = this.page;
    await page.getByRole('tab', { name: 'Designs' }).click();
    const designsTable = page.getByRole('table', { name: 'sidepabel-tab-content' });
    await expect(designsTable).toBeVisible();
    await designsTable.getByTitle(`designs-table-row-${name}`).click();
  }
}

export class ConnectionsPage {
  /**
   * @param {PlaywrightPage} page
   */
  constructor(page) {
    this.page = page;
  }
}

export const test = base.test.extend({
  designsPage: async ({ page }, use) => {
    const designerPage = new DesignsPage(page);
    await designerPage.loadDefault();
    await use(designerPage);
  },
});

export const expect = base.expect;
