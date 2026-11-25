import { describe, expect, it } from 'vitest';
import PascalCaseToKebab from '../PascalCaseToKebab';

describe('PascalCaseToKebab', () => {
  it('converts a PascalCase identifier to kebab-case', () => {
    expect(PascalCaseToKebab('HelloWorld')).toBe('hello-world');
  });

  it('converts a single-word PascalCase identifier to lowercase', () => {
    expect(PascalCaseToKebab('Hello')).toBe('hello');
  });

  it('keeps camelCase as-is but kebab-separated', () => {
    expect(PascalCaseToKebab('helloWorld')).toBe('hello-world');
  });

  it('handles multiple uppercase words', () => {
    expect(PascalCaseToKebab('FooBarBaz')).toBe('foo-bar-baz');
  });

  it('returns lowercase single-word input unchanged', () => {
    expect(PascalCaseToKebab('hello')).toBe('hello');
  });

  it('inserts a hyphen before each capital letter inside the word', () => {
    // First char lowercased, then each remaining capital prefixed with '-'
    expect(PascalCaseToKebab('ABC')).toBe('a-b-c');
  });

  it('handles strings beginning with a digit (no leading hyphen)', () => {
    expect(PascalCaseToKebab('1FooBar')).toBe('1-foo-bar');
  });

  it('safely handles nullish, non-string, or empty inputs without throwing', () => {
    // @ts-expect-error testing invalid inputs
    expect(PascalCaseToKebab(null)).toBe('');
    // @ts-expect-error testing invalid inputs
    expect(PascalCaseToKebab(undefined)).toBe('');
    // @ts-expect-error testing invalid inputs
    expect(PascalCaseToKebab(123)).toBe('');
    // @ts-expect-error testing invalid inputs
    expect(PascalCaseToKebab({})).toBe('');
    // @ts-expect-error testing invalid inputs
    expect(PascalCaseToKebab([])).toBe('');
    expect(PascalCaseToKebab('')).toBe('');
  });
});
