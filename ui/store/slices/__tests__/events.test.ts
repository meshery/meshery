import { describe, it, expect } from 'vitest';
import eventsReducer, {
  pushEvent,
  setEvents,
  toggleNotificationCenter,
  closeNotificationCenter,
  openNotificationCenter,
  isSerializable,
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

  it('pushEvent orders revisions at the sub-millisecond precision the server sends', () => {
    let state = eventsReducer(undefined, { type: 'unknown' });
    state = eventsReducer(
      state,
      pushEvent(makeEvent({ occurrenceCount: 3, updatedAt: '2025-01-01T07:48:29.932900Z' })),
    );
    // Both times fall within one millisecond, which Date.parse cannot tell apart.
    state = eventsReducer(
      state,
      pushEvent(makeEvent({ occurrenceCount: 2, updatedAt: '2025-01-01T07:48:29.932100Z' })),
    );
    expect(state.entities['evt-1'].occurrenceCount).toBe(3);

    state = eventsReducer(
      state,
      pushEvent(makeEvent({ occurrenceCount: 4, updatedAt: '2025-01-01T07:48:29.932900001Z' })),
    );
    expect(state.entities['evt-1'].occurrenceCount).toBe(4);
  });

  it('pushEvent compares update times across zones and fraction lengths', () => {
    let state = eventsReducer(undefined, { type: 'unknown' });
    state = eventsReducer(
      state,
      pushEvent(makeEvent({ occurrenceCount: 3, updatedAt: '2025-01-01T09:48:29.932900+02:00' })),
    );
    state = eventsReducer(
      state,
      pushEvent(makeEvent({ occurrenceCount: 2, updatedAt: '2025-01-01T07:48:29.9321Z' })),
    );
    expect(state.entities['evt-1'].occurrenceCount).toBe(3);
  });

  it('pushEvent breaks a tie between equal update times on lastSeenAt, then occurrenceCount', () => {
    const updatedAt = '2025-01-01T07:48:29.932900Z';
    let state = eventsReducer(undefined, { type: 'unknown' });
    state = eventsReducer(
      state,
      pushEvent(
        makeEvent({ occurrenceCount: 3, updatedAt, lastSeenAt: '2025-01-01T07:48:29.000000002Z' }),
      ),
    );
    state = eventsReducer(
      state,
      pushEvent(
        makeEvent({ occurrenceCount: 5, updatedAt, lastSeenAt: '2025-01-01T07:48:29.000000001Z' }),
      ),
    );
    expect(state.entities['evt-1'].occurrenceCount).toBe(3);

    state = eventsReducer(
      state,
      pushEvent(
        makeEvent({ occurrenceCount: 2, updatedAt, lastSeenAt: '2025-01-01T07:48:29.000000002Z' }),
      ),
    );
    expect(state.entities['evt-1'].occurrenceCount).toBe(3);

    state = eventsReducer(
      state,
      pushEvent(
        makeEvent({ occurrenceCount: 4, updatedAt, lastSeenAt: '2025-01-01T07:48:29.000000002Z' }),
      ),
    );
    expect(state.entities['evt-1'].occurrenceCount).toBe(4);
  });

  it('pushEvent accepts a revision without updatedAt even when lastSeenAt or occurrenceCount is lower', () => {
    let state = eventsReducer(undefined, { type: 'unknown' });
    state = eventsReducer(
      state,
      pushEvent(
        makeEvent({
          occurrenceCount: 5,
          updatedAt: '2025-01-01T07:48:29.932900Z',
          lastSeenAt: '2025-01-01T07:48:29.000000002Z',
        }),
      ),
    );
    state = eventsReducer(
      state,
      pushEvent(makeEvent({ occurrenceCount: 2, lastSeenAt: '2025-01-01T07:48:29.000000001Z' })),
    );
    expect(state.entities['evt-1'].occurrenceCount).toBe(2);

    state = eventsReducer(
      state,
      pushEvent(
        makeEvent({
          occurrenceCount: 1,
          updatedAt: 'not a time',
          lastSeenAt: '2025-01-01T07:48:29.000000000Z',
        }),
      ),
    );
    expect(state.entities['evt-1'].occurrenceCount).toBe(1);
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

describe('isSerializable', () => {
  it('accepts primitives, plain objects and arrays', () => {
    expect(isSerializable('text')).toBe(true);
    expect(isSerializable(42)).toBe(true);
    expect(isSerializable(true)).toBe(true);
    expect(isSerializable(null)).toBe(true);
    expect(isSerializable(undefined)).toBe(true);
    expect(isSerializable({ title: 'X', nested: { page: 0, flags: [true, 'a'] } })).toBe(true);
    expect(isSerializable([1, 'a', { ok: true }])).toBe(true);
    expect(isSerializable(Object.create(null))).toBe(true);
  });

  it('rejects functions, symbols and bigints', () => {
    expect(isSerializable(() => null)).toBe(false);
    expect(isSerializable(Symbol('s'))).toBe(false);
    expect(isSerializable(10n)).toBe(false);
  });

  it('rejects Date, Map, Set and class instances', () => {
    class Custom {}
    expect(isSerializable(new Date())).toBe(false);
    expect(isSerializable(new Map())).toBe(false);
    expect(isSerializable(new Set())).toBe(false);
    expect(isSerializable(new Custom())).toBe(false);
    expect(isSerializable({ at: new Date() })).toBe(false);
    expect(isSerializable([new Map()])).toBe(false);
  });

  it('rejects circular objects', () => {
    const value = { title: 'X' };
    value.self = value;
    expect(isSerializable(value)).toBe(false);
    const list = [];
    list.push(list);
    expect(isSerializable(list)).toBe(false);
  });

  it('accepts a shared reference that is not circular', () => {
    const shared = { page: 0 };
    expect(isSerializable({ a: shared, b: shared })).toBe(true);
  });
});

describe('openNotificationCenter', () => {
  const open = (ui) => {
    const initial = eventsReducer(undefined, { type: 'unknown' });
    return eventsReducer(initial, openNotificationCenter({ ui }));
  };

  it('drops an icon passed as a component and keeps serializable siblings', () => {
    const state = open({ icon: () => null, title: 'Kanvas' });
    expect(state.ui.icon).toBeUndefined();
    expect(state.ui.title).toBe('Kanvas');
    expect(state.isNotificationCenterOpen).toBe(true);
  });

  it('keeps an icon passed as a name', () => {
    const state = open({ icon: 'alert', title: 'Alerts' });
    expect(state.ui.icon).toBe('alert');
    expect(state.ui.title).toBe('Alerts');
  });

  it('drops values holding a Date, Map, Set, class instance, circular object or function', () => {
    class Custom {}
    const circular = { title: 'loop' };
    circular.self = circular;
    const state = open({
      at: new Date(),
      lookup: new Map(),
      tags: new Set(),
      custom: new Custom(),
      circular,
      onClose: () => {},
      title: 'Kept',
    });
    expect(state.ui.at).toBeUndefined();
    expect(state.ui.lookup).toBeUndefined();
    expect(state.ui.tags).toBeUndefined();
    expect(state.ui.custom).toBeUndefined();
    expect(state.ui.circular).toBeUndefined();
    expect(state.ui.onClose).toBeUndefined();
    expect(state.ui.title).toBe('Kept');
    expect(() => JSON.stringify(state.ui)).not.toThrow();
  });

  it('screens a raw action object, which is how extensions dispatch it', () => {
    const initial = eventsReducer(undefined, { type: 'unknown' });
    const state = eventsReducer(initial, {
      type: 'events/openNotificationCenter',
      payload: { ui: { icon: () => null, title: 'Extension' } },
    });
    expect(state.ui.icon).toBeUndefined();
    expect(state.ui.title).toBe('Extension');
  });

  it('preserves the defaults the payload does not override', () => {
    const state = open({ icon: 'bell' });
    expect(state.ui.title).toBe('Notifications');
    expect(state.ui.empty_message).toBe('No notifications found');
  });

  it('drops a view_to_open that is not serializable', () => {
    const initial = eventsReducer(undefined, { type: 'unknown' });
    const state = eventsReducer(
      initial,
      openNotificationCenter({
        view_to_open: { page: 1, at: new Date(), onSelect: () => {} },
      }),
    );
    expect(state.isNotificationCenterOpen).toBe(true);
    expect(state.view_to_fetch_on_open).toEqual(initial.view_to_fetch_on_open);
    expect(state.current_view).toEqual(initial.current_view);
  });

  it('applies a serializable view_to_open', () => {
    const initial = eventsReducer(undefined, { type: 'unknown' });
    const view = { page: 2, filters: { status: 'read' } };
    const state = eventsReducer(initial, openNotificationCenter({ view_to_open: view }));
    expect(state.view_to_fetch_on_open).toEqual(view);
    expect(state.current_view).toEqual(view);
  });

  it('applies a view_to_open given as a name', () => {
    const initial = eventsReducer(undefined, { type: 'unknown' });
    const state = eventsReducer(initial, openNotificationCenter({ view_to_open: 'history' }));
    expect(state.view_to_fetch_on_open).toBe('history');
    expect(state.current_view).toBe('history');
  });

  it('closeNotificationCenter resets ui to the defaults', () => {
    const opened = open({ icon: 'error', title: 'Errors' });
    const closed = eventsReducer(opened, closeNotificationCenter());
    expect(closed.ui.icon).toBeUndefined();
    expect(closed.ui.title).toBe('Notifications');
  });
});
