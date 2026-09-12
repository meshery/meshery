import { describe, expect, it } from 'vitest';
import { groupWorkloadByVersion, getUnformattedName, findWorkloadByName } from '../workloadFilter';

const makeWorkload = (overrides: Record<string, unknown> = {}) => ({
  workload: {
    oam_definition: {
      metadata: { name: 'WorkloadA' },
      spec: { metadata: { meshVersion: 'v1.0.0' } },
    },
    ...overrides,
  },
});

describe('groupWorkloadByVersion', () => {
  it('groups workloads by meshVersion', () => {
    const w1 = makeWorkload();
    const w2 = {
      workload: {
        oam_definition: {
          metadata: { name: 'B' },
          spec: { metadata: { meshVersion: 'v1.0.0' } },
        },
      },
    };
    const w3 = {
      workload: {
        oam_definition: {
          metadata: { name: 'C' },
          spec: { metadata: { meshVersion: 'v2.0.0' } },
        },
      },
    };
    const result = groupWorkloadByVersion([w1, w2, w3]);
    expect(Object.keys(result)).toEqual(expect.arrayContaining(['v1.0.0', 'v2.0.0']));
    expect(result['v1.0.0']).toHaveLength(2);
    expect(result['v2.0.0']).toHaveLength(1);
  });

  it('falls back to spec.metadata.version when meshVersion is missing', () => {
    const w = {
      workload: {
        oam_definition: {
          metadata: { name: 'A' },
          spec: { metadata: { version: '0.9.0' } },
        },
      },
    };
    const result = groupWorkloadByVersion([w]);
    expect(result['0.9.0']).toHaveLength(1);
  });

  it('falls back to "Meshery" when both meshVersion and version are missing', () => {
    const w = {
      workload: {
        oam_definition: { metadata: { name: 'A' }, spec: { metadata: {} } },
      },
    };
    const result = groupWorkloadByVersion([w]);
    expect(result.Meshery).toHaveLength(1);
  });

  it('falls back to "Meshery" when oam_definition is missing entirely', () => {
    const w = { workload: {} };
    const result = groupWorkloadByVersion([w]);
    expect(result.Meshery).toHaveLength(1);
  });

  it('returns an empty object for an empty workload list', () => {
    expect(groupWorkloadByVersion([])).toEqual({});
  });
});

describe('getUnformattedName', () => {
  it('returns the name from the oam_definition metadata', () => {
    expect(getUnformattedName(makeWorkload())).toBe('WorkloadA');
  });

  it('returns "Un-Named" when the name is missing', () => {
    const w = {
      workload: { oam_definition: { metadata: {}, spec: { metadata: {} } } },
    };
    expect(getUnformattedName(w)).toBe('Un-Named');
  });

  it('returns "Un-Named" when oam_definition is missing', () => {
    expect(getUnformattedName({ workload: {} })).toBe('Un-Named');
  });
});

describe('findWorkloadByName', () => {
  it('returns the workload whose unformatted name matches', () => {
    const w1 = makeWorkload();
    const w2 = {
      workload: {
        oam_definition: {
          metadata: { name: 'WorkloadB' },
          spec: { metadata: { meshVersion: 'v1.0.0' } },
        },
      },
    };
    expect(findWorkloadByName('WorkloadB', [w1, w2])).toBe(w2);
  });

  it('returns undefined when no match is found', () => {
    expect(findWorkloadByName('missing', [makeWorkload()])).toBeUndefined();
  });

  it('returns undefined when workloads is null/undefined', () => {
    expect(findWorkloadByName('whatever', null)).toBeUndefined();
    expect(findWorkloadByName('whatever', undefined)).toBeUndefined();
  });
});
