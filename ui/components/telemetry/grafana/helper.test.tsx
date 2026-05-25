import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import grafanaDateRangeToDate from './helper';

// Pin time so the relative-offset cases produce deterministic dates.
const FIXED_NOW = new Date('2024-07-15T12:34:56.789Z');

describe('grafanaDateRangeToDate', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.setSystemTime(FIXED_NOW);
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('returns now() when the token is "now"', () => {
    expect(grafanaDateRangeToDate('now', false).getTime()).toBe(FIXED_NOW.getTime());
  });

  it.each([
    ['now-2d', 2 * 24 * 60 * 60 * 1000],
    ['now-7d', 7 * 24 * 60 * 60 * 1000],
    ['now-30d', 30 * 24 * 60 * 60 * 1000],
    ['now-90d', 90 * 24 * 60 * 60 * 1000],
    ['now-5m', 5 * 60 * 1000],
    ['now-15m', 15 * 60 * 1000],
    ['now-30m', 30 * 60 * 1000],
    ['now-1h', 60 * 60 * 1000],
    ['now-3h', 3 * 60 * 60 * 1000],
    ['now-6h', 6 * 60 * 60 * 1000],
    ['now-12h', 12 * 60 * 60 * 1000],
    ['now-24h', 24 * 60 * 60 * 1000],
  ])('subtracts the relative offset for %s', (token, deltaMs) => {
    const result = grafanaDateRangeToDate(token, false);
    expect(FIXED_NOW.getTime() - result.getTime()).toBe(deltaMs);
  });

  it('handles month-relative tokens (now-6M)', () => {
    const result = grafanaDateRangeToDate('now-6M', false);
    expect(result.getMonth()).toBe((FIXED_NOW.getMonth() + 12 - 6) % 12);
  });

  it.each([
    ['now-1y', 1],
    ['now-2y', 2],
    ['now-5y', 5],
  ])('subtracts the right number of years for %s', (token, years) => {
    const result = grafanaDateRangeToDate(token, false);
    expect(result.getFullYear()).toBe(FIXED_NOW.getFullYear() - years);
  });

  it('returns start-of-day for now-1d/d when startDate is truthy', () => {
    const result = grafanaDateRangeToDate('now-1d/d', true);
    expect(result.getHours()).toBe(0);
    expect(result.getMinutes()).toBe(0);
    expect(result.getSeconds()).toBe(0);
    expect(result.getMilliseconds()).toBe(0);
  });

  it('returns end-of-day for now-1d/d when startDate is falsy', () => {
    const result = grafanaDateRangeToDate('now-1d/d', false);
    expect(result.getHours()).toBe(23);
    expect(result.getMinutes()).toBe(59);
    expect(result.getSeconds()).toBe(59);
    expect(result.getMilliseconds()).toBe(999);
  });

  it('parses unknown tokens as a numeric timestamp via parseFloat', () => {
    const result = grafanaDateRangeToDate('1700000000000', false);
    expect(result.getTime()).toBe(1700000000000);
  });

  it('handles end-of-month tokens (now/M)', () => {
    const start = grafanaDateRangeToDate('now/M', true);
    expect(start.getDate()).toBe(1);
    expect(start.getHours()).toBe(0);
  });

  it('handles end-of-year tokens (now/y) returning end of year when not startDate', () => {
    const end = grafanaDateRangeToDate('now/y', false);
    // The branch sets month to 12 then date to 0, normalising to Dec 31.
    expect(end.getMonth()).toBe(11);
    expect(end.getDate()).toBe(31);
    expect(end.getHours()).toBe(23);
  });
});
