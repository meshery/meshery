import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';

// Hand-rolled EventSource mock. We can't import the real one in jsdom and
// we don't want to ship eventsource-polyfill just for tests, so we install
// a stub on globalThis and track every instance created during the test.
//
// The mock implements just enough of the EventSource surface that sseClient
// touches:
//   - constructor(url, init)
//   - addEventListener(type, cb)
//   - close()
//   - helper methods on the mock instance to dispatch open/message/error
//     synthetically from the test side.
//
// Multiple sseSubscribe calls will create multiple instances; we keep them
// all in a registry so a test can interrogate "the second EventSource"
// after a rebind.
interface FakeEventSourceInit {
  withCredentials?: boolean;
}

type ListenerMap = Record<string, Array<(ev: Event) => void>>;

class FakeEventSource {
  static instances: FakeEventSource[] = [];
  url: string;
  withCredentials: boolean;
  readyState = 0; // CONNECTING
  closed = false;
  private listeners: ListenerMap = {};

  constructor(url: string, init?: FakeEventSourceInit) {
    this.url = url;
    this.withCredentials = !!init?.withCredentials;
    FakeEventSource.instances.push(this);
  }

  addEventListener(type: string, cb: (ev: Event) => void): void {
    (this.listeners[type] ||= []).push(cb);
  }

  close(): void {
    this.closed = true;
    this.readyState = 2; // CLOSED
  }

  // --- test-side dispatchers --------------------------------------------
  dispatchOpen(): void {
    this.readyState = 1; // OPEN
    for (const cb of this.listeners.open ?? []) cb(new Event('open'));
  }

  dispatchMessage(name: string, data: string): void {
    const ev = new MessageEvent(name, { data });
    for (const cb of this.listeners[name] ?? []) cb(ev);
  }

  dispatchError(): void {
    for (const cb of this.listeners.error ?? []) cb(new Event('error'));
  }

  static reset(): void {
    FakeEventSource.instances = [];
  }
}

// Install BEFORE importing sseClient — the module under test doesn't
// import the global eagerly (it constructs new EventSource() inside
// sseSubscribe), so this assignment is picked up at call time. Vitest's
// hoisting of vi.mock doesn't help here because EventSource is a global,
// not a module.
(globalThis as unknown as { EventSource: typeof FakeEventSource }).EventSource = FakeEventSource;

import { sseSubscribe } from '../sseClient';
import { _resetForTests, _activeCountForTests, on as onWsStatus } from '../wsStatus';

describe('sseSubscribe', () => {
  beforeEach(() => {
    FakeEventSource.reset();
    _resetForTests();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('opens an EventSource with withCredentials and the encoded URL', () => {
    sseSubscribe({
      path: '/api/stream',
      params: { id: 'abc', tags: ['x', 'y'] },
      onMessage: () => {},
    });

    expect(FakeEventSource.instances).toHaveLength(1);
    const es = FakeEventSource.instances[0]!;
    expect(es.withCredentials).toBe(true);
    expect(es.url).toBe('/api/stream?id=abc&tags=x&tags=y');
  });

  it('drops undefined params from the URL', () => {
    sseSubscribe({
      path: '/api/stream',
      params: { id: 'abc', skip: undefined },
      onMessage: () => {},
    });
    expect(FakeEventSource.instances[0]!.url).toBe('/api/stream?id=abc');
  });

  it('preserves an existing query string when appending params', () => {
    sseSubscribe({
      path: '/api/stream?fixed=1',
      params: { extra: 'two' },
      onMessage: () => {},
    });
    expect(FakeEventSource.instances[0]!.url).toBe('/api/stream?fixed=1&extra=two');
  });

  it('calls onMessage with JSON-parsed payload when a default message arrives', () => {
    const onMessage = vi.fn();
    sseSubscribe({ path: '/api/stream', onMessage });
    const es = FakeEventSource.instances[0]!;
    es.dispatchMessage('message', JSON.stringify({ hello: 'world' }));
    expect(onMessage).toHaveBeenCalledWith({ hello: 'world' });
  });

  it('honors a custom eventName instead of the default "message"', () => {
    const onMessage = vi.fn();
    sseSubscribe({ path: '/api/stream', eventName: 'tick', onMessage });
    const es = FakeEventSource.instances[0]!;
    es.dispatchMessage('message', JSON.stringify({ wrong: true }));
    expect(onMessage).not.toHaveBeenCalled();
    es.dispatchMessage('tick', JSON.stringify({ ok: true }));
    expect(onMessage).toHaveBeenCalledWith({ ok: true });
  });

  it('logs and swallows JSON parse failures (does not throw)', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    const onMessage = vi.fn();
    sseSubscribe({ path: '/api/stream', onMessage, subscriptionName: 'Bad' });
    const es = FakeEventSource.instances[0]!;
    es.dispatchMessage('message', '{not json');
    expect(onMessage).not.toHaveBeenCalled();
    expect(consoleSpy).toHaveBeenCalledWith(
      expect.stringContaining('[SSE: Bad] failed to parse message:'),
      expect.any(Error),
    );
  });

  it('dispose() closes the EventSource and is idempotent', () => {
    const sub = sseSubscribe({ path: '/api/stream', onMessage: () => {} });
    const es = FakeEventSource.instances[0]!;
    expect(es.closed).toBe(false);
    sub.dispose();
    expect(es.closed).toBe(true);
    // Idempotent — calling dispose again should not blow up and should not
    // create a new EventSource.
    sub.dispose();
    expect(FakeEventSource.instances).toHaveLength(1);
  });

  it('rebind() closes the old EventSource and opens a new one with new params', () => {
    const onMessage = vi.fn();
    const sub = sseSubscribe({
      path: '/api/stream',
      params: { ctx: 'a' },
      onMessage,
    });
    const first = FakeEventSource.instances[0]!;
    expect(first.url).toBe('/api/stream?ctx=a');

    sub.rebind({ ctx: 'b' });

    expect(first.closed).toBe(true);
    expect(FakeEventSource.instances).toHaveLength(2);
    const second = FakeEventSource.instances[1]!;
    expect(second.url).toBe('/api/stream?ctx=b');

    // Handlers carry over: a message on the second EventSource invokes the
    // same onMessage closure.
    second.dispatchMessage('message', JSON.stringify({ from: 'second' }));
    expect(onMessage).toHaveBeenCalledWith({ from: 'second' });
  });

  it('rebind() after dispose() is a no-op', () => {
    const sub = sseSubscribe({ path: '/api/stream', onMessage: () => {} });
    sub.dispose();
    sub.rebind({ next: '1' });
    expect(FakeEventSource.instances).toHaveLength(1);
  });

  it('default onError logs to console without throwing', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    sseSubscribe({
      path: '/api/stream',
      subscriptionName: 'Default',
      onMessage: () => {},
    });
    const es = FakeEventSource.instances[0]!;
    expect(() => es.dispatchError()).not.toThrow();
    expect(consoleSpy).toHaveBeenCalledWith(
      expect.stringContaining('[SSE: Default] Error:'),
      expect.any(Event),
    );
  });

  it('user-supplied onError replaces the default handler', () => {
    const onError = vi.fn();
    sseSubscribe({
      path: '/api/stream',
      onMessage: () => {},
      onError,
    });
    const es = FakeEventSource.instances[0]!;
    es.dispatchError();
    expect(onError).toHaveBeenCalled();
  });

  it('updates wsStatus active count on open and dispose', () => {
    const sub = sseSubscribe({ path: '/api/stream', onMessage: () => {} });
    expect(_activeCountForTests()).toBe(0);
    FakeEventSource.instances[0]!.dispatchOpen();
    expect(_activeCountForTests()).toBe(1);
    sub.dispose();
    expect(_activeCountForTests()).toBe(0);
  });

  it('emits wsStatus connected and closed across the 0->1->0 edge', () => {
    const connected = vi.fn();
    const closed = vi.fn();
    onWsStatus('connected', connected);
    onWsStatus('closed', closed);

    const sub = sseSubscribe({ path: '/api/stream', onMessage: () => {} });
    const es = FakeEventSource.instances[0]!;
    es.dispatchOpen();
    expect(connected).toHaveBeenCalledTimes(1);
    expect(closed).not.toHaveBeenCalled();
    sub.dispose();
    expect(closed).toHaveBeenCalledTimes(1);
  });

  it('marks the subscription closed on error so the next open re-fires connected', () => {
    // Regression for the case where EventSource fires error -> open (the
    // browser reconnected) but the wrapper had left its internal countedOpen
    // flag set, so it never observed the 0->1 edge and the connection
    // indicator stayed stuck in the "reconnecting" state forever.
    const connected = vi.fn();
    const closed = vi.fn();
    const errorSpy = vi.fn();
    onWsStatus('connected', connected);
    onWsStatus('closed', closed);
    onWsStatus('error', errorSpy);
    // Suppress the default console.error from the SSE error handler so the
    // test output stays clean.
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    sseSubscribe({ path: '/api/stream', onMessage: () => {} });
    const es = FakeEventSource.instances[0]!;
    es.dispatchOpen();
    expect(connected).toHaveBeenCalledTimes(1);

    // The transport hiccups: error fires, the EventSource enters its
    // built-in reconnect.
    es.dispatchError();
    expect(closed).toHaveBeenCalledTimes(1);
    expect(errorSpy).toHaveBeenCalledTimes(1);
    expect(_activeCountForTests()).toBe(0);

    // Then the browser reconnects and re-fires open on the same instance.
    es.dispatchOpen();
    expect(connected).toHaveBeenCalledTimes(2);
    expect(_activeCountForTests()).toBe(1);

    consoleSpy.mockRestore();
  });

  it('drops messages that arrive after dispose()', () => {
    // The browser can deliver buffered MessageEvents after close() returns
    // but before the underlying connection is fully torn down. onMessage
    // must not fire in that window or callers see callbacks after dispose,
    // which is a footgun for components that have unmounted.
    const onMessage = vi.fn();
    const sub = sseSubscribe({ path: '/api/stream', onMessage });
    const es = FakeEventSource.instances[0]!;
    es.dispatchOpen();
    sub.dispose();
    es.dispatchMessage('message', '{"after":"dispose"}');
    expect(onMessage).not.toHaveBeenCalled();
  });
});
