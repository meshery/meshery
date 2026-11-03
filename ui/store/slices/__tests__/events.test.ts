import { describe, it, expect, vi } from 'vitest';
import eventsReducer, {
  pushEvent,
  pushEvents,
  setEvents,
  setLoadError,
  toggleNotificationCenter,
  closeNotificationCenter,
  loadEvents,
  loadNextPage,
  retryCurrentPage,
} from '../events';

const makeEvent = (overrides = {}) => ({
  id: 'evt-1',
  severity: 'info',
  status: 'unread',
  description: 'Test event',
  createdAt: '2025-01-01T00:00:00Z',
  ...overrides,
});

describe('events slice', () => {
  it('returns initial state', () => {
    const state = eventsReducer(undefined, { type: 'unknown' });
    expect(state.isNotificationCenterOpen).toBe(false);
    expect(state.current_view.page).toBe(0);
  });

  it('pushEvent adds a single event', () => {
    const initial = eventsReducer(undefined, { type: 'unknown' });
    const event = makeEvent();
    const state = eventsReducer(initial, pushEvent(event));
    expect(state.ids).toContain('evt-1');
    expect(state.entities['evt-1'].description).toBe('Test event');
  });

  it('pushEvent trims severity and status', () => {
    const initial = eventsReducer(undefined, { type: 'unknown' });
    const event = makeEvent({ severity: '  warning  ', status: '  read  ' });
    const state = eventsReducer(initial, pushEvent(event));
    expect(state.entities['evt-1'].severity).toBe('warning');
    expect(state.entities['evt-1'].status).toBe('read');
  });

  it('setEvents replaces all events', () => {
    let state = eventsReducer(undefined, { type: 'unknown' });
    state = eventsReducer(state, pushEvent(makeEvent({ id: 'old' })));
    state = eventsReducer(
      state,
      setEvents([makeEvent({ id: 'new-1' }), makeEvent({ id: 'new-2' })]),
    );
    expect(state.ids).not.toContain('old');
    expect(state.ids).toContain('new-1');
    expect(state.ids).toContain('new-2');
  });

  it('toggleNotificationCenter toggles open state', () => {
    let state = eventsReducer(undefined, { type: 'unknown' });
    expect(state.isNotificationCenterOpen).toBe(false);
    state = eventsReducer(state, toggleNotificationCenter());
    expect(state.isNotificationCenterOpen).toBe(true);
    state = eventsReducer(state, toggleNotificationCenter());
    expect(state.isNotificationCenterOpen).toBe(false);
  });

  it('closeNotificationCenter resets view state', () => {
    let state = eventsReducer(undefined, { type: 'unknown' });
    state = eventsReducer(state, toggleNotificationCenter());
    state = eventsReducer(state, closeNotificationCenter());
    expect(state.isNotificationCenterOpen).toBe(false);
    expect(state.current_view.page).toBe(0);
  });

  it('setLoadError flags the error without touching has_more', () => {
    let state = eventsReducer(undefined, { type: 'unknown' });
    expect(state.current_view.has_more).toBe(true);
    state = eventsReducer(state, setLoadError());
    expect(state.current_view.load_error).toBe(true);
    expect(state.current_view.has_more).toBe(true);
  });

  it('a successful push clears a previous load error', () => {
    let state = eventsReducer(undefined, { type: 'unknown' });
    state = eventsReducer(state, setLoadError());
    expect(state.current_view.load_error).toBe(true);
    state = eventsReducer(state, pushEvents([makeEvent()]));
    expect(state.current_view.load_error).toBe(false);
  });
});

describe('events thunks – failed page handling', () => {
  const makeState = (overrides = {}) => ({
    events: {
      current_view: { page: 1, filters: { status: 'unread' }, has_more: true, load_error: false },
      ...overrides,
    },
  });

  it('a failed request flags the error and does not mark the list exhausted', async () => {
    const dispatched: any[] = [];
    const dispatch = (action: any) =>
      typeof action === 'function' ? action(dispatch, getState) : (dispatched.push(action), action);
    const getState = () => makeState();
    const fetch = vi.fn().mockResolvedValue({ error: { status: 500 } });

    await loadEvents(fetch, 2, { status: 'unread' })(dispatch, getState);

    expect(fetch).toHaveBeenCalledWith({ page: 2, filters: { status: 'unread' } });
    // only setLoadError dispatched; no setCurrentView / pushEvents advancing the page
    expect(dispatched.map((a) => a.type)).toEqual(['events/setLoadError']);
    // the failed page + filters are recorded for retry
    expect(dispatched[0].payload).toEqual({ page: 2, filters: { status: 'unread' } });
  });

  it('loadNextPage does not advance while a load error is set', async () => {
    const dispatch = vi.fn();
    const getState = () => makeState({ current_view: { page: 1, filters: {}, load_error: true } });

    await loadNextPage((fetch) => fetch)(dispatch, getState);

    expect(dispatch).not.toHaveBeenCalled();
  });

  it('retryCurrentPage re-requests the failed page, not the last loaded page', async () => {
    const dispatched: any[] = [];
    const dispatch = (action: any) =>
      typeof action === 'function' ? action(dispatch, getState) : (dispatched.push(action), action);
    // current page is 3 (last loaded), but page 4 failed and is pending retry
    const getState = () =>
      makeState({
        current_view: {
          page: 3,
          filters: { status: 'read' },
          load_error: true,
          failed_request: { page: 4, filters: { status: 'unread' } },
        },
      });
    const fetch = vi.fn().mockResolvedValue({ data: { events: [makeEvent()] } });

    await retryCurrentPage(fetch)(dispatch, getState);

    // retries the failed page 4 with its filters, not current page 3
    expect(fetch).toHaveBeenCalledWith({ page: 4, filters: { status: 'unread' } });
  });

  it('retryCurrentPage is a no-op when no request is pending', async () => {
    const dispatch = vi.fn();
    const getState = () =>
      makeState({ current_view: { page: 2, filters: {}, failed_request: null } });

    await retryCurrentPage(vi.fn())(dispatch, getState);

    expect(dispatch).not.toHaveBeenCalled();
  });

  it('a successful empty final page still stops pagination', async () => {
    const dispatched: any[] = [];
    const dispatch = (action: any) =>
      typeof action === 'function' ? action(dispatch, getState) : (dispatched.push(action), action);
    const getState = () => makeState();
    const fetch = vi.fn().mockResolvedValue({ data: { events: [] } });

    await loadEvents(fetch, 2, {})(dispatch, getState);

    const pushed = dispatched.find((a) => a.type === 'events/pushEvents');
    expect(pushed).toBeTruthy();
    // empty page -> has_more false via reducer
    const state = eventsReducer(eventsReducer(undefined, { type: 'x' }), pushed);
    expect(state.current_view.has_more).toBe(false);
  });
});
