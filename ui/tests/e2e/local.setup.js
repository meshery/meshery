import { test as setup } from './fixtures/project';
import { ProviderSelectionPage } from './pages/ProviderSelectionPage';
import { waitForAuthRedirection } from './pages/LoginPage';
import { ENV } from './env';

const PROVIDERS = {
  LOCAL: 'None',
  MESHERY: 'Layer5',
};

// if project is none provider then only select the none provider and save the auth state
setup('authenticate as None provider', async ({ page, provider }) => {
  console.log(`Provider in auth setup: ${provider}`);

  // Perform authentication steps. Replace these actions with your own.

  const providerSelectionPage = new ProviderSelectionPage(page);
  await providerSelectionPage.navigateToProviderSelection();
  await providerSelectionPage.selectProvider(PROVIDERS.LOCAL);

  await waitForAuthRedirection(page);

  await page.context().storageState({ path: ENV.AUTHFILELOCALPROVIDER });
});
