import { expect, test } from '@playwright/test';
import { DesignPage } from './pages/DesignPage/DesignPage';
import { waitForSnackBar } from './utils/waitForSnackBar';
import { mockEnvironmentsApi, mockConnectionsApi } from './pages/DesignPage/utils/mockApiRoutes';

const IMPORT_SOURCES = [
  {
    type: 'File',
    pathOrUrl: 'ui/tests/e2e/assets/GuestBook App.yml',
    designName: 'GuestBook App',
  },
  {
    type: 'URL',
    pathOrUrl: 'https://example.com/guestbook-app.yml',
    designName: 'GuestBook App',
  },
];

const DESIGN_TYPES = [
  {
    type: 'published',
    getCard: (designPage) => designPage.getPublishedDesign(),
  },
  // TODO add private/unpublised design cards
];

test.describe('Design Page Tests', () => {
  let designPage;

  test.beforeEach(async ({ page }) => {
    await page.route('**/api/pattern/import', async (route) => await route.fulfill());
    designPage = new DesignPage(page);
    await designPage.navigateTo();
  });

  test('renders design page UI', async () => {
    await expect(designPage.createDesignBtn).toBeVisible();
    await expect(designPage.importDesignBtn).toBeVisible();
    await expect(designPage.searchBar).toBeVisible();
    await expect(designPage.designGrid).toBeVisible();
    expect(await designPage.designCards.count()).toBeGreaterThan(0);
  });

  DESIGN_TYPES.forEach(({ type, getCard }) => {
    test(`displays ${type} design card correctly`, async () => {
      const card = getCard(designPage);

      const visibleElements = [...Object.values(card.display)];

      for (const el of visibleElements) {
        await expect(el).toBeVisible();
      }

      await card.actionToggleBtn.click();

      const visibleActions = [...Object.values(card.actionElements)];

      for (const el of visibleActions) {
        await expect(el).toBeVisible();
      }
    });
  });

  IMPORT_SOURCES.forEach(({ type, pathOrUrl, designName }) => {
    test(`imports design via ${type}`, async () => {
      await designPage.clickImportDesignButton();
      const ImportModal = designPage.ImportModal;

      await expect(ImportModal.modal).toBeVisible();
      await ImportModal.clickImportType(type);

      await ImportModal.importDesign(type, pathOrUrl, designName);

      await waitForSnackBar(designPage.page, `${designName}" design uploaded`);
    });
  });

  test('deletes a published design from the list', async () => {
    await designPage.page.route('**/api/pattern/*', async (route) => await route.fulfill());
    const card = designPage.getPublishedDesign();
    await card.display.name.click();
    await designPage.deleteDesignBtn.click();
    await expect(designPage.deleteDesignModalHeader).toHaveText('Delete 1 Design?');
    await designPage.deleteConfirmationBtn.click();
  });

  test('deploys a published design to a connected cluster', async () => {
    await designPage.page.route(
      '**/api/environments?page=0&pagesize=all&orgID=*',
      mockEnvironmentsApi,
    );
    await designPage.page.route(
      '**/api/environments/*/connections?page=0&pagesize=all',
      mockConnectionsApi,
    );
    await designPage.page.route(
      '**/api/pattern/deploy?contexts=*',
      async (route) => await route.fulfill(),
    );

    const card = designPage.getPublishedDesign();
    await card.actionToggleBtn.click();
    await card.actionElements.deploy.click();

    const deployModal = designPage.DeployModal;

    await expect(deployModal.getStepLabel(0)).toHaveText('Validate Design');
    await deployModal.goToNextStep();

    await expect(deployModal.getStepLabel(1)).toHaveText('Identify Environments');
    await deployModal.selectEnvironment('Test');
    await deployModal.goToNextStep();

    await expect(deployModal.getStepLabel(2)).toHaveText('Dry Run');
    await deployModal.waitForLoader();

    await deployModal.deployWithBypass();

    await expect(deployModal.modalContent).toContainText('Finalize Deployment');
    await expect(deployModal.nextButton).toHaveText('Deploy');

    await deployModal.goToNextStep();

    await expect(deployModal.loader).toBeVisible();

    await deployModal.goToNextStep();
  });
});
