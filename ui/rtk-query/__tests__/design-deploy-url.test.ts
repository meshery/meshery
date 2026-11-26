import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { configureStore } from '@reduxjs/toolkit';

// The query string of the design deploy and undeploy requests, as the endpoints build it with
// the real ctxUrl. A design checked while no Kubernetes context is selected has no contexts to
// start the query string with.

vi.mock('@/utils/utils', () => ({ urlEncodeParams: () => '' }));
vi.mock('../../store', () => ({ store: { dispatch: vi.fn() } }));

import { designsApi } from '../design';

// fetchBaseQuery builds a Request from a relative URL, which Node refuses; this keeps the URL.
class RecordedRequest {
  url: string;
  method: string;
  headers: Headers;
  signal: AbortSignal | undefined;
  constructor(
    url: string,
    init: { method?: string; headers?: HeadersInit; signal?: AbortSignal } = {},
  ) {
    this.url = url;
    this.method = init.method ?? 'GET';
    this.headers = new Headers(init.headers);
    this.signal = init.signal;
  }
  clone() {
    return this;
  }
}

describe('design deploy and undeploy URLs', () => {
  let requests: string[];

  beforeEach(() => {
    requests = [];
    vi.stubGlobal('Request', RecordedRequest);
    vi.stubGlobal('fetch', async (req: RecordedRequest) => {
      requests.push(`${req.method} ${req.url}`);
      return new Response('{}', { status: 200, headers: { 'content-type': 'application/json' } });
    });
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  const send = async (
    endpoint: 'deployPattern' | 'undeployPattern',
    args: Record<string, unknown>,
  ) => {
    const store = configureStore({
      reducer: { [designsApi.reducerPath]: designsApi.reducer },
      middleware: (getDefault) => getDefault().concat(designsApi.middleware),
    });
    await store.dispatch(
      designsApi.endpoints[endpoint].initiate({
        patternFile: 'f',
        patternId: 'p',
        ...args,
      }) as never,
    );
    return requests;
  };

  it('starts the query string with the flags when no context is selected', async () => {
    expect(
      await send('deployPattern', { selectedK8sContexts: [], dryRun: true, skipCRD: true }),
    ).toEqual(['POST /api/pattern/deploy?dryRun=true&skipCRD=true']);
  });

  it('puts the flags after the contexts', async () => {
    expect(
      await send('deployPattern', { selectedK8sContexts: ['a', 'b'], verify: true, dryRun: true }),
    ).toEqual(['POST /api/pattern/deploy?contexts=a&contexts=b&verify=true&dryRun=true']);
  });

  it('adds no query string when there is nothing to add', async () => {
    expect(await send('deployPattern', { selectedK8sContexts: undefined })).toEqual([
      'POST /api/pattern/deploy',
    ]);
  });

  it('does the same for an undeploy', async () => {
    expect(await send('undeployPattern', { selectedK8sContexts: [], dryRun: true })).toEqual([
      'DELETE /api/pattern/deploy?dryRun=true',
    ]);
  });
});
