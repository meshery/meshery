import { describe, it, expect } from 'vitest';
import eventsReducer, {
  pushEvent,
  setEvents,
  toggleNotificationCenter,
  closeNotificationCenter,
  openNotificationCenter,
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
});

// Extensions dispatch openNotificationCenter into this store over the event bus
// (store/index.ts), so its payload is not under this repo's control. A React
// icon component used to arrive in `ui` and sit in Redux state, which is what
// RTK's serializable-state check flagged at `events.ui.icon` (issue #16873).

const hasNonSerializable = (value) => {
  const type = typeof value;
  if (value === null || type === 'string' || type === 'number' || type === 'boolean') {
    return false;
  }
  if (type !== 'object') return true;
  return Object.values(value).some(hasNonSerializable);
};

describe('openNotificationCenter keeps the store serializable', () => {
  const open = (ui) => {
    const initial = eventsReducer(undefined, { type: 'unknown' });
    return eventsReducer(initial, openNotificationCenter({ ui }));
  };

  it('drops an icon passed as a React component', () => {
    const IconComponent = () => null;
    const state = open({ icon: IconComponent, title: 'Kanvas' });

    expect(state.ui.icon).toBeUndefined();
    expect(typeof state.ui.icon).not.toBe('function');
    // the serializable siblings still land
    expect(state.ui.title).toBe('Kanvas');
    expect(state.isNotificationCenterOpen).toBe(true);
  });

  it('keeps an icon passed as a registry name', () => {
    const state = open({ icon: 'alert', title: 'Alerts' });

    expect(state.ui.icon).toBe('alert');
    expect(state.ui.title).toBe('Alerts');
  });

  it('drops a React element and other non-serializable values', () => {
    const element = { $$typeof: Symbol.for('react.element'), type: 'svg', props: {} };
    const state = open({ icon: element, onClick: () => {}, history_mode: true });

    expect(state.ui.icon).toBeUndefined();
    expect(state.ui.onClick).toBeUndefined();
    expect(state.ui.history_mode).toBe(true);
  });

  it('leaves no non-serializable value anywhere under events.ui', () => {
    const state = open({ icon: () => null, nested: { deep: () => null }, title: 'X' });

    expect(hasNonSerializable(state.ui)).toBe(false);
    expect(() => JSON.stringify(state.ui)).not.toThrow();
  });

  it('preserves the defaults an extension does not override', () => {
    const state = open({ icon: 'bell' });

    expect(state.ui.title).toBe('Notifications');
    expect(state.ui.empty_message).toBe('No notifications found');
  });

  it('screens the raw action object too, which is how extensions dispatch it', () => {
    // store/index.ts forwards whatever an extension puts on the event bus
    // straight to dispatch, so the action never goes through the creator.
    const initial = eventsReducer(undefined, { type: 'unknown' });
    const state = eventsReducer(initial, {
      type: 'events/openNotificationCenter',
      payload: { ui: { icon: () => null, title: 'Extension' } },
    });

    expect(state.ui.icon).toBeUndefined();
    expect(state.ui.title).toBe('Extension');
  });

  it('closeNotificationCenter resets ui back to the serializable defaults', () => {
    const opened = open({ icon: 'error', title: 'Errors' });
    const closed = eventsReducer(opened, closeNotificationCenter());

    expect(closed.ui.icon).toBeUndefined();
    expect(closed.ui.title).toBe('Notifications');
  });
});
