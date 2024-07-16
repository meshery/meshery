const { test, expect, MesheryDashboardPage } = require('../fixtures/pages');

test('Verifies the logout functionality', async ({ page }) => {
  const dashboardPage = new MesheryDashboardPage(page);

  await test.step('When I visit the dashboard page', async () => {
    await dashboardPage.goTo();
  });

  await test.step('And I hover over the profile and click `logout` button', async () => {
    const profileButton = dashboardPage.page
      .locator('[data-test="profile-button"]')
      .getByRole('button');
    await profileButton.hover();
    await page.waitForTimeout(500); // Adding a delay to ensure the menu appears
    await expect(dashboardPage.page.getByRole('menuitem', { name: 'Logout' })).toBeVisible();
    await dashboardPage.page.getByRole('menuitem', { name: 'Logout' }).click();
  });

  await test.step('And I see successfully logged out', async () => {
    await expect(page.getByLabel('Select Provider')).toBeVisible();
  });
});
