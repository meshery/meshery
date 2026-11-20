import React from 'react';
import { afterAll, afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { act, renderHook, waitFor } from '@testing-library/react';

// `../utils` imports the app store at module scope; stub it before loading
// RTK Query modules so these tests can create an isolated store per setup.
vi.mock('../../store', () => ({ store: { dispatch: vi.fn() } }));

const { previousEndpointPrefix } = vi.hoisted(() => {
  const previous = process.env.RTK_MESHERY_ENDPOINT_PREFIX;
  process.env.RTK_MESHERY_ENDPOINT_PREFIX = 'http://localhost';
  return { previousEndpointPrefix: previous };
});

afterAll(() => {
  if (previousEndpointPrefix === undefined) {
    delete process.env.RTK_MESHERY_ENDPOINT_PREFIX;
  } else {
    process.env.RTK_MESHERY_ENDPOINT_PREFIX = previousEndpointPrefix;
  }
});

const okJsonResponse = (body: unknown) => ({
  ok: true,
  status: 200,
  redirected: false,
  headers: new Headers({ 'content-type': 'application/json' }),
  url: '',
  text: () => Promise.resolve(JSON.stringify(body)),
  json: () => Promise.resolve(body),
  clone() {
    return this;
  },
});

const setup = async () => {
  vi.resetModules();
  const apiMod = await import('../index');
  const meshModelMod = await import('../meshModel');
  const { configureStore } = await import('@reduxjs/toolkit');
  const { Provider } = await import('react-redux');

  const store = configureStore({
    reducer: { [apiMod.api.reducerPath]: apiMod.api.reducer },
    middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(apiMod.api.middleware),
  });

  const wrapper = ({ children }: { children: React.ReactNode }) =>
    React.createElement(Provider, { store }, children);

  return { apiMod, meshModelMod, store, wrapper };
};

describe('useGetCategoriesSummary', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn(async (input: RequestInfo | URL) => {
      const url =
        typeof input === 'string'
          ? input
          : typeof input === 'object' && input !== null && 'url' in input
            ? String((input as { url: string }).url)
            : String(input);
      if (url.includes('/api/registry/categories/Networking/models')) {
        return okJsonResponse({ totalCount: 5 });
      }
      if (url.includes('/api/registry/categories/Security/models')) {
        return okJsonResponse({ totalCount: 3 });
      }
      if (url.includes('/api/registry/categories')) {
        return okJsonResponse({
          categories: [{ name: 'Networking' }, { name: 'Security' }],
        });
      }
      throw new Error(`Unexpected request URL in test: ${url}`);
    });
    vi.spyOn(globalThis, 'fetch').mockImplementation(fetchMock as unknown as typeof fetch);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('fetches categories and per-category model counts by default (skip omitted)', async () => {
    const { meshModelMod, wrapper } = await setup();

    const { result } = renderHook(() => meshModelMod.useGetCategoriesSummary(), { wrapper });

    await waitFor(() => {
      expect(result.current).toEqual({ Networking: 5, Security: 3 });
    });
    expect(fetchMock).toHaveBeenCalledTimes(3);
  });

  it('issues zero category or model requests when { skip: true }', async () => {
    const { meshModelMod, wrapper } = await setup();

    const { result } = renderHook(() => meshModelMod.useGetCategoriesSummary({ skip: true }), {
      wrapper,
    });

    // Allow effects to settle and verify no fetch was initiated.
    await new Promise((resolve) => setTimeout(resolve, 20));
    expect(result.current).toEqual({});
    expect(fetchMock).not.toHaveBeenCalled();
  });

  it('clears cached category counts and skips refetches when skip flips from false to true', async () => {
    const { meshModelMod, wrapper } = await setup();

    const { result, rerender } = renderHook(
      ({ skip }: { skip: boolean }) => meshModelMod.useGetCategoriesSummary({ skip }),
      {
        initialProps: { skip: false },
        wrapper,
      },
    );

    await waitFor(() => {
      expect(result.current).toEqual({ Networking: 5, Security: 3 });
    });
    const callCountBeforeDeny = fetchMock.mock.calls.length;

    rerender({ skip: true });

    await waitFor(() => {
      expect(result.current).toEqual({});
    });
    expect(fetchMock).toHaveBeenCalledTimes(callCountBeforeDeny);
  });

  it('discards in-flight per-category model responses when skip flips from false to true', async () => {
    const { meshModelMod, wrapper } = await setup();
    let resolveNetworking!: (value: ReturnType<typeof okJsonResponse>) => void;
    let resolveSecurity!: (value: ReturnType<typeof okJsonResponse>) => void;
    const networkingPromise = new Promise<ReturnType<typeof okJsonResponse>>((resolve) => {
      resolveNetworking = resolve;
    });
    const securityPromise = new Promise<ReturnType<typeof okJsonResponse>>((resolve) => {
      resolveSecurity = resolve;
    });

    fetchMock.mockImplementation(async (input: RequestInfo | URL) => {
      const url =
        typeof input === 'string'
          ? input
          : typeof input === 'object' && input !== null && 'url' in input
            ? String((input as { url: string }).url)
            : String(input);
      if (url.includes('/api/registry/categories/Networking/models')) {
        return networkingPromise;
      }
      if (url.includes('/api/registry/categories/Security/models')) {
        return securityPromise;
      }
      if (url.includes('/api/registry/categories')) {
        return okJsonResponse({
          categories: [{ name: 'Networking' }, { name: 'Security' }],
        });
      }
      throw new Error(`Unexpected request URL in test: ${url}`);
    });

    const { result, rerender } = renderHook(
      ({ skip }: { skip: boolean }) => meshModelMod.useGetCategoriesSummary({ skip }),
      {
        initialProps: { skip: false },
        wrapper,
      },
    );

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledTimes(3);
    });

    await act(async () => {
      rerender({ skip: true });
      resolveNetworking(okJsonResponse({ totalCount: 5 }));
      resolveSecurity(okJsonResponse({ totalCount: 3 }));
    });

    await new Promise((resolve) => setTimeout(resolve, 20));
    expect(result.current).toEqual({});
    expect(fetchMock).toHaveBeenCalledTimes(3);
  });

  it('unsubscribes category queries so Registry_Registry invalidation does not refetch while skipped', async () => {
    const { apiMod, meshModelMod, store, wrapper } = await setup();

    const { result, rerender } = renderHook(
      ({ skip }: { skip: boolean }) => meshModelMod.useGetCategoriesSummary({ skip }),
      {
        initialProps: { skip: false },
        wrapper,
      },
    );

    await waitFor(() => {
      expect(result.current).toEqual({ Networking: 5, Security: 3 });
    });
    expect(fetchMock).toHaveBeenCalledTimes(3);

    act(() => {
      rerender({ skip: true });
    });
    await waitFor(() => {
      expect(result.current).toEqual({});
    });

    await act(async () => {
      store.dispatch(apiMod.api.util.invalidateTags(['Registry_Registry']));
      await new Promise((resolve) => setTimeout(resolve, 20));
    });

    expect(result.current).toEqual({});
    expect(fetchMock).toHaveBeenCalledTimes(3);

    act(() => {
      rerender({ skip: false });
    });
    await waitFor(() => {
      expect(result.current).toEqual({ Networking: 5, Security: 3 });
    });
    expect(fetchMock).toHaveBeenCalledTimes(6);
  });
});
