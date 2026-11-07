import { describe, it, expect } from 'vitest';
import eventsReducer, {
  pushEvent,
  setEvents,
  toggleNotificationCenter,
  closeNotificationCenter,
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

  it('pushEvent updates an event already shown in place', () => {
    let state = eventsReducer(undefined, { type: 'unknown' });
    state = eventsReducer(
      state,
      pushEvent(makeEvent({ severity: 'error', description: 'unreachable', occurrenceCount: 1 })),
    );
    state = eventsReducer(
      state,
      pushEvent(
        makeEvent({
          severity: 'success',
          description: 'reachable again',
          occurrenceCount: 2,
          resolvedAt: '2025-01-01T00:05:00Z',
        }),
      ),
    );
    expect(state.ids).toEqual(['evt-1']);
    expect(state.entities['evt-1']).toMatchObject({
      severity: 'success',
      description: 'reachable again',
      occurrenceCount: 2,
      resolvedAt: '2025-01-01T00:05:00Z',
    });
  });

  it('pushEvent keeps the read status of an event already shown', () => {
    let state = eventsReducer(undefined, { type: 'unknown' });
    state = eventsReducer(state, pushEvent(makeEvent({ status: 'read', occurrenceCount: 1 })));
    state = eventsReducer(state, pushEvent(makeEvent({ status: 'unread', occurrenceCount: 2 })));
    expect(state.entities['evt-1']).toMatchObject({ status: 'read', occurrenceCount: 2 });
  });

  it('pushEvent ignores an older revision delivered after a newer one', () => {
    let state = eventsReducer(undefined, { type: 'unknown' });
    state = eventsReducer(
      state,
      pushEvent(makeEvent({ occurrenceCount: 3, updatedAt: '2025-01-01T00:03:00Z' })),
    );
    state = eventsReducer(
      state,
      pushEvent(makeEvent({ occurrenceCount: 2, updatedAt: '2025-01-01T00:02:00Z' })),
    );
    expect(state.entities['evt-1'].occurrenceCount).toBe(3);

    state = eventsReducer(
      state,
      pushEvent(
        makeEvent({
          severity: 'success',
          resolvedAt: '2025-01-01T00:05:00Z',
          updatedAt: '2025-01-01T00:05:00Z',
        }),
      ),
    );
    // A delayed occurrence never reopens the resolved condition, even when it
    // carries no update time.
    state = eventsReducer(state, pushEvent(makeEvent({ severity: 'error', occurrenceCount: 4 })));
    expect(state.entities['evt-1']).toMatchObject({
      severity: 'success',
      resolvedAt: '2025-01-01T00:05:00Z',
    });
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
});
