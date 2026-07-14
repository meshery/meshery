import { useCallback, useEffect, useRef, useState } from 'react';
import { ConsoleSocket, type ConsoleClosed, type ConsoleStatus } from './consoleSocket';
import type { ConsoleCapabilities } from './protocol';

export interface UseConsoleSocketOptions {
  /** The console URL, from `consoleSocketUrl`. A null url keeps the socket shut. */
  url: string | null;
  /** Receives every binary frame. Kept in a ref, so it may change each render. */
  onData?: (chunk: Uint8Array) => void;
  /** Set false to hold the socket closed, e.g. while a panel is hidden. */
  enabled?: boolean;
}

export interface UseConsoleSocket {
  /** `open` means the server has attached to the target, not merely that TCP is up. */
  status: ConsoleStatus;
  /** Resolved by the server's `ready` frame; null until then. */
  capabilities: ConsoleCapabilities | null;
  /** Why the console ended; null while it is live. */
  closed: ConsoleClosed | null;
  send: (bytes: Uint8Array) => void;
  resize: (cols: number, rows: number) => void;
  close: () => void;
}

/**
 * Owns one console socket for the life of `url`.
 *
 * Changing the url — a different container, a different tail depth — tears the
 * old console down and opens a new one, which is the honest behaviour: those
 * parameters are fixed when the remote stream is created and cannot be
 * renegotiated on a live socket.
 */
export const useConsoleSocket = ({
  url,
  onData,
  enabled = true,
}: UseConsoleSocketOptions): UseConsoleSocket => {
  const [status, setStatus] = useState<ConsoleStatus>('connecting');
  const [capabilities, setCapabilities] = useState<ConsoleCapabilities | null>(null);
  const [closed, setClosed] = useState<ConsoleClosed | null>(null);

  const socketRef = useRef<ConsoleSocket | null>(null);

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
    // its successor's state, reporting a dead console over a working shell.
    let current = true;

    const socket = new ConsoleSocket(url, {
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
