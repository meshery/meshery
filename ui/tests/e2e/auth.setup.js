import { expect, test as setup } from './fixtures/project';
import { ENV } from './env';

setup.describe.configure({ mode: 'serial' });

setup('authenticate as Meshery provider', async ({ page }) => {
  // Perform authentication steps. Replace these actions with your own.

  await page.goto(ENV.PROVIDER_SELECTION_URL);
  await page.getByLabel('Select Provider').click();
  await page.getByRole('menuitem', { name: 'Meshery' }).click();

  await page.getByLabel('E-Mail').fill(ENV.REMOTE_PROVIDER_USER.email);
  await page.getByLabel('Password').fill(ENV.REMOTE_PROVIDER_USER.password);
  await page.getByRole('button', { name: 'Sign in', exact: true }).click();

  // Wait until the page receives the cookies.
  // Sometimes login flow sets cookies in the process of several redirects.
  // Wait for the final URL to ensure that the cookies are actually set.

  await expect(async () => {
    const url = page.url();

    const redirect_urls = new Set([
      ENV.MESHERY_SERVER_URL + '/',
      ENV.REMOTE_PROVIDER_URL + '/',
      ENV.REMOTE_PROVIDER_URL + '/dashboard',
    ]);
    const redirected = redirect_urls.has(url);
    return expect(redirected).toBeTruthy();
  }).toPass();
  // End of authentication steps.
  await page.context().storageState({ path: ENV.AUTHFILEMESHERYPROVIDER });
});

setup('authenticate as None provider', async ({ page }) => {
  // Perform authentication steps. Replace these actions with your own.

  await page.goto(ENV.PROVIDER_SELECTION_URL);
  await page.getByLabel('Select Provider').click();
  await page.getByRole('menuitem', { name: 'None' }).click();

  // Wait until the page receives the cookies.
  // Sometimes login flow sets cookies in the process of several redirects.
  // Wait for the final URL to ensure that the cookies are actually set.

  await expect(async () => {
    const url = page.url();

    const redirect_urls = new Set([
      ENV.MESHERY_SERVER_URL + '/',
      ENV.REMOTE_PROVIDER_URL + '/',
      ENV.REMOTE_PROVIDER_URL + '/dashboard',
    ]);
    const redirected = redirect_urls.has(url);
    return expect(redirected).toBeTruthy();
  }).toPass();
  // End of authentication steps.
  await page.context().storageState({ path: ENV.AUTHFILELOCALPROVIDER });
});
