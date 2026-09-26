import { beforeEach, describe, expect, it } from 'vitest';
import { loadCachedUserKeys } from '../userKeys';

const validKeys = [
  {
    id: '11111111-1111-1111-1111-111111111111',
    owner: '22222222-2222-2222-2222-222222222222',
    function: 'View Designs',
    category: 'Designs',
    subcategory: 'View',
    description: 'View designs',
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
  },
];

describe('loadCachedUserKeys', () => {
  beforeEach(() => {
    window.sessionStorage.clear();
  });

  it('returns the cached keys when the entry is a valid key array', () => {
    window.sessionStorage.setItem('keys', JSON.stringify(validKeys));
    expect(loadCachedUserKeys()).toEqual(validKeys);
    expect(window.sessionStorage.getItem('keys')).toBe(JSON.stringify(validKeys));
  });

  it('returns an empty array unchanged', () => {
    window.sessionStorage.setItem('keys', '[]');
    expect(loadCachedUserKeys()).toEqual([]);
  });

  it('returns null when nothing is cached', () => {
    expect(loadCachedUserKeys()).toBeNull();
  });

  it.each([
    ['the literal string "undefined"', 'undefined'],
    ['malformed JSON', 'undefined{broken]'],
    ['an object', '{"some_object": true}'],
    ['null', 'null'],
    ['an array with a null element', '[null]'],
    ['an array of primitives', '[1, 2]'],
    ['an array of objects without an id', '[{"function": "View Designs"}]'],
    ['an array of objects without a function', '[{"id": "read"}]'],
    ['a key with a blank id', '[{"id": "  ", "function": "View Designs"}]'],
    ['a key with a blank function', '[{"id": "read", "function": ""}]'],
  ])('returns null and discards the entry when it is %s', (_label, raw) => {
    window.sessionStorage.setItem('keys', raw);
    expect(loadCachedUserKeys()).toBeNull();
    expect(window.sessionStorage.getItem('keys')).toBeNull();
  });
});
