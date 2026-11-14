import React from 'react';
import '@testing-library/jest-dom';
import { configureStore } from '@reduxjs/toolkit';
import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { Provider } from 'react-redux';
import { SistentThemeProvider } from '@sistent/sistent';
import { afterAll, afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

// The schemas client reads its base URL at module load, so this is set before
// any import of it, and restored afterwards so suite order cannot start to matter.
const { previousEndpointPrefix, canViewSpy } = vi.hoisted(() => {
  const previous = process.env.RTK_MESHERY_ENDPOINT_PREFIX;
  process.env.RTK_MESHERY_ENDPOINT_PREFIX = 'http://localhost';
  return { previousEndpointPrefix: previous, canViewSpy: vi.fn(() => true) };
});

afterAll(() => {
  if (previousEndpointPrefix === undefined) {
    delete process.env.RTK_MESHERY_ENDPOINT_PREFIX;
  } else {
    process.env.RTK_MESHERY_ENDPOINT_PREFIX = previousEndpointPrefix;
  }
});

vi.mock('@sistent/sistent', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@sistent/sistent')>();
  return { ...actual, useHasPermission: () => canViewSpy() };
});

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
  useNotificationHandlers: () => ({ notifyApiError: vi.fn() }),
}));
vi.mock('../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
  useNotificationHandlers: () => ({ notifyApiError: vi.fn() }),
}));

// The create/edit modal and the delete prompt are not under test and pull in
// the RJSF form stack, so they are cut off here. DefaultError reaches the org
// switcher's Next.js-only imports; its stub still marks where it rendered.
vi.mock('./general/error-404/index', () => ({
  default: () => <div data-testid="credentials-access-denied" />,
}));
vi.mock('./shared/Modal/Modal', () => ({ default: () => null }));
vi.mock('./general/PromptComponent', () => ({
  default: React.forwardRef(function PromptStub() {
    return null;
  }),
}));

import { api } from '@/rtk-query/index';
import MesheryCredentialComponent from './MesheryCredentialComponent';

const TOTAL_CREDENTIALS = 12;
const SERVER_DEFAULT_PAGE_SIZE = 10;

const seedCredentials = () =>
  Array.from({ length: TOTAL_CREDENTIALS }, (_, i) => {
    const n = String(i + 1).padStart(2, '0');
    return {
      id: `00000000-0000-0000-0000-0000000000${n}`,
      name: `cred-${n}`,
      type: 'prometheus',
      createdAt: '2026-11-11T10:00:00Z',
      updatedAt: '2026-11-11T10:00:00Z',
    };
  });

// The server's current credentials; a test can shrink it to mimic a delete.
let storedCredentials = seedCredentials();

// Answers like GET /api/integrations/credentials: zero-based `page`, a page size
// that defaults to 10 when none is sent, and the total across all pages.
const credentialsResponse = (url: URL) => {
  const page = Number(url.searchParams.get('page') ?? 0);
  const pageSize = Number(url.searchParams.get('pageSize') ?? SERVER_DEFAULT_PAGE_SIZE);
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

const renderPage = () => {
  const store = configureStore({
    reducer: {
      [api.reducerPath]: api.reducer,
      ui: () => ({ connectionMetadataState: {} }),
    },
    middleware: (getDefault) =>
      getDefault({ serializableCheck: false, immutableCheck: false }).concat(api.middleware),
  });
  render(
    <Provider store={store}>
      <SistentThemeProvider>
        <MesheryCredentialComponent />
      </SistentThemeProvider>
    </Provider>,
  );
  return store;
};

describe('MesheryCredentialComponent pagination', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    canViewSpy.mockReturnValue(true);
    storedCredentials = seedCredentials();
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

  it('requests the first page explicitly and counts every credential', async () => {
    renderPage();

    expect(await screen.findByText('cred-01')).toBeInTheDocument();

    const [first] = credentialRequests(fetchMock);
    expect(first.searchParams.get('page')).toBe('0');
    expect(first.searchParams.get('pageSize')).toBe('10');

    // The total comes from the server, not from the rows that were fetched.
    expect(screen.getByText(/1-10 of 12/)).toBeInTheDocument();
    expect(screen.queryByText('cred-11')).not.toBeInTheDocument();
  });

  it('fetches the next page from the server when the user pages forward', async () => {
    renderPage();
    expect(await screen.findByText('cred-01')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: /next page/i }));

    expect(await screen.findByText('cred-12')).toBeInTheDocument();
    expect(screen.getByText('cred-11')).toBeInTheDocument();
    expect(screen.queryByText('cred-01')).not.toBeInTheDocument();
    expect(screen.getByText(/11-12 of 12/)).toBeInTheDocument();

    await waitFor(() => {
      const pages = credentialRequests(fetchMock).map((url) => url.searchParams.get('page'));
      expect(pages).toContain('1');
    });
  });

  it('offers no page size above the 25 the endpoint serves', async () => {
    renderPage();
    expect(await screen.findByText('cred-01')).toBeInTheDocument();

    fireEvent.mouseDown(screen.getByRole('combobox'));

    const sizes = (await screen.findAllByRole('option')).map((option) => option.textContent);
    expect(sizes).toEqual(['10', '20', '25']);
  });

  it('steps back a page when the last page empties after a delete', async () => {
    const store = renderPage();
    expect(await screen.findByText('cred-01')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: /next page/i }));
    expect(await screen.findByText('cred-12')).toBeInTheDocument();

    // Both credentials on page 2 are deleted; the delete mutation invalidates
    // this tag, which is what refetches the list.
    storedCredentials = storedCredentials.slice(0, 10);
    act(() => {
      store.dispatch(api.util.invalidateTags(['credential_credentials']));
    });

    expect(await screen.findByText('cred-01')).toBeInTheDocument();
    expect(screen.getByText(/1-10 of 10/)).toBeInTheDocument();
    expect(screen.queryByText('cred-12')).not.toBeInTheDocument();
  });

  it('issues no request and shows the access error when the view key is missing', async () => {
    canViewSpy.mockReturnValue(false);
    renderPage();

    expect(screen.getByTestId('credentials-access-denied')).toBeInTheDocument();
    expect(screen.queryByRole('table')).not.toBeInTheDocument();
    expect(credentialRequests(fetchMock)).toHaveLength(0);
  });
});

describe('MesheryCredentialComponent dates', () => {
  const NOW = new Date('2026-11-13T12:00:00Z');
  let previousWidth: number;

  beforeEach(() => {
    canViewSpy.mockReturnValue(true);
    // Freeze only Date, so a cell showing "now" is unambiguous while RTK Query
    // and testing-library keep their real timers.
    vi.useFakeTimers({ toFake: ['Date'] });
    vi.setSystemTime(NOW);
    // The date columns are only shown from the xl breakpoint (1140px).
    previousWidth = window.innerWidth;
    window.innerWidth = 1280;
    storedCredentials = [
      {
        id: '00000000-0000-0000-0000-000000000099',
        name: 'cred-dated',
        type: 'prometheus',
        createdAt: '2025-11-13T12:00:00Z',
        updatedAt: '2026-11-10T12:00:00Z',
      },
    ];
    vi.spyOn(globalThis, 'fetch').mockImplementation(async (req) =>
      credentialsResponse(new URL((req as Request).url)),
    );
  });

  afterEach(() => {
    vi.useRealTimers();
    window.innerWidth = previousWidth;
    vi.restoreAllMocks();
  });

  it("shows the credential's own creation and update dates, not the current time", async () => {
    renderPage();

    const row = (await screen.findByText('cred-dated')).closest('tr') as HTMLElement;
    expect(within(row).getByText('a year ago')).toBeInTheDocument();
    expect(within(row).getByText('3 days ago')).toBeInTheDocument();
  });
});
