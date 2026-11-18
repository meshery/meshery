import { describe, expect, it } from 'vitest';
import { analyzeDesignImpact } from './designImpactAnalysis';

describe('designImpactAnalysis', () => {
  it('returns safe empty result when design or target is missing', () => {
    const res1 = analyzeDesignImpact({ design: null, target: null });
    expect(res1.targetType).toBe('unknown');
    expect(res1.totalAffectedCount).toBe(0);

    const res2 = analyzeDesignImpact({ design: { components: [] }, target: null });
    expect(res2.targetType).toBe('unknown');
    expect(res2.totalAffectedCount).toBe(0);
  });

  it('handles a component with no relationships producing an empty impact result', () => {
    const design = {
      components: [
        {
          id: 'isolated-pod',
          name: 'isolated-pod',
          component: { kind: 'Pod' },
          model: { name: 'kubernetes' },
        },
      ],
      relationships: [],
    };

    const result = analyzeDesignImpact({
      design,
      target: 'isolated-pod',
      modificationType: 'update',
    });

    expect(result.targetType).toBe('component');
    expect(result.targetId).toBe('isolated-pod');
    expect(result.targetKind).toBe('Pod');
    expect(result.directlyAffected).toHaveLength(0);
    expect(result.indirectlyAffected).toHaveLength(0);
    expect(result.affectedRelationships).toHaveLength(0);
    expect(result.totalAffectedCount).toBe(0);
    expect(result.summary).toContain('no dependent components or relationships found');
  });

  it('identifies direct relationships of a modified component (one relationship)', () => {
    const design = {
      components: [
        {
          id: 'comp-svc',
          name: 'frontend-svc',
          component: { kind: 'Service' },
          model: { name: 'kubernetes' },
        },
        {
          id: 'comp-deploy',
          name: 'frontend-deploy',
          component: { kind: 'Deployment' },
          model: { name: 'kubernetes' },
        },
      ],
      relationships: [
        {
          id: 'rel-svc-deploy',
          kind: 'Edge',
          type: 'network',
          subType: 'ServiceToDeployment',
          selectors: [
            {
              allow: {
                from: [{ id: 'comp-svc' }],
                to: [{ id: 'comp-deploy' }],
              },
            },
          ],
        },
      ],
    };

    const result = analyzeDesignImpact({
      design,
      target: 'comp-svc',
      modificationType: 'update',
    });

    expect(result.targetType).toBe('component');
    expect(result.directlyAffected).toHaveLength(1);
    expect(result.directlyAffected[0].id).toBe('comp-deploy');
    expect(result.directlyAffected[0].name).toBe('frontend-deploy');
    expect(result.directlyAffected[0].kind).toBe('Deployment');
    expect(result.directlyAffected[0].model).toBe('kubernetes');
    expect(result.affectedRelationships).toHaveLength(1);
    expect(result.totalAffectedCount).toBe(1);
  });

  it('handles a component with multiple direct relationships', () => {
    const design = {
      components: [
        { id: 'comp-svc', name: 'my-svc', component: { kind: 'Service' } },
        { id: 'comp-deploy', name: 'my-deploy', component: { kind: 'Deployment' } },
        { id: 'comp-ingress', name: 'my-ingress', component: { kind: 'Ingress' } },
      ],
      relationships: [
        {
          id: 'rel-1',
          kind: 'Edge',
          type: 'network',
          selectors: [
            {
              allow: {
                from: [{ id: 'comp-svc' }],
                to: [{ id: 'comp-deploy' }],
              },
            },
          ],
        },
        {
          id: 'rel-2',
          kind: 'Edge',
          type: 'network',
          selectors: [
            {
              allow: {
                from: [{ id: 'comp-ingress' }],
                to: [{ id: 'comp-svc' }],
              },
            },
          ],
        },
      ],
    };

    const result = analyzeDesignImpact({
      design,
      target: 'comp-svc',
      modificationType: 'update',
    });

    expect(result.directlyAffected).toHaveLength(2);
    const affectedIds = result.directlyAffected.map((c) => c.id);
    expect(affectedIds).toContain('comp-deploy');
    expect(affectedIds).toContain('comp-ingress');
    expect(result.affectedRelationships).toHaveLength(2);
    expect(result.totalAffectedCount).toBe(2);
  });

  it('traverses multi-level relationship chains (indirect dependencies)', () => {
    // Ingress -> Service -> Deployment -> Pod -> PVC
    const design = {
      components: [
        { id: 'c-ingress', name: 'ingress', component: { kind: 'Ingress' } },
        { id: 'c-svc', name: 'service', component: { kind: 'Service' } },
        { id: 'c-deploy', name: 'deployment', component: { kind: 'Deployment' } },
        { id: 'c-pvc', name: 'storage', component: { kind: 'PersistentVolumeClaim' } },
      ],
      relationships: [
        {
          id: 'rel-ing-svc',
          selectors: [{ allow: { from: [{ id: 'c-ingress' }], to: [{ id: 'c-svc' }] } }],
        },
        {
          id: 'rel-svc-dep',
          selectors: [{ allow: { from: [{ id: 'c-svc' }], to: [{ id: 'c-deploy' }] } }],
        },
        {
          id: 'rel-dep-pvc',
          selectors: [{ allow: { from: [{ id: 'c-deploy' }], to: [{ id: 'c-pvc' }] } }],
        },
      ],
    };

    const result = analyzeDesignImpact({
      design,
      target: 'c-ingress',
      modificationType: 'update',
    });

    expect(result.directlyAffected).toHaveLength(1);
    expect(result.directlyAffected[0].id).toBe('c-svc');

    expect(result.indirectlyAffected).toHaveLength(2);
    const indirectIds = result.indirectlyAffected.map((c) => c.id);
    expect(indirectIds).toContain('c-deploy');
    expect(indirectIds).toContain('c-pvc');
    expect(result.totalAffectedCount).toBe(3);
  });

  it('resolves selector endpoints carrying both id and name to a single canonical identity', () => {
    const design = {
      components: [
        { id: 'comp-a', name: 'alpha', component: { kind: 'Service' } },
        { id: 'comp-b', name: 'beta', component: { kind: 'Deployment' } },
        { id: 'comp-c', name: 'gamma', component: { kind: 'Pod' } },
      ],
      relationships: [
        {
          id: 'rel-a-b',
          kind: 'Edge',
          selectors: [
            {
              allow: {
                from: [{ id: 'comp-a', name: 'alpha' }],
                // Endpoint has both id and name
                to: [{ id: 'comp-b', name: 'beta' }],
              },
            },
          ],
        },
        {
          id: 'rel-b-c',
          kind: 'Edge',
          selectors: [
            {
              allow: {
                // Next hop references B by name only, connecting seamlessly to gamma by ID
                from: ['beta'],
                to: [{ id: 'comp-c' }],
              },
            },
          ],
        },
      ],
    };

    const result = analyzeDesignImpact({
      design,
      target: 'comp-a',
      modificationType: 'update',
    });

    // comp-b must appear only once in directlyAffected, and comp-c in indirectlyAffected
    expect(result.directlyAffected).toHaveLength(1);
    expect(result.directlyAffected[0].id).toBe('comp-b');
    expect(result.indirectlyAffected).toHaveLength(1);
    expect(result.indirectlyAffected[0].id).toBe('comp-c');
    expect(result.totalAffectedCount).toBe(2);
  });

  it('covers component removal (modificationType: delete)', () => {
    const design = {
      components: [
        { id: 'c-vol', name: 'data-vol', component: { kind: 'PersistentVolume' } },
        { id: 'c-pvc', name: 'data-claim', component: { kind: 'PersistentVolumeClaim' } },
      ],
      relationships: [
        {
          id: 'rel-mount',
          kind: 'Edge',
          type: 'mount',
          selectors: [{ allow: { from: [{ id: 'c-vol' }], to: [{ id: 'c-pvc' }] } }],
        },
      ],
    };

    const result = analyzeDesignImpact({
      design,
      target: 'c-vol',
      modificationType: 'delete',
    });

    expect(result.modificationType).toBe('delete');
    expect(result.summary).toContain('Removing PersistentVolume/data-vol');
    expect(result.directlyAffected).toHaveLength(1);
    expect(result.directlyAffected[0].id).toBe('c-pvc');
  });

  it('covers relationship removal (target is a relationship)', () => {
    const design = {
      components: [
        { id: 'c1', name: 'cm', component: { kind: 'ConfigMap' } },
        { id: 'c2', name: 'deploy', component: { kind: 'Deployment' } },
      ],
      relationships: [
        {
          id: 'rel-binding',
          kind: 'Edge',
          type: 'reference',
          selectors: [{ allow: { from: [{ id: 'c1' }], to: [{ id: 'c2' }] } }],
        },
      ],
    };

    const result = analyzeDesignImpact({
      design,
      target: 'rel-binding',
      modificationType: 'delete',
    });

    expect(result.targetType).toBe('relationship');
    expect(result.targetId).toBe('rel-binding');
    expect(result.directlyAffected).toHaveLength(2);
    const ids = result.directlyAffected.map((c) => c.id);
    expect(ids).toContain('c1');
    expect(ids).toContain('c2');
    expect(result.summary).toContain('Removing relationship');
  });

  it('honors explicit targetType option when provided', () => {
    const design = {
      components: [{ id: 'target-rel-id', name: 'ambiguous', component: { kind: 'Service' } }],
      relationships: [{ id: 'target-rel-id', kind: 'Edge', selectors: [] }],
    };

    const resRel = analyzeDesignImpact({
      design,
      target: 'target-rel-id',
      targetType: 'relationship',
    });
    expect(resRel.targetType).toBe('relationship');

    const resComp = analyzeDesignImpact({
      design,
      target: 'target-rel-id',
      targetType: 'component',
    });
    expect(resComp.targetType).toBe('component');
  });
});
