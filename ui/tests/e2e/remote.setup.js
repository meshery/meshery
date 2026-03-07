import { test as setup } from './fixtures/project';
import { ProviderSelectionPage } from './pages/ProviderSelectionPage';
import { LoginPage } from './pages/LoginPage';
import { ENV } from './env';

const PROVIDERS = {
  LOCAL: 'None',
  MESHERY: 'Layer5',
  // Extension Point: Add other providers as needed
};

setup('authenticate with Remote Provider', async ({ page }) => {
  const baseURL = ENV.MESHERY_SERVER_URL;
  const token = ENV.PROVIDER_TOKEN;
  const email = ENV.REMOTE_PROVIDER_USER.email;
  const password = ENV.REMOTE_PROVIDER_USER.password;
  const loginPage = new LoginPage(page);

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
