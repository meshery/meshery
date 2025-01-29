import { expect, test } from '@playwright/test';
import { ENV } from './env';

// Extensions Section Tests
test.describe('Extensions Section Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Visit The Meshery Serve URL
    await page.goto(ENV.MESHERY_SERVER_URL);
    const extensionsNav = page.locator('[data-cy="Extensions"]');
    await extensionsNav.click();
  });

  test('Verify Kanvas Snapshot using data-testid', async ({ page }) => {
    await expect(page.getByTestId('kanvas-snapshot-heading')).toBeVisible();
    await expect(page.getByTestId('kanvas-snapshot-description')).toBeVisible();

    const enableButton = page.getByTestId('kanvas-snapshot-enable-btn');
    await expect(enableButton).toBeVisible();
    await expect(enableButton).toBeDisabled();

    await expect(page.getByTestId('kanvas-snapshot-image')).toBeVisible();
  });

  test('Verify Performance Analysis Details', async ({ page }) => {
    await expect(page.getByTestId('performance-analysis-heading')).toBeVisible();
    const performanceEnableButton = page.getByTestId('performance-analysis-enable-btn');
    await expect(performanceEnableButton).toBeVisible();
    await expect(performanceEnableButton).toBeDisabled();
  });

  test('Verify Kanvas Details', async ({ page, context }) => {
    await expect(page.getByTestId('kanvas-signup-heading')).toBeVisible();
    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      page.getByTestId('kanvas-signup-btn').click(),
    ]);
    await expect(newPage).toHaveURL('https://docs.layer5.io/kanvas/');
    await newPage.close();
  });

  test('Verify Meshery Docker Extension Details', async ({ page, context }) => {
    await expect(page.getByTestId('docker-extension-heading')).toBeVisible();
    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      page.getByTestId('docker-extension-download-btn').click(),
    ]);
    await expect(newPage).toHaveURL(
      'https://hub.docker.com/extensions/meshery/docker-extension-meshery',
    );
    await newPage.close();
  });

  test('Verify Meshery Design Embed Details', async ({ page, context }) => {
    await expect(page.getByTestId('design-embed-learn-more-btn')).toBeVisible();
    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      page.getByTestId('design-embed-learn-more-btn').click(),
    ]);
    await expect(newPage).toHaveURL('https://docs.layer5.io/kanvas/designer/embedding-designs/');
    await newPage.close();
  });

  test('Verify Meshery Catalog Section Details', async ({ page, context }) => {
    await expect(page.getByTestId('catalog-section-heading')).toBeVisible();
    const toggleButton = page.getByTestId('catalog-toggle-switch');
    await toggleButton.click();
    const catalogLink = page.locator('a[href="https://meshery.io/catalog"]');
    const [newPage] = await Promise.all([context.waitForEvent('page'), catalogLink.click()]);
    await expect(newPage).toHaveURL('https://meshery.io/catalog');
    await newPage.close();
  });

  test('Verify Meshery Adapter for Istio Section', async ({ page, context }) => {
    // Test the "Open Adapter docs" link
    const [docsPage] = await Promise.all([
      context.waitForEvent('page'),
      page.getByRole('link', { name: 'Open Adapter docs' }).click(),
    ]);
    await expect(docsPage).toHaveURL('https://docs.meshery.io/concepts/architecture/adapters');
    await docsPage.close();
  });
});
