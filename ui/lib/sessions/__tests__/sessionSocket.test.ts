import { beforeEach, describe, expect, it, vi } from 'vitest';
import { SessionSocket, type SessionClosed } from '../sessionSocket';

/** Minimal stand-in for the browser's WebSocket, driven by the test. */
class FakeWebSocket {
  static CONNECTING = 0;
  static OPEN = 1;
  static CLOSING = 2;
  static CLOSED = 3;

  static instances: FakeWebSocket[] = [];

  readyState = FakeWebSocket.CONNECTING;
  binaryType = '';
  onmessage: ((event: MessageEvent) => void) | null = null;
  onclose: ((event: CloseEvent) => void) | null = null;
  onerror: (() => void) | null = null;

  closeCalls: Array<{ code?: number; reason?: string }> = [];

  constructor(public url: string) {
    FakeWebSocket.instances.push(this);
  }

  close(code?: number, reason?: string) {
    this.closeCalls.push({ code, reason });
    this.readyState = FakeWebSocket.CLOSED;
  }

  /** Fires the close event the browser would deliver, asynchronously in reality. */
  emitClose(code: number, reason = '') {
    this.readyState = FakeWebSocket.CLOSED;
    this.onclose?.({ code, reason } as CloseEvent);
  }

  emitText(payload: string) {
    this.onmessage?.({ data: payload } as MessageEvent);
  }
}

beforeEach(() => {
  FakeWebSocket.instances = [];
  vi.stubGlobal('WebSocket', FakeWebSocket);
});

const openSocket = (handlers: { onClose?: (info: SessionClosed) => void } = {}) => {
  const socket = new SessionSocket('ws://localhost/session', handlers);
  return { socket, fake: FakeWebSocket.instances.at(-1)! };
};

describe('SessionSocket close reporting', () => {
  it('reports a client-initiated close of a CONNECTING socket as graceful', () => {
    // The regression: a socket torn down before its handshake completes cannot
    // send a close frame, so the browser reports 1006 (abnormal). Without
    // remembering that we asked for the close, that reads as "connection lost"
    // and gets painted into the terminal of the socket that replaced it.
    let closed: SessionClosed | undefined;
    const { socket, fake } = openSocket({ onClose: (info) => (closed = info) });

    expect(fake.readyState).toBe(FakeWebSocket.CONNECTING);
    socket.close();
    fake.emitClose(1006);

    expect(closed).toEqual({ graceful: true });
  });

  it('still reports an unexpected close as a lost connection', () => {
    let closed: SessionClosed | undefined;
    const { fake } = openSocket({ onClose: (info) => (closed = info) });

    fake.emitClose(1006);

    expect(closed?.graceful).toBe(false);
    expect(closed?.message).toContain('connection to the session was lost');
  });

  it('prefers the server’s terminal control frame over the close code', () => {
    let closed: SessionClosed | undefined;
    const { fake } = openSocket({ onClose: (info) => (closed = info) });

    fake.emitText(JSON.stringify({ type: 'error', code: 'meshery-server-1447', message: 'boom' }));
    fake.emitClose(4001);

    expect(closed).toEqual({ graceful: false, code: 'meshery-server-1447', message: 'boom' });
  });

  it('reports a non-zero exit as graceful, carrying the status', () => {
    let closed: SessionClosed | undefined;
    const { fake } = openSocket({ onClose: (info) => (closed = info) });

    fake.emitText(JSON.stringify({ type: 'exit', exitCode: 130 }));
    fake.emitClose(4000);

    expect(closed).toEqual({ graceful: true, exitCode: 130 });
  });

  it('notifies of a close exactly once', () => {
    const onClose = vi.fn();
    const { fake } = openSocket({ onClose });

    fake.emitClose(1006);
    fake.emitClose(1006);

    expect(onClose).toHaveBeenCalledTimes(1);
  });
});
