import { expect } from '@playwright/test';
import * as allure from 'allure-js-commons';
import { test } from './fixtures/project';
import { DashboardPage } from './pages/DashboardPage';
import {
  BUILT_IN,
  CONNECTIONS_API,
  SCHEMA_VERSION,
  SYSTEM_CONFIG_API,
  TEST_GROUP,
  annotateControllersCase as annotateCase,
  brokerGrid,
  chooseOption,
  clearServerDefaults,
  controllersTags as caseTags,
  fieldByLabel,
  fieldIn,
  getConnectionConfig,
  getServerDefaults,
  kubernetesConnections,
  meshsyncGrid,
  openControllersSettings,
  operatorGrid,
  putConnectionConfig,
  putServerDefaults,
  saveDefaults,
  waitForDefaultsLoad,
  type ConfigDoc,
} from './controllers-config.testmap';

// End-to-end coverage for the Operator, MeshSync & Broker settings - the layered
// `MesheryControllersConfig` document mapped row by row in
// docs/content/en/project/contributing/contributing-controllers-config.md.
//
// What is proven here, and where:
//
//   - the rendered server-wide editor (Settings -> Operator, MeshSync & Broker),
//     driven in a real browser: every setting present, storage of only the fields
//     the user set, the Inherit round-trip, the LoadBalancer-only service fields,
//     and what the form states about the deployment mode;
//   - the layering itself, driven through the two APIs the editors use:
//     GET/PUT `/api/system/controllers/config` (layer 2) and
//     GET/PUT `/api/integrations/connections/{id}/controllers/config` (layer 1),
//     which return `override`, `default` and `effective` separately.
//
// None of that needs a Kubernetes cluster: storage, layering, validation and
// rendering are all decided in Meshery Server and the browser. Only actual
// propagation to cluster objects does, and the one case that asserts it
// self-skips when no connected cluster is in reach - mirroring connections.spec.ts,
// so an infra-less run degrades to "skipped" and never to a false pass.
//
// Each case names the exact wire path(s) it covers, both as an inline comment and
// as a `wirePath` annotation, so the doc's settings tables can be reconciled
// against this spec mechanically.

/**
 * Every row of the doc's three settings tables: its label in the editor, its
 * wire path, and the section it lives in. The rendering case asserts each of
 * these is actually on the form, so a setting added to the document without an
 * editor control - or renamed in one place only - fails here.
 *
 * `loadBalancerOnly` rows are rendered only while the effective
 * `broker.service.type` is `LoadBalancer`; they have their own case below.
 */
const SETTING_ROWS: {
  label: string;
  path: string;
  section: 'operator' | 'meshsync' | 'broker';
  loadBalancerOnly?: true;
}[] = [
  { label: 'Deployment mode', path: 'operator.deploymentMode', section: 'operator' },
  { label: 'Operator version', path: 'operator.version', section: 'operator' },

  { label: 'MeshSync version', path: 'meshsync.version', section: 'meshsync' },
  { label: 'Replicas', path: 'meshsync.replicas', section: 'meshsync' },
  { label: 'Watched resources (discovery scope)', path: 'meshsync.watchList', section: 'meshsync' },
  { label: 'Output namespaces', path: 'meshsync.outputNamespaces', section: 'meshsync' },
  { label: 'Output resources', path: 'meshsync.outputResources', section: 'meshsync' },
  { label: 'Secret redaction', path: 'meshsync.redactSecrets', section: 'meshsync' },
  { label: 'Broker content dedup', path: 'meshsync.brokerContentDedup', section: 'meshsync' },
  { label: 'Debug logging', path: 'meshsync.debugLogging', section: 'meshsync' },

  { label: 'Broker version', path: 'broker.version', section: 'broker' },
  { label: 'Replicas', path: 'broker.replicas', section: 'broker' },
  { label: 'Service type', path: 'broker.service.type', section: 'broker' },
  { label: 'Service annotations', path: 'broker.service.annotations', section: 'broker' },
  {
    label: 'External endpoint override',
    path: 'broker.service.externalEndpointOverride',
    section: 'broker',
  },
  {
    label: 'Load balancer class',
    path: 'broker.service.loadBalancerClass',
    section: 'broker',
    loadBalancerOnly: true,
  },
  {
    label: 'Load balancer source ranges',
    path: 'broker.service.loadBalancerSourceRanges',
    section: 'broker',
    loadBalancerOnly: true,
  },
];

test.describe.serial('Operator, MeshSync & Broker Settings', () => {
  // Navigating the dashboard shell and mounting the settings page runs several
  // 120s visibility waits inside the page objects; the default BASE_TIMEOUT is
  // not enough headroom for that under load (same reasoning as connections.spec.ts).
  test.describe.configure({ timeout: 180_000 });

  test.beforeEach(async ({ page, request }) => {
    // One call here tags every case in this describe with its Test Plan Test
    // Group, so the Test-Group-keyed meshery/qa report picks the whole file up.
    await allure.label('testGroup', TEST_GROUP);

    // Each case starts from a cleared server-wide layer, so precedence
    // assertions are about what the case sets and not about leftovers.
    await clearServerDefaults(request);
    // `page` is required so Playwright builds the browser context for cases that
    // only use the API fixtures too - keeping the reset identical across cases.
    expect(page).toBeTruthy();
  });

  // `test.afterAll` receives only worker-scoped fixtures, so the test-scoped
  // `request` context is not available here - taking it would make this hook
  // throw and silently leave the server-wide defaults set for whatever runs
  // next. Build a request context explicitly instead.
  test.afterAll(async ({ playwright, baseURL }) => {
    // Leave the server as the run found it: the defaults are server-wide state
    // that outlives this spec and fans out to every tracked connection.
    const cleanup = await playwright.request.newContext({ baseURL });
    try {
      await clearServerDefaults(cleanup);
    } finally {
      await cleanup.dispose();
    }
  });

  // Covers: the entry point itself - Settings -> "Operator, MeshSync & Broker".
  // The docs point users here for the server-wide layer, so a tab that does not
  // open makes every setting below unreachable regardless of what it stores.
  test(
    'Reach the server-wide controllers editor from the header menu',
    { tag: caseTags('UI/Settings') },
    async ({ page }, testInfo) => {
      annotateCase(testInfo, {
        feature: 'Server-wide defaults editor',
        story: 'Open Settings -> Operator, MeshSync & Broker',
        wirePaths: [],
      });
      test.slow();

      const dashboardPage = new DashboardPage(page);
      await dashboardPage.navigateToDashboard();
      await dashboardPage.navigateToSettings();
      await page.waitForURL(/\/settings/);

      const loaded = waitForDefaultsLoad(page);
      await page.getByTestId('settings-tab-controllers').click();
      await loaded;

      await expect(page).toHaveURL(/settingsCategory=Controllers/);
      await expect(
        page.getByRole('heading', { name: 'Meshery Operator, MeshSync & Broker' }),
      ).toBeVisible();
      await expect(page.getByTestId('controllers-config-save')).toBeVisible();
    },
  );

  // Covers every unconditionally-rendered row of SETTING_ROWS - i.e. all of the
  // doc's three settings tables except the two LoadBalancer-only rows, which
  // have their own case.
  test(
    'Render a control for every setting in the three sections',
    { tag: caseTags('UI/Controllers Editor') },
    async ({ page }, testInfo) => {
      annotateCase(testInfo, {
        feature: 'Server-wide defaults editor',
        story: 'Every documented setting has a control',
        wirePaths: SETTING_ROWS.filter((row) => !row.loadBalancerOnly).map((row) => row.path),
      });

      await openControllersSettings(page);

      await expect(page.getByTestId('controllers-config-section-operator')).toContainText(
        'Meshery Operator',
      );
      await expect(page.getByTestId('controllers-config-section-meshsync')).toContainText(
        'MeshSync',
      );
      await expect(page.getByTestId('controllers-config-section-broker')).toContainText(
        'Meshery Broker',
      );

      const grids = {
        operator: operatorGrid(page),
        meshsync: meshsyncGrid(page),
        broker: brokerGrid(page),
      };

      for (const row of SETTING_ROWS) {
        if (row.loadBalancerOnly) continue;
        // The watch-scope control lives outside the MeshSync grid (it renders a
        // mode select plus either a whitelist table or a blacklist textarea), so
        // it is located by its unique label rather than within a section.
        const field =
          row.path === 'meshsync.watchList'
            ? fieldByLabel(page, row.label)
            : fieldIn(grids[row.section], row.label);
        await expect(field, `no control rendered for ${row.path} ("${row.label}")`).toHaveCount(1);
      }
    },
  );

  // Covers: broker.version, broker.replicas, broker.service.type
  //
  // The stored document must carry only the fields the user actually set - it is
  // the "merged" document that propagates to the cluster, and a field written
  // where the user asked to inherit pins that layer forever.
  test(
    'Store only the fields set through the server-wide editor',
    { tag: caseTags('UI/Controllers Editor') },
    async ({ page, request }, testInfo) => {
      annotateCase(testInfo, {
        feature: 'Server-wide defaults editor',
        story: 'Saving stores only explicitly-set fields',
        wirePaths: ['broker.version', 'broker.replicas', 'broker.service.type'],
      });

      await openControllersSettings(page);
      const broker = brokerGrid(page);

      await fieldIn(broker, 'Broker version').locator('input').fill('2.10.24');
      await fieldIn(broker, 'Replicas').locator('input').fill('3');
      await chooseOption(page, fieldIn(broker, 'Service type').getByRole('combobox'), 'NodePort');

      const stored = await saveDefaults(page);

      expect(stored).toEqual({
        schemaVersion: SCHEMA_VERSION,
        broker: { version: '2.10.24', replicas: 3, service: { type: 'NodePort' } },
      });
      // ...and it is what a fresh read returns, not just what the PUT echoed.
      expect(await getServerDefaults(request)).toEqual(stored);
    },
  );

  // Covers: meshsync.version, meshsync.replicas
  //
  // The inherit round-trip: a field set and then returned to Inherit must leave
  // the stored document entirely, so the next precedence layer applies again.
  // Clearing the last field of a section leaves no empty section behind.
  test(
    'Return a field to Inherit and it leaves the stored document',
    { tag: caseTags('UI/Controllers Editor') },
    async ({ page, request }, testInfo) => {
      annotateCase(testInfo, {
        feature: 'Server-wide defaults editor',
        story: 'Inherit round-trip removes the field from the document',
        wirePaths: ['meshsync.version', 'meshsync.replicas'],
      });

      await openControllersSettings(page);
      const meshsync = meshsyncGrid(page);
      const version = fieldIn(meshsync, 'MeshSync version').locator('input');
      const replicas = fieldIn(meshsync, 'Replicas').locator('input');

      await version.fill('v0.8.6');
      await replicas.fill('4');
      expect(await saveDefaults(page)).toEqual({
        schemaVersion: SCHEMA_VERSION,
        meshsync: { version: 'v0.8.6', replicas: 4 },
      });
      // The reloaded form shows both as overrides, not as Inherit.
      await expect(version).toHaveValue('v0.8.6');
      await expect(replicas).toHaveValue('4');

      // Emptying a control is how the editor expresses "Inherit".
      await replicas.fill('');
      expect(await saveDefaults(page)).toEqual({
        schemaVersion: SCHEMA_VERSION,
        meshsync: { version: 'v0.8.6' },
      });
      // Only the field returned to Inherit left; its sibling is untouched.
      await expect(replicas).toHaveValue('');
      await expect(version).toHaveValue('v0.8.6');

      await version.fill('');
      const cleared = await saveDefaults(page);
      expect(cleared).toEqual({ schemaVersion: SCHEMA_VERSION });
      expect(await getServerDefaults(request)).toEqual(cleared);

      // Back on Inherit, the placeholder states the value that now applies -
      // the built-in default for replicas.
      await expect(replicas).toHaveAttribute(
        'placeholder',
        `Inherit (${BUILT_IN.meshsyncReplicas})`,
      );
    },
  );

  // Covers: broker.service.type, broker.service.loadBalancerClass,
  //         broker.service.loadBalancerSourceRanges
  //
  // The two LoadBalancer-only fields are valid only while the effective service
  // type is LoadBalancer - the server rejects them otherwise. The form therefore
  // renders them only then, and clears them when the type moves away, so a stale
  // value cannot trip a server-side rejection the user can neither see nor clear.
  test(
    'Expose the LoadBalancer-only service fields only while the service type is LoadBalancer',
    { tag: caseTags('UI/Controllers Editor') },
    async ({ page }, testInfo) => {
      annotateCase(testInfo, {
        feature: 'Server-wide defaults editor',
        story: 'LoadBalancer-only fields follow broker.service.type',
        wirePaths: [
          'broker.service.type',
          'broker.service.loadBalancerClass',
          'broker.service.loadBalancerSourceRanges',
        ],
      });

      await openControllersSettings(page);
      const broker = brokerGrid(page);
      const serviceType = fieldIn(broker, 'Service type').getByRole('combobox');

      // Inheriting the built-in ClusterIP: neither field exists.
      await expect(fieldIn(broker, 'Load balancer class')).toHaveCount(0);
      await expect(fieldIn(broker, 'Load balancer source ranges')).toHaveCount(0);

      await chooseOption(page, serviceType, 'LoadBalancer');
      await expect(fieldIn(broker, 'Load balancer class')).toHaveCount(1);
      await expect(fieldIn(broker, 'Load balancer source ranges')).toHaveCount(1);

      await fieldIn(broker, 'Load balancer class').locator('input').fill('service.k8s.aws/nlb');
      await fieldIn(broker, 'Load balancer source ranges').locator('input').fill('10.0.0.0/8');

      expect(await saveDefaults(page)).toEqual({
        schemaVersion: SCHEMA_VERSION,
        broker: {
          service: {
            type: 'LoadBalancer',
            loadBalancerClass: 'service.k8s.aws/nlb',
            loadBalancerSourceRanges: ['10.0.0.0/8'],
          },
        },
      });

      // The reloaded form still has them, so what follows edits the saved state.
      await expect(fieldIn(broker, 'Load balancer class').locator('input')).toHaveValue(
        'service.k8s.aws/nlb',
      );

      // Moving off LoadBalancer hides both fields and drops their stored values,
      // rather than leaving a document the server would reject.
      await chooseOption(page, serviceType, 'ClusterIP (cluster-internal)');
      await expect(fieldIn(broker, 'Load balancer class')).toHaveCount(0);
      await expect(fieldIn(broker, 'Load balancer source ranges')).toHaveCount(0);

      expect(await saveDefaults(page)).toEqual({
        schemaVersion: SCHEMA_VERSION,
        broker: { service: { type: 'ClusterIP' } },
      });
    },
  );

  // Covers: operator.deploymentMode
  //
  // Mode gating on the server-wide layer. The mode decides which of the other
  // settings can reach anything, so the editor states it before the save rather
  // than letting the server report it afterwards as a skip. On this layer the
  // statement is conditional ("reach only connections running in Operator mode"),
  // not inert: a connection that overrides the mode to Operator uses every value
  // stored here.
  test(
    'State in Embedded mode which settings reach only Operator-mode connections',
    { tag: caseTags('UI/Controllers Editor') },
    async ({ page }, testInfo) => {
      annotateCase(testInfo, {
        feature: 'Deployment mode gating',
        story: 'Embedded mode states the reduced reach of MeshSync and Broker settings',
        wirePaths: ['operator.deploymentMode'],
      });

      await openControllersSettings(page);

      // The built-in default mode is embedded, so the statements are present
      // before anything is touched.
      const banner = page.getByTestId('controllers-config-mode-banner');
      await expect(banner).toContainText('Default deployment mode: Embedded (in Meshery Server)');
      await expect(banner).toContainText('Set by the built-in default.');

      await expect(
        page.getByText(
          'MeshSync version, replicas, watch scope and pod environment settings reach only connections running in Operator mode.',
        ),
      ).toBeVisible();
      await expect(
        page.getByText('Meshery Broker settings reach only connections running in Operator mode.'),
      ).toBeVisible();

      // Choosing embedded explicitly says the same thing, now attributed to the
      // server-wide default rather than to the built-in one.
      await chooseOption(
        page,
        fieldByLabel(page, 'Deployment mode').getByRole('combobox'),
        'Embedded (in Meshery Server)',
      );
      await expect(banner).toContainText('Default deployment mode: Embedded (in Meshery Server)');
      await expect(banner).toContainText('Set by this server-wide default.');
      await expect(
        page.getByText('Meshery Broker settings reach only connections running in Operator mode.'),
      ).toBeVisible();

      expect(await saveDefaults(page)).toEqual({
        schemaVersion: SCHEMA_VERSION,
        operator: { deploymentMode: 'embedded' },
      });
    },
  );

  // Covers: operator.deploymentMode
  //
  // The other half of the gating contract: in Operator mode every setting on the
  // form applies, so the statements about reduced reach must be gone. A notice
  // that survives the mode change is as wrong as one that never appears.
  test(
    'Drop the reduced-reach statements in Operator mode',
    { tag: caseTags('UI/Controllers Editor') },
    async ({ page }, testInfo) => {
      annotateCase(testInfo, {
        feature: 'Deployment mode gating',
        story: 'Operator mode removes the reduced-reach statements',
        wirePaths: ['operator.deploymentMode'],
      });

      await openControllersSettings(page);

      await chooseOption(
        page,
        fieldByLabel(page, 'Deployment mode').getByRole('combobox'),
        'Operator (in-cluster)',
      );

      const banner = page.getByTestId('controllers-config-mode-banner');
      await expect(banner).toContainText('Default deployment mode: Operator (in-cluster)');
      await expect(banner).toContainText(
        'Connections that do not override the mode run Meshery Operator, which applies every setting below.',
      );

      await expect(
        page.getByText(
          'MeshSync version, replicas, watch scope and pod environment settings reach only connections running in Operator mode.',
        ),
      ).toHaveCount(0);
      await expect(
        page.getByText('Meshery Broker settings reach only connections running in Operator mode.'),
      ).toHaveCount(0);

      expect(await saveDefaults(page)).toEqual({
        schemaVersion: SCHEMA_VERSION,
        operator: { deploymentMode: 'operator' },
      });
    },
  );

  // Covers all seventeen wire paths of SETTING_ROWS at the storage boundary: a
  // document carrying every setting must survive the write/read round-trip byte
  // for byte, with no field silently dropped and no default silently injected.
  test(
    'Round-trip every documented wire path through the server-wide defaults',
    { tag: caseTags('API/System Controllers Config') },
    async ({ request }, testInfo) => {
      annotateCase(testInfo, {
        feature: 'Layered document storage',
        story: 'Every setting survives write and read unchanged',
        wirePaths: SETTING_ROWS.map((row) => row.path),
      });

      const full = {
        operator: { deploymentMode: 'operator', version: 'v0.7.9' },
        meshsync: {
          version: 'v0.8.6',
          replicas: 2,
          watchList: { blacklist: ['secrets.v1.', 'events.v1.'] },
          outputNamespaces: ['kube-system', 'default'],
          outputResources: ['pod', 'deployment'],
          redactSecrets: true,
          brokerContentDedup: true,
          debugLogging: true,
        },
        broker: {
          version: '2.10.24',
          replicas: 3,
          service: {
            type: 'LoadBalancer',
            annotations: { 'service.beta.kubernetes.io/aws-load-balancer-internal': 'true' },
            loadBalancerClass: 'service.k8s.aws/nlb',
            loadBalancerSourceRanges: ['10.0.0.0/8', '192.168.0.0/16'],
            externalEndpointOverride: 'broker.example.com:4222',
          },
        },
      };

      const expected = { schemaVersion: SCHEMA_VERSION, ...full };
      expect(await putServerDefaults(request, full)).toEqual(expected);
      expect(await getServerDefaults(request)).toEqual(expected);
    },
  );

  // Covers: operator.deploymentMode, meshsync.replicas, broker.replicas,
  //         broker.service.type
  //
  // Precedence: per-connection override (layer 1) over server-wide default
  // (layer 2) over built-in (layer 3), field by field. Every field layers
  // independently, so a field left unset at layer 1 still resolves through
  // layer 2 while a sibling is overridden.
  test(
    'Resolve each field through override, then server-wide default, then built-in',
    { tag: caseTags('API/Connection Controllers Config') },
    async ({ request }, testInfo) => {
      annotateCase(testInfo, {
        feature: 'Layered precedence',
        story: 'Override beats server-wide default beats built-in',
        wirePaths: [
          'operator.deploymentMode',
          'meshsync.replicas',
          'broker.replicas',
          'broker.service.type',
        ],
      });

      const connections = await kubernetesConnections(request);
      if (connections.length === 0) {
        test.skip(true, 'No Kubernetes connection registered to carry a per-connection override.');
        return;
      }
      const connectionId = connections[0].id;

      try {
        // Layer 3 only: nothing set anywhere.
        //
        // deploymentMode is deliberately excluded here. Unlike every other
        // setting, it has a layer BENEATH the document that the client cannot
        // see - the server's MESHSYNC_DEFAULT_DEPLOYMENT_MODE - so the resolved
        // value with nothing set is the server's env default, not the compiled
        // built-in. Asserting the built-in made this test pass in CI (where the
        // variable is unset) and fail for anyone running the repo's own
        // `make server`, which sets it to `operator` by default. The mode's
        // full precedence, including that env layer, is covered by
        // TestB3ServerWideDeploymentModeReachesInheritingConnection.
        const builtIn = await getConnectionConfig(request, connectionId);
        expect(builtIn.override).toBeUndefined();
        expect(builtIn.effective).toMatchObject({
          meshsync: { replicas: BUILT_IN.meshsyncReplicas },
          broker: {
            replicas: BUILT_IN.brokerReplicas,
            service: { type: BUILT_IN.brokerServiceType },
          },
        });
        expect(['operator', 'embedded']).toContain(builtIn.effective.operator?.deploymentMode);

        // Layer 2: the server-wide default reaches a connection that overrides
        // nothing.
        await putServerDefaults(request, {
          meshsync: { replicas: 5 },
          broker: { replicas: 3, service: { type: 'NodePort' } },
        });
        const defaulted = await getConnectionConfig(request, connectionId);
        expect(defaulted.override).toBeUndefined();
        expect(defaulted.default).toMatchObject({
          meshsync: { replicas: 5 },
          broker: { replicas: 3, service: { type: 'NodePort' } },
        });
        expect(defaulted.effective).toMatchObject({
          meshsync: { replicas: 5 },
          broker: { replicas: 3, service: { type: 'NodePort' } },
        });
        // deploymentMode is excluded for the same reason as above: with neither
        // editable layer setting it, the resolved value is the server's
        // MESHSYNC_DEFAULT_DEPLOYMENT_MODE, which the client cannot observe.
        expect(['operator', 'embedded']).toContain(defaulted.effective.operator?.deploymentMode);

        // Layer 1: the override wins for the fields it sets, and only those.
        const overridden = await putConnectionConfig(request, connectionId, {
          broker: { replicas: 7 },
        });
        expect(overridden.override).toEqual({
          schemaVersion: SCHEMA_VERSION,
          broker: { replicas: 7 },
        });
        expect(overridden.effective).toMatchObject({
          broker: {
            replicas: 7,
            // Not overridden - still the server-wide default.
            service: { type: 'NodePort' },
          },
          // Not overridden - still the server-wide default.
          meshsync: { replicas: 5 },
        });
      } finally {
        await putConnectionConfig(request, connectionId, {});
      }
    },
  );

  // Covers: meshsync.outputNamespaces, meshsync.outputResources
  //
  // The inherit round-trip at layer 1: clearing a per-connection override must
  // remove the field from the stored override document, so the server-wide
  // default applies to that connection again. Collections merge whole, so the
  // override replaces the default's list rather than being unioned with it.
  test(
    'Clear a per-connection override and the server-wide default applies again',
    { tag: caseTags('API/Connection Controllers Config') },
    async ({ request }, testInfo) => {
      annotateCase(testInfo, {
        feature: 'Layered precedence',
        story: 'Clearing an override restores the server-wide default',
        wirePaths: ['meshsync.outputNamespaces', 'meshsync.outputResources'],
      });

      const connections = await kubernetesConnections(request);
      if (connections.length === 0) {
        test.skip(true, 'No Kubernetes connection registered to carry a per-connection override.');
        return;
      }
      const connectionId = connections[0].id;

      try {
        await putServerDefaults(request, {
          meshsync: {
            outputNamespaces: ['kube-system', 'default'],
            outputResources: ['pod', 'deployment'],
          },
        });

        // The override replaces the whole collection - not element-wise, which
        // would produce a scope neither layer asked for.
        const overridden = await putConnectionConfig(request, connectionId, {
          meshsync: { outputNamespaces: ['meshery'] },
        });
        expect(overridden.effective).toMatchObject({
          meshsync: {
            outputNamespaces: ['meshery'],
            outputResources: ['pod', 'deployment'],
          },
        });

        // Back to Inherit: the field leaves the override document entirely, and
        // with nothing left in it the layer itself stops being reported - an
        // empty override and no override are the same state.
        const cleared = await putConnectionConfig(request, connectionId, {});
        expect(cleared.override).toBeUndefined();
        expect(cleared.effective).toMatchObject({
          meshsync: {
            outputNamespaces: ['kube-system', 'default'],
            outputResources: ['pod', 'deployment'],
          },
        });
      } finally {
        await putConnectionConfig(request, connectionId, {});
      }
    },
  );

  // Covers: meshsync.replicas, broker.replicas, meshsync.watchList,
  //         broker.service.loadBalancerClass, broker.service.loadBalancerSourceRanges,
  //         broker.service.type, operator.deploymentMode
  //
  // Every guardrail rejected with 400 and nothing persisted. A rejection that
  // half-writes is worse than no validation: it leaves a stored document the
  // editor cannot show and the cluster would be asked to apply.
  test(
    'Reject invalid documents and persist nothing',
    { tag: caseTags('API/System Controllers Config') },
    async ({ request }, testInfo) => {
      annotateCase(testInfo, {
        feature: 'Validation',
        story: 'Guardrails reject with 400 and leave the stored document untouched',
        wirePaths: [
          'meshsync.replicas',
          'broker.replicas',
          'meshsync.watchList',
          'broker.service.loadBalancerClass',
          'broker.service.loadBalancerSourceRanges',
          'broker.service.type',
          'operator.deploymentMode',
        ],
      });

      // A known-good document to prove nothing is clobbered by a rejection.
      const good = { meshsync: { replicas: 2 } };
      const stored = await putServerDefaults(request, good);

      const rejections: { why: string; body: ConfigDoc; because: string }[] = [
        {
          why: 'meshsync.replicas above the allowed range',
          body: { meshsync: { replicas: 11 } },
          because: 'meshsync.replicas must be between 1 and 10',
        },
        {
          why: 'meshsync.replicas below the allowed range',
          body: { meshsync: { replicas: 0 } },
          because: 'meshsync.replicas must be between 1 and 10',
        },
        {
          why: 'broker.replicas above the allowed range',
          body: { broker: { replicas: 11 } },
          because: 'broker.replicas must be between 1 and 10',
        },
        {
          why: 'meshsync.watchList setting both a whitelist and a blacklist',
          body: {
            meshsync: {
              watchList: {
                whitelist: [{ resource: 'pods.v1.', events: ['ADDED'] }],
                blacklist: ['events.v1.'],
              },
            },
          },
          because: 'whitelist and blacklist are mutually exclusive',
        },
        {
          why: 'broker.service.loadBalancerClass without broker.service.type LoadBalancer',
          body: { broker: { service: { type: 'ClusterIP', loadBalancerClass: 'nlb' } } },
          because:
            'broker.service.loadBalancerClass is only valid when broker.service.type is LoadBalancer',
        },
        {
          why: 'broker.service.loadBalancerSourceRanges without broker.service.type LoadBalancer',
          body: {
            broker: { service: { type: 'ClusterIP', loadBalancerSourceRanges: ['10.0.0.0/8'] } },
          },
          because:
            'broker.service.loadBalancerSourceRanges is only valid when broker.service.type is LoadBalancer',
        },
        {
          why: 'an unknown operator.deploymentMode',
          body: { operator: { deploymentMode: 'in-cluster-ish' } },
          because: 'operator.deploymentMode must be either "operator" or "embedded"',
        },
      ];

      for (const rejection of rejections) {
        const response = await request.put(SYSTEM_CONFIG_API, { data: rejection.body });
        const text = await response.text();
        expect(response.status(), `${rejection.why} was not rejected - body: ${text}`).toBe(400);
        // The rule is named in the MeshKit error's longDescription. Asserting on
        // the parsed field rather than the raw body keeps the expectation in the
        // server's own words - the wire form escapes the quotes around the
        // allowed deploymentMode values.
        const failed: { longDescription?: string[] } = JSON.parse(text);
        expect(
          (failed.longDescription ?? []).join('\n'),
          `${rejection.why} was rejected without naming the rule`,
        ).toContain(rejection.because);
        // Nothing was written: the previously stored document still stands.
        expect(
          await getServerDefaults(request),
          `${rejection.why} clobbered the stored document`,
        ).toEqual(stored);
      }
    },
  );

  // Covers: meshsync.watchList, broker.service.annotations
  //
  // The only case here that needs a reachable cluster: the per-connection PUT
  // applies the resolved document to that connection's cluster, and a watch-scope
  // change additionally rolls the MeshSync pods. Everything the browser and the
  // API can decide without a cluster is asserted above; this asserts that the
  // apply against a live connection succeeds and leaves the connection connected.
  // Self-skips when no connected Kubernetes cluster is in reach, so an
  // infra-less run degrades to "skipped" rather than to a false pass.
  test(
    'Apply a watch-scope and broker service override to a connected cluster',
    { tag: caseTags('API/Connection Controllers Config') },
    async ({ page, request }, testInfo) => {
      annotateCase(testInfo, {
        feature: 'Cluster propagation',
        story: 'A per-connection override applies to a connected cluster',
        wirePaths: ['meshsync.watchList', 'broker.service.annotations'],
      });
      test.slow();

      const connected = (await kubernetesConnections(request)).filter(
        (connection) => connection.status === 'connected',
      );
      if (connected.length === 0) {
        test.skip(true, 'No connected Kubernetes cluster in reach; skipping cluster propagation.');
        return;
      }
      const connection = connected[0];

      try {
        const applied = await putConnectionConfig(request, connection.id, {
          meshsync: { watchList: { blacklist: ['events.v1.'] } },
          broker: {
            service: { annotations: { 'meshery.io/managed-by': 'controllers-config-e2e' } },
          },
        });
        expect(applied.effective).toMatchObject({
          meshsync: { watchList: { blacklist: ['events.v1.'] } },
          broker: {
            service: { annotations: { 'meshery.io/managed-by': 'controllers-config-e2e' } },
          },
        });

        // The apply is asynchronous with respect to the cluster and must not
        // knock the connection out of `connected` - a configuration change is
        // not a lifecycle transition.
        const after = await request.get(`${CONNECTIONS_API}/${connection.id}`);
        expect(after.status()).toBe(200);
        expect((await after.json()).status).toBe('connected');

        // And the editor's own surface still loads against that connection.
        await openControllersSettings(page);
      } finally {
        await putConnectionConfig(request, connection.id, {});
      }
    },
  );
});
