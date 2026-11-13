import React from 'react';
import { configureStore } from '@reduxjs/toolkit';
import { renderHook, waitFor } from '@testing-library/react';
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

// The step registry pulls in the whole step UI; this suite only needs the
// credentials the wizard hands its steps.
vi.mock('./registry', () => ({
  buildSteps: () => [{ id: 'select', label: 'Choose Connection', Component: () => null }],
}));
vi.mock('./useKindPermission', () => ({ useKindPermission: () => () => true }));
vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

import { api } from '@/rtk-query/index';
import { filterCredentialsForKind } from '../ConnectionWizard.helpers';
import { useConnectionWizard } from './useConnectionWizard';

// Newest first, as the server orders them: ten Grafana credentials created
// after the one Prometheus credential, which is therefore the 11th newest.
const storedCredentials = [
  ...Array.from({ length: 10 }, (_, i) => ({
    id: `00000000-0000-0000-0000-0000000000${String(10 - i).padStart(2, '0')}`,
    name: `grafana-${String(10 - i).padStart(2, '0')}`,
    type: 'grafana',
  })),
  { id: '00000000-0000-0000-0000-000000000099', name: 'prometheus-cred', type: 'prometheus' },
];

// GET /api/integrations/credentials as served by a server that applies its
// default page size of 10 whatever size is requested.
const credentialsResponse = (url: URL) => {
  const page = Number(url.searchParams.get('page') ?? 0);
  const pageSize = 10;
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

const credentialRequests = (fetchMock: ReturnType<typeof vi.fn>) =>
  fetchMock.mock.calls
    .map(([req]) => new URL((req as Request).url))
    .filter((url) => url.pathname === '/api/integrations/credentials');

const renderWizard = (isOpen: boolean) => {
  const store = configureStore({
    reducer: { [api.reducerPath]: api.reducer },
    middleware: (getDefault) =>
      getDefault({ serializableCheck: false, immutableCheck: false }).concat(api.middleware),
  });
  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <Provider store={store}>{children}</Provider>
  );
  return renderHook(() => useConnectionWizard({ mode: 'create', isOpen, availableKinds: [] }), {
    wrapper,
  });
};

describe('useConnectionWizard existing credentials', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi
      .spyOn(globalThis, 'fetch')
      .mockImplementation(async (req) =>
        credentialsResponse(new URL((req as Request).url)),
      ) as unknown as ReturnType<typeof vi.fn>;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('offers a credential that is not among the 10 newest', async () => {
    const { result } = renderWizard(true);

    await waitFor(() => expect(result.current.ctx.services.credentials).toHaveLength(11));

    const prometheus = filterCredentialsForKind(
      result.current.ctx.services.credentials,
      'prometheus',
    );
    expect(prometheus.map((credential) => credential.name)).toEqual(['prometheus-cred']);
  });

  it('fetches no credentials while the wizard is closed', async () => {
    const { result } = renderWizard(false);

    await new Promise((resolve) => setTimeout(resolve, 50));
    expect(result.current.ctx.services.credentials).toEqual([]);
    expect(credentialRequests(fetchMock)).toHaveLength(0);
  });
});
