import { describe, expect, it } from 'vitest';
import { safeJsonParse } from './json-parse';

describe('safeJsonParse', () => {
  it('parses a valid JSON string into an object', () => {
    const result = safeJsonParse('{"hostIP":"10.0.0.1"}');
    expect(result).toEqual({ hostIP: '10.0.0.1' });
  });

  it('returns the fallback for null input', () => {
    const result = safeJsonParse(null);
    expect(result).toEqual({});
  });

  it('returns the fallback for undefined input', () => {
    const result = safeJsonParse(undefined);
    expect(result).toEqual({});
  });

  it('returns the fallback for malformed JSON strings', () => {
    const result = safeJsonParse('{invalid json');
    expect(result).toEqual({});
  });

  it('returns the fallback for an empty string', () => {
    const result = safeJsonParse('');
    expect(result).toEqual({});
  });

  it('returns a pre-parsed object as-is', () => {
    const obj = { cpu: '4', memory: '8Gi' };
    const result = safeJsonParse(obj);
    expect(result).toBe(obj);
  });

  it('returns the fallback for numeric input', () => {
    const result = safeJsonParse(42);
    expect(result).toEqual({});
  });

  it('returns the fallback for boolean input', () => {
    const result = safeJsonParse(true);
    expect(result).toEqual({});
  });

  it('returns the fallback when JSON.parse yields a primitive (string)', () => {
    const result = safeJsonParse('"just a string"');
    expect(result).toEqual({});
  });

  it('returns the fallback when JSON.parse yields a primitive (number)', () => {
    const result = safeJsonParse('123');
    expect(result).toEqual({});
  });

  it('returns the fallback when JSON.parse yields null', () => {
    const result = safeJsonParse('null');
    expect(result).toEqual({});
  });

  it('accepts a custom fallback value', () => {
    const fallback = { status: 'unknown' };
    const result = safeJsonParse('{bad', fallback);
    expect(result).toBe(fallback);
  });

  it('parses a valid JSON array string', () => {
    const result = safeJsonParse('[1,2,3]');
    expect(result).toEqual([1, 2, 3]);
  });

  it('returns a pre-parsed array as-is', () => {
    const arr = [{ name: 'a' }];
    const result = safeJsonParse(arr);
    expect(result).toBe(arr);
  });

  it('handles deeply nested valid JSON', () => {
    const json = '{"spec":{"template":{"spec":{"restartPolicy":"Always"}}}}';
    const result = safeJsonParse(json);
    expect(result).toEqual({
      spec: { template: { spec: { restartPolicy: 'Always' } } },
    });
  });

  it('returns the fallback for the string literal "undefined"', () => {
    const result = safeJsonParse('undefined');
    expect(result).toEqual({});
  });
});
