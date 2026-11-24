import { describe, expect, it, vi } from 'vitest';

const hoisted = vi.hoisted(() => ({
  eventBus: { publish: vi.fn() },
  useSelectorMock: vi.fn(),
}));

vi.mock('../../lib/trueRandom', () => ({
  trueRandom: vi.fn(() => 0.5),
}));

// Avoid evaluating `pages/_app` (which pulls notistack/material/relay/etc.).
vi.mock('../../pages/_app', () => ({
  mesheryExtensionRoute: '/extension/meshmap',
}));

vi.mock('../eventBus', () => ({
  mesheryEventBus: hoisted.eventBus,
}));

vi.mock('react-redux', () => ({
  useSelector: (sel: (s: unknown) => unknown) => hoisted.useSelectorMock(sel),
}));

import {
  camelcaseToSnakecase,
  getColumnValue,
  getComponentFromDesign,
  getVisibilityColums,
  isEmptyArr,
  isEqualArr,
  urlEncodeParams,
} from '../utils';

describe('isEmptyArr defensive guards', () => {
  it('returns true strictly for empty arrays', () => {
    expect(isEmptyArr([])).toBe(true);
  });

  it('returns false for non-empty arrays', () => {
    expect(isEmptyArr([1, 2])).toBe(false);
  });

  it('returns false for null, undefined, empty strings, numbers, and plain objects', () => {
    expect(isEmptyArr(null as never)).toBe(false);
    expect(isEmptyArr(undefined as never)).toBe(false);
    expect(isEmptyArr('' as never)).toBe(false);
    expect(isEmptyArr(0 as never)).toBe(false);
    expect(isEmptyArr({ length: 0 } as never)).toBe(false);
  });
});

describe('isEqualArr defensive guards and immutability', () => {
  it('does not mutate input arrays when orderMatters is false', () => {
    const arr1 = [3, 1, 2];
    const arr2 = [2, 3, 1];
    expect(isEqualArr(arr1, arr2, false)).toBe(true);
    expect(arr1).toEqual([3, 1, 2]);
    expect(arr2).toEqual([2, 3, 1]);
  });

  it('safely compares frozen arrays without throwing TypeError', () => {
    const frozen1 = Object.freeze([3, 1, 2]);
    const frozen2 = Object.freeze([2, 3, 1]);
    expect(isEqualArr(frozen1, frozen2, false)).toBe(true);
    expect(isEqualArr(frozen1, frozen2, true)).toBe(false);
  });

  it('returns false when non-array arguments are passed', () => {
    expect(isEqualArr('abc' as never, 'abc' as never)).toBe(false);
    expect(isEqualArr({ length: 2 } as never, [1, 2])).toBe(false);
    expect(isEqualArr([1, 2], { length: 2 } as never)).toBe(false);
    expect(isEqualArr(null as never, [1])).toBe(false);
    expect(isEqualArr([1], undefined as never)).toBe(false);
  });
});

describe('getColumnValue defensive guards', () => {
  const columns = [{ name: 'id' }, { name: 'status' }];

  it('returns column value when valid rowData and column match', () => {
    expect(getColumnValue(['row-1', 'active'], 'status', columns)).toBe('active');
  });

  it('returns undefined safely when columns or rowData is nullish or malformed', () => {
    expect(getColumnValue(null, 'status', columns)).toBeUndefined();
    expect(getColumnValue(undefined, 'status', columns)).toBeUndefined();
    expect(getColumnValue(['row-1'], 'status', null as never)).toBeUndefined();
    expect(getColumnValue(['row-1'], 'status', undefined as never)).toBeUndefined();
    expect(getColumnValue(['row-1'], 'status', [{}] as never)).toBeUndefined();
  });

  it('returns undefined when column is not present in columns array', () => {
    expect(getColumnValue(['row-1', 'active'], 'nonexistent', columns)).toBeUndefined();
  });
});

describe('getVisibilityColums defensive guards', () => {
  it('returns empty array when columns is nullish or not an array', () => {
    expect(getVisibilityColums(null as never)).toEqual([]);
    expect(getVisibilityColums(undefined as never)).toEqual([]);
    expect(getVisibilityColums({} as never)).toEqual([]);
  });

  it('filters columns based on options.display safely', () => {
    const cols = [
      { name: 'col1', options: { display: false } },
      { name: 'col2' },
      { name: 'col3', options: { display: true } },
    ];
    expect(getVisibilityColums(cols)).toEqual([
      { name: 'col2' },
      { name: 'col3', options: { display: true } },
    ]);
  });
});

describe('camelcaseToSnakecase defensive guards', () => {
  it('passes through undefined input', () => {
    expect(camelcaseToSnakecase(undefined as never)).toBeUndefined();
  });

  it('returns empty string for non-string values without throwing', () => {
    expect(camelcaseToSnakecase(123 as never)).toBe('');
    expect(camelcaseToSnakecase(null as never)).toBe('');
    expect(camelcaseToSnakecase({} as never)).toBe('');
    expect(camelcaseToSnakecase(true as never)).toBe('');
  });

  it('converts camelCase strings properly', () => {
    expect(camelcaseToSnakecase('camelCaseText')).toBe('camel_case_text');
  });
});

describe('getComponentFromDesign defensive guards', () => {
  it('returns the matching component by id', () => {
    const design = { components: [{ id: 'comp-1' }, { id: 'comp-2' }] };
    expect(getComponentFromDesign(design, 'comp-2')).toEqual({ id: 'comp-2' });
  });

  it('returns undefined when design or components is nullish or malformed', () => {
    expect(getComponentFromDesign(null as never, 'comp-1')).toBeUndefined();
    expect(getComponentFromDesign(undefined as never, 'comp-1')).toBeUndefined();
    expect(getComponentFromDesign({}, 'comp-1')).toBeUndefined();
    expect(getComponentFromDesign({ components: null } as never, 'comp-1')).toBeUndefined();
    expect(getComponentFromDesign({ components: 'not-array' } as never, 'comp-1')).toBeUndefined();
  });

  it('returns undefined when componentId is empty or not found', () => {
    const design = { components: [{ id: 'comp-1' }] };
    expect(getComponentFromDesign(design, '')).toBeUndefined();
    expect(getComponentFromDesign(design, 'comp-99')).toBeUndefined();
  });
});

describe('urlEncodeParams defensive guards', () => {
  it('returns empty string when params is nullish or not an object', () => {
    expect(urlEncodeParams(null as never)).toBe('');
    expect(urlEncodeParams(undefined as never)).toBe('');
    expect(urlEncodeParams('string' as never)).toBe('');
  });

  it('encodes valid object parameters and skips nil values', () => {
    expect(urlEncodeParams({ a: '1', b: undefined, c: null, d: '2' })).toBe('a=1&d=2');
  });
});
