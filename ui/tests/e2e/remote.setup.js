import { test as setup } from './fixtures/project';
import { ProviderSelectionPage } from './pages/ProviderSelectionPage';
import { LoginPage } from './pages/LoginPage';
import { ENV } from './env';

const PROVIDERS = {
  LOCAL: 'Local',
  MESHERY: 'Meshery',
  // Extension Point: Add other providers as needed
};

setup('authenticate with Remote Provider', async ({ page }) => {
  const baseURL = ENV.MESHERY_SERVER_URL;
  const token = ENV.PROVIDER_TOKEN;
  const email = ENV.REMOTE_PROVIDER_USER.email;
  const password = ENV.REMOTE_PROVIDER_USER.password;
  const loginPage = new LoginPage(page);

  // Fail, do not skip. Playwright collapses a dependent project when its setup
  // *fails* - it does not when the setup *skips*, because a skip is a normal
  // outcome and the dependents stay scheduled. They then all die on the storage
  // state file this setup never wrote. Both shapes are on record for the same
  // defect: run 31039121068 (failing setup) reported 1 failure and ran 0
  // chromium-meshery-provider tests; run 31701664917 (skipping setup) reported
  // 62 failures and 23 skips, every one of them "ENOENT ... user-meshery-
  // provider.json" rather than anything about authentication.
  //
  // Skipping is also unnecessary. Only `push` selects this project
  // (npm test:e2e:ci:full); pull_request runs test:e2e:ci:local, so fork PRs -
  // the case the skip was added to protect - never reach this setup at all.
  if (!token && !(email && password)) {
    throw new Error(
      'No remote-provider credentials. Set PROVIDER_TOKEN, or REMOTE_PROVIDER_USER_EMAIL and ' +
        'REMOTE_PROVIDER_USER_PASSWORD, to run the chromium-meshery-provider project. In CI this ' +
        'comes from the REMOTE_PROVIDER_TEST_USER_TOKEN org secret (.github/workflows/test-e2e.yml).',
    );
  }

  if (token) {
    console.log('Using token-based authentication');
    await loginPage.loginWithToken(token, baseURL);
  } else {
    console.log('Using form-based authentication');
    const providerSelectionPage = new ProviderSelectionPage(page);
    await providerSelectionPage.navigateToProviderSelection();
    await providerSelectionPage.selectProvider(PROVIDERS.MESHERY);
    await loginPage.loginWithEmail(email, password);
  }

  await loginPage.waitForRedirection();

  await page.context().storageState({ path: ENV.AUTHFILEMESHERYPROVIDER });
});
