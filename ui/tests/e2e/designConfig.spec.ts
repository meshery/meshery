import { expect, test, Page, Route, Request } from '@playwright/test';
import { DesignConfiguratorPage } from './pages/DesignConfiguratorPage';

const DESIGN_ID = 'test-design';
const DESIGN_PATTERN = `name: Sample Design
schemaVersion: designs.meshery.io/v1beta1
components: []`;

const MOCK_CATEGORY = 'Database';
const MOCK_MODEL = 'azure-db-for-mysql';
const MOCK_MODEL_VERSION = '1.0.0';

interface DesignResponse {
  id: string;
  patternFile: string;
}

interface Category {
  name: string;
}

interface CategoriesResponse {
  categories: Category[];
}

interface Model {
  name: string;
  displayName: string;
  version: string;
}

interface ModelsResponse {
  models: Model[];
}

interface ComponentModel {
  name: string;
  version: string;
  category: { name: string; metadata: null };
  registrant: { hostname: string; kind: string };
  metadata: Record<string, unknown>;
  model: { version: string };
}

interface Component {
  id: string;
  component: { kind: string; version: string; schema: string };
  displayName: string;
  model: ComponentModel;
  metadata: Record<string, unknown>;
}

interface ComponentsResponse {
  components: Component[];
}

interface UpdateDesignPayload {
  id: string;
  name: string;
}

test.describe('Design Configurator Tests', () => {
  let designConfigPage: DesignConfiguratorPage;

  test.beforeEach(async ({ page }: { page: Page }) => {
    await page.route(`**/api/pattern/${DESIGN_ID}`, async (route: Route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: DESIGN_ID,
          patternFile: DESIGN_PATTERN,
        } as DesignResponse),
      });
    });

    await page.route('**/api/pattern', async (route: Route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([{ id: DESIGN_ID }] as DesignResponse[]),
      });
    });

    await page.route('**/api/registry/categories**', async (route: Route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          categories: [{ name: MOCK_CATEGORY }],
        } as CategoriesResponse),
      });
    });

    await page.route(
      `**/api/registry/categories/${MOCK_CATEGORY}/models**`,
      async (route: Route) => {
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
          } as ModelsResponse),
        });
      },
    );

    await page.route(`**/api/registry/models/${MOCK_MODEL}/components**`, async (route: Route) => {
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
        } as ComponentsResponse),
      });
    });

    designConfigPage = new DesignConfiguratorPage(page);
    await designConfigPage.navigateTo(DESIGN_ID);
  });

  test('should verify Design Configurator page elements', async () => {
    await expect(designConfigPage.appBar).toBeVisible();
    await expect(designConfigPage.codeEditor).toBeVisible();
  });

  test('should edit design in Design Configurator', async ({ page }: { page: Page }) => {
    await designConfigPage.categorySelector.click();
    await designConfigPage.databaseCategory.click();

    await designConfigPage.modelSelector.click();
    await designConfigPage.modelAzure.click();

    await expect(designConfigPage.modelContainer).toBeVisible();

    const updateRequest = page.waitForRequest(
      (request: Request) => request.url().includes('/api/pattern') && request.method() === 'POST',
    );

    await designConfigPage.updateDesign();

    const updatePayload = (await updateRequest).postDataJSON() as UpdateDesignPayload;
    await expect(updatePayload).toMatchObject({
      id: DESIGN_ID,
      name: 'Sample Design',
    });
  });
});
