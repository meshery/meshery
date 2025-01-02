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
    // Verify the heading
    const kanvasHeading = page.locator('text=GitHub Action: Kanvas Snapshot');
    await kanvasHeading.isVisible();
    console.log('Kanvas Snapshot heading is visible');

    // Verify the "Enable" button
    const enableButton = page.locator('button:has-text("Enable")').first();
    await expect(enableButton).toBeVisible();
    await expect(enableButton).toBeEnabled();
    console.log('Enable button is visible and enabled');

    // Verify the Kanvas Image
    const kanvasImage = page.locator('img[src*="meshmap-snapshot-logo"]');
    await expect(kanvasImage).toBeVisible();
    console.log('Kanvas Image is visible');
  });

  test('Verify Performance Analysis Details', async ({ page }) => {
    // Verify the Performance Analysis Heading
    const performanceHeading = page.locator('text=GitHub Action: Performance Analysis');
    await performanceHeading.isVisible();
    console.log('Performance Analysis Heading is Visible');

    // Verify the "Enable" button works as expected or not
    const performanceEnableButton = page.locator('button:has-text("Enable")').nth(1);
    await expect(performanceEnableButton).toBeVisible();
    await expect(performanceEnableButton).toBeEnabled();
    console.log('Performance Enable button is visible and works fine');
  });

  // Kanvas Component Tests
  test('Verify Kanvas Details', async ({ page }) => {
    // Verify Kanvas Heading
    const kanvasHeading = page.locator('text=Kanvas');
    await kanvasHeading.isVisible();
    console.log('Kanvas Heading is Visible');

    // Verify "Sign Up" Button works or not
    const signUpButton = page.locator('button.MuiButtonBase-root:has-text("Sign Up")');

    // Verify button is visible
    await expect(signUpButton).toBeVisible();
    console.log('Sign Up button is visible');

    // Verify button is enabled
    await expect(signUpButton).toBeEnabled();
    console.log('Sign Up button is enabled');

    // Click the button
    await Promise.all([
      // Wait for the new page/tab to open
      page.waitForEvent('popup'),
      // Click the button
      signUpButton.click(),
    ]);

    // Get the new page/tab
    const newPage = await page.context().pages()[1];

    // Verify the URL
    await expect(newPage).toHaveURL('https://docs.layer5.io/kanvas/');
    console.log('Successfully navigated to Kanvas documentation');

    // Verify the new page has loaded
    await newPage.waitForLoadState('networkidle');

    // Close the new page/tab
    await newPage.close();
    console.log('Sign Up Button is working Fine');
  });

  // Meshery Docker Extension Tests
  test('Verify Meshery Docker Extension Details', async ({ page }) => {
    // Verify the Heading
    const dockerHeading = await page.locator('text=Meshery Docker Extension');
    await expect(dockerHeading).toBeVisible();
    console.log('Meshery Docker Extension Heading is Visible');

    // Verify the Download Button works or not
    const downloadButton = page.locator('button.MuiButtonBase-root:has-text("Download")');

    // Verify the button is visible
    await expect(downloadButton).toBeVisible();
    console.log('Download button is visible');

    // Verify the Button is enabled
    await expect(downloadButton).toBeEnabled();
    console.log('Download button is enabled');

    // Click the Button and Handle new tab
    await Promise.all([page.waitForEvent('popup'), downloadButton.click()]);
    const newPage = await page.context().pages()[1];

    // Verify the URL
    await expect(newPage).toHaveURL(
      'https://hub.docker.com/extensions/meshery/docker-extension-meshery',
    );
    console.log('Successfully navigated to Meshery Docker Extension');
    await newPage.waitForLoadState('networkidle');
    await newPage.close();
    console.log('Download Button is working Fine');
  });

  test('Verify Meshery Desing Embed Details', async ({ page }) => {
    // Verify the Leanr More Button works or not
    const learnMoreButton = page.locator('button.MuiButtonBase-root:has-text("Learn More")');
    await expect(learnMoreButton).toBeVisible();
    console.log('Learn More button is visible');

    await expect(learnMoreButton).toBeEnabled();
    console.log('Learn More button is enabled');

    await Promise.all([page.waitForEvent('popup'), learnMoreButton.click()]);
    const newPage = await page.context().pages()[1];

    // Verify URL
    await expect(newPage).toHaveURL('https://docs.layer5.io/kanvas/designer/embedding-designs/');
    console.log('Successfully navigated to Meshery Design Embed');
    await newPage.waitForLoadState('networkidle');
    await newPage.close();
    console.log('Learn More Button is working Fine');
  });

  test('Verify Meshery Catalog Section Details', async ({ page }) => {
    // First verify the toggle button functionality
    const toggleButton = page
      .locator('div:has-text("Meshery Catalog") >> input[type="checkbox"]')
      .nth(0);
    // console.log(await toggleButton.count());
    await page.waitForLoadState('networkidle');
    const istoggleVisible = await toggleButton.isVisible();

    if (istoggleVisible) {
      console.log('Toggle button is visible');

      // Check whether the Toggle Button is checked or Not BEFORE Clicking it
      const isCheckedBefore = await toggleButton.isChecked();
      console.log(`Toggle button Status: ${isCheckedBefore ? 'Enable' : 'Disable'}`);

      // Click the Button
      await toggleButton.click();
      console.log('Toggle button clicked');

      // Check whether the Toggle Button is checked or Not AFTER Clicking it
      const isCheckedAfter = await toggleButton.isChecked();
      console.log(`Toggle button Status: ${isCheckedAfter ? 'Enable' : 'Disable'}`);
    } else {
      console.log('Toggle button is not visible');
      throw new Error('Toggle button not found');
    }

    // Verify the Meshert Catalog Link
    const catalogLink = page.locator('a[href="https://meshery.io/catalog"]');

    // Verify the Link is visible
    await expect(catalogLink).toBeVisible();
    console.log('Meshery Catalog Link is visible');

    // Verify the link text
    const linkText = await page.locator('text=Explore the Meshery Catalog').isVisible();
    expect(linkText).toBeTruthy();
    console.log('Link text is correct');

    // Click the link and verify navigation
    await Promise.all([
      // Wait for the new page/tab to open
      page.waitForEvent('popup'),
      // Click the link
      catalogLink.click(),
    ]);

    // Get the new page/tab
    const newPage = await page.context().pages()[1];

    // Verify the URL
    await expect(newPage).toHaveURL('https://meshery.io/catalog');
    console.log('Successfully navigated to Meshery Catalog page');

    // Wait for the new page to load
    await newPage.waitForLoadState('networkidle');

    // Close the new page/tab
    await newPage.close();
    console.log('Meshery Catalog Link is working Fine');
  });

  // Meshery Adapter for Istio Section Tests
  test('Verify Meshery Adapter for Istio Section', async ({ page }) => {
    // Find the Istio section container to scope our selectors
    const istioSection = page.locator('div', {
      has: page.locator('text=Meshery Adapter for Istio'),
    });

    // Test the "Open Adapter docs" link
    const adapterDocsLink = istioSection.locator('a:has-text("Open Adapter docs")').first();
    await expect(adapterDocsLink).toBeVisible();
    console.log('Open Adapter docs link is visible');

    // Click the link and verify navigation
    await Promise.all([page.waitForEvent('popup'), adapterDocsLink.click()]);

    // Get and verify the new page
    const newPage = await page.context().pages()[1];
    await expect(newPage).toHaveURL('https://docs.meshery.io/concepts/architecture/adapters');
    console.log('Successfully navigated to Adapter documentation page');
    await newPage.close();

    // First verify the toggle button functionality
    const toggleButton = page
      .locator('div:has-text("Meshery Catalog") >> input[type="checkbox"]')
      .nth(1);
    // console.log(await toggleButton.count());
    await page.waitForLoadState('networkidle');
    const istoggleVisible = await toggleButton.isVisible();

    if (istoggleVisible) {
      console.log('Toggle button is visible');

      // Check whether the Toggle Button is checked or Not BEFORE Clicking it
      const isCheckedBefore = await toggleButton.isChecked();
      console.log(`Toggle button Status: ${isCheckedBefore ? 'Enable' : 'Disable'}`);

      // Click the Button
      await toggleButton.click();
      console.log('Toggle button clicked');

      // Check whether the Toggle Button is checked or Not AFTER Clicking it
      const isCheckedAfter = await toggleButton.isChecked();
      console.log(`Toggle button Status: ${isCheckedAfter ? 'Enable' : 'Disable'}`);
    } else {
      console.log('Toggle button is not visible');
      throw new Error('Toggle button not found');
    }
  });
});
