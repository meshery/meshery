import { describe, expect, it, vi, afterEach, afterAll } from 'vitest';
import { configureStore } from '@reduxjs/toolkit';

// Pin the absolute base every rtk-query suite uses before the client loads, and
// restore it afterwards so suite order cannot start to matter.
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

// view.ts pulls in @/utils/utils (SVG-in-JS transitively) and, through
// ../utils, the app store. Mock the small surface it uses.
vi.mock('@/utils/utils', () => ({
  urlEncodeArrayParam: (key: string, array: unknown) =>
    Array.isArray(array) ? array.map((item) => `${key}=${item}`).join('&') : '',
  urlEncodeParams: (params: Record<string, unknown>) =>
    new URLSearchParams(
      Object.entries(params)
        .filter(([, v]) => v != null)
        .map(([k, v]) => [k, String(v)]),
    ).toString(),
}));
vi.mock('../../store', () => ({ store: { dispatch: vi.fn() } }));

import { api } from '../index';
import { viewsApi } from '../view';

const VIEW_URL = /\/api\/extensions\/api\/content\/views\/([^/?]+)$/;
const LIST_URL = '/api/extensions/api/content/views?';
const PAGE_SIZE = 2;

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), {
    status,
    headers: { 'content-type': 'application/json' },
  });

// A tiny in-memory remote provider: PUT and DELETE change what later GETs
// return, so the assertions below are about what the cache ends up showing.
const setup = () => {
  const server: Record<string, Record<string, unknown>> = {};
  ['v-1', 'v-2', 'v-3', 'v-4'].forEach((id) => {
    server[id] = { id, name: id, visibility: 'public', metadata: { notes: 'old' } };
  });

  const fetchMock = vi.fn(async (request: Request) => {
    const url = new URL(request.url);
    const match = url.pathname.match(VIEW_URL);
    if (match && request.method === 'PUT') {
      Object.assign(server[match[1]], await request.json());
      return json(server[match[1]]);
    }
    if (match && request.method === 'DELETE') {
      delete server[match[1]];
      return json({});
    }
    if (match) {
      return server[match[1]] ? json(server[match[1]]) : json({ error: 'not found' }, 404);
    }
    const page = Number(url.searchParams.get('page'));
    const all = Object.values(server);
    return json({
      page,
      pageSize: PAGE_SIZE,
      totalCount: all.length,
      views: all.slice(page * PAGE_SIZE, (page + 1) * PAGE_SIZE),
    });
  });
  vi.stubGlobal('fetch', fetchMock);

  const store = configureStore({
    reducer: { [api.reducerPath]: api.reducer },
    middleware: (getDefault) => getDefault().concat(api.middleware),
  });

  const requests = (method: string, path: string | RegExp) =>
    fetchMock.mock.calls.filter(([request]) => {
      const { pathname, search } = new URL(request.url);
      const target = pathname + (search ? '?' : '');
      return (
        request.method === method &&
        (typeof path === 'string' ? target.endsWith(path) : path.test(target))
      );
    }).length;

  const getViewData = (viewId: string) =>
    viewsApi.endpoints.getView.select({ viewId })(store.getState()).data as
      | Record<string, unknown>
      | undefined;

  return { store, requests, getViewData };
};

// Let invalidation-driven refetches dispatched by the middleware settle.
const settle = () => new Promise((resolve) => setTimeout(resolve, 0));

describe('view mutations keep the cached views fresh', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('refetches the open view after a visibility change, and only that view', async () => {
    const { store, requests, getViewData } = setup();
    const open = store.dispatch(viewsApi.endpoints.getView.initiate({ viewId: 'v-1' }));
    const other = store.dispatch(viewsApi.endpoints.getView.initiate({ viewId: 'v-2' }));
    await Promise.all([open, other]);

    await store
      .dispatch(
        viewsApi.endpoints.updateViewVisibility.initiate({
          id: 'v-1',
          body: { visibility: 'private' },
        }),
      )
      .unwrap();
    await settle();

    expect(requests('GET', '/views/v-1')).toBe(2);
    expect(getViewData('v-1')?.visibility).toBe('private');
    expect(requests('GET', '/views/v-2')).toBe(1);

    open.unsubscribe();
    other.unsubscribe();
  });

  it('refetches the open view after its notes are saved', async () => {
    const { store, requests, getViewData } = setup();
    const open = store.dispatch(viewsApi.endpoints.getView.initiate({ viewId: 'v-1' }));
    await open;

    await store
      .dispatch(
        viewsApi.endpoints.updateViewVisibility.initiate({
          id: 'v-1',
          body: { metadata: { notes: 'new' } },
        }),
      )
      .unwrap();
    await settle();

    expect(requests('GET', '/views/v-1')).toBe(2);
    expect(getViewData('v-1')?.metadata).toEqual({ notes: 'new' });

    open.unsubscribe();
  });

  it('stops serving a deleted view from the cache', async () => {
    const { store, requests } = setup();
    const before = store.dispatch(viewsApi.endpoints.getView.initiate({ viewId: 'v-2' }));
    await before;
    before.unsubscribe();

    await store.dispatch(viewsApi.endpoints.deleteView.initiate({ id: 'v-2' })).unwrap();
    await settle();

    const after = store.dispatch(viewsApi.endpoints.getView.initiate({ viewId: 'v-2' }));
    const result = await after;

    expect(requests('GET', '/views/v-2')).toBe(2);
    expect(result.isError).toBe(true);

    after.unsubscribe();
  });

  it('leaves the paged views list to its caller so a refetch cannot duplicate rows', async () => {
    // fetchViews keeps every page in one merged cache entry. If a view
    // mutation invalidated it, RTK would refetch the current page (page 1
    // here) and the merge would append that page a second time.
    const { store, requests } = setup();
    const list = store.dispatch(
      viewsApi.endpoints.fetchViews.initiate({ page: 0, pagesize: PAGE_SIZE }),
    );
    await list;
    const nextPage = store.dispatch(
      viewsApi.endpoints.fetchViews.initiate({ page: 1, pagesize: PAGE_SIZE }),
    );
    await nextPage;
    const listRequests = requests('GET', LIST_URL);

    await store
      .dispatch(
        viewsApi.endpoints.updateViewVisibility.initiate({
          id: 'v-3',
          body: { visibility: 'private' },
        }),
      )
      .unwrap();
    await settle();

    const views = (
      viewsApi.endpoints.fetchViews.select({})(store.getState()).data as {
        views: { id: string }[];
      }
    ).views;
    expect(requests('GET', LIST_URL)).toBe(listRequests);
    expect(views.map((view) => view.id)).toEqual(['v-1', 'v-2', 'v-3', 'v-4']);

    list.unsubscribe();
    nextPage.unsubscribe();
  });
});
