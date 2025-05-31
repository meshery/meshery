import { expect, test } from './fixtures/project';

const resourceCategories = [
  { name: 'Node', testId: 'Node' },
  { name: 'Namespace', testId: 'Namespace' },
  { name: 'Workload', testId: 'Workload' },
  { name: 'Configuration', testId: 'Configuration' },
  { name: 'Security', testId: 'Security' },
  { name: 'Storage', testId: 'Storage' },
  { name: 'CRDS', testId: 'CRDS' },
];

const widgets = [
  { name: 'Cluster Resource Overview', testId: 'cluster-res-overview' },
  { name: 'Getting Started', testId: 'getting-started' },
  { name: 'Help Center', testId: 'help-center' },
  { name: 'My Recent Designs', testId: 'my-recent-designs' },
  { name: 'Workspace Activity', testId: 'workspace-activity' },
  { name: 'Kubernetes Cluster Status', testId: 'kubernetes-cluster-status' },
  { name: 'Latest Blogs', testId: 'latest-blogs' },
];

test.describe('Index Page UI Tests', () => {

  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('Left Navigation Panel is visible', async ({ page }) => {
    await expect(page.getByTestId('navigation')).toBeVisible();
  });

  test('Settings button is visible', async ({ page }) => {
    await expect(page.getByTestId('settings-button')).toBeVisible();
  });

  test('Notification button is visible', async ({ page }) => {
    await expect(page.getByTestId('notification-button')).toBeVisible();
  });

  test('Profile button is visible', async ({ page }) => {
    await expect(page.getByTestId('profile-button')).toBeVisible();
  });

  test('Getting Started modal opens on Start click', async ({ page }) => {
    const getStartedComponent =  page.getByTestId('getting-started');
    await getStartedComponent.toBeVisible();
    await getStartedComponent.getByRole('button', { name: 'Start' }).click();
    await expect(page.getByTestId('get-started-modal')).toBeVisible();
  });

  test('Connect Clusters button works', async ({ page }) => {
   const connectClustersComponent = page.getByTestId('connect-clusters');
  await expect(connectClustersComponent).toBeVisible();
  await connectClustersComponent.getByRole('button', { name: 'Connect Clusters' }).click();//or get by testi d
  await expect(page.getByText('Cluster Connection')).toBeVisible();
  });

  resourceCategories.forEach(({ name, testId }, index) => {
    test(`${name} resource category is visible and clickable`, async ({ page }) => {
      await expect(page.getByTestId(testId)).toBeVisible();
      await page.getByTestId(testId).click();
      await expect(page.getByTestId(`${name}-${index}`)).toBeVisible();

    });
  });

  widgets.forEach(({ name, testId }) => {
    test(`${name} widget is visible`, async ({ page }) => {
      await expect(page.getByTestId(testId)).toBeVisible();
    });
  });

});

/*
tests 
text before dashboard is displayed and can be clicked
overview workload bars are checked 
on clicking cluster overview multiple options are visible 
connect cluster button is working 
on clicking start left panel is displayed and closed 
all workspaces button is working 
see all designs is working 
kubernetes cluster status is shown
edit button is show and clickable
overview panel testing
 */