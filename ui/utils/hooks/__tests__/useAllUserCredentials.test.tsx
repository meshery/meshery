import React from 'react';
import { configureStore } from '@reduxjs/toolkit';
import { act, renderHook, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { afterAll, afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

// The schemas client reads its base URL at module load, so this is set before
// any import of it, and restored afterwards so suite order cannot start to matter.
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

import { api } from '@/rtk-query/index';
import { CREDENTIALS_PAGE_SIZE, useAllUserCredentials } from '../useAllUserCredentials';

const makeCredentials = (count: number) =>
  Array.from({ length: count }, (_, i) => {
    const n = String(i + 1).padStart(3, '0');
    return {
      id: `00000000-0000-0000-0000-000000000${n}`,
      name: `cred-${n}`,
      type: 'grafana',
    };
  });

// The server's credentials, newest first, and the most it serves per page. A
// server that ignores the requested pageSize serves its default of 10.
let storedCredentials = makeCredentials(0);
let servedPageSize = CREDENTIALS_PAGE_SIZE;

// Answers like GET /api/integrations/credentials: zero-based `page`, and a page
// size of at most `servedPageSize`, reported back along with the total.
const credentialsResponse = (url: URL) => {
  const page = Number(url.searchParams.get('page') ?? 0);
  const requested = Number(url.searchParams.get('pageSize') ?? 10);
  const pageSize = Math.min(requested, servedPageSize);
  const body = {
    credentials: storedCredentials.slice(page * pageSize, (page + 1) * pageSize),
    totalCount: storedCredentials.length,
    page,
    pageSize,
  };
  return new Response(JSON.stringify(body), {
    status: 200,
    headers: { 'content-type': 'application/json' },
  });
};

const requestedPages = (fetchMock: ReturnType<typeof vi.fn>) =>
  fetchMock.mock.calls
    .map(([req]) => new URL((req as Request).url))
    .filter((url) => url.pathname === '/api/integrations/credentials')
    .map((url) => `${url.searchParams.get('page')}/${url.searchParams.get('pageSize')}`);

const renderCredentials = (props: { skip?: boolean } = {}) => {
  const store = configureStore({
    reducer: { [api.reducerPath]: api.reducer },
    middleware: (getDefault) =>
      getDefault({ serializableCheck: false, immutableCheck: false }).concat(api.middleware),
  });
  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <Provider store={store}>{children}</Provider>
  );
  const hook = renderHook(() => useAllUserCredentials(props), { wrapper });
  return { ...hook, store };
};

const names = (credentials: { name?: string }[]) => credentials.map((c) => c.name);

describe('useAllUserCredentials', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    storedCredentials = makeCredentials(0);
    servedPageSize = CREDENTIALS_PAGE_SIZE;
    // spyOn so restoreAllMocks reverts it rather than leaking into later suites.
    fetchMock = vi
      .spyOn(globalThis, 'fetch')
      .mockImplementation(async (req) =>
        credentialsResponse(new URL((req as Request).url)),
      ) as unknown as ReturnType<typeof vi.fn>;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('returns credentials beyond the first page', async () => {
    storedCredentials = makeCredentials(30);
    const { result } = renderCredentials();

    await waitFor(() => expect(result.current.credentials).toHaveLength(30));
    expect(result.current.isLoading).toBe(false);
    expect(names(result.current.credentials)).toContain('cred-030');
    expect(requestedPages(fetchMock).sort()).toEqual(['0/25', '1/25']);
  });

  it('pages by the size the server served when it ignores the requested size', async () => {
    // The 11th-newest credential is the one a server returning 10 hides.
    storedCredentials = makeCredentials(11);
    servedPageSize = 10;
    const { result } = renderCredentials();

    await waitFor(() => expect(result.current.credentials).toHaveLength(11));
    expect(names(result.current.credentials)).toContain('cred-011');
    expect(requestedPages(fetchMock).sort()).toEqual(['0/25', '1/25']);
  });

  it('makes a single request when everything fits on one page', async () => {
    storedCredentials = makeCredentials(3);
    const { result } = renderCredentials();

    await waitFor(() => expect(result.current.credentials).toHaveLength(3));
    expect(requestedPages(fetchMock)).toEqual(['0/25']);
  });

  it('picks up a change on a later page when the credential tag is invalidated', async () => {
    storedCredentials = makeCredentials(30);
    const { result, store } = renderCredentials();
    await waitFor(() => expect(result.current.credentials).toHaveLength(30));

    // The oldest credential, on page 2, is deleted; the delete mutation
    // invalidates this tag, which is what refetches the list.
    storedCredentials = storedCredentials.slice(0, 29);
    act(() => {
      store.dispatch(api.util.invalidateTags(['credential_credentials']));
    });

    await waitFor(() => expect(result.current.credentials).toHaveLength(29));
    expect(names(result.current.credentials)).not.toContain('cred-030');
  });

  it('issues no request while skipped', async () => {
    storedCredentials = makeCredentials(30);
    const { result } = renderCredentials({ skip: true });

    // Give a skipped query the chance to fire before asserting it did not.
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 50));
    });
    expect(result.current.credentials).toEqual([]);
    expect(requestedPages(fetchMock)).toHaveLength(0);
  });
});
