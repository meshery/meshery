import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import {
  Box,
  Button,
  CachedIcon,
  CircularProgress,
  FormControl,
  InputLabel,
  MenuItem,
  PlayArrowIcon,
  Select,
  TimerIcon,
  ToggleButton,
  ToggleButtonGroup,
  Tooltip,
  Typography,
} from '@sistent/sistent';
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

export interface LogSessionProps {
  connectionId: string;
  target: SessionTarget;
  /** Containers the user may switch between; from the capabilities endpoint. */
  containers?: string[];
  /** Rendered at the right of the toolbar, e.g. the host's fullscreen control. */
  toolbarEnd?: React.ReactNode;
  active?: boolean;
  /** Whether this session is the one on show, and so owns the panel's header. */
  focused?: boolean;
}

/** Tail depths offered in the toolbar. */
const TAIL_OPTIONS = [100, 500, 1000, 5000];

/**
 * A log tail for the target container.
 *
 * Changing any stream parameter — container, follow, tail depth, previous
 * instance — rebuilds the URL, which tears down the old stream and opens a new
 * one. That is not a shortcut: the Kubernetes API fixes these when the log
 * stream is created and offers no way to renegotiate them on a live stream.
 */
const LogSession: React.FC<LogSessionProps> = ({
  connectionId,
  target,
  containers = [],
  toolbarEnd,
  active = true,
  focused = true,
}) => {
  const { resource, namespace, name } = target;

  const [container, setContainer] = useState(target.container ?? '');
  const [follow, setFollow] = useState(true);
  const [timestamps, setTimestamps] = useState(false);
  const [previous, setPrevious] = useState(false);
  const [tailLines, setTailLines] = useState(500);
  const [query, setQuery] = useState('');
  const [paused, setPaused] = useState(false);

  // Whether the next stream should replay history. It must not after a resume:
  // the pane already holds everything the previous stream sent, and re-tailing
  // would duplicate the last `tailLines` of it beneath itself.
  const [replayHistory, setReplayHistory] = useState(true);

  const terminalWriteRef = useRef<((chunk: Uint8Array) => void) | null>(null);
  const onData = useCallback((chunk: Uint8Array) => terminalWriteRef.current?.(chunk), []);

  // Identifies the stream the *user* asked for. A pause is not a new stream, so
  // it is deliberately absent here: only a change to one of these wipes the pane
  // and replays from the top.
  const streamKey = `${container}|${follow}|${timestamps}|${previous}|${tailLines}`;

  const url = useMemo(
    () =>
      active && !paused
        ? sessionSocketUrl(
            connectionId,
            'logs',
            { resource, namespace, name, container: container || undefined },
            // Resuming asks for zero history: pausing tore the upstream stream
            // down, so anything produced while paused is gone either way, and
            // replaying the tail would duplicate what is already on screen.
            { follow, timestamps, previous, tailLines: replayHistory ? tailLines : 0 },
          )
        : null,
    [
      active,
      paused,
      connectionId,
      resource,
      namespace,
      name,
      container,
      follow,
      timestamps,
      previous,
      tailLines,
      replayHistory,
    ],
  );

  const { status, closed } = useSessionSocket({ url, onData });

  const { containerRef, write, clear, search, ready } = useXTerm({
    readOnly: true,
    // Logs are the reason someone opens this pane; keep far more history than a
    // shell needs.
    scrollback: 20000,
  });
  terminalWriteRef.current = write;

  // A genuinely new stream replays history from the top, so the pane must start
  // empty or the replay would land beneath the previous stream's output. Keyed
  // on the user's parameters, not on the URL, so a pause/resume cycle — which
  // does change the URL — keeps its scrollback.
  useEffect(() => {
    if (!ready) return;
    clear();
    setReplayHistory(true);
  }, [ready, streamKey, clear]);

  const togglePaused = useCallback(() => {
    setPaused((current) => {
      // Closing the socket makes the server cancel the upstream log stream, so
      // this is real backpressure rather than a frozen viewport.
      if (!current) setReplayHistory(false);
      return !current;
    });
  }, []);

  const activeToggles = useMemo(
    () =>
      [follow && 'follow', timestamps && 'timestamps', previous && 'previous'].filter(
        Boolean,
      ) as string[],
    [follow, timestamps, previous],
  );

  const handleToggle = useCallback((_event: React.MouseEvent, next: string[]) => {
    setFollow(next.includes('follow'));
    setTimestamps(next.includes('timestamps'));
    setPrevious(next.includes('previous'));
  }, []);

  const statusText = () => {
    if (paused) return 'Paused — the upstream stream is closed';
    if (status === 'connecting') return 'Opening the log stream…';
    if (status === 'open') return follow ? 'Streaming…' : 'Loaded';
    if (closed && !closed.graceful) return closed.message ?? 'The log stream was lost.';
    return follow ? 'Stream ended' : 'End of logs';
  };

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

        <FormControl size="small" sx={{ minWidth: 120 }}>
          <InputLabel id="log-tail-label">Tail</InputLabel>
          <Select
            labelId="log-tail-label"
            label="Tail"
            value={tailLines}
            onChange={(event) => setTailLines(Number(event.target.value))}
          >
            {TAIL_OPTIONS.map((lines) => (
              <MenuItem key={lines} value={lines}>
                {lines} lines
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        {/*
         * Icon toggles rather than switch-plus-label: three labelled switches
         * pushed the search box and the pause control off the end of a panel
         * that the user can make narrow.
         */}
        {/*
         * A controlled multi-select group, not three standalone toggles:
         * ToggleButtonGroup clones its children and replaces their `onChange`
         * with its own, so per-button handlers would silently never fire.
         */}
        <ToggleButtonGroup size="small" value={activeToggles} onChange={handleToggle}>
          <ToggleButton value="follow" title="Follow the stream">
            <PlayArrowIcon width={16} height={16} />
          </ToggleButton>
          <ToggleButton value="timestamps" title="Show timestamps">
            <TimerIcon width={16} height={16} />
          </ToggleButton>
          <ToggleButton value="previous" title="Logs of the previous container instance">
            <CachedIcon width={16} height={16} />
          </ToggleButton>
        </ToggleButtonGroup>

        <Tooltip
          title={
            paused
              ? 'Resume streaming (output produced while paused is not backfilled)'
              : 'Pause: close the upstream log stream'
          }
        >
          {/* Pausing is meaningless without follow: there is no live stream to stop. */}
          <Box component="span" sx={{ display: 'inline-flex' }}>
            <Button size="small" variant="outlined" onClick={togglePaused} disabled={!follow}>
              {paused ? 'Resume' : 'Pause'}
            </Button>
          </Box>
        </Tooltip>
        <Button size="small" variant="outlined" onClick={clear}>
          Clear
        </Button>

        <StatusText $error={Boolean(closed && !closed.graceful)}>{statusText()}</StatusText>
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
            <Typography variant="body2">Opening the log stream…</Typography>
          </Box>
        )}
      </SessionSurface>
    </Box>
  );
};

export default LogSession;
