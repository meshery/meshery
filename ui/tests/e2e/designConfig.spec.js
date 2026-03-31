import { expect, test } from '@playwright/test';
import { DesignConfiguratorPage } from './pages/DesignConfiguratorPage';

const DESIGN_ID = 'test-design';
const DESIGN_PATTERN = `name: Sample Design
schemaVersion: designs.meshery.io/v1beta1
components: []
`;

test.describe('Design Configurator Tests', () => {
  let designConfigPage;

  test.beforeEach(async ({ page }) => {
    await page.route(`**/api/pattern/${DESIGN_ID}`, async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: DESIGN_ID,
          pattern_file: DESIGN_PATTERN,
        }),
      });
    });

    await page.route('**/api/pattern', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([{ id: DESIGN_ID }]),
      });
    });

    designConfigPage = new DesignConfiguratorPage(page);
    await designConfigPage.navigateTo(DESIGN_ID);
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
