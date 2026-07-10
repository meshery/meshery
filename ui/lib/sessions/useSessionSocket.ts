import { useCallback, useEffect, useRef, useState } from 'react';
import { SessionSocket, type SessionClosed, type SessionStatus } from './sessionSocket';
import type { SessionCapabilities } from './protocol';

export interface UseSessionSocketOptions {
  /** The session URL, from `sessionSocketUrl`. A null url keeps the socket shut. */
  url: string | null;
  /** Receives every binary frame. Kept in a ref, so it may change each render. */
  onData?: (chunk: Uint8Array) => void;
  /** Set false to hold the socket closed, e.g. while a panel is hidden. */
  enabled?: boolean;
}

export interface UseSessionSocket {
  /** `open` means the server has attached to the target, not merely that TCP is up. */
  status: SessionStatus;
  /** Resolved by the server's `ready` frame; null until then. */
  capabilities: SessionCapabilities | null;
  /** Why the session ended; null while it is live. */
  closed: SessionClosed | null;
  send: (bytes: Uint8Array) => void;
  resize: (cols: number, rows: number) => void;
  close: () => void;
}

/**
 * Owns one session socket for the life of `url`.
 *
 * Changing the url — a different container, a different tail depth — tears the
 * old session down and opens a new one, which is the honest behaviour: those
 * parameters are fixed when the remote stream is created and cannot be
 * renegotiated on a live socket.
 */
export const useSessionSocket = ({
  url,
  onData,
  enabled = true,
}: UseSessionSocketOptions): UseSessionSocket => {
  const [status, setStatus] = useState<SessionStatus>('connecting');
  const [capabilities, setCapabilities] = useState<SessionCapabilities | null>(null);
  const [closed, setClosed] = useState<SessionClosed | null>(null);

  const socketRef = useRef<SessionSocket | null>(null);

  // Held in a ref so a caller may pass an inline callback without forcing the
  // socket to be torn down and reopened on every render.
  const onDataRef = useRef(onData);
  onDataRef.current = onData;

  useEffect(() => {
    if (!enabled || !url) {
      return undefined;
    }

    setStatus('connecting');
    setCapabilities(null);
    setClosed(null);

    // A socket this effect has torn down must not speak for the one that
    // replaced it. Its `close` event arrives asynchronously — after the next
    // effect has already opened a live socket — so under React's StrictMode
    // double-mount the discarded socket's onClose would otherwise land on top of
    // its successor's state, reporting a dead session over a working shell.
    let current = true;

    const socket = new SessionSocket(url, {
      onReady: (caps) => {
        if (!current) return;
        setCapabilities(caps);
        setStatus('open');
      },
      onData: (chunk) => {
        if (!current) return;
        onDataRef.current?.(chunk);
      },
      onClose: (info) => {
        if (!current) return;
        setStatus('closed');
        setClosed(info);
      },
    });
    socketRef.current = socket;

    return () => {
      current = false;
      socket.close();
      socketRef.current = null;
    };
  }, [url, enabled]);

  const send = useCallback((bytes: Uint8Array) => socketRef.current?.send(bytes), []);
  const resize = useCallback(
    (cols: number, rows: number) => socketRef.current?.resize(cols, rows),
    [],
  );
  const close = useCallback(() => socketRef.current?.close(), []);

  return { status, capabilities, closed, send, resize, close };
};
