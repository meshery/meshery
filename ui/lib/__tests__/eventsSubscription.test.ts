import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { subscribeToEvents } from '../eventsSubscription';

class FakeEventSource {
  static CONNECTING = 0;
  static OPEN = 1;
  static CLOSED = 2;

  readyState = FakeEventSource.CONNECTING;
  onopen: (() => void) | null = null;
  onmessage: ((message: { data: string }) => void) | null = null;
  onerror: ((error: Event) => void) | null = null;
  closed = false;

  constructor() {
    streams.push(this);
  }

  close() {
    this.closed = true;
    this.readyState = FakeEventSource.CLOSED;
  }
}

const streams: FakeEventSource[] = [];

describe('subscribeToEvents', () => {
  beforeEach(() => {
    streams.length = 0;
    vi.stubGlobal('EventSource', FakeEventSource);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('invokes onOpen on the initial open and on each native reconnect', () => {
    const onOpen = vi.fn();
    const subscription = subscribeToEvents(vi.fn(), undefined, onOpen);
    expect(onOpen).not.toHaveBeenCalled();

    const source = streams[0];
    source.onopen?.();
    source.onopen?.();

    expect(onOpen).toHaveBeenCalledTimes(2);

    subscription.dispose();
    expect(source.closed).toBe(true);
  });

  it('still forwards stream messages after an open callback', () => {
    const onNext = vi.fn();
    const subscription = subscribeToEvents(onNext, undefined, vi.fn());
    const source = streams[0];

    source.onopen?.();
    source.onmessage?.({ data: '{"id":"evt-1"}' });

    expect(onNext).toHaveBeenCalledWith({ id: 'evt-1' });
    subscription.dispose();
  });
});
