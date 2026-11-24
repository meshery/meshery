// Which Meshery Operator / MeshSync / Broker settings can actually take effect,
// and why.
//
// Meshery Operator manages MeshSync and Meshery Broker. In `embedded` mode
// Meshery Server runs MeshSync in-process and installs nothing into the
// cluster: there is no Meshery Operator release, no MeshSync Deployment, no
// Broker and no `meshery.io` custom resources. Meshery Server skips every
// setting that targets one of those and reports it in
// `ControllersConfigApplyResult.Skipped` (server/models/controllers_config_apply.go).
//
// This module is the single client-side statement of that dependency structure,
// so the editors cannot offer a control that reaches nothing without saying so.

import { getPath, type FieldPath } from './fieldPath';
import type { ControllersConfigDoc } from './ControllersConfigForm';

export type DeploymentMode = 'operator' | 'embedded';

export const DEPLOYMENT_MODE_LABEL: Record<DeploymentMode, string> = {
  operator: 'Operator (in-cluster)',
  embedded: 'Embedded (in Meshery Server)',
};

/** The compiled-in default, mirrored from `connections.BuiltInControllersConfig()`. */
export const BUILT_IN_DEPLOYMENT_MODE: DeploymentMode = 'embedded';

export type ConfigSection = 'operator' | 'meshsync' | 'broker';

// The settings that still reach something when MeshSync runs embedded:
//
//   - operator.deploymentMode is what selects the mode itself;
//   - the output filters are passed to the in-process libmeshsync run
//     (models.meshsyncDataHandlersStartLibMeshsyncRun), which restarts on save.
//
// Everything else configures an in-cluster object that embedded mode never
// creates. Note that secret redaction, broker content dedup and debug logging
// are environment variables on the MeshSync Deployment: embedded MeshSync takes
// them from the Meshery Server process environment instead, so setting them
// here changes nothing in that mode.
const EMBEDDED_LIVE_PATHS = new Set([
  'operator.deploymentMode',
  'meshsync.outputNamespaces',
  'meshsync.outputResources',
]);

/** True when a setting at `path` can take effect in `mode`. */
export const takesEffectIn = (mode: DeploymentMode, path: FieldPath): boolean =>
  mode === 'operator' || EMBEDDED_LIVE_PATHS.has(path.join('.'));

/**
 * Why a section's settings cannot take effect in embedded mode, in the words
 * shown to the user. Stated in the form itself rather than in a tooltip: a
 * control that reaches nothing has to say so where it is read.
 */
export const EMBEDDED_INERT_REASON: Record<ConfigSection, string> = {
  operator:
    'Embedded mode installs no Meshery Operator release on this cluster, so the operator chart version is not applied.',
  meshsync:
    'Embedded mode runs MeshSync inside Meshery Server, so there is no MeshSync Deployment or custom resource to configure. Only the output filters below take effect; secret redaction, broker content dedup and debug logging follow the Meshery Server process environment in this mode.',
  broker:
    'Meshery Broker is deployed by Meshery Operator. Embedded mode installs no Broker on this cluster, so nothing in this section is applied.',
};

/**
 * How a section's settings relate to the server-wide default mode. The
 * server-wide layer configures every managed cluster, and a connection may
 * override the mode, so nothing here is inert - it is conditional.
 */
export const SERVER_DEFAULT_SCOPE_NOTE: Record<ConfigSection, string> = {
  operator: '',
  meshsync:
    'MeshSync version, replicas, watch scope and pod environment settings reach only connections running in Operator mode. The output filters apply in both modes.',
  broker: 'Meshery Broker settings reach only connections running in Operator mode.',
};

/** Every editable path of a section, so dormant stored values can be found and cleared. */
export const SECTION_PATHS: Record<ConfigSection, FieldPath[]> = {
  operator: [
    ['operator', 'deploymentMode'],
    ['operator', 'version'],
  ],
  meshsync: [
    ['meshsync', 'version'],
    ['meshsync', 'replicas'],
    ['meshsync', 'redactSecrets'],
    ['meshsync', 'brokerContentDedup'],
    ['meshsync', 'debugLogging'],
    ['meshsync', 'outputNamespaces'],
    ['meshsync', 'outputResources'],
    ['meshsync', 'watchList'],
  ],
  broker: [
    ['broker', 'version'],
    ['broker', 'replicas'],
    ['broker', 'service', 'type'],
    ['broker', 'service', 'externalEndpointOverride'],
    ['broker', 'service', 'loadBalancerClass'],
    ['broker', 'service', 'loadBalancerSourceRanges'],
    ['broker', 'service', 'annotations'],
  ],
};

/**
 * The deployment mode that governs an editor, and where it came from.
 *
 * `scope` decides how the form treats settings that cannot take effect in
 * `mode`: on a connection the mode is a fact about one cluster, so they are
 * rendered inert; on the server-wide defaults the mode is only what inheriting
 * connections get, and a connection that overrides to Operator mode still uses
 * every value stored here, so they stay editable and are annotated instead.
 */
export type DeploymentModeGovernance = {
  mode: DeploymentMode;
  /** Where `mode` comes from, phrased for the user. */
  sourceLabel: string;
  /** True when `mode` is an unsaved draft edit rather than the persisted state. */
  unsaved: boolean;
  scope: 'connection' | 'serverDefault';
};

const modeOf = (doc?: ControllersConfigDoc | null): DeploymentMode | undefined =>
  (doc?.operator?.deploymentMode as DeploymentMode | undefined) ?? undefined;

export type ConnectionDeploymentModeInput = {
  /** The override currently being edited. */
  draft: ControllersConfigDoc;
  /** The override as persisted, to tell an unsaved mode change from the current one. */
  persistedOverride?: ControllersConfigDoc | null;
  /** The server-wide default document. */
  serverDefault?: ControllersConfigDoc | null;
  /**
   * The effective document the server returned. Its `operator.deploymentMode`
   * is the mode the connection actually runs, resolved through layers the
   * client cannot see (the materialized `meshsync_deployment_mode` cache and
   * `MESHSYNC_DEFAULT_DEPLOYMENT_MODE`).
   */
  serverEffective?: ControllersConfigDoc | null;
};

/**
 * Resolves the mode governing the per-connection editor, mirroring the server's
 * precedence (`connections.ResolveDeploymentMode`) over the layers the client
 * holds: the draft override, then the server-wide default, then the mode the
 * server reports the connection running.
 *
 * The last of those is why the server's effective document has to name the real
 * mode: when neither editable layer sets one, it is the only thing that
 * distinguishes a cluster running Meshery Operator from one that is not.
 */
export const connectionDeploymentMode = ({
  draft,
  persistedOverride,
  serverDefault,
  serverEffective,
}: ConnectionDeploymentModeInput): DeploymentModeGovernance => {
  const draftMode = modeOf(draft);
  const persistedMode = modeOf(persistedOverride);
  const defaultMode = modeOf(serverDefault);
  const effectiveMode = modeOf(serverEffective) ?? BUILT_IN_DEPLOYMENT_MODE;
  const unsaved = draftMode !== persistedMode;

  if (draftMode) {
    return {
      mode: draftMode,
      sourceLabel: "this connection's override",
      unsaved,
      scope: 'connection',
    };
  }
  if (defaultMode) {
    return {
      mode: defaultMode,
      sourceLabel: 'the server-wide default',
      unsaved,
      scope: 'connection',
    };
  }
  // Neither editable layer sets a mode. The server's effective mode is the
  // truth while nothing is unsaved; once the user clears an override the new
  // mode is whatever Meshery Server resolves on save, and saying so is more
  // honest than guessing at a floor the client cannot see.
  return {
    mode: effectiveMode,
    sourceLabel: unsaved
      ? 'the Meshery Server default, resolved when you save'
      : 'the Meshery Server default',
    unsaved,
    scope: 'connection',
  };
};

/** Resolves the mode governing the server-wide defaults editor. */
export const serverDefaultDeploymentMode = (
  draft: ControllersConfigDoc,
): DeploymentModeGovernance => {
  const draftMode = modeOf(draft);
  return {
    mode: draftMode ?? BUILT_IN_DEPLOYMENT_MODE,
    sourceLabel: draftMode ? 'this server-wide default' : 'the built-in default',
    unsaved: false,
    scope: 'serverDefault',
  };
};

/**
 * True when a setting cannot take effect and the form must render it inert.
 *
 * Only the per-connection editor does this: there the mode is a fact about one
 * cluster. On the server-wide defaults the mode is only what inheriting
 * connections get, and a connection that overrides it to Operator mode uses
 * every value stored there, so nothing on that layer is dead.
 */
export const isInertIn = (
  governance: DeploymentModeGovernance | undefined,
  path: FieldPath,
): boolean => governance?.scope === 'connection' && !takesEffectIn(governance.mode, path);

/**
 * The settings of a section that carry a stored value the current mode cannot
 * apply. They are kept and shown rather than dropped - the user chose them, and
 * they become live again in Operator mode - but the form names them as dormant
 * and offers to clear them.
 */
export const dormantPathsIn = (
  governance: DeploymentModeGovernance | undefined,
  doc: ControllersConfigDoc,
  section: ConfigSection,
): FieldPath[] =>
  SECTION_PATHS[section].filter(
    (path) => isInertIn(governance, path) && getPath(doc, path) !== undefined,
  );
