import { describe, expect, it } from 'vitest';
import getMostRecentVersion, {
  versionSortComparatorFn,
  sortByVersionInDecreasingOrder,
  getGreaterVersion,
  sortAndGroupVersionsInModel,
} from '../versionSort';
import { WILDCARD_V } from '../hooks/useMeshModelComponents';

describe('getMostRecentVersion', () => {
  it('returns undefined when called without a list', () => {
    expect(getMostRecentVersion(undefined)).toBeUndefined();
  });

  it('prefers stable (v1-v9) over beta over alpha', () => {
    expect(getMostRecentVersion(['v1alpha1', 'v1beta1', 'v1'])).toBe('v1');
    expect(getMostRecentVersion(['v1alpha1', 'v1beta1'])).toBe('v1beta1');
    expect(getMostRecentVersion(['v1alpha1'])).toBe('v1alpha1');
  });

  it('returns the highest stable version when multiple stables exist', () => {
    expect(getMostRecentVersion(['v1', 'v2', 'v3'])).toBe('v3');
  });

  it('handles multi-digit stable versions correctly', () => {
    expect(getMostRecentVersion(['v2', 'v10'])).toBe('v10');
    expect(getMostRecentVersion(['v1', 'v9', 'v10'])).toBe('v10');
  });

  it('handles semantic versions for stable releases', () => {
    expect(getMostRecentVersion(['v1.0.0', 'v2.1.0', 'v1.9.0'])).toBe('v2.1.0');
    expect(getMostRecentVersion(['1.0.0', '2.0.0'])).toBe('2.0.0');
  });

  it('returns the first entry when no stable/beta/alpha can be classified', () => {
    expect(getMostRecentVersion(['unknown', 'mystery'])).toBe('unknown');
  });

  it('handles an empty list (falls back to undefined via optional chaining)', () => {
    expect(getMostRecentVersion([])).toBeUndefined();
  });
});

describe('versionSortComparatorFn', () => {
  it('returns undefined when either argument is missing', () => {
    expect(versionSortComparatorFn(undefined, '1.0.0')).toBeUndefined();
    expect(versionSortComparatorFn('1.0.0', undefined)).toBeUndefined();
    expect(versionSortComparatorFn(null, null)).toBeUndefined();
  });

  it('satisfies comparator contract for WILDCARD_V (anti-symmetry and reflexivity)', () => {
    expect(versionSortComparatorFn(WILDCARD_V, '1.0.0')).toBe(-1);
    expect(versionSortComparatorFn('1.0.0', WILDCARD_V)).toBe(1);
    expect(versionSortComparatorFn(WILDCARD_V, WILDCARD_V)).toBe(0);
  });

  it('compares simple semver-style versions numerically', () => {
    expect(versionSortComparatorFn('1.0.0', '2.0.0')).toBeLessThan(0);
    expect(versionSortComparatorFn('2.0.0', '1.0.0')).toBeGreaterThan(0);
  });

  it('returns 0 when versions are equal', () => {
    expect(versionSortComparatorFn('1.2.3', '1.2.3')).toBe(0);
    expect(versionSortComparatorFn('v1.0.0', 'v1.0.0')).toBe(0);
  });

  it('handles a leading v prefix on either operand', () => {
    expect(versionSortComparatorFn('v1.0.0', '2.0.0')).toBeLessThan(0);
    expect(versionSortComparatorFn('v2.0.0', 'v1.0.0')).toBeGreaterThan(0);
  });

  it('compares the second segment when the first is equal', () => {
    expect(versionSortComparatorFn('1.10.0', '1.2.0')).toBeGreaterThan(0);
  });

  it('resolves unequal version segment lengths correctly', () => {
    expect(versionSortComparatorFn('1.0', '1.0.1')).toBeLessThan(0);
    expect(versionSortComparatorFn('1.0.1', '1.0')).toBeGreaterThan(0);
    expect(versionSortComparatorFn('2.0', '2.0.0')).toBeLessThan(0);
    expect(versionSortComparatorFn('2.0.0', '2.0')).toBeGreaterThan(0);
  });
});

describe('sortByVersionInDecreasingOrder', () => {
  it('returns undefined when given undefined / null', () => {
    expect(sortByVersionInDecreasingOrder(undefined)).toBeUndefined();
    expect(sortByVersionInDecreasingOrder(null)).toBeUndefined();
  });

  it('does not inject a wildcard for a single-version list', () => {
    expect(sortByVersionInDecreasingOrder(['1.0.0'])).toEqual(['1.0.0']);
  });

  it('prepends WILDCARD_V when there are multiple versions', () => {
    const result = sortByVersionInDecreasingOrder(['1.0.0', '2.0.0']);
    expect(result?.[0]).toBe(WILDCARD_V);
  });

  it('keeps exactly one WILDCARD_V when the input already contains it', () => {
    expect(sortByVersionInDecreasingOrder([WILDCARD_V, '2.0.0', '1.0.0'])).toEqual([
      WILDCARD_V,
      '2.0.0',
      '1.0.0',
    ]);
    expect(sortByVersionInDecreasingOrder([WILDCARD_V, '2.0.0'])).toEqual([WILDCARD_V, '2.0.0']);
    expect(sortByVersionInDecreasingOrder([WILDCARD_V])).toEqual([WILDCARD_V]);
  });

  it('sorts in decreasing order numerically', () => {
    expect(sortByVersionInDecreasingOrder(['2.2.1', '2.10.11', '10.1.2', '10.1.1'])).toEqual([
      WILDCARD_V,
      '10.1.2',
      '10.1.1',
      '2.10.11',
      '2.2.1',
    ]);
  });

  it('does not mutate the input array', () => {
    const input = ['1.0.0', '2.0.0'];
    sortByVersionInDecreasingOrder(input);
    expect(input).toEqual(['1.0.0', '2.0.0']);
  });
});

describe('getGreaterVersion', () => {
  it('returns the higher of two versions', () => {
    expect(getGreaterVersion('2.0.0', '1.0.0')).toBe('2.0.0');
    expect(getGreaterVersion('1.0.0', '2.0.0')).toBe('2.0.0');
  });

  it('returns the first when versions are equal (comparator returns 0)', () => {
    const result = getGreaterVersion('1.0.0', '1.0.0');
    expect(result).toBe('1.0.0');
  });

  it('treats v-prefixed versions consistently with bare numbers', () => {
    expect(getGreaterVersion('v2.0.0', 'v1.0.0')).toBe('v2.0.0');
  });
});

describe('sortAndGroupVersionsInModel', () => {
  it('groups duplicate model names and de-duplicates versions', () => {
    const models = [
      { name: 'foo', version: '1.0.0' },
      { name: 'foo', version: '2.0.0' },
      { name: 'foo', version: '1.0.0' },
      { name: 'bar', version: '0.5.0' },
    ];

    const grouped = sortAndGroupVersionsInModel(models);
    expect(grouped).toHaveLength(2);

    const foo = grouped.find((m) => m.name === 'foo');
    expect(foo).toBeDefined();
    // wildcard appears at the front and duplicates are removed
    expect(foo!.version[0]).toBe(WILDCARD_V);
    expect(foo!.version).toContain('1.0.0');
    expect(foo!.version).toContain('2.0.0');
    // de-duplication: should not contain two '1.0.0' entries
    expect(foo!.version.filter((v) => v === '1.0.0')).toHaveLength(1);

    const bar = grouped.find((m) => m.name === 'bar');
    expect(bar).toBeDefined();
    expect(bar!.version).toEqual(['0.5.0']); // single version, no wildcard
  });

  it('does not mutate input model objects in-place', () => {
    const original = [
      { name: 'foo', version: '1.0.0' },
      { name: 'foo', version: '2.0.0' },
    ];

    sortAndGroupVersionsInModel(original);
    expect(original[0].version).toBe('1.0.0');
    expect(original[1].version).toBe('2.0.0');
  });

  it('returns an empty array when given null/undefined', () => {
    expect(sortAndGroupVersionsInModel(null)).toEqual([]);
    expect(sortAndGroupVersionsInModel(undefined)).toEqual([]);
  });

  it('handles an empty array', () => {
    expect(sortAndGroupVersionsInModel([])).toEqual([]);
  });
});

describe('getMostRecentVersion – defensive type guards', () => {
  it('returns undefined for non-array inputs', () => {
    expect(getMostRecentVersion(null)).toBeUndefined();
    expect(getMostRecentVersion(123)).toBeUndefined();
    expect(getMostRecentVersion('v1.0.0')).toBeUndefined();
    expect(getMostRecentVersion({})).toBeUndefined();
    expect(getMostRecentVersion(true)).toBeUndefined();
  });

  it('safely skips nullish and non-string elements inside array', () => {
    expect(getMostRecentVersion(['v1', null, undefined, 42, {}])).toBe('v1');
    expect(getMostRecentVersion([null, undefined, 42])).toBeUndefined();
    expect(getMostRecentVersion(['v1alpha1', null, 'v1'])).toBe('v1');
    expect(getMostRecentVersion([null, '', 'mystery'])).toBe('mystery');
    expect(getMostRecentVersion([''])).toBeUndefined();
  });
});

describe('versionSortComparatorFn – defensive type guards', () => {
  it('returns undefined when either argument is not a string', () => {
    expect(versionSortComparatorFn(42, '1.0.0')).toBeUndefined();
    expect(versionSortComparatorFn('1.0.0', 42)).toBeUndefined();
    expect(versionSortComparatorFn({}, '1.0.0')).toBeUndefined();
    expect(versionSortComparatorFn('1.0.0', [])).toBeUndefined();
    expect(versionSortComparatorFn(null, '1.0.0')).toBeUndefined();
    expect(versionSortComparatorFn('1.0.0', null)).toBeUndefined();
    expect(versionSortComparatorFn('', '1.0.0')).toBeUndefined();
    expect(versionSortComparatorFn('1.0.0', '')).toBeUndefined();
  });
});

describe('sortByVersionInDecreasingOrder – defensive type guards', () => {
  it('returns undefined for non-array inputs', () => {
    expect(sortByVersionInDecreasingOrder(123)).toBeUndefined();
    expect(sortByVersionInDecreasingOrder('1.0.0')).toBeUndefined();
    expect(sortByVersionInDecreasingOrder({})).toBeUndefined();
    expect(sortByVersionInDecreasingOrder(true)).toBeUndefined();
  });

  it('filters out non-string elements when sorting', () => {
    expect(sortByVersionInDecreasingOrder(['1.0.0', null, undefined, 42])).toEqual(['1.0.0']);
    expect(sortByVersionInDecreasingOrder([null, undefined])).toEqual([]);
    expect(sortByVersionInDecreasingOrder(['1.0.0', '2.0.0', null])).toEqual([
      WILDCARD_V,
      '2.0.0',
      '1.0.0',
    ]);
  });
});

describe('getGreaterVersion – defensive type guards', () => {
  it('falls back gracefully when non-string arguments are provided', () => {
    expect(getGreaterVersion(42, '1.0.0')).toBe('1.0.0');
    expect(getGreaterVersion('1.0.0', 42)).toBe('1.0.0');
    expect(getGreaterVersion('', '1.0.0')).toBe('1.0.0');
    expect(getGreaterVersion('1.0.0', '')).toBe('1.0.0');
    expect(getGreaterVersion(null, undefined)).toBeUndefined();
  });
});

describe('sortAndGroupVersionsInModel – defensive type guards', () => {
  it('returns empty array for non-array inputs', () => {
    expect(sortAndGroupVersionsInModel(123)).toEqual([]);
    expect(sortAndGroupVersionsInModel('models')).toEqual([]);
    expect(sortAndGroupVersionsInModel({})).toEqual([]);
  });

  it('de-duplicates duplicate versions in a single model entry without adding a wildcard', () => {
    const singleModelWithDuplicates = [{ name: 'foo', version: ['1.0.0', '1.0.0'] }];
    const result = sortAndGroupVersionsInModel(singleModelWithDuplicates);
    expect(result).toHaveLength(1);
    expect(result[0].name).toBe('foo');
    expect(result[0].version).toEqual(['1.0.0']);
  });

  it('safely ignores null or invalid elements in array', () => {
    const mixed = [
      null,
      undefined,
      42,
      { name: 'foo', version: '1.0.0' },
      { notAModel: true },
      { name: '', version: '2.0.0' },
    ];
    const result = sortAndGroupVersionsInModel(mixed);
    expect(result).toHaveLength(1);
    expect(result[0].name).toBe('foo');
    expect(result[0].version).toEqual(['1.0.0']);
  });
});
