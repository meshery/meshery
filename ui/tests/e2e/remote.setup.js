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

  // No credentials is an environment fact, not a defect: forks and PRs from
  // forks cannot read REMOTE_PROVIDER_TOKEN, and a contributor running the
  // suite locally has neither. Failing here made every remote-provider run
  // report a hard failure that said "Email is required for login", which reads
  // as a broken test rather than an unconfigured environment - and once the
  // E2E job can fail the build, it would fail every fork PR for a reason the
  // contributor cannot fix. Skip cleanly instead, so the whole
  // chromium-meshery-provider project reports "skipped".
  if (!token && !(email && password)) {
    setup.skip(
      true,
      'No remote-provider credentials: set PROVIDER_TOKEN, or REMOTE_PROVIDER_USER_EMAIL and REMOTE_PROVIDER_USER_PASSWORD, to run the Meshery-provider project.',
    );
    return;
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
