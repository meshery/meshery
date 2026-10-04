import { describe, expect, it } from 'vitest';
import { formatValue } from './time';

describe('formatValue', () => {
  it('scales unitless values with K/M/B suffixes', () => {
    expect(formatValue(42)).toBe('42.00');
    expect(formatValue(1500)).toBe('1.50K');
    expect(formatValue(2_500_000)).toBe('2.50M');
    expect(formatValue(2_500_000_000, 'short')).toBe('2.50B');
    expect(formatValue(0.1234)).toBe('0.1234');
  });

  it('moves to the next suffix when rounding reaches it', () => {
    expect(formatValue(999_999)).toBe('1.00M');
    expect(formatValue(999.999)).toBe('1.00K');
  });

  it('keeps percent units', () => {
    expect(formatValue(42, 'percent')).toBe('42.00%');
    expect(formatValue(0.42, 'percentunit')).toBe('42.00%');
  });

  it('formats bytes with IEC prefixes', () => {
    expect(formatValue(512, 'bytes')).toBe('512.00 B');
    expect(formatValue(512 * 1024 ** 2, 'bytes')).toBe('512.00 MiB');
    expect(formatValue(1.5 * 1024 ** 3, 'bytes')).toBe('1.50 GiB');
    expect(formatValue(1024 ** 2 - 1, 'bytes')).toBe('1.00 MiB');
  });

  it('formats decbytes with SI prefixes', () => {
    expect(formatValue(2e9, 'decbytes')).toBe('2.00 GB');
    expect(formatValue(1500, 'decbytes')).toBe('1.50 kB');
  });

  it('formats data rates', () => {
    expect(formatValue(2e6, 'Bps')).toBe('2.00 MB/s');
    expect(formatValue(2 * 1024 ** 2, 'binBps')).toBe('2.00 MiB/s');
  });

  it('scales seconds to the largest fitting time unit', () => {
    expect(formatValue(0.25, 's')).toBe('250.00 ms');
    expect(formatValue(42, 's')).toBe('42.00 s');
    expect(formatValue(90, 's')).toBe('1.50 min');
    expect(formatValue(7200, 's')).toBe('2.00 hour');
    expect(formatValue(3 * 86400, 's')).toBe('3.00 day');
    expect(formatValue(59.999, 's')).toBe('1.00 min');
  });

  it('scales milliseconds', () => {
    expect(formatValue(250, 'ms')).toBe('250.00 ms');
    expect(formatValue(1500, 'ms')).toBe('1.50 s');
  });

  it('keeps rate units', () => {
    expect(formatValue(1500, 'reqps')).toBe('1.50K req/s');
    expect(formatValue(12, 'ops')).toBe('12.00 ops/s');
  });

  it('applies custom suffix and prefix units', () => {
    expect(formatValue(3, 'suffix: pods')).toBe('3.00 pods');
    expect(formatValue(3, 'prefix:$')).toBe('$3.00');
  });

  it('keeps an unknown unit instead of dropping it', () => {
    expect(formatValue(3, 'widgets')).toBe('3.00 widgets');
  });

  it('formats non-finite values as a dash', () => {
    expect(formatValue(Number.NaN, 'bytes')).toBe('—');
    expect(formatValue(Number.POSITIVE_INFINITY, 's')).toBe('—');
  });
});
