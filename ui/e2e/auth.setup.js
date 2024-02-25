import { test as setup, expect } from '@playwright/test';

export const authFile = 'playwright/.auth/user.json';

setup('authenticate', async ({ page }) => {
  // Perform authentication steps.

  await page.goto('/');
  await page.getByRole('button').click();
  await page.getByRole('menuitem', { name: 'Meshery' }).click();
  // await page.getByLabel('Password').fill('password');
  // await page.getByRole('button', { name: 'Sign in' }).click();
  // Wait until the page receives the cookies.
  //
  // Sometimes login flow sets cookies in the process of several redirects.
  // Wait for the final URL to ensure that the cookies are actually set.
  await page.waitForURL('https://meshery.layer5.io/**');
  // Alternatively, you can wait until the page reaches a state where all cookies are set.
  // await page.getByRole('textbox', { name: 'email' }).fill('my-user');
  await page.locator('input[name="identifier"]').fill(process.env.USER_EMAIL);
  await page.locator('input[name="password"]').fill(process.env.USER_PASSWORD);
  await page.getByRole('button', { name: 'Sign in', exact: true }).click();
  await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible();

  // End of authentication steps.

  await page.context().storageState({ path: authFile });
});