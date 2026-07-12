import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Box, CircularProgress, Typography } from '@sistent/sistent';
import { sessionSocketUrl, type SessionTarget } from 'lib/sessions/protocol';
import { useSessionSocket } from 'lib/sessions/useSessionSocket';
import { useXTerm } from './useXTerm';
import SessionControls from './SessionControls';
import {
  SessionSurface,
  SessionToolbar,
  StatusText,
  TerminalMount,
  ToolbarEnd,
} from './session-styles';

export interface TerminalSessionProps {
  connectionId: string;
  target: SessionTarget;
  /** Containers the user may attach to; from the capabilities endpoint. */
  containers?: string[];
  /** Rendered at the right of the toolbar, e.g. the host's fullscreen control. */
  toolbarEnd?: React.ReactNode;
  /** Held false while the pane is hidden, so a background tab holds no shell. */
  active?: boolean;
  /** Whether this session is the one on show, and so owns the panel's header. */
  focused?: boolean;
}

const encoder = new TextEncoder();

/**
 * An interactive shell in the target container.
 *
 * The socket is opened once per (connection, target) pair and closed when the
 * pane unmounts or goes inactive. There is no reconnect: the remote shell dies
 * with its socket, and quietly replacing it would give the user a fresh shell
 * still showing the dead one's scrollback.
 */
const TerminalSession: React.FC<TerminalSessionProps> = ({
  connectionId,
  target,
  containers = [],
  toolbarEnd,
  active = true,
  focused = true,
}) => {
  const { resource, namespace, name } = target;

  // Switching container opens a shell in a different process namespace, so it is
  // a new session, not a reconfiguration of the running one.
  const [container, setContainer] = useState(target.container ?? '');
  const [query, setQuery] = useState('');

  // Depend on the target's fields, not on the object: callers pass an object
  // literal, and a new identity each render would tear the shell down and open
  // a new one on every parent re-render.
  const url = useMemo(
    () =>
      active
        ? sessionSocketUrl(connectionId, 'terminal', {
            resource,
            namespace,
            name,
            container: container || undefined,
          })
        : null,
    [active, connectionId, resource, namespace, name, container],
  );

  // The socket is created first so its stable `send`/`resize` can back the
  // terminal's input handlers. `write` is likewise stable, so passing it as the
  // socket's data sink below does not churn the socket.
  const terminalRef = useRef<((chunk: Uint8Array) => void) | null>(null);
  const onData = useCallback((chunk: Uint8Array) => terminalRef.current?.(chunk), []);

  const { status, capabilities, closed, send, resize } = useSessionSocket({ url, onData });

  const handleInput = useCallback((data: string) => send(encoder.encode(data)), [send]);
  const handleResize = useCallback(
    ({ cols, rows }: { cols: number; rows: number }) => resize(cols, rows),
    [resize],
  );

  const { containerRef, write, writeln, fit, getGeometry, search, ready } = useXTerm({
    onInput: handleInput,
    onResize: handleResize,
    scrollback: 5000,
  });
  terminalRef.current = write;

  // The remote pty is created at a default 80x24, so it must be told the real
  // geometry once the session attaches. fit() alone is not enough: it only
  // emits onResize when the geometry actually changes, and the terminal was
  // already fitted when it was created. Send the geometry explicitly.
  useEffect(() => {
    if (!ready || status !== 'open') return;
    fit();
    const geometry = getGeometry();
    if (geometry) {
      resize(geometry.cols, geometry.rows);
    }
  }, [ready, status, fit, getGeometry, resize]);

  // Report the outcome inside the terminal, where the user is already looking,
  // not only in a status bar they may have scrolled past.
  useEffect(() => {
    if (!closed || !ready) return;
    if (closed.graceful) {
      writeln(`\r\n\x1b[2m— session ended (exit code ${closed.exitCode ?? 0}) —\x1b[0m`);
    } else {
      writeln(`\r\n\x1b[31m— session failed: ${closed.message ?? 'connection lost'} —\x1b[0m`);
    }
  }, [closed, ready, writeln]);

  const attachedContainer = capabilities?.defaultContainer ?? container;

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%', minHeight: 0 }}>
      <SessionToolbar>
        {/* Container and search render in the panel's header; see SessionControls. */}
        <SessionControls
          containers={containers}
          container={container}
          onContainerChange={setContainer}
          query={query}
          onQueryChange={setQuery}
          onSearch={search}
          focused={focused}
        />

        <StatusText $error={Boolean(closed && !closed.graceful)}>
          {closed && !closed.graceful
            ? `${closed.message ?? 'connection lost'}${closed.code ? ` (${closed.code})` : ''}`
            : status === 'open' && attachedContainer
              ? `Attached to ${attachedContainer}`
              : status === 'closed'
                ? `Exited${closed?.exitCode !== undefined ? ` (${closed.exitCode})` : ''}`
                : 'Connecting…'}
        </StatusText>

        <ToolbarEnd>{toolbarEnd}</ToolbarEnd>
      </SessionToolbar>

      <SessionSurface>
        <TerminalMount ref={containerRef} />
        {status === 'connecting' && (
          <Box
            sx={{
              position: 'absolute',
              inset: 0,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              gap: '0.75rem',
            }}
          >
            <CircularProgress size={20} />
            <Typography variant="body2">Opening a shell…</Typography>
          </Box>
        )}
      </SessionSurface>
    </Box>
  );
};

export default TerminalSession;
