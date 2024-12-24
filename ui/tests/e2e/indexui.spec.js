import { expect, test } from '@playwright/test';
import { ENV } from './env';

test.describe('Index Page UI Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Visit Index Page
    await page.goto(ENV.MESHERY_SERVER_URL);
  });

  test('Test if Left Navigation Panel is displayed', async ({ page }) => {
    await expect(page.locator('[data-test=navigation]')).toBeVisible();
  });

  test('Test if Settings button is displayed', async ({ page }) => {
    await expect(page.locator('[data-test=settings-button]')).toBeVisible();
  });

  test('Test if Notification button is displayed', async ({ page }) => {
    await expect(page.locator('[data-test=notification-button]')).toBeVisible();
  });

  test('Test if Profile button is displayed', async ({ page }) => {
    await expect(page.locator('[data-test=profile-button]')).toBeVisible();
  });

  test('Test if Dashboard is displayed', async ({ page }) => {
    await expect(page.locator('[data-cy=dashboard]')).toBeVisible();
  });

  test('Test if Lifecycle is displayed', async ({ page }) => {
    await expect(page.locator('[data-cy=lifecycle]')).toBeVisible();
  });

  test('Test if Configuration is displayed', async ({ page }) => {
    await expect(page.locator('[data-cy=configuration]')).toBeVisible();
  });

  test('Test if Performance is displayed', async ({ page }) => {
    await expect(page.locator('[data-cy=performance]')).toBeVisible();
  });

  test('Test if Extensions is displayed', async ({ page }) => {
    await expect(page.locator('[data-cy=Extensions]')).toBeVisible();
  });
});

// Dashboard Section Tests
test.describe('Dashboard Section Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Visit the Meshery Server URL
    await page.goto(ENV.MESHERY_SERVER_URL);

    // Navigate to the Dashboard sectionj
    const dashboardNav = page.locator('[data-test="navigation"] >> text=Dashboard');
    await dashboardNav.click();
  });

  // Verifying if the Dashboard section has the elements present in navbar
  test.describe('Navbar Items Tests', () => {
    test('Test if all Navbar items are displayed', async ({ page }) => {
      const navbarItems = [
        'Overview',
        'Node',
        'Namespace',
        'Workload',
        'Configuration',
        'Network',
        'Security',
        'Storage',
      ];

      for (const item of navbarItems) {
        const element = await page.getByRole('tab', { name: item });
        await expect(element).toBeVisible();
      }
    });
  });

  // Navigating to Node in Navbar
  test('Test if default section is Overview and navigates to Node', async ({ page }) => {
    // Check if the default tab (Overview) is selected
    const overviewTab = await page.getByRole('tab', { name: 'Overview' });
    await expect(overviewTab).toHaveAttribute('aria-selected', 'true'); // Verify Overview is selected
    // console.log('Overview tab is selected');

    // Click on the Node tab
    const nodeTab = await page.getByRole('tab', { name: 'Node' });
    await nodeTab.click();
    // console.log('Node button clicked');

    // Verify the Node tab is now selected
    await expect(nodeTab).toHaveAttribute('aria-selected', 'true'); // Ensure Node is selected
    // console.log('Node tab Verified');
  });
});

// Lifecycle Section Tests
test.describe('Lifecycle Section Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Visit the Meshery Server URL
    await page.goto(ENV.MESHERY_SERVER_URL);

    // Navigate to the Lifecycle section
    const lifecycleNav = page.locator('[data-test="navigation"] >> text=Lifecycle');
    await lifecycleNav.click();
  });

  // Test to verify the "Connections" tab is present and functional
  test('Verify Connections tab under Lifecycle', async ({ page }) => {
    const connectionsTab = page.locator('button[role="tab"] >> text=Connections');
    await expect(connectionsTab).toBeVisible();
    await connectionsTab.click();

    // Verify the tab content is displayed
    const connectionsContent = page.locator('text="Add Cluster"');
    await expect(connectionsContent).toBeVisible();
  });

  // Test to verify the "MeshSync" tab is present and functional
  test('Verify MeshSync tab under Lifecycle', async ({ page }) => {
    const meshsyncTab = page.locator('button[role="tab"] >> text=MeshSync');
    await expect(meshsyncTab).toBeVisible();
    await meshsyncTab.click();

    // Verify the tab content for MeshSync is displayed
    const meshsyncContent = page.locator('text="MeshSync"');
    await expect(meshsyncContent).toBeVisible();
  });
});

// Configuration Sections Tests
test.describe('Configuration Section Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Visit the Meshery Server URL
    await page.goto(ENV.MESHERY_SERVER_URL);

    // Navigate to the Configuration section
    const configurationNav = page.locator('[data-test="navigation"] >> text=Configuration');
    await configurationNav.click();
    // console.log('Configuration section clicked');
  });

  // Test to verify the "Meshery Operator" tab is present and functional
  test('Verify Create Design buttons', async ({ page }) => {
    // Ensure the "Design" button is visible and clickable
    const designNav = page.locator('span.MuiTypography-root', { hasText: 'Design' });

    // Wait for the collapsible menu to expand and the "Design" button to appear
    await designNav.waitFor();
    await expect(designNav).toBeVisible();

    // Click the "Design" button
    await designNav.click();
    // console.log('Design button clicked');

    // Verify that the URL Matches the Design URL
    await expect(page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/configuration/designs`);
    // console.log('URL Verified');

    // Verify "Create Design Button"
    const createDesignButton = page.locator('button:has-text("Create Design")');

    // Wait for potential loading states
    await page.waitForLoadState('networkidle');
    const isButtonVisible = await createDesignButton.isVisible();

    if (isButtonVisible) {
      console.log('Create Design button found, clicking it...');
      await createDesignButton.click();

      // Verify the code editor is visible
      const codeEditor = page.locator('.react-codemirror2');
      await expect(codeEditor).toBeVisible();
      console.log('Code Edior is visible');

      // Verify default YAML content is present
      const expectedContent = [
        'name: Untitled Design',
        'components: []',
        'schemaVersion: designs.meshery.io/v1beta1',
      ];

      // Get the text content and verify it contains expected lines
      const editorContent = await page.locator('.CodeMirror-lines').textContent();
      for (const line of expectedContent) {
        expect(editorContent).toContain(line);
      }
    } else {
      console.log('Create Design button not found');
      throw new Error('Create Design button not found on the page');
    }
  });
});

// Performance Section Tests
test.describe('Performance Section Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Visit the Meshery Servevr URL
    await page.goto(ENV.MESHERY_SERVER_URL);

    // Navigate to Performance Section
    const performanceNav = page.locator('[data-test="navigation"] >> text=Performance');
    await performanceNav.click();
    // console.log('Performance section clicked'); });
  });

  // Test to verify the Week adn Month Button is present and functional
  test('Verify Week and Month Button', async ({ page }) => {
    const weekButton = page.locator('button:has-text("Week")');
    const monthButton = page.locator('button:has-text("Month")');

    // Wait for the button to appear
    await weekButton.isVisible();
    await weekButton.click();
    // console.log('Week button is visible and Clicked');

    // Wait for the button to appear
    await monthButton.isVisible();
    await monthButton.click();
    // console.log('Month button is visible and Clicked');
  });

  // Test to verify the Configure Metrics Button is present and functional
  test('Verify Configure Metrics Button', async ({ page }) => {
    const configureMetricsButton = page.locator('button:has-text("Configure Metrics")');

    // Wait for Button and Click it
    await page.waitForLoadState('networkidle');
    const isButtonVisible = await configureMetricsButton.isVisible();

    if (isButtonVisible) {
      console.log('Configure Metrics button found, clicking it...');
      await configureMetricsButton.click();

      // Verify the URL after clicking
      await expect(page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/settings#metrics`);
      console.log('Settings metrics URL verified');

      // Verify the Setting tabs are visible
      const tabstoVerify = ['Adapters', 'Metrics', 'Registry', 'Reset'];
      for (const tabName of tabstoVerify) {
        const tab = page.locator(`button[role="tab"]:has-text("${tabName}")`);
        await expect(tab).toBeVisible();
        console.log(`${tabName} tab is visible`);
      }

      // Verify Mesh Adapter elements
      const meshAdapterLabel = page.locator('label.MuiFormLabel-root:has-text("Mesh Adapter URL")');
      await expect(meshAdapterLabel).toBeVisible();

      const meshAdapterSelect = page.locator('#react-select-2-placeholder');
      await expect(meshAdapterSelect).toBeVisible();

      // Verify action buttons
      const undeployButton = page.locator('button:has-text("Undeploy")');
      await expect(undeployButton).toBeVisible();

      const connectButton = page.locator('button:has-text("Connect")');
      await expect(connectButton).toBeVisible();

      // Verify Available Adapters section
      const availableAdaptersText = page.locator(
        'label.MuiFormLabel-root:has-text("Available Adapters")',
      );
      await expect(availableAdaptersText).toBeVisible();

      const deployButton = page.locator(
        'button:has-text("Deploy"):below(label:has-text("Available Adapters"))',
      );
      await expect(deployButton).toBeVisible();
    } else {
      console.log('Configure Metrics button not found');
      throw new Error('Configure Metrics button not found on the page');
    }
  });
});

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

    // Verify the Description text
    const description = page.locator(
      'text=Connect Kanvas to your GitHub repo and see changes pull request-to-pull request. Get snapshots of your infrastructure directly in your PRs.',
    );
    await expect(description).toBeVisible();
    console.log('Description text is visible');

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

    // Verify the Performance Analysis Description
    const performanceDescription = page.locator(
      "text=Characterize the performance of your services using Meshery's performance analysis GitHub Action to benchmark and visually compare percentiles (e.g. P99) over time.",
    );
    await performanceDescription.isVisible();
    console.log('Performance Analysis Description is Visible');

    // Verify the "Enable" button
    const performanceEnableButton = page.locator('button:has-text("Enable")').nth(1);
    await expect(performanceEnableButton).toBeVisible();
    await expect(performanceEnableButton).toBeEnabled();
    console.log('Performance Enable button is visible');

    // Verify The Performance Analysis Image
    const performanceImage = page.locator('img[src*="smp-dark"]');
    await expect(performanceImage).toBeVisible();
    console.log('Performance Analysis Image is visible');
  });

  test('Verify and Click Toggle Button in Meshery Catalog Section', async ({ page }) => {
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
  });
});
