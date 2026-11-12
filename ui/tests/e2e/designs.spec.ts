import { test, expect } from './fixtures/project';
import { DesignPage } from './pages/DesignPage/DesignPage';
import { waitForSnackBar } from './utils/waitForSnackBar';
import {
  mockEnvironmentsApi,
  mockConnectionsApi,
  mockPatternsApi,
} from './pages/DesignPage/utils/mockApiRoutes';

const IMPORT_SOURCES = [
  {
    type: 'File',
    pathOrUrl: 'tests/e2e/fixtures/relationships/meshery-design-fixture.json',
    designName: 'Meshery Design',
  },
  {
    type: 'URL',
    pathOrUrl:
      'https://raw.githubusercontent.com/meshery/meshery/refs/heads/master/ui/tests/e2e/fixtures/relationships/meshery-design-fixture.json',
    designName: 'GuestBook App',
  },
];

const DESIGN_TYPES = [
  {
    type: 'published',
  },
  {
    type: 'private',
  },
];

test.describe('Design Page Tests', () => {
  let designPage: DesignPage;

  test.beforeEach(async ({ page }) => {
    // Intercept and mock API routes for deterministic, resilient E2E execution
    await page.route('**/api/pattern?*', mockPatternsApi);
    await page.route('**/api/pattern', async (route) => {
      if (route.request().method() === 'POST') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ message: 'Pattern saved successfully' }),
        });
      } else {
        await mockPatternsApi(route);
      }
    });

    await page.route('**/api/pattern/import', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([{ id: 'imported-design-1', name: 'Imported Design' }]),
      });
    });

    await page.route('**/api/pattern/deploy*', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'Deployment initiated' }),
      });
    });

    await page.route('**/api/pattern/catalog/unpublish', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'Design Unpublished' }),
      });
    });

    await page.route('**/api/patterns/delete', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'Design Deleted' }),
      });
    });

    await page.route('**/api/pattern/*', async (route) => {
      if (route.request().method() === 'DELETE') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ message: 'Design Deleted' }),
        });
      } else {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            id: 'design-test-published-1',
            name: 'Published Sample Design',
            patternFile:
              'name: Published Sample Design\nschemaVersion: designs.meshery.io/v1beta1\ncomponents: []',
          }),
        });
      }
    });

    await page.route('**/api/environments?*', mockEnvironmentsApi);
    await page.route('**/api/environments/*/connections?*', mockConnectionsApi);

    designPage = new DesignPage(page);
    await designPage.navigateTo();
  });

  test('renders design page UI elements', async () => {
    await expect(designPage.createDesignBtn).toBeVisible();
    await expect(designPage.importDesignBtn).toBeVisible();
    await expect(designPage.searchBar).toBeVisible();
    await expect(designPage.designGrid).toBeVisible();
    expect(await designPage.designCards.count()).toBeGreaterThan(0);
  });

  DESIGN_TYPES.forEach(({ type }) => {
    test(`displays ${type} design card correctly`, async ({ provider }) => {
      test.skip(
        provider === 'Local' && type === 'public',
        `Skipping public design test for provider: ${provider}`,
      );

      await designPage.applyVisibilityFilter(type);
      const card = designPage.getFirstCardByVisibilityBadge(type);
      const cardElements = designPage.getCardElements(card, type);

      await expect(cardElements.display.name).toBeVisible();
      await expect(cardElements.display.visibility).toBeVisible();
      await expect(cardElements.display.actionsBar).toBeVisible();

      await cardElements.actionToggleBtn.click();
      await expect(cardElements.actionElements.validate).toBeVisible();
      await expect(cardElements.actionElements.deploy).toBeVisible();
    });
  });

  IMPORT_SOURCES.forEach(({ type, pathOrUrl, designName }) => {
    test(`imports design via ${type}`, async () => {
      await designPage.clickImportDesignButton();
      const importModal = designPage.ImportModal;

      await expect(importModal.modal).toBeVisible();
      await importModal.clickImportType(type);
      await importModal.importDesign(type, pathOrUrl, designName);

      await waitForSnackBar(designPage.page, `${designName}" design uploaded`);
    });
  });

  test('opens details of design', async () => {
    const card = designPage.getFirstCardByVisibilityBadge('published');
    const cardElements = designPage.getCardElements(card, 'published');

    await cardElements.actionElements.info.click();
    await expect(
      designPage.page.locator('[role="dialog"]').or(designPage.page.getByTestId('info-modal')),
    ).toBeVisible();
  });

  test('opens download design modal and displays export options', async () => {
    const card = designPage.getFirstCardByVisibilityBadge('published');
    const cardElements = designPage.getCardElements(card, 'published');

    await cardElements.actionElements.download.click();
    await expect(designPage.page.getByRole('dialog')).toBeVisible();
    await expect(designPage.page.getByText('Meshery YAML')).toBeVisible();
  });

  test('validates a design', async () => {
    const card = designPage.getFirstCardByVisibilityBadge('published');
    const cardElements = designPage.getCardElements(card, 'published');

    await cardElements.actionToggleBtn.click();
    await cardElements.actionElements.validate.click();

    await expect(designPage.page.getByRole('dialog')).toBeVisible();
    await expect(designPage.page.getByText(/Validate design/i)).toBeVisible();
  });

  test('unpublishes a published design', async () => {
    await designPage.applyVisibilityFilter('published');
    const card = designPage.getFirstCardByVisibilityBadge('published');
    const cardElements = designPage.getCardElements(card, 'published');

    if (await cardElements.actionElements.unpublish.isVisible()) {
      await cardElements.actionElements.unpublish.click();
      await expect(designPage.page.getByText(/Unpublish Catalog item\?/i)).toBeVisible();
      const unpublishConfirmBtn = designPage.page.locator('button:has-text("UNPUBLISH")').first();
      await expect(unpublishConfirmBtn).toBeVisible();
      await unpublishConfirmBtn.click();
      await waitForSnackBar(designPage.page, 'Design Unpublished');
    }
  });

  test('deletes a published design from the list', async () => {
    await designPage.applyVisibilityFilter('published');
    const card = designPage.getFirstCardByVisibilityBadge('published');
    const cardElements = designPage.getCardElements(card, 'published');

    await cardElements.display.name.click();
    await designPage.deleteDesignBtn.click();
    await expect(designPage.deleteDesignModalHeader).toBeVisible();
    await designPage.deleteConfirmationBtn.click();
  });

  test('deploys a published design to a connected cluster', async () => {
    await designPage.applyVisibilityFilter('published');
    const card = designPage.getFirstCardByVisibilityBadge('published');
    const cardElements = designPage.getCardElements(card, 'published');

    await cardElements.actionToggleBtn.click();
    await cardElements.actionElements.deploy.click();

    const deployModal = designPage.DeployModal;

    await expect(deployModal.getStepLabel(0, 'Validate Design')).toBeVisible();
    await deployModal.goToNextStep();

    await expect(deployModal.getStepLabel(1, 'Identify Environments')).toBeVisible();
    await deployModal.selectEnvironment('Test');
    await deployModal.goToNextStep();

    await expect(deployModal.getStepLabel(2, 'Dry Run')).toBeVisible();
    await deployModal.waitForLoader();

    await deployModal.deployWithBypass();

    await expect(deployModal.modalContent).toContainText('Finalize Deployment');
    await expect(deployModal.nextButton).toHaveText('Deploy');

    await deployModal.goToNextStep();
    await deployModal.goToNextStep();
  });

  test('undeploys a published design from a connected cluster', async () => {
    await designPage.applyVisibilityFilter('published');
    const card = designPage.getFirstCardByVisibilityBadge('published');
    const cardElements = designPage.getCardElements(card, 'published');

    await cardElements.actionToggleBtn.click();
    await cardElements.actionElements.unDeploy.click();

    const deployModal = designPage.DeployModal;
    await expect(deployModal.modalContent).toBeVisible();
    await expect(designPage.page.getByText(/Undeploy design/i)).toBeVisible();
  });

  test('edits a design in the configurator', async () => {
    await designPage.applyVisibilityFilter('published');
    const card = designPage.getFirstCardByVisibilityBadge('published');
    const cardElements = designPage.getCardElements(card, 'published');

    if (await cardElements.actionElements.edit.isVisible()) {
      await cardElements.actionElements.edit.click();
      await expect(designPage.page).toHaveURL(/.*pattern.*/);
    }
  });
});
