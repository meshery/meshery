import { expect, test } from '@playwright/test';
import { DesignConfiguratorPage } from './pages/DesignConfiguratorPage';

const DESIGN_ID = 'test-design';
const DESIGN_PATTERN = `name: Sample Design
schemaVersion: designs.meshery.io/v1beta1
components: []
`;

const MOCK_CATEGORY = 'Database';
const MOCK_MODEL = 'azure-db-for-mysql';
const MOCK_MODEL_VERSION = '1.0.0';

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
          patternFile: DESIGN_PATTERN,
        }),
      });
    });

    await page.route('**/api/pattern', async (route) => {
      if (route.request().method() === 'GET') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify([{ id: DESIGN_ID }]),
        });
        return;
      }

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
              version: MOCK_MODEL_VERSION,
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
              model: {
                name: MOCK_MODEL,
                version: MOCK_MODEL_VERSION,
                category: { name: MOCK_CATEGORY, metadata: null },
                registrant: { hostname: 'artifacthub', kind: 'artifacthub' },
                metadata: {},
                model: { version: MOCK_MODEL_VERSION },
              },
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

  test('should edit design in Design Configurator', async ({ page }) => {
    await designConfigPage.categorySelector.click();
    await expect(designConfigPage.databaseCategory).toBeVisible();
    await designConfigPage.databaseCategory.click();

    await designConfigPage.modelSelector.click();
    await designConfigPage.modelAzure.click();

    await expect(designConfigPage.modelContainer).toBeVisible();

    const updateRequest = page.waitForRequest(
      (request) => request.url().includes('/api/pattern') && request.method() === 'POST',
    );

    await designConfigPage.updateDesign();

    const updatePayload = (await updateRequest).postDataJSON();
    await expect(updatePayload).toMatchObject({
      id: DESIGN_ID,
      name: 'Sample Design',
    });
  });
});
