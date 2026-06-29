import { describe, expect, it, vi, beforeEach } from 'vitest';

let createThemeCalls = 0;
let createCacheCalls = 0;

vi.mock('@sistent/sistent', () => ({
  createTheme: () => {
    createThemeCalls += 1;
    return { __theme: true };
  },
}));

vi.mock('@emotion/cache', () => ({
  default: (opts: any) => {
    createCacheCalls += 1;
    return { __cache: true, ...opts };
  },
}));

describe('getPageContext', () => {
  beforeEach(async () => {
    vi.resetModules();
    createThemeCalls = 0;
    createCacheCalls = 0;
  });

  it('returns an object with theme and emotionCache on the server', async () => {
    const originalWindow = global.window;
    // Simulate the server: delete window before importing the module.
    // The module reads `typeof window === 'undefined'` at call time.
    // @ts-expect-error - intentionally removing window for server simulation
    delete (global as any).window;

    const { default: getPageContext } = await import('./PageContext');
    const ctx = getPageContext();

    expect(ctx).toHaveProperty('theme');
    expect(ctx).toHaveProperty('emotionCache');
    expect(ctx.theme).toEqual({ __theme: true });

    (global as any).window = originalWindow;
  });

  it('creates a new context for every server-side request', async () => {
    const originalWindow = global.window;
    // @ts-expect-error - intentionally removing window for server simulation
    delete (global as any).window;

    const { default: getPageContext } = await import('./PageContext');
    getPageContext();
    getPageContext();

    // On the server every call recreates the cache.
    expect(createCacheCalls).toBeGreaterThanOrEqual(2);

    (global as any).window = originalWindow;
  });

  it('reuses the same context across calls on the client', async () => {
    // jsdom provides a window by default in this test, so we just import.
    const { default: getPageContext } = await import('./PageContext');
    const first = getPageContext();
    const second = getPageContext();

    expect(first).toBe(second);
    // Only one client cache should have been created.
    expect(createCacheCalls).toBe(1);
  });
});
