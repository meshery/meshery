import { test as setup } from './fixtures/project';
import { ProviderSelectionPage } from './pages/ProviderSelectionPage';
import { LoginPage } from './pages/LoginPage';
import { ENV } from './env';

const PROVIDERS = {
  LOCAL: 'None',
  MESHERY: 'Layer5',
}

setup('authenticate as Meshery provider', async ({ page }) => {
  // Perform authentication steps. Replace these actions with your own.
  const providerSelectionPage = new ProviderSelectionPage(page);
  await providerSelectionPage.navigateToProviderSelection();
  await providerSelectionPage.selectProvider(PROVIDERS.MESHERY);

  const loginPage = new LoginPage(page);
  await loginPage.login(ENV.REMOTE_PROVIDER_USER.email, ENV.REMOTE_PROVIDER_USER.password);

  await loginPage.waitForRedirection();

  await page.context().storageState({ path: ENV.AUTHFILEMESHERYPROVIDER });
});

setup('authenticate as None provider', async ({ page }) => {
  // Perform authentication steps. Replace these actions with your own.

  const providerSelectionPage = new ProviderSelectionPage(page);
  await providerSelectionPage.navigateToProviderSelection();
  await providerSelectionPage.selectProvider(PROVIDERS.LOCAL);

  const loginPage = new LoginPage(page);
  await loginPage.waitForRedirection();

  await page.context().storageState({ path: ENV.AUTHFILELOCALPROVIDER });
});
