import { expect, test } from '@playwright/test';
import { ENV } from './env';

// Extensions Section Tests
test.describe('Extensions Section Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Visit The Meshery Serve URL
    await page.goto(ENV.MESHERY_SERVER_URL);

    const extensionsNav = page.locator('[data-test="navigation"] >> text=Extensions');
    await extensionsNav.click();
  });

  test('Verify Kanvas Snapshot using data-testid', async ({ page }) => {
    // Verify heading using data-testid
    await expect(page.getByTestId('kanvas-snapshot-heading')).toBeVisible();
    
    // Verify description using data-testid
    await expect(page.getByTestId('kanvas-snapshot-description')).toBeVisible();
    
    // Verify enable button using data-testid
    const enableButton = page.getByTestId('kanvas-snapshot-enable-btn');
    await expect(enableButton).toBeVisible();
    await expect(enableButton).toBeDisabled();

    // Verify snapshot logo
    await expect(page.locator('img[src="/static/img/meshmap-snapshot-logo.svg"]')).toBeVisible();
  });

  test('Verify Performance Analysis Details', async ({ page }) => {
    // Verify heading using data-testid
    await expect(page.getByTestId('performance-analysis-heading')).toBeVisible();
    
    // Verify enable button using data-testid
    const performanceEnableButton = page.getByTestId('performance-analysis-enable-btn');
    await expect(performanceEnableButton).toBeVisible();
    await expect(performanceEnableButton).toBeDisabled();
  });

  test('Verify Kanvas Details', async ({ page, context }) => {
    // Verify Kanvas heading using data-testid
    await expect(page.getByTestId('kanvas-signup-heading')).toBeVisible();

    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      page.getByTestId('kanvas-signup-btn').click(),
    ]);

    await expect(newPage).toHaveURL('https://docs.layer5.io/kanvas/');
    await newPage.close();
  });

  test('Verify Meshery Docker Extension Details', async ({ page, context }) => {
    // Verify Docker Extension heading using data-testid
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
    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      page.getByTestId('design-embed-learn-more-btn').click(),
    ]);

    await expect(newPage).toHaveURL('https://docs.layer5.io/kanvas/designer/embedding-designs/');
    await newPage.close();
  });

  test('Verify Meshery Catalog Section Details', async ({ page, context }) => {
    // Verify Catalog section heading using data-testid
    await expect(page.getByTestId('catalog-section-heading')).toBeVisible();

    // First verify the toggle button functionality
    const toggleButton = page
      .locator(
        'div:has-text("Explore the Meshery Catalog") >> .MuiSwitch-root input[type="checkbox"]',
      )
      .nth(0);
    await expect(toggleButton).toBeVisible();

    // Click toggle Button and Verify the State changes
    const initialState = await toggleButton.isChecked();
    await toggleButton.click();
    await expect(toggleButton).toBeChecked(!initialState);

    // Verify the Meshery Catalog Link
    const catalogLink = page.locator('a[href="https://meshery.io/catalog"]');
    const [newPage] = await Promise.all([context.waitForEvent('page'), catalogLink.click()]);

    // Verify the URL
    await expect(newPage).toHaveURL('https://meshery.io/catalog');
    await newPage.close();
  });

  test('Verify Meshery Adapter for Istio Section', async ({ page, context }) => {
    // Find the Istio section container to scope our selectors
    const istioSection = page.locator('div', {
      has: page.locator('text=Meshery Adapter for Istio'),
    });

    // Test the "Open Adapter docs" link
    const adapterDocsLink = istioSection.locator('a:has-text("Open Adapter docs")').first();
    await expect(adapterDocsLink).toBeVisible();

    const [docsPage] = await Promise.all([context.waitForEvent('page'), adapterDocsLink.click()]);

    await expect(docsPage).toHaveURL('https://docs.meshery.io/concepts/architecture/adapters');
    await docsPage.close();

    const toggleButton = page
      .locator('div:has-text("Meshery Catalog") >> input[type="checkbox"]')
      .nth(1);
    await expect(toggleButton).toBeVisible();

    const initialState = await toggleButton.isChecked();
    await toggleButton.click();
    await expect(toggleButton).toBeChecked(!initialState);
  });
});
