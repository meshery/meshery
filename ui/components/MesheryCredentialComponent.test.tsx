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
import customValidator from '../utils/rjsfValidator';

vi.mock('./general/error-404/index', () => ({
  default: () => <div data-testid="credentials-access-denied" />,
}));
vi.mock('./shared/Modal/Modal', () => ({
  default: (props: any) => {
    const [localData, setLocalData] = React.useState<Record<string, any>>(props.formData || {});
    React.useEffect(() => {
      setLocalData(props.formData || {});
    }, [props.formData, props.type]);
    if (!props.open) return null;
    const handleSave = () => {
      const validation = customValidator.validateFormData(localData, props.schema);
      if (validation.errors.length === 0) {
        props.handleSubmit?.(localData);
      }
    };
    const updateLocal = (next: Record<string, any>) => {
      setLocalData(next);
      props.onChange?.(next);
    };
    return (
      <div
        data-testid="credential-modal"
        data-title={props.title}
        data-type={props.type}
        data-schema-count={props.schema_array?.length}
        data-has-schema={Boolean(props.schema)}
        data-submit-disabled={Boolean(props.submitDisabled)}
        data-close-on-submit={String(props.closeOnSubmit)}
        data-form-data={JSON.stringify(props.formData || {})}
      >
        <button data-testid="modal-save-btn" disabled={props.submitDisabled} onClick={handleSave}>
          {props.submitBtnText || 'Save'}
        </button>
        <button data-testid="modal-close-btn" onClick={(e) => props.handleClose?.(e)}>
          Close
        </button>
        <button
          data-testid="modal-switch-grafana-btn"
          onClick={() => props.schemaChangeHandler?.('grafana')}
        >
          Switch Grafana
        </button>
        <button
          data-testid="modal-switch-prom-btn"
          onClick={() => props.schemaChangeHandler?.('prometheus')}
        >
          Switch Prometheus
        </button>
        <button
          data-testid="modal-fill-grafana-no-key-btn"
          onClick={() =>
            updateLocal({
              name: localData?.name || 'prom-cred',
              secret: { grafanaURL: 'https://grafana.new.example' },
            })
          }
        >
          Fill Grafana URL only
        </button>
        <button
          data-testid="modal-fill-grafana-with-key-btn"
          onClick={() =>
            updateLocal({
              name: localData?.name || 'prom-cred',
              secret: {
                grafanaURL: 'https://grafana.new.example',
                grafanaAPIKey: 'new-grafana-key',
              },
            })
          }
        >
          Fill Grafana URL and Key
        </button>
      </div>
    );
  },
}));
vi.mock('./general/PromptComponent', () => ({
  default: React.forwardRef(function PromptStub() {
    return null;
  }),
}));

import { api } from '@/rtk-query/index';
import MesheryCredentialComponent, {
  cleanEmptyPasswordFields,
  sanitizeInitialData,
} from './MesheryCredentialComponent';
import { CredentialRjsfSchemasV1Beta1, CredentialRjsfUiSchemasV1Beta1 } from '@meshery/schemas';

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

describe('MesheryCredentialComponent password and secret sanitation', () => {
  it('redacts writeOnly / password fields from initialData when editing', () => {
    const k8sData = {
      clusterName: 'production',
      clusterServerURL: 'https://k8s.example.com',
      auth: {
        clusterUserName: 'admin',
        clusterToken: 'super-secret-token',
        clusterClientCertificateData: 'cert',
        clusterClientKeyData: 'key',
        clusterCertificateAuthorityData: 'ca',
      },
    };
    const sanitized = sanitizeInitialData(
      k8sData,
      CredentialRjsfSchemasV1Beta1['kubernetes'].properties.secret,
      CredentialRjsfUiSchemasV1Beta1['kubernetes'].secret,
    );

    expect(sanitized.clusterName).toBe('production');
    expect(sanitized.clusterServerURL).toBe('https://k8s.example.com');
    expect(sanitized.auth.clusterUserName).toBe('admin');
    expect(sanitized.auth.clusterToken).toBeUndefined();
    expect(sanitized.auth.clusterClientCertificateData).toBeUndefined();
    expect(sanitized.auth.clusterClientKeyData).toBeUndefined();
  });

  it('omits empty password fields on edit so stored secrets are not blanked out', () => {
    const editPayload = {
      clusterName: 'production-v2',
      clusterServerURL: 'https://k8s.example.com',
      auth: {
        clusterUserName: 'admin',
        clusterToken: '', // left blank by user
        clusterClientCertificateData: '',
      },
    };
    const cleaned = cleanEmptyPasswordFields(
      editPayload,
      CredentialRjsfSchemasV1Beta1['kubernetes'].properties.secret,
      CredentialRjsfUiSchemasV1Beta1['kubernetes'].secret,
    );

    expect(cleaned.clusterName).toBe('production-v2');
    expect(cleaned.auth.clusterUserName).toBe('admin');
    expect(cleaned.auth.clusterToken).toBeUndefined();
    expect(cleaned.auth.clusterClientCertificateData).toBeUndefined();
  });

  it('retains password fields when the user explicitly provides a new value', () => {
    const editPayload = {
      grafanaURL: 'https://grafana.example.com',
      grafanaAPIKey: 'new-api-key-999',
    };
    const cleaned = cleanEmptyPasswordFields(
      editPayload,
      CredentialRjsfSchemasV1Beta1['grafana'].properties.secret,
      CredentialRjsfUiSchemasV1Beta1['grafana'].secret,
    );

    expect(cleaned.grafanaURL).toBe('https://grafana.example.com');
    expect(cleaned.grafanaAPIKey).toBe('new-api-key-999');
  });
});

describe('MesheryCredentialComponent schema forms and toolbar actions', () => {
  beforeEach(() => {
    canViewSpy.mockReturnValue(true);
    storedCredentials = seedCredentials();
    vi.spyOn(globalThis, 'fetch').mockImplementation(async (req) =>
      credentialsResponse(new URL((req as Request).url)),
    );
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('renders the Create button and opens modal with all canonical schemas', async () => {
    renderPage();
    expect(await screen.findByText('cred-01')).toBeInTheDocument();

    const createBtn = screen.getByRole('button', { name: /create/i });
    expect(createBtn).toBeInTheDocument();

    fireEvent.click(createBtn);

    const modal = screen.getByTestId('credential-modal');
    expect(modal).toBeInTheDocument();
    expect(modal).toHaveAttribute('data-title', 'Create Credential');
    expect(modal).toHaveAttribute('data-type', 'prometheus');
    const schemaCount = Number(modal.getAttribute('data-schema-count'));
    expect(schemaCount).toBeGreaterThanOrEqual(21);
    expect(modal).toHaveAttribute('data-has-schema', 'true');
    expect(modal).toHaveAttribute('data-submit-disabled', 'false');

    // Create schema still enforces required password fields (unlike edit mode).
    fireEvent.click(screen.getByTestId('modal-switch-grafana-btn'));
    expect(screen.getByTestId('credential-modal')).toHaveAttribute('data-type', 'grafana');
    const putOrPostCountBefore = storedCredentials.length;
    fireEvent.click(screen.getByTestId('modal-save-btn'));
    expect(storedCredentials.length).toBe(putOrPostCountBefore);
  });

  it('renders Edit button on each credential row and opens Edit modal', async () => {
    renderPage();
    expect(await screen.findByText('cred-01')).toBeInTheDocument();

    const editButtons = screen.getAllByRole('button', { name: /edit/i });
    expect(editButtons.length).toBeGreaterThan(0);

    fireEvent.click(editButtons[0]);

    const modal = screen.getByTestId('credential-modal');
    expect(modal).toBeInTheDocument();
    expect(modal).toHaveAttribute('data-title', 'Edit Credential');
    expect(modal).toHaveAttribute('data-close-on-submit', 'false');
  });

  it('submits an edited Grafana credential with unchanged password omitted through schema validation', async () => {
    const putPayloads: any[] = [];
    storedCredentials = [
      {
        id: '00000000-0000-0000-0000-000000000101',
        name: 'my-grafana',
        type: 'grafana',
        secret: {
          grafanaURL: 'https://grafana.example.com',
          grafanaAPIKey: 'stored-secret-key',
        },
        createdAt: '2026-11-11T10:00:00Z',
        updatedAt: '2026-11-11T10:00:00Z',
      } as any,
    ];
    vi.restoreAllMocks();
    vi.spyOn(globalThis, 'fetch').mockImplementation(async (req) => {
      const request = req as Request;
      if (request.method === 'PUT') {
        putPayloads.push(await request.clone().json());
        return new Response(JSON.stringify({}), {
          status: 200,
          headers: { 'content-type': 'application/json' },
        });
      }
      return credentialsResponse(new URL(request.url));
    });

    renderPage();
    expect(await screen.findByText('my-grafana')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: /edit/i }));
    const modal = screen.getByTestId('credential-modal');
    const initialForm = JSON.parse(modal.getAttribute('data-form-data') || '{}');
    expect(initialForm.secret.grafanaAPIKey).toBeUndefined();
    expect(initialForm.secret.grafanaURL).toBe('https://grafana.example.com');

    fireEvent.click(screen.getByTestId('modal-save-btn'));

    await waitFor(() => {
      expect(putPayloads).toHaveLength(1);
    });
    expect(putPayloads[0]).toEqual({
      id: '00000000-0000-0000-0000-000000000101',
      name: 'my-grafana',
      type: 'grafana',
      secret: { grafanaURL: 'https://grafana.example.com' },
    });
  });

  it('clears the old secret while keeping the credential name when switching types during edit', async () => {
    storedCredentials = [
      {
        id: '00000000-0000-0000-0000-000000000102',
        name: 'prom-cred',
        type: 'prometheus',
        secret: { prometheusURL: 'https://prom.example.com' },
        createdAt: '2026-11-11T10:00:00Z',
        updatedAt: '2026-11-11T10:00:00Z',
      } as any,
    ];
    renderPage();
    expect(await screen.findByText('prom-cred')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: /edit/i }));
    fireEvent.click(screen.getByTestId('modal-switch-grafana-btn'));

    const modal = screen.getByTestId('credential-modal');
    expect(modal).toHaveAttribute('data-type', 'grafana');
    const switchedForm = JSON.parse(modal.getAttribute('data-form-data') || '{}');
    expect(switchedForm).toEqual({ name: 'prom-cred', secret: {} });
  });

  it('requires new-type password fields after switching credential type during edit', async () => {
    const putPayloads: any[] = [];
    storedCredentials = [
      {
        id: '00000000-0000-0000-0000-000000000104',
        name: 'prom-to-grafana',
        type: 'prometheus',
        secret: { prometheusURL: 'https://prom.example.com' },
        createdAt: '2026-11-11T10:00:00Z',
        updatedAt: '2026-11-11T10:00:00Z',
      } as any,
    ];
    vi.restoreAllMocks();
    vi.spyOn(globalThis, 'fetch').mockImplementation(async (req) => {
      const request = req as Request;
      if (request.method === 'PUT') {
        putPayloads.push(await request.clone().json());
        return new Response(JSON.stringify({}), {
          status: 200,
          headers: { 'content-type': 'application/json' },
        });
      }
      return credentialsResponse(new URL(request.url));
    });

    renderPage();
    expect(await screen.findByText('prom-to-grafana')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: /edit/i }));
    fireEvent.click(screen.getByTestId('modal-switch-grafana-btn'));

    // Switching type retains full creation requirements for Grafana: URL without API key must fail validation.
    fireEvent.click(screen.getByTestId('modal-fill-grafana-no-key-btn'));
    fireEvent.click(screen.getByTestId('modal-save-btn'));
    expect(putPayloads).toHaveLength(0);
    expect(screen.getByTestId('credential-modal')).toBeInTheDocument();

    // Providing both URL and the required Grafana API key passes validation and submits.
    fireEvent.click(screen.getByTestId('modal-fill-grafana-with-key-btn'));
    fireEvent.click(screen.getByTestId('modal-save-btn'));

    await waitFor(() => {
      expect(putPayloads).toHaveLength(1);
    });
    expect(putPayloads[0]).toEqual({
      id: '00000000-0000-0000-0000-000000000104',
      name: 'prom-to-grafana',
      type: 'grafana',
      secret: {
        grafanaURL: 'https://grafana.new.example',
        grafanaAPIKey: 'new-grafana-key',
      },
    });
  });

  it('keeps the modal open on failed save and prevents a slow save from closing a newly opened modal', async () => {
    let resolveSlowSave!: (res: Response) => void;
    let callCount = 0;
    storedCredentials = [
      {
        id: '00000000-0000-0000-0000-000000000103',
        name: 'prom-edit',
        type: 'prometheus',
        secret: { prometheusURL: 'https://prom.example.com' },
        createdAt: '2026-11-11T10:00:00Z',
        updatedAt: '2026-11-11T10:00:00Z',
      } as any,
    ];
    vi.restoreAllMocks();
    vi.spyOn(globalThis, 'fetch').mockImplementation(async (req) => {
      const request = req as Request;
      if (request.method === 'PUT') {
        callCount += 1;
        if (callCount === 1) {
          return new Response(JSON.stringify({ message: 'save failed' }), {
            status: 500,
            headers: { 'content-type': 'application/json' },
          });
        }
        return await new Promise<Response>((resolve) => {
          resolveSlowSave = resolve;
        });
      }
      return credentialsResponse(new URL(request.url));
    });

    renderPage();
    expect(await screen.findByText('prom-edit')).toBeInTheDocument();

    // 1. Failed save keeps the edit dialog open with its data intact.
    fireEvent.click(screen.getByRole('button', { name: /edit/i }));
    fireEvent.click(screen.getByTestId('modal-save-btn'));
    await waitFor(() => {
      expect(callCount).toBe(1);
    });
    expect(screen.getByTestId('credential-modal')).toBeInTheDocument();

    // 2. Start a slow save, close that modal, open a Create modal, then resolve the old save.
    fireEvent.click(screen.getByTestId('modal-save-btn'));
    await waitFor(() => {
      expect(callCount).toBe(2);
    });
    fireEvent.click(screen.getByTestId('modal-close-btn'));
    expect(screen.queryByTestId('credential-modal')).not.toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: /create/i }));
    expect(screen.getByTestId('credential-modal')).toHaveAttribute(
      'data-title',
      'Create Credential',
    );

    await act(async () => {
      resolveSlowSave(
        new Response(JSON.stringify({}), {
          status: 200,
          headers: { 'content-type': 'application/json' },
        }),
      );
    });

    // The newly opened Create dialog must remain open after the old save finishes.
    expect(screen.getByTestId('credential-modal')).toHaveAttribute(
      'data-title',
      'Create Credential',
    );
  });
});
