import { describe, expect, it } from 'vitest';
import { isEmptyAtAllDepths, findNestedObject, filterEmptyFields } from '../objects';

describe('isEmptyAtAllDepths', () => {
  it('treats nullish / empty primitives as empty', () => {
    expect(isEmptyAtAllDepths('')).toBe(true);
    expect(isEmptyAtAllDepths(null)).toBe(true);
    expect(isEmptyAtAllDepths(undefined)).toBe(true);
  });

  it('treats lodash-non-empty primitives (numbers, booleans) as empty', () => {
    // lodash's `isEmpty` returns true for numbers and booleans by design
    expect(isEmptyAtAllDepths(0)).toBe(true);
    expect(isEmptyAtAllDepths(42)).toBe(true);
    expect(isEmptyAtAllDepths(true)).toBe(true);
  });

  it('treats populated strings as non-empty', () => {
    expect(isEmptyAtAllDepths('hello')).toBe(false);
  });

  it('returns true for an empty array', () => {
    expect(isEmptyAtAllDepths([])).toBe(true);
  });

  it('returns true for an empty object', () => {
    expect(isEmptyAtAllDepths({})).toBe(true);
  });

  it('returns true for nested empty structures', () => {
    expect(isEmptyAtAllDepths([[], {}, [{}]])).toBe(true);
    expect(isEmptyAtAllDepths({ a: {}, b: [], c: { d: [] } })).toBe(true);
    expect(isEmptyAtAllDepths({ a: '' })).toBe(true);
  });

  it('returns false when any leaf has a non-empty string', () => {
    expect(isEmptyAtAllDepths({ a: { b: 'value' } })).toBe(false);
    expect(isEmptyAtAllDepths([{ deep: [{ leaf: 'x' }] }])).toBe(false);
  });
});

describe('findNestedObject', () => {
  it('returns the first object that satisfies the condition (DFS, stack-based)', () => {
    const tree = {
      id: 1,
      child: { id: 2, leaf: { id: 3, kind: 'target' } },
    };
    const result = findNestedObject(tree, (obj) => obj?.kind === 'target');
    expect(result).toEqual({ id: 3, kind: 'target' });
  });

  it('returns the root when the root itself matches', () => {
    const root = { kind: 'target', value: 1 };
    expect(findNestedObject(root, (obj) => obj?.kind === 'target')).toBe(root);
  });

  it('returns null when nothing matches', () => {
    const tree = { a: 1, b: { c: 'two' } };
    expect(findNestedObject(tree, (obj) => obj?.kind === 'missing')).toBeNull();
  });

  it('does not descend into the `models` property of an object', () => {
    const tree = {
      kind: 'root',
      models: { kind: 'target' }, // should be skipped
      child: { kind: 'other' },
    };
    expect(findNestedObject(tree, (obj) => obj?.kind === 'target')).toBeNull();
  });

  it('searches inside arrays', () => {
    const tree = {
      list: [{ id: 1 }, { id: 2, kind: 'needle' }],
    };
    expect(findNestedObject(tree, (obj) => obj?.kind === 'needle')).toEqual({
      id: 2,
      kind: 'needle',
    });
  });

  it('returns the input itself if condition matches a non-object', () => {
    expect(findNestedObject('hello', (val) => val === 'hello')).toBe('hello');
  });
});

describe('filterEmptyFields', () => {
  it('returns an empty object when given falsy input', () => {
    expect(filterEmptyFields(null)).toEqual({});
    expect(filterEmptyFields(undefined)).toEqual({});
    expect(filterEmptyFields(0)).toEqual({});
    expect(filterEmptyFields(false)).toEqual({});
  });

  it('removes empty string and undefined properties', () => {
    const data = { a: 'value', b: '', c: undefined, d: null, e: 0, f: false };
    expect(filterEmptyFields(data)).toEqual({ a: 'value', d: null, e: 0, f: false });
  });

  it('keeps all properties when none are empty', () => {
    expect(filterEmptyFields({ a: 1, b: 'b' })).toEqual({ a: 1, b: 'b' });
  });

  it('keeps array and object values verbatim', () => {
    const data = { list: [1, 2], obj: { k: 'v' } };
    expect(filterEmptyFields(data)).toEqual(data);
  });

  it('returns a new object (does not mutate the input)', () => {
    const data = { a: '', b: 'keep' };
    const result = filterEmptyFields(data);
    expect(result).not.toBe(data);
    expect(data).toEqual({ a: '', b: 'keep' });
  });
});
