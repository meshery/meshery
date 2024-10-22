import { expect, test as base } from '@playwright/test';
import { EnvironmentPage } from './fixtures/environmentsPage';
import { ENV } from './env';

export const test = base.extend({
  environmentPage: async ({ page }, use) => {
    const envPage = new EnvironmentPage(page);
    await envPage.navigate();
    await use(envPage);
  },
});

test.describe('Lifecycle Environments Tests', () => {
  test.beforeAll(async () => {
    console.log('Setting up lifecycle environments for tests...');
  });

  const environmentName = 'Sample-playwright-test';

  test('Create a Environment', async ({ environmentPage }) => {
    await environmentPage.createEnvironmentProfile(environmentName);
    await expect(environmentPage.page).toHaveURL(
      `${ENV.MESHERY_SERVER_URL}/management/environments`,
    );
  });

  test('Assign Connections', async ({ environmentPage }) => {
    await environmentPage.assignConnections(environmentName);
    await expect(environmentPage.page).toHaveURL(
      `${ENV.MESHERY_SERVER_URL}/management/environments`,
    );
  });

  test('Move Connections From Assigned To Available', async ({ environmentPage }) => {
    await environmentPage.moveAssignedConnectionsToAvailable(environmentName);
    await expect(environmentPage.page).toHaveURL(
      `${ENV.MESHERY_SERVER_URL}/management/environments`,
    );
  });

  test('Assign Environment To Connection', async ({ environmentPage }) => {
    await environmentPage.goToConnections();
    await environmentPage.addEnvironmentsToConnection(environmentName);
    await expect(environmentPage.page).toHaveURL(
      `${ENV.MESHERY_SERVER_URL}/management/connections`,
    );
  });

  test('Edit Environment Card', async ({ environmentPage }) => {
    await environmentPage.editEnvironmentCard(environmentName);
    await expect(environmentPage.page).toHaveURL(
      `${ENV.MESHERY_SERVER_URL}/management/environments`,
    );
  });

  test('Delete Environment profile', async ({ environmentPage }) => {
    await environmentPage.deleteEnvironmentCard(environmentName);
    await expect(environmentPage.page).toHaveURL(
      `${ENV.MESHERY_SERVER_URL}/management/environments`,
    );
  });

  test.afterAll(async () => {
    console.log('Cleaning up test environment...');
  });
});
