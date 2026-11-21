import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { subscribeToEvents } from '../eventsSubscription';

// jsdom has no EventSource; this fake lets a test fire `open` like the browser.
class FakeEventSource {
  static CONNECTING = 0;
  static OPEN = 1;
  static CLOSED = 2;
  static latest: FakeEventSource;

  readyState = FakeEventSource.CONNECTING;
  onopen: (() => void) | null = null;
  onmessage: ((message: { data: string }) => void) | null = null;
  onerror: ((error: Event) => void) | null = null;
  close = vi.fn();

  constructor() {
    FakeEventSource.latest = this;
  }

  open() {
    this.readyState = FakeEventSource.OPEN;
    this.onopen?.();
  }
}

describe('subscribeToEvents onOpen', () => {
  beforeEach(() => {
    vi.stubGlobal('EventSource', FakeEventSource);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('calls onOpen each time the stream (re)connects', () => {
    const onOpen = vi.fn();
    subscribeToEvents(vi.fn(), vi.fn(), onOpen);

    FakeEventSource.latest.open();
    FakeEventSource.latest.open();

    expect(onOpen).toHaveBeenCalledTimes(2);
  });

  it('does not call onOpen after dispose', () => {
    const onOpen = vi.fn();
    const subscription = subscribeToEvents(vi.fn(), vi.fn(), onOpen);

    subscription.dispose();
    FakeEventSource.latest.open();

    expect(onOpen).not.toHaveBeenCalled();
  });

  it('works without an onOpen callback', () => {
    subscribeToEvents(vi.fn());

    expect(() => FakeEventSource.latest.open()).not.toThrow();
  });
});
