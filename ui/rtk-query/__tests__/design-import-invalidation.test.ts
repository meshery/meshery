import { describe, expect, it, vi, beforeEach, afterEach, afterAll } from 'vitest';
import { configureStore } from '@reduxjs/toolkit';

// The schemas client reads its base URL from the environment at module-load
// time, and the app relies on a relative base in the browser. Under the test
// runner a relative base has no origin to resolve against, so pin an absolute
// one before the client is imported (vi.hoisted runs ahead of the imports).
//
// Use the same value every other rtk-query suite uses, and restore whatever was
// there afterwards: test files sharing a worker process share `process.env`, so
// a divergent value left behind here is a way for suite order to start mattering.
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

// design.ts pulls in @/utils/utils + @/utils/multi-ctx, which transitively
// import non-test-safe modules (SVG-in-JS). Mock the small surface it uses.
vi.mock('@/utils/utils', () => ({
  urlEncodeParams: (params: Record<string, unknown>) => {
    const sp = new URLSearchParams();
    Object.entries(params).forEach(([k, v]) => {
      if (v != null) sp.append(k, String(v));
    });
    return sp.toString();
  },
}));
vi.mock('@/utils/multi-ctx', () => ({ ctxUrl: (url: string) => url }));
const { dispatch } = vi.hoisted(() => ({ dispatch: vi.fn() }));
vi.mock('../../store', () => ({ store: { dispatch } }));

import { api } from '../index';
import { designsApi } from '../design';
import { appendInvalidatesTags } from '../utils';

// Nothing in the generated client provides the schemas-side `Design_designs`
// tag, so this stands in for any future consumer of it: it is the only way to
// observe from the outside that the enhancement kept the generated tag.
const schemasTaggedApi = api.injectEndpoints({
  endpoints: (builder) => ({
    designsTaggedBySchemas: builder.query({
      query: () => '/api/test/schemas-tagged-designs',
      providesTags: ['Design_designs'],
    }),
  }),
});

const jsonResponse = (body: unknown) =>
  Promise.resolve(
    new Response(JSON.stringify(body), {
      status: 200,
      headers: { 'content-type': 'application/json' },
    }),
  );

// fetchBaseQuery calls fetch() with a Request object, so both the URL and the
// method come off argument 0 rather than an init bag.
const urlOf = (call: unknown[]) =>
  typeof call[0] === 'string' ? call[0] : (call[0] as Request).url;

const methodOf = (call: unknown[]) =>
  typeof call[0] === 'string'
    ? ((call[1] as RequestInit)?.method ?? 'GET')
    : (call[0] as Request).method;

/**
 * `design.ts` attaches the local `designs` cache tag to the schemas-generated
 * `importDesign` endpoint through the callback form of `enhanceEndpoints`, which
 * is handed the live definition from `getEndpointDefinition(context, name)`. That
 * lookup has no fallback, so a missing endpoint fails at module load instead of
 * enhancing a throwaway object the way the object form
 * (`Object.assign(getEndpointDefinition(...) || {}, partial)`) silently would.
 * These tests pin the three parts of that wiring: that `importDesign` is the
 * generated endpoint and no local duplicate has come back, that an import
 * invalidates the local `designs` tag the design lists provide, and that the tags
 * schemas itself declares are still invalidated alongside it.
 */
describe('importDesign — schemas-generated endpoint with local cache tag', () => {
  const setup = () => {
    const fetchMock = vi.fn((input: RequestInfo | URL) => {
      const url = typeof input === 'string' ? input : (input as Request).url;
      if (url.includes('/api/pattern/import')) {
        return jsonResponse([{ id: 'design-1', name: 'Imported design' }]);
      }
      return jsonResponse({ patterns: [], total_count: 0 });
    });
    vi.stubGlobal('fetch', fetchMock);

    const store = configureStore({
      reducer: { [api.reducerPath]: api.reducer },
      middleware: (getDefault) => getDefault().concat(api.middleware),
    });

    return { fetchMock, store };
  };

  const importDesign = (store: ReturnType<typeof setup>['store']) =>
    store
      .dispatch(
        designsApi.endpoints.importDesign.initiate({
          body: { url: 'https://example.com/design.yaml', name: 'Imported design' },
        }),
      )
      .unwrap();

  beforeEach(() => {
    vi.restoreAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
  });

  it('is the schemas-generated endpoint, never a locally re-declared one', () => {
    // `enhanceEndpoints` only ever mutates an existing definition, so a defined
    // `importDesign` is the generated one design.ts enhanced.
    expect(api.endpoints.importDesign).toBeDefined();
    // The local duplicate that used to declare POST /api/pattern/import.
    expect(designsApi.endpoints.importPattern).toBeUndefined();
  });

  it('names the endpoint and the package if schemas stops generating it', () => {
    const attachTag = appendInvalidatesTags('importDesign', { type: 'designs' });
    expect(() => attachTag(undefined)).toThrow(/importDesign/);
    expect(() => attachTag(undefined)).toThrow(/@meshery\/schemas/);
    expect(() => attachTag(undefined)).toThrow(/renamed or removed/);
  });

  it('POSTs the generated /api/pattern/import request and refetches the designs list', async () => {
    const { fetchMock, store } = setup();

    // Subscribe to the designs list so RTK will act on the invalidation.
    const listSub = store.dispatch(
      designsApi.endpoints.getPatterns.initiate({ page: 0, pagesize: 10 }),
    );
    await listSub;
    const fetchesAfterList = fetchMock.mock.calls.length;

    await importDesign(store);

    const importCall = fetchMock.mock.calls.find((call) =>
      urlOf(call).includes('/api/pattern/import'),
    );
    expect(importCall, 'import request was issued').toBeDefined();
    expect(methodOf(importCall!)).toBe('POST');

    // Let the invalidation-driven refetch settle.
    await new Promise((resolve) => setTimeout(resolve, 0));

    const listRefetches = fetchMock.mock.calls
      .slice(fetchesAfterList)
      .filter((call) => urlOf(call).includes('/api/pattern?'));
    expect(listRefetches.length, 'designs list refetched after import').toBeGreaterThan(0);

    listSub.unsubscribe();
  });

  it('keeps the tags the generated definition declares, not just the local one', async () => {
    const { fetchMock, store } = setup();

    const schemasTagSub = store.dispatch(
      schemasTaggedApi.endpoints.designsTaggedBySchemas.initiate(undefined),
    );
    await schemasTagSub;
    const fetchesAfterSubscribe = fetchMock.mock.calls.length;

    await importDesign(store);
    await new Promise((resolve) => setTimeout(resolve, 0));

    const refetches = fetchMock.mock.calls
      .slice(fetchesAfterSubscribe)
      .filter((call) => urlOf(call).includes('/api/test/schemas-tagged-designs'));
    expect(refetches.length, 'Design_designs consumer refetched after import').toBeGreaterThan(0);

    schemasTagSub.unsubscribe();
  });
});
