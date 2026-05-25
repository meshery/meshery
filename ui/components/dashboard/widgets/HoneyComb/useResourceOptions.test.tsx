import { renderHook } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('../../resources/config', () => ({
  ResourceMenuConfig: {
    Node: [],
    Workload: ['Pod', 'Deployment'],
    Configuration: ['ConfigMap'],
  },
  getAllCustomResourceDefinitionsKinds: (kinds: Array<{ Kind: string }>) =>
    kinds.filter((k) => k.Kind === 'ExoticCrd'),
}));

import {
  DEFAULT_GROUP_BY,
  SORT_DIRECTIONS,
  useResourceFiltering,
  useResourceOptions,
} from './useResourceOptions';

describe('useResourceOptions', () => {
  it('returns "all", filtered category entries (skipping empty arrays), and CRDs', () => {
    const { result } = renderHook(() => useResourceOptions());
    const values = result.current.map((o) => o.value);
    expect(values).toEqual(['all', 'workload', 'configuration', 'crds']);
    expect(result.current[0]).toEqual({ value: 'all', label: 'All Resources' });
    expect(result.current.at(-1)).toEqual({ value: 'crds', label: 'Custom Resources' });
  });
});

describe('useResourceFiltering', () => {
  const kinds = [
    { Kind: 'Pod', Count: 4 },
    { Kind: 'Deployment', Count: 2 },
    { Kind: 'ConfigMap', Count: 5 },
    { Kind: 'ExoticCrd', Count: 1 },
  ];

  it('returns an empty array when kinds is undefined', () => {
    const { result } = renderHook(() => useResourceFiltering(undefined, DEFAULT_GROUP_BY, null));
    expect(result.current).toEqual([]);
  });

  it('returns a copy of all kinds when grouping by "all"', () => {
    const { result } = renderHook(() => useResourceFiltering(kinds, DEFAULT_GROUP_BY, null));
    expect(result.current).toEqual(kinds);
    expect(result.current).not.toBe(kinds);
  });

  it('filters by category when groupBy matches a known category', () => {
    const { result } = renderHook(() => useResourceFiltering(kinds, 'workload', null));
    expect(result.current.map((k) => k.Kind)).toEqual(['Pod', 'Deployment']);
  });

  it('uses getAllCustomResourceDefinitionsKinds when groupBy is "crds"', () => {
    const { result } = renderHook(() => useResourceFiltering(kinds, 'crds', null));
    expect(result.current.map((k) => k.Kind)).toEqual(['ExoticCrd']);
  });

  it('sorts ascending by Count when sortDirection is "asc"', () => {
    const { result } = renderHook(() =>
      useResourceFiltering(kinds, DEFAULT_GROUP_BY, SORT_DIRECTIONS.ASC),
    );
    expect(result.current.map((k) => k.Count)).toEqual([1, 2, 4, 5]);
  });

  it('sorts descending by Count when sortDirection is "desc"', () => {
    const { result } = renderHook(() =>
      useResourceFiltering(kinds, DEFAULT_GROUP_BY, SORT_DIRECTIONS.DESC),
    );
    expect(result.current.map((k) => k.Count)).toEqual([5, 4, 2, 1]);
  });

  it('returns an empty array when no kinds match the category', () => {
    const { result } = renderHook(() => useResourceFiltering(kinds, 'unknown-cat', null));
    expect(result.current).toEqual([]);
  });

  it('treats missing Count as zero when sorting', () => {
    const sparse = [{ Kind: 'A' }, { Kind: 'B', Count: 3 }];
    const { result } = renderHook(() =>
      useResourceFiltering(sparse, DEFAULT_GROUP_BY, SORT_DIRECTIONS.ASC),
    );
    expect(result.current.map((k) => k.Kind)).toEqual(['A', 'B']);
  });
});
