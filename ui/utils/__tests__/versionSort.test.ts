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

  it('treats WILDCARD_V as the smallest element (sorts to the front)', () => {
    expect(versionSortComparatorFn(WILDCARD_V, '1.0.0')).toBe(-1);
    expect(versionSortComparatorFn('1.0.0', WILDCARD_V)).toBe(-1);
  });

  it('compares simple semver-style versions numerically', () => {
    expect(versionSortComparatorFn('1.0.0', '2.0.0')).toBeLessThan(0);
    expect(versionSortComparatorFn('2.0.0', '1.0.0')).toBeGreaterThan(0);
  });

  it('returns 0 (falls through) when versions are equal', () => {
    // The function continues looping while equal segments are found, so equal
    // versions of the same length return `undefined` (no comparison made).
    const result = versionSortComparatorFn('1.2.3', '1.2.3');
    expect(result === 0 || result === undefined).toBe(true);
  });

  it('handles a leading v prefix on either operand', () => {
    expect(versionSortComparatorFn('v1.0.0', '2.0.0')).toBeLessThan(0);
    expect(versionSortComparatorFn('v2.0.0', 'v1.0.0')).toBeGreaterThan(0);
  });

  it('compares the second segment when the first is equal', () => {
    expect(versionSortComparatorFn('1.10.0', '1.2.0')).toBeGreaterThan(0);
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

  it('returns the first when versions are equal (comparator returns 0/undefined)', () => {
    // When the comparator returns `undefined`, `>= 0` is false, so the function returns v2.
    // This pins down current behaviour rather than the ideal.
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

  it('returns an empty array when given null/undefined', () => {
    expect(sortAndGroupVersionsInModel(null)).toEqual([]);
    expect(sortAndGroupVersionsInModel(undefined)).toEqual([]);
  });

  it('handles an empty array', () => {
    expect(sortAndGroupVersionsInModel([])).toEqual([]);
  });
});
