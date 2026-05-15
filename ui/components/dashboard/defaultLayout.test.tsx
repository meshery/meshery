import { describe, expect, it } from 'vitest';
import { DEFAULT_LAYOUT, LOCAL_PROVIDER_LAYOUT, OVERVIEW_LAYOUT } from './defaultLayout';

const BREAKPOINTS = ['lg', 'md', 'sm', 'xs', 'xxs'] as const;

describe('defaultLayout', () => {
  describe('DEFAULT_LAYOUT', () => {
    it('exposes layouts for every breakpoint', () => {
      for (const bp of BREAKPOINTS) {
        expect(Array.isArray(DEFAULT_LAYOUT[bp])).toBe(true);
        expect(DEFAULT_LAYOUT[bp].length).toBeGreaterThan(0);
      }
    });

    it('includes the canonical widget keys at every breakpoint', () => {
      const expectedKeys = [
        'OVERVIEW',
        'GETTING_STARTED',
        'HELP_CENTER',
        'MY_DESIGNS',
        'WORKSPACE_ACTIVITY',
        'CONNECTIONS_STATUS_CHART',
      ];
      for (const bp of BREAKPOINTS) {
        const keys = DEFAULT_LAYOUT[bp].map((item: { i: string }) => item.i);
        for (const expected of expectedKeys) {
          expect(keys).toContain(expected);
        }
      }
    });

    it('each layout item exposes integer-like positional fields', () => {
      for (const bp of BREAKPOINTS) {
        for (const item of DEFAULT_LAYOUT[bp]) {
          expect(typeof item.x).toBe('number');
          expect(typeof item.y).toBe('number');
          expect(typeof item.w).toBe('number');
          expect(typeof item.h).toBe('number');
        }
      }
    });
  });

  describe('LOCAL_PROVIDER_LAYOUT', () => {
    it('exposes layouts for every breakpoint', () => {
      for (const bp of BREAKPOINTS) {
        expect(Array.isArray(LOCAL_PROVIDER_LAYOUT[bp])).toBe(true);
        expect(LOCAL_PROVIDER_LAYOUT[bp].length).toBe(2);
      }
    });

    it('contains OVERVIEW and LATEST_BLOGS widgets', () => {
      for (const bp of BREAKPOINTS) {
        const keys = LOCAL_PROVIDER_LAYOUT[bp].map((item: { i: string }) => item.i);
        expect(keys).toEqual(expect.arrayContaining(['OVERVIEW', 'LATEST_BLOGS']));
      }
    });
  });

  describe('OVERVIEW_LAYOUT', () => {
    it('exposes only an OVERVIEW widget per breakpoint', () => {
      for (const bp of BREAKPOINTS) {
        expect(OVERVIEW_LAYOUT[bp]).toHaveLength(1);
        expect(OVERVIEW_LAYOUT[bp][0].i).toBe('OVERVIEW');
      }
    });
  });
});
