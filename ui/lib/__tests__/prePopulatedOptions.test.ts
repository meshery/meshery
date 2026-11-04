import { describe, expect, it } from 'vitest';
import { durationOptions } from '../prePopulatedOptions';

describe('durationOptions', () => {
  it('is an array', () => {
    expect(Array.isArray(durationOptions)).toBe(true);
  });

  it('contains 12 entries', () => {
    expect(durationOptions).toHaveLength(12);
  });

  it('contains all expected duration strings in order', () => {
    expect(durationOptions).toEqual([
      '15s',
      '30s',
      '1m',
      '3m',
      '5m',
      '10m',
      '30m',
      '1h',
      '2h',
      '5h',
      '10h',
      '1d',
    ]);
  });

  it('contains only string values', () => {
    durationOptions.forEach((opt) => {
      expect(typeof opt).toBe('string');
    });
  });

  it('every entry ends with a recognised time unit suffix', () => {
    const allowedSuffixes = /[smhd]$/;
    durationOptions.forEach((opt) => {
      expect(opt).toMatch(allowedSuffixes);
    });
  });

  it('entries follow the "<number><unit>" pattern', () => {
    const pattern = /^\d+[smhd]$/;
    durationOptions.forEach((opt) => {
      expect(opt).toMatch(pattern);
    });
  });

  it('contains no duplicates', () => {
    expect(new Set(durationOptions).size).toBe(durationOptions.length);
  });

  it('starts with the smallest duration', () => {
    expect(durationOptions[0]).toBe('15s');
  });

  it('ends with the largest duration', () => {
    expect(durationOptions[durationOptions.length - 1]).toBe('1d');
  });
});
