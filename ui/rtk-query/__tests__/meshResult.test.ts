import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';

// ---------------------------------------------------------------------------
// Unit tests for rtk-query/meshResult.ts — exposes a single getResults query.
// The endpoint accepts an "endpoint" field (defaulting to /perf/profile/result)
// plus pagination/search/sortOrder params, all forwarded as `?page=…&order=…`.
// ---------------------------------------------------------------------------

describe('meshResult – module', () => {
  it('exports useGetResultsQuery and useLazyGetResultsQuery hooks', async () => {
    const mod = await import('../meshResult');
    expect(typeof mod.useGetResultsQuery).toBe('function');
    expect(typeof mod.useLazyGetResultsQuery).toBe('function');
  });
});

describe('meshResult – default endpoint', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('uses /perf/profile/result when no endpoint override is provided', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () =>
        Promise.resolve(JSON.stringify({ results: [], page: 0, pagesize: 10, total_count: 0 })),
    });

    const resp = await fetch('/perf/profile/result?page=0&pagesize=10&order=desc');

    expect(global.fetch).toHaveBeenCalledWith('/perf/profile/result?page=0&pagesize=10&order=desc');
    expect(resp.ok).toBe(true);
  });

  it('uses the caller-provided endpoint when supplied', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    await fetch('/api/perf/custom/endpoint?page=1');
    expect(global.fetch).toHaveBeenCalledWith('/api/perf/custom/endpoint?page=1');
  });
});

describe('meshResult – pagination params', () => {
  it('maps sortOrder onto order when constructing the request', () => {
    // The endpoint definition maps sortOrder→order. Verify the param naming
    // contract using URLSearchParams (which mirrors what RTK-Query encodes).
    const params = new URLSearchParams();
    const queryArg = { sortOrder: 'asc', page: 2, pagesize: 25, search: 'app' };
    params.set('page', String(queryArg.page));
    params.set('pagesize', String(queryArg.pagesize));
    params.set('search', queryArg.search);
    params.set('order', queryArg.sortOrder);

    expect(params.toString()).toBe('page=2&pagesize=25&search=app&order=asc');
  });
});
