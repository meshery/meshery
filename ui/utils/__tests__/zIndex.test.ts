import { describe, expect, it } from 'vitest';
import { ziCalc } from '../zIndex';

describe('ziCalc', () => {
  it('returns "9" by default (power = 1)', () => {
    expect(ziCalc()).toBe('9');
  });

  it('returns "9" for explicit power 1', () => {
    expect(ziCalc(1)).toBe('9');
  });

  it('returns a string of nines whose length matches the power', () => {
    expect(ziCalc(2)).toBe('99');
    expect(ziCalc(3)).toBe('999');
    expect(ziCalc(5)).toBe('99999');
  });

  it('returns undefined when the power is below 1 (e.g. 0)', () => {
    // The function intentionally short-circuits when `p < 1` and does NOT
    // return the literal '0' (that line is an unreferenced expression).
    expect(ziCalc(0)).toBeUndefined();
    expect(ziCalc(-1)).toBeUndefined();
  });

  it('rounds-down-style: fractional powers between 0 and 1 return undefined', () => {
    expect(ziCalc(0.5)).toBeUndefined();
  });

  it('fractional powers >= 1 produce ceil(p) nines', () => {
    // The for-loop iterates `i = 0..p-1` with `i < p`. For p=1.7 that means
    // i=0 (0 < 1.7) and i=1 (1 < 1.7) -- two nines. For p=2.5 you get three.
    expect(ziCalc(1.7)).toBe('99');
    expect(ziCalc(2.5)).toBe('999');
  });
});
