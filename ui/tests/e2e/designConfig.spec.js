import { expect, test } from '@playwright/test';
import { DesignConfiguratorPage } from './pages/DesignConfiguratorPage';

test.describe('Design Configurator Tests', () => {
  let designConfigPage;

  test.beforeEach(async ({ page }) => {
    await page.route('**/api/pattern', async (route) => {
      await route.fulfill();
    });

    designConfigPage = new DesignConfiguratorPage(page);
    await designConfigPage.navigateTo();
  });

  test('should verify Design Configurator page elements', async () => {
    await expect(designConfigPage.appBar).toBeVisible();
    await expect(designConfigPage.codeEditor).toBeVisible();
  });

  test('should edit design in Design Configurator', async () => {
    await designConfigPage.categorySelector.click();
    await designConfigPage.databaseCategory.click();
    await designConfigPage.modelSelector.click();
    await designConfigPage.modelAzure.click();

    await expect(designConfigPage.modelContainer).toBeVisible();

    await designConfigPage.updateDesign();
  });
});
