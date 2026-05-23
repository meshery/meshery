import { test as setup } from './fixtures/project';
import { ProviderSelectionPage } from './pages/ProviderSelectionPage';
import { waitForAuthRedirection } from './pages/LoginPage';
import { ENV } from './env';

const PROVIDERS = {
  LOCAL: 'Local',
  MESHERY: 'Meshery',
};

// Authenticate as the built-in Local provider and persist the auth state.
setup('authenticate as Local provider', async ({ page, provider }) => {
  // Increase timeout for authentication setup (default 60s → 120s)
  setup.setTimeout(120000);

  console.log(`Provider in auth setup: ${provider}`);

  // Perform authentication steps. Replace these actions with your own.

  const providerSelectionPage = new ProviderSelectionPage(page);
  await providerSelectionPage.navigateToProviderSelection();
  await providerSelectionPage.selectProvider(PROVIDERS.LOCAL);

  await waitForAuthRedirection(page);

  await page.context().storageState({ path: ENV.AUTHFILELOCALPROVIDER });
});
