import { describe, it, expect, vi, beforeEach } from 'vitest';

describe('sessionInterceptor', () => {
  let mockDispatch: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    // Reset module state between tests so each test gets a fresh interceptor
    vi.resetModules();
    mockDispatch = vi.fn();
  });

  async function setupAndFetch(mockResponse: Partial<Response>) {
    // Set up the mock fetch BEFORE importing the interceptor
    global.fetch = vi.fn().mockResolvedValue(mockResponse) as unknown as typeof fetch;

    // Dynamic import to get a fresh module (vi.resetModules clears the cache)
    const { installSessionInterceptor } = await import('../sessionInterceptor');
    installSessionInterceptor({ dispatch: mockDispatch });

    // Call the now-intercepted fetch
    return fetch('/api/test');
  }

  it('dispatches SESSION_EXPIRED on 401 response', async () => {
    await setupAndFetch({ status: 401, redirected: false });
    expect(mockDispatch).toHaveBeenCalledWith({ type: 'SESSION_EXPIRED' });
  });

  it('dispatches SESSION_EXPIRED on redirected response', async () => {
    await setupAndFetch({ status: 200, redirected: true });
    expect(mockDispatch).toHaveBeenCalledWith({ type: 'SESSION_EXPIRED' });
  });

  it('does not dispatch on successful response', async () => {
    await setupAndFetch({ status: 200, redirected: false });
    expect(mockDispatch).not.toHaveBeenCalled();
  });

  it('returns the original response unchanged', async () => {
    const mockResponse = { status: 200, redirected: false, body: 'data' };
    const result = await setupAndFetch(mockResponse as Partial<Response>);
    expect(result).toEqual(mockResponse);
  });
});
