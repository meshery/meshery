import { describe, expect, it } from 'vitest';
import { CamelCaseToSentanceCase } from '../camelCaseToSentanceCase';

// The implementation runs two regex replaces in sequence:
//   1) /([A-Z]+)/g  -> prefixes ANY capital run with a space
//   2) /([A-Z][a-z])/g -> prefixes a `Capital + lowercase` pair with a space again
// As a result, a simple "camel" -> "Case" boundary ends up doubly-spaced.
// These tests pin down the actual behaviour rather than the idealised one.

describe('CamelCaseToSentanceCase', () => {
  it('separates each capital-cased word with spaces (double-spaced at simple boundaries)', () => {
    expect(CamelCaseToSentanceCase('camelCase')).toBe('camel  Case');
    expect(CamelCaseToSentanceCase('helloWorldFoo')).toBe('hello  World  Foo');
  });

  it('handles PascalCase input by prepending leading spaces', () => {
    expect(CamelCaseToSentanceCase('PascalCase')).toBe('  Pascal  Case');
  });

  it('returns an empty string unchanged', () => {
    expect(CamelCaseToSentanceCase('')).toBe('');
  });

  it('returns lowercase-only words unchanged', () => {
    expect(CamelCaseToSentanceCase('hello')).toBe('hello');
  });

  it('keeps acronyms together while still separating subsequent TitleCase tokens', () => {
    // 'URLValidator' -> first regex inserts a space before the acronym,
    // second regex inserts another space at the `R` + `V` boundary
    const result = CamelCaseToSentanceCase('URLValidator');
    expect(result).toContain('URL');
    expect(result).toContain('Validator');
    // The acronym should still read as one contiguous token
    expect(result.replace(/\s+/g, ' ').trim()).toBe('URL Validator');
  });

  it('does not change numeric characters', () => {
    expect(CamelCaseToSentanceCase('version2Update')).toBe('version2  Update');
  });

  it('handles a single uppercase letter (one leading space from each regex)', () => {
    // Single uppercase matches the first regex only (no following lowercase)
    expect(CamelCaseToSentanceCase('A')).toBe(' A');
  });

  it('handles a single uppercase + lowercase pair (matches both regexes)', () => {
    expect(CamelCaseToSentanceCase('Ab')).toBe('  Ab');
  });
});
