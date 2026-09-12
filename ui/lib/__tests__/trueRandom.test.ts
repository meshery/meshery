import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { trueRandom } from '../trueRandom';

describe('trueRandom', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('returns a number', () => {
    const value = trueRandom();
    expect(typeof value).toBe('number');
  });

  it('returns a value in [0, 1)', () => {
    for (let i = 0; i < 100; i++) {
      const value = trueRandom();
      expect(value).toBeGreaterThanOrEqual(0);
      expect(value).toBeLessThan(1);
    }
  });

  it('returns 0 when crypto returns 0', () => {
    vi.spyOn(crypto, 'getRandomValues').mockImplementation((arr) => {
      (arr as Uint32Array)[0] = 0;
      return arr as Uint32Array;
    });
    expect(trueRandom()).toBe(0);
  });

  it('returns a value just below 1 when crypto returns max uint32', () => {
    vi.spyOn(crypto, 'getRandomValues').mockImplementation((arr) => {
      (arr as Uint32Array)[0] = 0xffffffff;
      return arr as Uint32Array;
    });
    const value = trueRandom();
    expect(value).toBeLessThan(1);
    expect(value).toBeGreaterThan(0.9999);
  });

  it('uses crypto.getRandomValues with a Uint32Array of length 1', () => {
    const spy = vi.spyOn(crypto, 'getRandomValues');
    trueRandom();
    expect(spy).toHaveBeenCalledTimes(1);
    const arg = spy.mock.calls[0][0];
    expect(arg).toBeInstanceOf(Uint32Array);
    expect((arg as Uint32Array).length).toBe(1);
  });

  it('produces varied output across many calls', () => {
    // Reasonable check that we are not getting the same value each time
    const values = new Set<number>();
    for (let i = 0; i < 50; i++) {
      values.add(trueRandom());
    }
    expect(values.size).toBeGreaterThan(1);
  });
});
