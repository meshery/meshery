import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import {
  ArticleIcon,
  Box,
  CloseIcon,
  IconButton,
  Tab,
  TerminalIcon,
  Tooltip,
} from '@sistent/sistent';
import type { SessionKind, SessionTarget } from 'lib/sessions/protocol';
import SessionPanel from './SessionPanel';
import SessionsShell, { type SessionsPanelMode } from './SessionsShell';
import { bindDock, publishDockState } from './dock-store';
import { SessionTabs, TabLabel, TabTitle } from './session-styles';

export interface OpenSessionRequest {
  connectionId: string;
  target: SessionTarget;
  kind: SessionKind;
}

interface OpenSession extends OpenSessionRequest {
  id: string;
  label: string;
}

export interface SessionsContextValue {
  /** Opens a session in the shared drawer, or focuses it if already open. */
  openSession: (request: OpenSessionRequest) => void;
  /** Closes a session, tearing its socket down. */
  closeSession: (request: OpenSessionRequest) => void;
  /**
   * Whether a session for this exact target is already open. Callers use it to
   * render an indicator on a resource that already has a shell attached.
   */
  hasSession: (request: OpenSessionRequest) => boolean;
  /** Every session currently open in the drawer. */
  getOpenSessions: () => OpenSessionRequest[];
  /**
   * Observes the set of open sessions. The listener fires immediately with the
   * current set and again on every change; the returned function unsubscribes.
   *
   * This exists so a non-React consumer — Kanvas mirrors it into its own store to
   * badge graph nodes that already have a session — can stay in sync with the
   * drawer as the source of truth, rather than tracking what it opened and
   * silently drifting when the user closes a session from the drawer itself.
   */
  subscribe: (listener: (open: OpenSessionRequest[]) => void) => () => void;
}

const SessionsContext = createContext<SessionsContextValue | null>(null);

const NO_SESSIONS: SessionsContextValue = {
  openSession: () => {},
  closeSession: () => {},
  hasSession: () => false,
  getOpenSessions: () => [],
  subscribe: () => () => {},
};

/**
 * The imperative handle on the shared session drawer.
 *
 * This is the API surface remote components (the Kanvas operator view) use to
 * launch a shell or a log tail from a graph node, so it is exported through
 * `remote-component.config.js` alongside the components themselves.
 *
 * Outside a provider it degrades to no-ops rather than throwing, so a component
 * that merely *offers* a session action can be dropped anywhere without the
 * tree it lands in having to know about sessions.
 */
export const useSessions = (): SessionsContextValue => useContext(SessionsContext) ?? NO_SESSIONS;

/**
 * The mounted provider's handle, for callers that are not React components.
 *
 * Kanvas dispatches its context-menu actions from plain functions, not from a
 * component, so it cannot read the context. This module-scope handle is bound
 * while a provider is mounted and reset when it unmounts, so a stale reference
 * can never resurrect a dead drawer.
 */
let mountedSessions: SessionsContextValue | null = null;

/** The open-session set, mirrored outside React so non-React callers can read it. */
let openSessionRequests: OpenSessionRequest[] = [];
const sessionListeners = new Set<(open: OpenSessionRequest[]) => void>();

const publishOpenSessions = (open: OpenSessionRequest[]) => {
  openSessionRequests = open;
  sessionListeners.forEach((listener) => listener(open));
};

/**
 * Subscribes to the open-session set. Defined at module scope rather than as a
 * closure over the provider, so both the imperative handle and the context value
 * hand out the same stable function — and so it does not collide with the
 * provider's own `sessions` state variable.
 */
const subscribeToSessions = (listener: (open: OpenSessionRequest[]) => void) => {
  sessionListeners.add(listener);
  // Fire immediately so a subscriber that attaches after sessions were opened
  // does not have to wait for the next change to learn about them.
  listener(openSessionRequests);
  return () => {
    sessionListeners.delete(listener);
  };
};

/**
 * Imperative session control for non-React callers. Falls back to no-ops when no
 * provider is mounted, so an extension loading before the app tree is ready
 * degrades rather than throwing.
 */
export const sessions: SessionsContextValue = {
  openSession: (request) => (mountedSessions ?? NO_SESSIONS).openSession(request),
  closeSession: (request) => (mountedSessions ?? NO_SESSIONS).closeSession(request),
  hasSession: (request) => (mountedSessions ?? NO_SESSIONS).hasSession(request),
  getOpenSessions: () => openSessionRequests,
  subscribe: subscribeToSessions,
};

/** Identifies a session, so re-requesting the same one focuses it rather than duplicating it. */
const sessionId = ({ connectionId, target, kind }: OpenSessionRequest): string =>
  [kind, connectionId, target.namespace ?? '', target.name, target.container ?? ''].join('/');

/**
 * The tab's label. Just the resource — the kind is carried by the tab's icon,
 * because a "Shell: " / "Logs: " prefix on every tab spent a third of the strip
 * repeating what the icon already says and pushed the strip into scroll.
 */
const sessionLabel = ({ target }: OpenSessionRequest): string =>
  target.container ? `${target.name}/${target.container}` : target.name;

/**
 * Hosts the sessions panel and hands its opener to the tree below.
 *
 * The panel docks to the foot of the content area, where it shares space with the
 * page rather than covering it, and detaches into a free-floating window that can
 * be dragged and resized anywhere over that area. Minimizing hides it and leaves
 * a badged control at the foot of the Navigator. The panel is mounted app-wide,
 * so a session survives navigation in any of those states.
 *
 * Every open session stays mounted for as long as it is open — including the
 * ones behind an inactive tab, and including while the panel is minimized. A
 * terminal is stateful on the remote end, so unmounting it to switch tabs, to
 * minimize, or to move between docked and floating would kill the shell. Hidden
 * panes are hidden with CSS, the shell changes style rather than structure, and
 * xterm refits from its ResizeObserver when a pane is shown or resized again.
 */
export const SessionsProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [sessions, setSessions] = useState<OpenSession[]>([]);
  const [activeId, setActiveId] = useState<string | null>(null);
  const [minimized, setMinimized] = useState(false);
  const [mode, setMode] = useState<SessionsPanelMode>('docked');

  const openSession = useCallback((request: OpenSessionRequest) => {
    const id = sessionId(request);
    setSessions((current) =>
      current.some((session) => session.id === id)
        ? current
        : [...current, { ...request, id, label: sessionLabel(request) }],
    );
    setActiveId(id);
    // Opening a session while minimized should show it, not silently queue it.
    setMinimized(false);
  }, []);

  const closeSessionById = useCallback((id: string) => {
    setSessions((current) => {
      const remaining = current.filter((session) => session.id !== id);
      setActiveId((active) => {
        if (active !== id) return active;
        return remaining.length ? remaining[remaining.length - 1].id : null;
      });
      return remaining;
    });
  }, []);

  const closeSession = useCallback(
    (request: OpenSessionRequest) => closeSessionById(sessionId(request)),
    [closeSessionById],
  );

  // Reads `sessions` and so must not be memoized against an empty dep list, or
  // callers would see a stale answer after a session opened.
  const hasSession = useCallback(
    (request: OpenSessionRequest) => {
      const id = sessionId(request);
      return sessions.some((session) => session.id === id);
    },
    [sessions],
  );

  const closeAll = useCallback(() => {
    setSessions([]);
    setActiveId(null);
    setMinimized(false);
  }, []);

  const openRequests = useMemo(
    () => sessions.map(({ connectionId, target, kind }) => ({ connectionId, target, kind })),
    [sessions],
  );

  const getOpenSessions = useCallback(() => openRequests, [openRequests]);

  const value = useMemo(
    () => ({
      openSession,
      closeSession,
      hasSession,
      getOpenSessions,
      subscribe: subscribeToSessions,
    }),
    [openSession, closeSession, hasSession, getOpenSessions],
  );

  // Publish to the non-React mirror whenever the set changes.
  useEffect(() => {
    publishOpenSessions(openRequests);
  }, [openRequests]);

  // Feed the Navigator's restore control, which sits outside this context.
  useEffect(() => {
    publishDockState({ count: sessions.length, minimized });
  }, [sessions.length, minimized]);

  useEffect(() => bindDock(() => setMinimized(false)), []);

  // Bind the imperative handle to *this* provider for as long as it is mounted.
  // The guard on unbind matters under StrictMode's double-mount, where the new
  // provider binds before the old one's cleanup runs.
  useEffect(() => {
    mountedSessions = value;
    return () => {
      if (mountedSessions === value) mountedSessions = null;
    };
  }, [value]);
  const open = sessions.length > 0;

  return (
    <SessionsContext.Provider value={value}>
      {children}

      <SessionsShell
        open={open}
        minimized={minimized}
        mode={mode}
        onModeChange={setMode}
        onMinimize={() => setMinimized(true)}
        onClose={closeAll}
        title="Sessions"
        tabs={
          <SessionTabs
            value={activeId ?? false}
            onChange={(_, id) => setActiveId(id)}
            variant="scrollable"
            scrollButtons="auto"
            allowScrollButtonsMobile
          >
            {sessions.map((session) => (
              <Tab
                key={session.id}
                value={session.id}
                disableRipple
                icon={
                  session.kind === 'terminal' ? (
                    <TerminalIcon width={15} height={15} />
                  ) : (
                    <ArticleIcon width={15} height={15} />
                  )
                }
                iconPosition="start"
                label={
                  <TabLabel>
                    <Tooltip title={session.label}>
                      <TabTitle>{session.label}</TabTitle>
                    </Tooltip>
                    <IconButton
                      size="small"
                      component="span"
                      aria-label={`Close ${session.label}`}
                      onClick={(event) => {
                        // The close control lives inside the tab, so stop the
                        // click from also selecting the tab being closed.
                        event.stopPropagation();
                        closeSessionById(session.id);
                      }}
                    >
                      <CloseIcon width={12} height={12} />
                    </IconButton>
                  </TabLabel>
                }
              />
            ))}
          </SessionTabs>
        }
      >
        {sessions.map((session) => (
          <Box
            key={session.id}
            sx={{
              flex: 1,
              minHeight: 0,
              // Hidden rather than unmounted: see the note on this component.
              display: session.id === activeId ? 'flex' : 'none',
              flexDirection: 'column',
              paddingTop: '0.5rem',
            }}
          >
            <SessionPanel
              connectionId={session.connectionId}
              target={session.target}
              kind={session.kind}
              // Every session stays `active` — and so keeps its socket — while
              // only the one on show owns the header's controls.
              focused={session.id === activeId}
            />
          </Box>
        ))}
      </SessionsShell>
    </SessionsContext.Provider>
  );
};

export default SessionsProvider;
