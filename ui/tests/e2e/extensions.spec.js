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

  test('Verify Kanvas Snapshot details', async ({ page }) => {
    await expect(
      page.locator('.MuiTypography-h5:has-text("GitHub Action: Kanvas Snapshot")'),
    ).toBeVisible();

    // Verify the "Enable" button
    const enableButton = page.locator('.MuiButton-label:has-text("Enable")').first();
    await enableButton.click();

    await expect(page.locator('img[src*="meshmap-snapshot-logo"]')).toBeVisible();
  });

  test('Verify Performance Analysis Details', async ({ page }) => {
    await expect(
      page.locator('.MuiTypography-h5:has-text("GitHub Action: Performance Analysis")'),
    ).toBeVisible();

    const performanceEnableButton = page.locator('.MuiButton-label:has-text("Enable")').nth(1);
    await expect(performanceEnableButton).toBeVisible();
  });

  // Kanvas Component Tests
  test('Verify Kanvas Details', async ({ page, context }) => {
    await expect(page.locator('.MuiTypography-h5:has-text("Kanvas")').first()).toBeVisible();

    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      page.locator('.MuiButton-label:has-text("Sign Up")').click(),
    ]);

    await expect(newPage).toHaveURL('https://docs.layer5.io/kanvas/');
    await newPage.close();
  });

  // Meshery Docker Extension Tests
  test('Verify Meshery Docker Extension Details', async ({ page, context }) => {
    await expect(
      page.locator('.MuiTypography-h5:has-text("Meshery Docker Extension")'),
    ).toBeVisible();

    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      page.locator('.MuiButton-label:has-text("Download")').click(),
    ]);

    await expect(newPage).toHaveURL(
      'https://hub.docker.com/extensions/meshery/docker-extension-meshery',
    );
    await newPage.close();
  });

  test('Verify Meshery Design Embed Details', async ({ page, context }) => {
    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      page.locator('.MuiButton-label:has-text("Learn More")').click(),
    ]);

    await expect(newPage).toHaveURL('https://docs.layer5.io/kanvas/designer/embedding-designs/');
    await newPage.close();
  });

  test('Verify Meshery Catalog Section Details', async ({ page, context }) => {
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

    // Verify the Meshert Catalog Link
    const catalogLink = page.locator('a[href="https://meshery.io/catalog"]');
    const [newPage] = await Promise.all([context.waitForEvent('page'), catalogLink.click()]);

    // Verify the URL
    await expect(newPage).toHaveURL('https://meshery.io/catalog');
    await newPage.close();
  });

  // Meshery Adapter for Istio Section Tests
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
