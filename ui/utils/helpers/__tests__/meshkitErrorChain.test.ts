import { configureStore } from '@reduxjs/toolkit';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { formatApiError } from '../meshkitError';

// The schemas baseQuery reads its endpoint prefix at module scope, and
// `fetchBaseQuery` builds an absolute `Request` before dispatching, so a
// relative default throws "Failed to parse URL". Set the prefix first, then
// pull in the real client - a dynamic import so this assignment wins over the
// hoisting of static imports.
process.env.RTK_MESHERY_ENDPOINT_PREFIX = 'http://localhost:9081';
const { mesheryApi } = await import('@meshery/schemas/mesheryApi');

// ---------------------------------------------------------------------------
// End-to-end coverage for the error chain the app actually runs:
//
//   server JSON envelope
//     -> real `@meshery/schemas` baseQuery (`withMeshkitErrorTransform`)
//       -> real `mesheryApi.endpoints.*`
//         -> `formatApiError`
//           -> snackbar markdown
//
// Deliberately NOT hand-constructing `error.meshkit`: driving the real wrapper
// is what catches a casing regression. The server emits camelCase
// (`probableCause`, ...); a casing mismatch in the wrapper once dropped the
// probable cause and the remediation list on the floor (meshery/schemas#1081,
// fixed in @meshery/schemas v1.3.37). The wrapper now reads camelCase with a
// snake_case fallback, and these tests guard that end-to-end so the gap cannot
// silently reopen on a future schemas bump.
// ---------------------------------------------------------------------------

/**
 * Verbatim shape of the `errorResponse` struct in
 * `server/models/httputil/httputil.go`, as emitted for a provider-rejected
 * Create Environment.
 */
const CREATE_ENVIRONMENT_403 = {
  error: 'Unable to create the environment',
  code: 'meshery-server-1448',
  severity: 'ALERT',
  probableCause: [
    'Your account does not have permission to create environments in this organization.',
    'An environment with the same name already exists in this organization.',
  ],
  suggestedRemediation: [
    'Confirm your role in the selected organization.',
    'The environment was NOT created - retry after resolving the cause.',
  ],
  longDescription: ['The remote provider rejected the create request.'],
};

/**
 * The same envelope emitted with the legacy snake_case detail keys. Some
 * producers still send these; `withMeshkitErrorTransform` reads camelCase
 * first and falls back to snake_case, so this MUST yield a MeshKit envelope
 * identical to the camelCase case above.
 */
const CREATE_ENVIRONMENT_403_SNAKE = {
  error: CREATE_ENVIRONMENT_403.error,
  code: CREATE_ENVIRONMENT_403.code,
  severity: CREATE_ENVIRONMENT_403.severity,
  probable_cause: CREATE_ENVIRONMENT_403.probableCause,
  suggested_remediation: CREATE_ENVIRONMENT_403.suggestedRemediation,
  long_description: CREATE_ENVIRONMENT_403.longDescription,
};

const makeStore = () =>
  configureStore({
    reducer: { [mesheryApi.reducerPath]: mesheryApi.reducer },
    middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(mesheryApi.middleware),
  });

/** Stub `fetch` so the real baseQuery parses a real JSON error response. */
const respondWith = (status: number, body: unknown) =>
  vi.stubGlobal(
    'fetch',
    vi.fn(
      async () =>
        new Response(JSON.stringify(body), {
          status,
          headers: { 'Content-Type': 'application/json' },
        }),
    ),
  );

/** Drive the real `createEnvironment` mutation and return its RTK error. */
const createEnvironmentError = async () => {
  const store = makeStore();
  const result = await store.dispatch(
    mesheryApi.endpoints.createEnvironment.initiate({
      body: { name: 'prod', description: 'production', organizationId: 'org-1' },
    }),
  );
  return (result as { error?: unknown }).error;
};

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('MeshKit error chain (real schemas client)', () => {
  it('populates the whole MeshKit envelope from a real camelCase 403', async () => {
    respondWith(403, CREATE_ENVIRONMENT_403);

    const error = (await createEnvironmentError()) as {
      status?: number;
      meshkit?: Record<string, unknown>;
    };

    // The wrapper attaches the full envelope: message, code and - since
    // @meshery/schemas v1.3.37 (meshery/schemas#1081) - the detail arrays,
    // read from the server's camelCase spellings. If a future schemas bump
    // reintroduces the casing mismatch, these assertions fail here rather than
    // silently dropping the probable cause and remediation list in production.
    expect(error?.status).toBe(403);
    expect(error?.meshkit?.message).toBe('Unable to create the environment');
    expect(error?.meshkit?.code).toBe('meshery-server-1448');
    expect(error?.meshkit?.suggestedRemediation).toEqual(
      CREATE_ENVIRONMENT_403.suggestedRemediation,
    );
    expect(error?.meshkit?.probableCause).toEqual(CREATE_ENVIRONMENT_403.probableCause);
    expect(error?.meshkit?.longDescription).toEqual(CREATE_ENVIRONMENT_403.longDescription);
  });

  it('populates the whole MeshKit envelope from a real snake_case 403', async () => {
    respondWith(403, CREATE_ENVIRONMENT_403_SNAKE);

    const error = (await createEnvironmentError()) as {
      status?: number;
      meshkit?: Record<string, unknown>;
    };

    // The wrapper reads camelCase first and falls back to the legacy
    // snake_case spellings, so a producer still emitting `probable_cause` /
    // `suggested_remediation` / `long_description` must yield the SAME complete
    // envelope as the camelCase case. Guards the fallback so a schemas
    // regression can't silently drop the detail arrays for such producers.
    expect(error?.status).toBe(403);
    expect(error?.meshkit?.message).toBe('Unable to create the environment');
    expect(error?.meshkit?.code).toBe('meshery-server-1448');
    expect(error?.meshkit?.suggestedRemediation).toEqual(
      CREATE_ENVIRONMENT_403.suggestedRemediation,
    );
    expect(error?.meshkit?.probableCause).toEqual(CREATE_ENVIRONMENT_403.probableCause);
    expect(error?.meshkit?.longDescription).toEqual(CREATE_ENVIRONMENT_403.longDescription);
  });

  it('renders the title, every remediation bullet and the code from a real 403', async () => {
    respondWith(403, CREATE_ENVIRONMENT_403);

    const formatted = formatApiError(
      await createEnvironmentError(),
      'Failed to create environment',
    );

    expect(formatted.message).toContain('**Unable to create the environment**');
    expect(formatted.message).toContain('*Try:*');
    for (const remediation of CREATE_ENVIRONMENT_403.suggestedRemediation) {
      expect(formatted.message).toContain(`- ${remediation}`);
    }
    expect(formatted.message).toContain('`meshery-server-1448`');

    // One bullet per remediation entry - not one run-on bullet.
    const bullets = formatted.message.split('\n').filter((line) => line.startsWith('- '));
    expect(bullets).toHaveLength(CREATE_ENVIRONMENT_403.suggestedRemediation.length);

    // The recovered envelope carries the details too, for non-markdown callers.
    expect(formatted.meshkit?.probableCause).toEqual(CREATE_ENVIRONMENT_403.probableCause);
    expect(formatted.meshkit?.longDescription).toEqual(CREATE_ENVIRONMENT_403.longDescription);
  });

  it('still recovers the details when a producer emits snake_case', async () => {
    respondWith(409, {
      error: 'Unable to create the workspace',
      code: 'meshery-server-1454',
      probable_cause: ['A workspace with the same name already exists.'],
      suggested_remediation: ['Pick a different workspace name.'],
    });

    const formatted = formatApiError(await createEnvironmentError());

    expect(formatted.message).toContain('**Unable to create the workspace**');
    expect(formatted.message).toContain('- Pick a different workspace name.');
    expect(formatted.meshkit?.probableCause).toEqual([
      'A workspace with the same name already exists.',
    ]);
  });

  it('degrades to a single-line message when the body is not a MeshKit envelope', async () => {
    respondWith(500, { message: 'internal server error' });

    const formatted = formatApiError(
      await createEnvironmentError(),
      'Failed to create environment',
    );

    expect(formatted.meshkit).toBeUndefined();
    expect(formatted.message).toBe('internal server error');
  });
});
