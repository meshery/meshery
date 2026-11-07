import { describe, expect, it } from 'vitest';
import {
  ExpectedResource,
  LiveResource,
  TutorialStep,
  checkResources,
  evaluateRecovery,
  extractExpectedResources,
  suggestRecovery,
  validateEnvironment,
} from './environmentRecovery';

const ns: ExpectedResource = { kind: 'Namespace', name: 'demo' };
const deploy: ExpectedResource = { kind: 'Deployment', name: 'web', namespace: 'demo' };
const svc: ExpectedResource = { kind: 'Service', name: 'web', namespace: 'demo' };
const expected = [ns, deploy, svc];

const live = (kind: string, name: string, namespace?: string, phase?: string): LiveResource => ({
  kind,
  metadata: { name, namespace },
  status: phase ? { phase } : null,
});

const steps: TutorialStep[] = [
  { id: 'ns', title: 'Create namespace', requires: [ns] },
  { id: 'deploy', title: 'Deployment configuration', requires: [deploy] },
  { id: 'svc', title: 'Service creation', requires: [svc] },
  { id: 'verify', title: 'Application verification', requires: [deploy, svc] },
];

describe('extractExpectedResources', () => {
  it('reads kind, name and namespace from design components and dedupes', () => {
    const design = {
      components: [
        {
          displayName: 'ignored',
          component: { kind: 'Deployment' },
          configuration: { metadata: { name: 'web', namespace: 'demo' } },
        },
        {
          component: { kind: 'Deployment' },
          configuration: { metadata: { name: 'web', namespace: 'demo' } },
        },
        { displayName: 'demo', component: { kind: 'Namespace' } },
        { displayName: 'no-kind', component: {} },
        { component: { kind: 'Service' } },
      ],
    };
    expect(extractExpectedResources(design)).toEqual([deploy, ns]);
  });

  it('only falls back to displayName when it is a valid resource name', () => {
    const component = (displayName: string) => ({ displayName, component: { kind: 'Service' } });
    expect(
      extractExpectedResources({
        components: [
          component('My Web Service'),
          component('Web'),
          component('web_svc'),
          component('-web'),
          component('a'.repeat(254)),
          component('web-svc.v1'),
        ],
      }),
    ).toEqual([{ kind: 'Service', name: 'web-svc.v1' }]);
  });

  it('prefers metadata.name over a valid displayName', () => {
    const design = {
      components: [
        {
          displayName: 'label',
          component: { kind: 'Service' },
          configuration: { metadata: { name: 'real-name' } },
        },
      ],
    };
    expect(extractExpectedResources(design)).toEqual([{ kind: 'Service', name: 'real-name' }]);
  });

  it.each([null, undefined, {}, { components: 'x' }])('returns [] for %j', (design) => {
    expect(extractExpectedResources(design)).toEqual([]);
  });
});

describe('checkResources', () => {
  it('classifies available, missing and unhealthy', () => {
    const checks = checkResources(expected, [
      live('Namespace', 'demo'),
      live('Deployment', 'web', 'demo', 'Failed'),
    ]);
    expect(checks.map((c) => c.state)).toEqual(['available', 'unhealthy', 'missing']);
  });

  it('matches kind case-insensitively and cluster-scoped resources in any namespace', () => {
    expect(checkResources([ns], [live('namespace', 'demo', 'whatever')])[0].state).toBe(
      'available',
    );
  });

  it('does not match the same name in another namespace', () => {
    expect(checkResources([deploy], [live('Deployment', 'web', 'other')])[0].state).toBe('missing');
  });

  it('treats a resource as available when any duplicate is healthy', () => {
    const checks = checkResources(
      [deploy],
      [live('Deployment', 'web', 'demo', 'Failed'), live('Deployment', 'web', 'demo', 'Running')],
    );
    expect(checks[0].state).toBe('available');
  });

  it('counts a present resource as available when it has no phase or is Pending', () => {
    const checks = checkResources(
      [deploy, { kind: 'Pod', name: 'p', namespace: 'demo' }],
      [live('Deployment', 'web', 'demo'), live('Pod', 'p', 'demo', 'Pending')],
    );
    expect(checks.map((c) => c.state)).toEqual(['available', 'available']);
  });

  it('reads the phase from status.attribute too', () => {
    const r: LiveResource = {
      kind: 'Pod',
      metadata: { name: 'p', namespace: 'demo' },
      status: { attribute: { phase: 'Failed' } },
    };
    expect(checkResources([{ kind: 'Pod', name: 'p', namespace: 'demo' }], [r])[0].state).toBe(
      'unhealthy',
    );
  });
});

describe('validateEnvironment', () => {
  it('reports an unreachable cluster without blaming individual resources', () => {
    const report = validateEnvironment({ reachability: 'unreachable', expected, live: [], steps });
    expect(report.status).toBe('environment-unavailable');
    expect(report.problems).toEqual([]);
    expect(report.affectedSteps).toEqual([]);
  });

  it('is ready when everything is present', () => {
    const report = validateEnvironment({
      reachability: 'reachable',
      expected,
      live: [
        live('Namespace', 'demo'),
        live('Deployment', 'web', 'demo'),
        live('Service', 'web', 'demo'),
      ],
      steps,
    });
    expect(report.status).toBe('ready');
    expect(report.affectedSteps).toEqual([]);
  });

  it('is partial and lists the affected steps when only some resources are missing', () => {
    const report = validateEnvironment({
      reachability: 'reachable',
      expected,
      live: [live('Namespace', 'demo')],
      steps,
    });
    expect(report.status).toBe('partial');
    expect(report.affectedSteps.map((s) => s.id)).toEqual(['deploy', 'svc', 'verify']);
  });

  it('is empty when the cluster is reachable but nothing expected exists', () => {
    const report = validateEnvironment({ reachability: 'reachable', expected, live: [], steps });
    expect(report.status).toBe('empty');
    expect(report.affectedSteps).toHaveLength(4);
  });

  it('still validates when reachability is unknown', () => {
    const report = validateEnvironment({ reachability: 'unknown', expected: [ns], live: [] });
    expect(report.status).toBe('empty');
  });

  it('is ready for a design that expects nothing', () => {
    expect(validateEnvironment({ reachability: 'reachable', expected: [], live: [] }).status).toBe(
      'ready',
    );
  });
});

describe('suggestRecovery', () => {
  it('suggests waiting when the environment is unavailable', () => {
    const report = validateEnvironment({ reachability: 'unreachable', expected, live: [] });
    expect(suggestRecovery(report).map((s) => s.action)).toEqual([
      'wait-for-environment',
      'retry-validation',
    ]);
  });

  it('suggests nothing when ready', () => {
    const report = validateEnvironment({ reachability: 'reachable', expected: [], live: [] });
    expect(suggestRecovery(report)).toEqual([]);
  });

  it('requires confirmation to redeploy missing resources and never redeploys unhealthy ones', () => {
    const report = validateEnvironment({
      reachability: 'reachable',
      expected,
      live: [live('Namespace', 'demo'), live('Deployment', 'web', 'demo', 'Failed')],
    });
    const suggestions = suggestRecovery(report);
    const redeploy = suggestions.find((s) => s.action === 'redeploy-design');
    const review = suggestions.find((s) => s.action === 'review-unhealthy');
    expect(redeploy).toEqual({
      action: 'redeploy-design',
      resources: [svc],
      requiresConfirmation: true,
    });
    expect(review?.resources).toEqual([deploy]);
    expect(suggestions.at(-1)?.action).toBe('retry-validation');
  });
});

describe('evaluateRecovery', () => {
  const before = validateEnvironment({ reachability: 'reachable', expected, live: [] });
  const report = (
    liveResources: LiveResource[],
    reachability: 'reachable' | 'unreachable' = 'reachable',
  ) => validateEnvironment({ reachability, expected, live: liveResources });

  it('recovered when everything is back', () => {
    const after = report([
      live('Namespace', 'demo'),
      live('Deployment', 'web', 'demo'),
      live('Service', 'web', 'demo'),
    ]);
    expect(evaluateRecovery({ before, after })).toBe('recovered');
  });

  it('partially recovered when fewer problems remain', () => {
    expect(evaluateRecovery({ before, after: report([live('Namespace', 'demo')]) })).toBe(
      'partially-recovered',
    );
  });

  it('failed when nothing improved', () => {
    expect(evaluateRecovery({ before, after: report([]) })).toBe('recovery-failed');
  });

  describe('when the cluster was unreachable before recovery', () => {
    const unavailableBefore = validateEnvironment({
      reachability: 'unreachable',
      expected,
      live: [],
    });

    it('is partially recovered when the cluster is back and some resources exist', () => {
      expect(
        evaluateRecovery({ before: unavailableBefore, after: report([live('Namespace', 'demo')]) }),
      ).toBe('partially-recovered');
    });

    it('is recovered when everything is back', () => {
      const after = report([
        live('Namespace', 'demo'),
        live('Deployment', 'web', 'demo'),
        live('Service', 'web', 'demo'),
      ]);
      expect(evaluateRecovery({ before: unavailableBefore, after })).toBe('recovered');
    });

    it('failed when the cluster is back but nothing expected exists', () => {
      expect(evaluateRecovery({ before: unavailableBefore, after: report([]) })).toBe(
        'recovery-failed',
      );
    });

    it('stays environment unavailable when the cluster is still unreachable', () => {
      expect(
        evaluateRecovery({ before: unavailableBefore, after: report([], 'unreachable') }),
      ).toBe('environment-unavailable');
    });
  });

  it('environment unavailable when the cluster drops during recovery', () => {
    expect(evaluateRecovery({ before, after: report([], 'unreachable') })).toBe(
      'environment-unavailable',
    );
  });
});
