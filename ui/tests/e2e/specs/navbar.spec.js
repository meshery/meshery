const { test, expect, MesheryDashboardPage } = require('../fixtures/pages');

const fs = require('fs').promises;
test.describe('Navbar Tests', () => {
  let dashboardPage;

  test.beforeEach(async ({ page }) => {
    dashboardPage = new MesheryDashboardPage(page);
    await test.step('When I visit the dashboard page', async () => {
      await dashboardPage.goTo();
    });
  });

  test('Verifies the download authentication token', async () => {
    await test.step('I Hover over the profile and click `Get Token`', async () => {
      await dashboardPage.hoverAndClickMenuItem('Get Token');
    });

    await test.step('And I see the file contains the token details', async () => {
      const [download] = await Promise.all([dashboardPage.page.waitForEvent('download')]);
      const path = await download.path();
      const content = await fs.readFile(path, 'utf-8');
      const jsonContent = JSON.parse(content);
      expect(jsonContent['meshery-provider']).toBe('Meshery');
      expect(jsonContent).toHaveProperty('token');
    });
  });

  test('Verifies the logout functionality', async () => {
    await test.step('I Hover over the profile and click `Logout`', async () => {
      await dashboardPage.hoverAndClickMenuItem('Logout');
    });

    await test.step('And I see successfully logged out', async () => {
      await expect(dashboardPage.page.getByLabel('Select Provider')).toBeVisible();
    });
  });
});
