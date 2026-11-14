/**
 * Tutorial-aware environment recovery: pure helpers (no I/O, no React).
 *
 * Compares the resources a saved Meshery design (and optionally a tutorial) expects with
 * the live state reported by MeshSync, and turns the difference into a status, the
 * affected tutorial steps and recovery suggestions. Fetching live state and redeploying
 * the design stay with the existing endpoints; this module only interprets their results.
 * See https://github.com/meshery/meshery/issues/22880
 */

export type ClusterReachability = 'reachable' | 'unreachable' | 'unknown';

export interface ExpectedResource {
  kind: string;
  name: string;
  /** Omitted for cluster-scoped resources, which then match in any namespace. */
  namespace?: string;
}

/** The subset of a MeshSync resource this module reads. */
export interface LiveResource {
  kind?: string;
  metadata?: { name?: string; namespace?: string };
  status?: { phase?: string; attribute?: { phase?: string } } | null;
}

export interface TutorialStep {
  id: string;
  title: string;
  /** Resources this step creates or depends on. */
  requires: ExpectedResource[];
}

export type ResourceState = 'available' | 'missing' | 'unhealthy';

export interface ResourceCheck {
  resource: ExpectedResource;
  state: ResourceState;
}

export type EnvironmentStatus =
  | 'environment-unavailable' // the cluster itself cannot be reached
  // Every expected resource is present and not known to be unhealthy. This means "present",
  // not "running": a resource with no reported phase, or a Pod still Pending, counts as
  // available. Only explicitly failing phases (see UNHEALTHY_PHASES) are flagged.
  | 'ready'
  | 'partial' // some, but not all, expected resources are present and healthy
  | 'empty'; // the cluster is reachable but none of the expected resources are usable

export interface ValidationReport {
  status: EnvironmentStatus;
  checks: ResourceCheck[];
  /** Resources that are missing or unhealthy. */
  problems: ResourceCheck[];
  affectedSteps: TutorialStep[];
}

export type RecoveryAction =
  | 'wait-for-environment'
  | 'redeploy-design'
  | 'review-unhealthy'
  | 'retry-validation';

export interface RecoverySuggestion {
  action: RecoveryAction;
  resources: ExpectedResource[];
  /** Redeploys change the cluster, so they must be confirmed by the user. */
  requiresConfirmation: boolean;
}

export type RecoveryOutcome =
  | 'recovered'
  | 'partially-recovered'
  | 'recovery-failed'
  | 'environment-unavailable';

const UNHEALTHY_PHASES = new Set(['failed', 'unknown', 'crashloopbackoff', 'error']);

const normalize = (value?: string) => (value ?? '').trim().toLowerCase();

// DNS-1123 subdomain: the name format most Kubernetes kinds require.
const RESOURCE_NAME_PATTERN = /^[a-z0-9]([-a-z0-9.]*[a-z0-9])?$/;

const isValidResourceName = (value?: string): value is string =>
  !!value && value.length <= 253 && RESOURCE_NAME_PATTERN.test(value);

const keyOf = (r: ExpectedResource) => `${normalize(r.kind)}/${normalize(r.namespace)}/${r.name}`;

/**
 * Reads the expected resources out of a saved design. Components without a kind or a name
 * (for example purely visual elements) are skipped. Duplicates are collapsed.
 */
export function extractExpectedResources(design: unknown): ExpectedResource[] {
  const components = (design as { components?: unknown[] } | null | undefined)?.components;
  if (!Array.isArray(components)) return [];

  const seen = new Set<string>();
  const result: ExpectedResource[] = [];

  for (const raw of components) {
    const component = raw as {
      displayName?: string;
      component?: { kind?: string };
      configuration?: { metadata?: { name?: string; namespace?: string } };
    };
    const kind = component?.component?.kind;
    const metadata = component?.configuration?.metadata;
    // displayName is a human label, not necessarily the object name, so it is only used as a
    // fallback when it is a valid Kubernetes resource name. Otherwise the component is skipped
    // rather than invented into an expected resource that could never match live state.
    const displayName = component?.displayName;
    const name = metadata?.name || (isValidResourceName(displayName) ? displayName : undefined);
    if (!kind || !name) continue;

    const resource: ExpectedResource = { kind, name };
    if (metadata?.namespace) resource.namespace = metadata.namespace;

    const key = keyOf(resource);
    if (seen.has(key)) continue;
    seen.add(key);
    result.push(resource);
  }
  return result;
}

const matches = (expected: ExpectedResource, live: LiveResource) => {
  if (normalize(live.kind) !== normalize(expected.kind)) return false;
  if (live.metadata?.name !== expected.name) return false;
  return (
    !expected.namespace || normalize(live.metadata?.namespace) === normalize(expected.namespace)
  );
};

/**
 * True only when the resource reports an explicitly failing phase. A missing phase or a
 * transient one such as Pending is not treated as unhealthy (see `EnvironmentStatus`).
 */
export function isUnhealthy(live: LiveResource): boolean {
  const phase = live.status?.phase ?? live.status?.attribute?.phase;
  return !!phase && UNHEALTHY_PHASES.has(normalize(phase));
}

/** Classifies each expected resource as available, missing or unhealthy. */
export function checkResources(
  expected: ExpectedResource[],
  live: LiveResource[],
): ResourceCheck[] {
  return expected.map((resource) => {
    const found = live.filter((l) => matches(resource, l));
    if (found.length === 0) return { resource, state: 'missing' as const };
    // Healthy if any matching object is healthy: one stale duplicate should not flag the resource.
    return {
      resource,
      state: found.some((l) => !isUnhealthy(l)) ? 'available' : 'unhealthy',
    } as ResourceCheck;
  });
}

/**
 * Builds the validation report. An unreachable cluster is reported as such and is never
 * mistaken for individual missing resources (the live list is meaningless in that case).
 */
export function validateEnvironment(input: {
  reachability: ClusterReachability;
  expected: ExpectedResource[];
  live: LiveResource[];
  steps?: TutorialStep[];
}): ValidationReport {
  const { reachability, expected, live, steps = [] } = input;

  if (reachability === 'unreachable') {
    return { status: 'environment-unavailable', checks: [], problems: [], affectedSteps: [] };
  }

  const checks = checkResources(expected, live);
  const problems = checks.filter((c) => c.state !== 'available');

  // A step is affected when any resource it requires is a problem. A requirement without a
  // namespace matches the problem resource in any namespace.
  const affectedSteps = steps.filter((step) =>
    step.requires.some((req) =>
      problems.some(
        (p) =>
          normalize(p.resource.kind) === normalize(req.kind) &&
          p.resource.name === req.name &&
          (!req.namespace || normalize(p.resource.namespace) === normalize(req.namespace)),
      ),
    ),
  );

  let status: EnvironmentStatus;
  if (problems.length === 0) status = 'ready';
  else if (problems.length === checks.length) status = 'empty';
  else status = 'partial';

  return { status, checks, problems, affectedSteps };
}

/** Explicit, user-controlled recovery suggestions for a report. */
export function suggestRecovery(report: ValidationReport): RecoverySuggestion[] {
  if (report.status === 'environment-unavailable') {
    return [
      { action: 'wait-for-environment', resources: [], requiresConfirmation: false },
      { action: 'retry-validation', resources: [], requiresConfirmation: false },
    ];
  }
  if (report.status === 'ready') return [];

  const missing = report.problems.filter((p) => p.state === 'missing').map((p) => p.resource);
  const unhealthy = report.problems.filter((p) => p.state === 'unhealthy').map((p) => p.resource);

  const suggestions: RecoverySuggestion[] = [];
  if (missing.length > 0) {
    suggestions.push({ action: 'redeploy-design', resources: missing, requiresConfirmation: true });
  }
  if (unhealthy.length > 0) {
    // Never redeploy over something that exists: ask the user to review it first.
    suggestions.push({
      action: 'review-unhealthy',
      resources: unhealthy,
      requiresConfirmation: false,
    });
  }
  suggestions.push({ action: 'retry-validation', resources: [], requiresConfirmation: false });
  return suggestions;
}

/**
 * Judges a recovery attempt by validating again afterwards.
 * `before` is the report that prompted recovery, `after` the re-validation.
 */
export function evaluateRecovery(input: {
  before: ValidationReport;
  after: ValidationReport;
}): RecoveryOutcome {
  const { before, after } = input;

  if (after.status === 'environment-unavailable') return 'environment-unavailable';
  if (after.status === 'ready') return 'recovered';
  // An unavailable `before` report carries no problems by construction, so there is no
  // count to compare against. If the cluster is back and some expected resources exist,
  // that is progress. If it is back but nothing expected exists ('empty'), nothing was
  // recovered, so that falls through to 'recovery-failed'.
  if (before.status === 'environment-unavailable') {
    return after.status === 'partial' ? 'partially-recovered' : 'recovery-failed';
  }
  // Otherwise progress is fewer problems than before.
  if (after.problems.length < before.problems.length) return 'partially-recovered';
  return 'recovery-failed';
}
