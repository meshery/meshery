import {
  expect,
  type APIRequestContext,
  type Locator,
  type Page,
  type TestInfo,
} from '@playwright/test';

/**
 * Shared contract and drivers for the Operator, MeshSync & Broker Playwright
 * suite (controllers-config.spec.ts).
 *
 * Two things live here, both of which the spec would otherwise repeat:
 *
 *  - the traceability contract - the Test Plan Test Group every case in the
 *    spec is keyed to, and the Allure labels each case emits. This mirrors
 *    connections.testmap.ts, which plays the same role for the Connection
 *    Lifecycle report;
 *  - the drivers for the surface under test: the two layered-configuration
 *    endpoints the editors use, and the locators for the rendered editor.
 *
 * The settings rows themselves deliberately stay in the spec: the wire path of
 * every row in
 * docs/content/en/project/contributing/contributing-controllers-config.md has
 * to be named there for the spec to be reconcilable against the doc on its own.
 */

// --- Traceability -----------------------------------------------------------

/**
 * Test Plan Test Group (col B) for every case in the spec. Emitted once as an
 * Allure `testGroup` label in the describe's beforeEach - the same mechanism
 * the Connection Lifecycle report uses (see connections.spec.ts and
 * mesheryctl/bats-to-allure.js). It must go through `allure.label` from
 * allure-js-commons: allure-playwright maps only epic/feature/story annotations
 * to labels, so a custom label set as a Playwright annotation never reaches the
 * report.
 */
export const TEST_GROUP = 'Operator, MeshSync & Broker Settings';

/** Allure epic shared by every case in this suite. */
export const CONTROLLERS_EPIC = 'Operator, MeshSync & Broker';

/** Client label for the shared cross-lane contract (UI vs CLI). */
export const CONTROLLERS_CLIENT = 'UI';

/**
 * Emit the shared label contract for a case: the wire paths it covers plus the
 * epic/feature/story grouping. The wire paths are the reconciliation key
 * between the doc's settings tables and the spec - a row whose path is named by
 * no case has no end-to-end coverage.
 */
export function annotateControllersCase(
  testInfo: TestInfo,
  { feature, story, wirePaths }: { feature: string; story: string; wirePaths: string[] },
): void {
  testInfo.annotations.push(
    { type: 'epic', description: CONTROLLERS_EPIC },
    { type: 'client', description: CONTROLLERS_CLIENT },
    { type: 'feature', description: feature },
    { type: 'story', description: story },
  );
  for (const path of wirePaths) {
    testInfo.annotations.push({ type: 'wirePath', description: path });
  }
}

/** Grep-able Playwright tags, mirroring connections.testmap.ts's `connTags`. */
export const controllersTags = (componentUnderTest: string): string[] => [
  '@controllers-config',
  `@cut:${componentUnderTest
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)/g, '')}`,
  `@client:${CONTROLLERS_CLIENT.toLowerCase()}`,
];

// --- The layered-configuration APIs -----------------------------------------

export const SYSTEM_CONFIG_API = '/api/system/controllers/config';
export const CONNECTIONS_API = '/api/integrations/connections';

/** The schema version every stored document carries back. */
export const SCHEMA_VERSION = 'controllers.meshery.io/v1alpha1';

/**
 * The compiled-in defaults, mirrored from `connections.BuiltInControllersConfig()`
 * and from `BUILT_IN_CONTROLLERS_CONFIG` in
 * ui/components/configuration/ControllersConfigForm.tsx. This is precedence
 * layer 3: what `effective` resolves to when neither editable layer sets a field.
 */
export const BUILT_IN = {
  deploymentMode: 'embedded',
  meshsyncReplicas: 1,
  brokerReplicas: 1,
  brokerServiceType: 'ClusterIP',
};

export type ConfigDoc = Record<string, unknown>;

/** The layered document for one connection: `override`, `default` and `effective`. */
export type LayeredConfig = {
  override?: ConfigDoc;
  default?: ConfigDoc;
  effective: ConfigDoc;
};

export async function getServerDefaults(request: APIRequestContext): Promise<ConfigDoc> {
  const response = await request.get(SYSTEM_CONFIG_API);
  expect(
    response.status(),
    `GET ${SYSTEM_CONFIG_API} returned ${response.status()}: ${await response.text()}`,
  ).toBe(200);
  return response.json();
}

export async function putServerDefaults(
  request: APIRequestContext,
  body: ConfigDoc,
): Promise<ConfigDoc> {
  const response = await request.put(SYSTEM_CONFIG_API, { data: body });
  expect(
    response.status(),
    `PUT ${SYSTEM_CONFIG_API} returned ${response.status()}: ${await response.text()}`,
  ).toBe(200);
  return response.json();
}

/** Clear the server-wide layer so each case starts from the built-in defaults. */
export const clearServerDefaults = (request: APIRequestContext) => putServerDefaults(request, {});

const connectionConfigApi = (connectionId: string) =>
  `${CONNECTIONS_API}/${connectionId}/controllers/config`;

export async function getConnectionConfig(
  request: APIRequestContext,
  connectionId: string,
): Promise<LayeredConfig> {
  const response = await request.get(connectionConfigApi(connectionId));
  expect(
    response.status(),
    `GET connection controllers config returned ${response.status()}: ${await response.text()}`,
  ).toBe(200);
  return response.json();
}

export async function putConnectionConfig(
  request: APIRequestContext,
  connectionId: string,
  body: ConfigDoc,
): Promise<LayeredConfig> {
  const response = await request.put(connectionConfigApi(connectionId), { data: body });
  expect(
    response.status(),
    `PUT connection controllers config returned ${response.status()}: ${await response.text()}`,
  ).toBe(200);
  return response.json();
}

export type KubernetesConnection = { id: string; name: string; status: string };

/**
 * Kubernetes connections known to this server. The per-connection layer is a row
 * in the connections table, so cases that exercise it need a registered
 * connection but not a reachable cluster - a connection in any state carries an
 * override. Only the propagation case additionally requires `connected`.
 */
export async function kubernetesConnections(
  request: APIRequestContext,
): Promise<KubernetesConnection[]> {
  const response = await request.get(`${CONNECTIONS_API}?page=0&pagesize=100`);
  if (response.status() !== 200) return [];
  const body = await response.json();
  const connections: Record<string, string>[] = body?.connections ?? [];
  return connections
    .filter((connection) => connection.kind === 'kubernetes')
    .map((connection) => ({
      id: connection.id,
      name: connection.name,
      status: connection.status,
    }));
}

// --- Locating controls in the rendered editor -------------------------------
//
// Field labels are Typography (a bare `<p>`), not `<label for>`, so they cannot
// be reached through getByLabel. Each label sits in a Box next to its control
// inside one grid item, and every field of a section shares one grid container -
// which is what tells the MeshSync "Replicas" field from the Broker one without
// depending on DOM order.

/** Quote a label for use as an XPath string literal. */
const xpathLiteral = (value: string): string => `"${value}"`;

export const fieldByLabel = (page: Page, label: string): Locator =>
  page.locator(`xpath=//p[normalize-space(text())=${xpathLiteral(label)}]/../..`);

/** The grid container holding every field of the section `uniqueLabel` is in. */
const sectionGrid = (page: Page, uniqueLabel: string): Locator =>
  fieldByLabel(page, uniqueLabel).locator('xpath=..');

export const fieldIn = (grid: Locator, label: string): Locator =>
  grid.locator(`xpath=.//p[normalize-space(text())=${xpathLiteral(label)}]/../..`);

export const operatorGrid = (page: Page) => sectionGrid(page, 'Operator version');
export const meshsyncGrid = (page: Page) => sectionGrid(page, 'MeshSync version');
export const brokerGrid = (page: Page) => sectionGrid(page, 'Broker version');

/** Pick an option from a Sistent/MUI select by its exact visible text. */
export async function chooseOption(page: Page, control: Locator, option: string): Promise<void> {
  await control.click();
  await page.getByRole('option', { name: option, exact: true }).click();
  // The listbox is a portal that intercepts pointer events until it unmounts.
  await expect(page.getByRole('listbox')).toHaveCount(0);
}

/** Wait for the server-wide defaults GET that hydrates the editor. */
export const waitForDefaultsLoad = (page: Page) =>
  page.waitForResponse(
    (response) =>
      response.request().method() === 'GET' &&
      response.url().includes(SYSTEM_CONFIG_API) &&
      response.status() === 200,
  );

/** Wait for the save PUT the editor issues. */
const waitForDefaultsSave = (page: Page) =>
  page.waitForResponse(
    (response) =>
      response.request().method() === 'PUT' && response.url().includes(SYSTEM_CONFIG_API),
  );

/** Deep-link straight to Settings -> Operator, MeshSync & Broker. */
export async function openControllersSettings(page: Page): Promise<void> {
  const loaded = waitForDefaultsLoad(page);
  await page.goto('/settings?settingsCategory=Controllers', { waitUntil: 'domcontentloaded' });
  await expect(
    page.getByRole('heading', { name: 'Meshery Operator, MeshSync & Broker' }),
  ).toBeVisible();
  await loaded;
  // The mode banner is the first thing the form renders, so its presence means
  // the editor is mounted with a resolved governing mode rather than mid-hydration.
  await expect(page.getByTestId('controllers-config-mode-banner')).toBeVisible();
}

/**
 * Save the editor and return the document the server stored.
 *
 * The save does not settle when the PUT resolves. `useControllersConfigDraft`
 * clears `dirty` at that point, which re-seeds the draft from the RTK query
 * cache - still holding the pre-save document - and the editor briefly shows
 * every field back on Inherit. The mutation invalidates
 * `Meshery_Controllers_Configuration_controllers`, so the refetch that follows
 * is what puts the saved values back on the form. Editing before it lands edits
 * a form that is about to be overwritten, and a `fill('')` into an already-blank
 * control then produces no change event at all - leaving the Save button
 * disabled with nothing to explain why. Waiting for the refetch is the real
 * settle point.
 */
export async function saveDefaults(page: Page): Promise<ConfigDoc> {
  const saved = waitForDefaultsSave(page);
  const resettled = waitForDefaultsLoad(page);
  await page.getByTestId('controllers-config-save').click();
  const response = await saved;
  expect(
    response.status(),
    `save PUT returned ${response.status()} ${response.statusText()} - body: ${await response.text()}`,
  ).toBe(200);
  // A test that saves more than once accumulates success toasts, and the
  // previous one need not have dismissed before the next arrives. An unscoped
  // getByText then matches both and fails on strict mode - not because the save
  // did not succeed, but because it succeeded twice. Assert the first match:
  // the presence of the confirmation is the signal, not its cardinality.
  await expect(
    page.getByText('Server-wide controllers configuration defaults saved').first(),
  ).toBeVisible();
  await resettled;
  return response.json();
}
