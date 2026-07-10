/**
 * Low-level client for a session WebSocket.
 *
 * Deliberately transport-only: it frames and unframes the protocol in
 * `protocol.ts` and reports lifecycle, and knows nothing about terminals, logs,
 * or React. That keeps it usable by a future session kind, and testable without
 * a DOM.
 *
 * It does not reconnect. A terminal is stateful on the remote end — a dropped
 * socket means the shell is gone, and silently opening a new one would hand the
 * user a fresh shell wearing the old one's scrollback. Reconnection, where it
 * makes sense at all, is the caller's decision.
 */

import {
  CLOSE_SESSION_ENDED,
  CLOSE_SESSION_FAILED,
  type ControlMessage,
  type SessionCapabilities,
} from './protocol';

export type SessionStatus = 'connecting' | 'open' | 'closed';

/** Why a session ended. */
export interface SessionClosed {
  /** True when the session ended as intended: command exited, or logs hit EOF. */
  graceful: boolean;
  /** The remote command's status, for a terminal that ran to completion. */
  exitCode?: number;
  /** MeshKit error code, when the server reported a failure. */
  code?: string;
  /** Human-readable failure description. */
  message?: string;
}

export interface SessionSocketHandlers {
  /** Fired once, when the server has attached to the target. */
  onReady?: (capabilities: SessionCapabilities) => void;
  /** Fired for every binary frame: stdout or log bytes. */
  onData?: (chunk: Uint8Array) => void;
  /** Fired once, when the session ends for any reason. */
  onClose?: (info: SessionClosed) => void;
}

export class SessionSocket {
  private socket: WebSocket;
  private handlers: SessionSocketHandlers;
  /** The last terminal control frame, which explains the close that follows it. */
  private finalMessage?: ControlMessage;
  private closeReported = false;
  /** Whether this end asked for the close, rather than suffering one. */
  private closedByClient = false;

  constructor(url: string, handlers: SessionSocketHandlers = {}) {
    this.handlers = handlers;
    this.socket = new WebSocket(url);
    // Binary frames arrive as ArrayBuffer rather than Blob, so output can reach
    // the terminal synchronously instead of via an async Blob read.
    this.socket.binaryType = 'arraybuffer';

    this.socket.onmessage = this.onMessage;
    this.socket.onclose = this.onSocketClose;
    this.socket.onerror = this.onSocketError;
  }

  get status(): SessionStatus {
    switch (this.socket.readyState) {
      case WebSocket.CONNECTING:
        return 'connecting';
      case WebSocket.OPEN:
        return 'open';
      default:
        return 'closed';
    }
  }

  /** Sends stdin bytes. Ignored unless the socket is open. */
  send(bytes: Uint8Array): void {
    if (this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(bytes);
    }
  }

  /** Reports a new terminal geometry to the remote end. */
  resize(cols: number, rows: number): void {
    if (this.socket.readyState !== WebSocket.OPEN) return;
    if (!Number.isFinite(cols) || !Number.isFinite(rows) || cols <= 0 || rows <= 0) return;
    this.socket.send(JSON.stringify({ type: 'resize', cols, rows }));
  }

  /** Closes the session. Safe to call more than once. */
  close(): void {
    // Recorded before the call, because closing a socket that is still
    // CONNECTING cannot send a close frame: the handshake never finished, so the
    // browser reports code 1006 (abnormal) and the close is indistinguishable
    // from a dropped connection unless we remember we asked for it.
    this.closedByClient = true;

    if (
      this.socket.readyState === WebSocket.OPEN ||
      this.socket.readyState === WebSocket.CONNECTING
    ) {
      this.socket.close(1000, 'client closed the session');
    }
  }

  private onMessage = (event: MessageEvent): void => {
    if (typeof event.data === 'string') {
      this.onControl(event.data);
      return;
    }
    if (event.data instanceof ArrayBuffer) {
      this.handlers.onData?.(new Uint8Array(event.data));
    }
  };

  private onControl(payload: string): void {
    let message: ControlMessage;
    try {
      message = JSON.parse(payload);
    } catch {
      // A frame we cannot parse is not worth killing a live session over.
      return;
    }

    switch (message.type) {
      case 'ready':
        this.handlers.onReady?.(message.capabilities ?? { terminal: false, logs: false });
        return;
      case 'error':
      case 'exit':
      case 'eof':
        // Remember it; the close frame that follows is what fires onClose, so
        // there is exactly one close notification per session.
        this.finalMessage = message;
        return;
      default:
        return;
    }
  }

  private onSocketError = (): void => {
    // `error` is always followed by `close`, which is where the session's
    // outcome is reported. Reporting here too would double-notify.
  };

  private onSocketClose = (event: CloseEvent): void => {
    if (this.closeReported) return;
    this.closeReported = true;

    this.handlers.onClose?.(this.describeClose(event));
  };

  /**
   * Explains a close, preferring the server's terminal control frame over the
   * close code — the frame carries an error code and a message; the code alone
   * carries neither.
   */
  private describeClose(event: CloseEvent): SessionClosed {
    const final = this.finalMessage;

    if (final?.type === 'error') {
      return { graceful: false, code: final.code, message: final.message };
    }
    if (final?.type === 'exit') {
      return { graceful: true, exitCode: final.exitCode ?? 0 };
    }
    if (final?.type === 'eof') {
      return { graceful: true };
    }

    // No terminal frame arrived: the socket died under us. Our own close() is
    // graceful whatever code the browser attaches to it; anything else is a
    // broken connection.
    if (
      this.closedByClient ||
      event.code === CLOSE_SESSION_ENDED ||
      event.code === 1000 ||
      event.code === 1005
    ) {
      return { graceful: true };
    }
    return {
      graceful: false,
      message:
        event.code === CLOSE_SESSION_FAILED && event.reason
          ? event.reason
          : 'The connection to the session was lost.',
    };
  }
}
