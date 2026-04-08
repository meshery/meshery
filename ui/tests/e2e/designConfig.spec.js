import { expect, test } from '@playwright/test';
import { DesignConfiguratorPage } from './pages/DesignConfiguratorPage';

const DESIGN_ID = 'test-design';
const DESIGN_PATTERN = `name: Sample Design
schemaVersion: designs.meshery.io/v1beta1
components: []
`;

const MOCK_CATEGORY = 'Database';
const MOCK_MODEL = 'azure-db-for-mysql';

test.describe('Design Configurator Tests', () => {
  let designConfigPage;

  test.beforeEach(async ({ page }) => {
    // Mock design fetch
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

    // Mock meshmodel categories
    await page.route('**/api/meshmodels/categories**', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          categories: [{ name: MOCK_CATEGORY }],
        }),
      });
    });

    // Mock meshmodel models for the selected category
    await page.route(`**/api/meshmodels/categories/${MOCK_CATEGORY}/models**`, async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          models: [
            {
              name: MOCK_MODEL,
              displayName: 'Azure DB for MySQL',
              version: '1.0.0',
              model: { version: '1.0.0' },
            },
          ],
        }),
      });
    });

    // Mock meshmodel components for the selected model
    await page.route(`**/api/meshmodels/models/${MOCK_MODEL}/components**`, async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          components: [
            {
              id: 'comp-1',
              component: { kind: 'FlexibleServer', version: 'v1', schema: '{}' },
              displayName: 'Flexible Server',
              model: { name: MOCK_MODEL, version: '1.0.0' },
              metadata: {},
            },
          ],
        }),
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
    await expect(designConfigPage.databaseCategory).toBeVisible();
    await designConfigPage.databaseCategory.click();

    await designConfigPage.modelSelector.click();
    await designConfigPage.modelAzure.click();

    await expect(designConfigPage.modelContainer).toBeVisible();
    await designConfigPage.updateDesign();
  });
});
